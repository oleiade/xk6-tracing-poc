package tracing

import (
	"fmt"
	"time"

	"github.com/dop251/goja"
	"go.k6.io/k6/js/common"
	"go.k6.io/k6/js/modules"
	"go.k6.io/k6/metrics"
)

// Tracing is the JS module instance that will be created for each VU.
type Tracing struct {
	vu modules.VU

	propagator Propagator
}

// InstrumentHTTP instruments the HTTP module with tracing headers.
//
// When used in the context of a k6 script, it will automatically replace
// the imported http module's methods with instrumented ones.
func (t *Tracing) InstrumentHTTP(options instrumentationOptions) {
	if err := t.configure(options); err != nil {
		common.Throw(t.vu.Runtime(), err)
	}

	// Explicitly inject the http module in the VU's runtime.
	// This allows us to later on override the http module's methods
	// with instrumented ones.
	httpModuleValue, err := t.vu.Runtime().RunString(`require('k6/http')`)
	if err != nil {
		common.Throw(t.vu.Runtime(), err)
	}

	httpModuleObj := httpModuleValue.ToObject(t.vu.Runtime())

	HTTPMethods := [...]k6HTTPMethodName{
		k6HTTPDeleteMethodName,
		k6HTTPGetMethodName,
		k6HTTPHeadMethodName,
		k6HTTPOptionsMethodName,
		k6HTTPPatchMethodName,
		k6HTTPPostMethodName,
		k6HTTPPutMethodName,
	}

	for _, method := range HTTPMethods {
		originalMethodFn, ok := goja.AssertFunction(httpModuleObj.Get(string(method)))
		if !ok {
			common.Throw(t.vu.Runtime(), fmt.Errorf("http.%s is not a function", method))
		}

		tracedMethodFn := t.instrumentHTTPMethod(method, originalMethodFn)

		// Inject the new get function, adding tracing headers
		// to the request in the HTTP module object.
		err = httpModuleObj.Set(string(method), tracedMethodFn)
		if err != nil {
			common.Throw(t.vu.Runtime(), err)
		}
	}

	// Inject the updated HTTP module object in the runtime,
	// overriding any previously imported one in the process.
	err = t.vu.Runtime().Set("http", httpModuleObj)
	if err != nil {
		common.Throw(t.vu.Runtime(), err)
	}
}

// configure configures the tracing module with the given options.
func (t *Tracing) configure(opts instrumentationOptions) error {
	switch opts.Propagator {
	case "w3c":
		t.propagator = &W3CPropagator{}
	case "b3":
		t.propagator = &B3Propagator{}
	case "jaeger":
		t.propagator = &JaegerPropagator{}
	default:
		return fmt.Errorf("unknown propagator: %s", opts.Propagator)
	}

	return nil
}

// instrumentationOptions are the options that can be passed to the
// tracing.instrument() method.
type instrumentationOptions struct {
	// Sampling is the sampling rate to use for the tracer.
	Sampling float64 `js:"sampling"`

	// Propagation is the propagation format to use for the tracer.
	Propagator string `js:"propagator"`

	// Baggage is a map of baggage items to add to the tracer.
	Baggage map[string]string `js:"baggage"`
}

// instrumentHTTPMethod returns a new function that wraps the original http
// method, adding tracing headers to the request.
//
// The function takes the HTTP method name as argument, as well as the exported
// method itself, as a goja.Callable.
//
// The produced resulting function is ready to be injected back into the http
// module, in place of the original method.
func (t *Tracing) instrumentHTTPMethod(methodName k6HTTPMethodName, methodFn goja.Callable) goja.Callable {
	return func(this goja.Value, args ...goja.Value) (goja.Value, error) {
		rt := t.vu.Runtime()

		// Ensure the arguments have a params object, in which
		// we can add the tracing headers at a later point in time.
		args, params, err := t.getOrCreateParams(methodName, args...)
		if err != nil {
			common.Throw(rt, fmt.Errorf("failed to normalize HTTP arguments: %w", err))
		}

		// Ensure that the params object contains a headers object.
		// Create it if it doesn't.
		headers, err := t.getOrCreateHeaders(params)
		if err != nil {
			common.Throw(rt, fmt.Errorf("failed to normalize HTTP headers: %w", err))
		}

		traceID := NewTraceID(k6Prefix, k6CloudCode, uint64(time.Now().UnixNano())/uint64(time.Millisecond))
		encodedTraceID, _, err := traceID.Encode()
		if err != nil {
			common.Throw(rt, fmt.Errorf("failed to encode trace ID: %w", err))
		}

		// Produce a trace header in the format defined by the
		// configured propagator.
		header, err := t.propagator.Propagate(encodedTraceID)
		if err != nil {
			common.Throw(rt, fmt.Errorf("failed to propagate trace ID: %w", err))
		}

		for key, value := range header {
			err = headers.Set(key, value)
			if err != nil {
				common.Throw(rt, err)
			}
		}

		vuState := t.vu.State()

		// Add the trace ID to the VU's state, so that it can be
		// used in the metrics emitted by the HTTP module.
		vuState.Tags.Modify(func(t *metrics.TagsAndMeta) {
			t.SetMetadata(metadataTraceIDKeyName, encodedTraceID)
		})

		// call the original http.get method, with overridden arguments
		args = append([]goja.Value{this}, args...)
		result, err := methodFn(goja.Undefined(), args...)
		if err != nil {
			common.Throw(rt, err)
		}

		// Remove the trace ID from the VU's state, so that it doesn't
		// leak into other requests.
		vuState.Tags.Modify(func(t *metrics.TagsAndMeta) {
			t.DeleteMetadata(metadataTraceIDKeyName)
		})

		return result, err
	}
}

// getOrCreateParams ensures that the HTTP method arguments list contains
// a params object. If it doesn't, it creates one.
//
// The method returns the normalized arguments list as well as its params
// object.
//
// Note that as of k6 v0.42.0 the HTTP API can turn out to be a bit inconsistent,
// as it can be called with 0, 1 or 2 arguments, and the second argument
// can be either a request's body, or a params object.
func (t *Tracing) getOrCreateParams(m k6HTTPMethodName, args ...goja.Value) ([]goja.Value, *goja.Object, error) {
	rt := t.vu.Runtime()
	params := rt.NewObject()

	switch len(args) {
	case 2:
		paramsValue := args[1]
		if !isNullish(paramsValue) {
			params = paramsValue.ToObject(rt)
			break
		}

		args[1] = params
	case 1:
		// The http.get and the http.head methods take an optional params
		// object as first argument. Whereas the other methods take a
		// request's body as first argument, and an optional params object
		// as second argument.
		if m != k6HTTPGetMethodName && m != k6HTTPHeadMethodName {
			args = append(args, params)
			break
		}

		paramsValue := args[0]
		if !isNullish(paramsValue) {
			params = paramsValue.ToObject(rt)
			break
		}

		args[0] = params
	case 0:
		args = []goja.Value{goja.Null(), params}
	default:
		return args, params, fmt.Errorf("unexpected number of arguments for http.%s method", m)
	}

	return args, params, nil
}

// getOrCreateHeaders ensures that a http method params object is properly
// formed, and has the expected properties set.
//
// This method modifies the params object in place, and returns the headers object.
func (t *Tracing) getOrCreateHeaders(params *goja.Object) (*goja.Object, error) {
	rt := t.vu.Runtime()

	headersValue := params.Get("headers")
	if !isNullish(headersValue) {
		return headersValue.ToObject(rt), nil
	}

	headers := rt.NewObject()
	if err := params.Set("headers", headers); err != nil {
		return nil, err
	}

	return headers, nil
}

// k6HTTPMethodName represents the name of a k6 HTTP method.
type k6HTTPMethodName string

// k6HTTPMethodName constants.
const (
	k6HTTPDeleteMethodName  k6HTTPMethodName = "del"
	k6HTTPGetMethodName     k6HTTPMethodName = "get"
	k6HTTPPostMethodName    k6HTTPMethodName = "post"
	k6HTTPPutMethodName     k6HTTPMethodName = "put"
	k6HTTPPatchMethodName   k6HTTPMethodName = "patch"
	k6HTTPHeadMethodName    k6HTTPMethodName = "head"
	k6HTTPOptionsMethodName k6HTTPMethodName = "options"
)

// HTTPMethods is a static list of all the k6 HTTP method names.
