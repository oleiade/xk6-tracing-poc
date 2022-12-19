package tracing

import (
	"fmt"
	"time"

	"github.com/dop251/goja"
	"go.k6.io/k6/js/common"
	"go.k6.io/k6/js/modules"
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

	httpModuleValue, err := t.vu.Runtime().RunString(`require('k6/http')`)
	if err != nil {
		common.Throw(t.vu.Runtime(), err)
	}

	HTTPMethods := []k6HTTPMethodName{
		k6HTTPDeleteMethodName,
		k6HTTPGetMethodName,
		k6HTTPHeadMethodName,
		k6HTTPOptionsMethodName,
		k6HTTPPatchMethodName,
		k6HTTPPostMethodName,
		k6HTTPPutMethodName,
	}

	httpModuleObj := httpModuleValue.ToObject(t.vu.Runtime())

	for _, method := range HTTPMethods {
		originalMethodFn, ok := goja.AssertFunction(httpModuleObj.Get(string(method)))
		if !ok {
			common.Throw(t.vu.Runtime(), fmt.Errorf("http.%s is not a function", method))
		}

		tracedMethodFn := t.withTracing(method, originalMethodFn)

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

// withTracing returns a new function that wraps the original http
// method, adding tracing headers to the request. The produced function
// is ready to be injected back into the http module, in place of the
// original method.
func (t *Tracing) withTracing(methodName k6HTTPMethodName, methodFn goja.Callable) goja.Callable {
	return func(url goja.Value, args ...goja.Value) (goja.Value, error) {
		rt := t.vu.Runtime()

		// normalize the arguments, and extract the params object
		// from the arguments.
		args, params := t.normalizeArgs(methodName, args...)

		// Ensure that the params object contains a headers object.
		// Create it if it doesn't.
		headers, err := t.normalizeParamsHeaders(params)
		if err != nil {
			common.Throw(rt, fmt.Errorf("failed to normalize HTTP headers: %w", err))
		}

		traceID := NewTraceID(k6Prefix, k6CloudCode, uint64(time.Now().UnixNano())/uint64(time.Millisecond))
		encodedTraceID, _, err := traceID.Encode()
		if err != nil {
			common.Throw(rt, fmt.Errorf("failed to encode trace ID: %w", err))
		}

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

		args = append([]goja.Value{url}, args...)

		// call the original http.get method, with overridden arguments
		result, err := methodFn(goja.Undefined(), args...)
		if err != nil {
			common.Throw(rt, err)
		}

		return result, err
	}
}

// normalizeArgs ensures that the http method arguments are in the
// expected format, and returns the normalized arguments and the params
// object.
//
// As of k6 v0.42.0 the HTTP API can turn out to be a bit inconsistent,
// as it can be called with 0, 1 or 2 arguments, and the second argument
// can be either a request's body, or a params object.
func (t *Tracing) normalizeArgs(forMethod k6HTTPMethodName, args ...goja.Value) ([]goja.Value, *goja.Object) {
	rt := t.vu.Runtime()
	var params *goja.Object

	switch len(args) {
	case 2:
		paramsValue := args[1]
		if params == nil || goja.IsUndefined(params) || goja.IsNull(params) {
			params = rt.NewObject()
			args[1] = params
		} else {
			params = paramsValue.ToObject(rt)
		}
	case 1:
		// The http.get and the http.head methods take an optional params
		// object as first argument. Whereas the other methods take a
		// request's body as first argument, and an optional params object
		// as second argument.
		if forMethod == k6HTTPGetMethodName || forMethod == k6HTTPHeadMethodName {
			params = args[0].ToObject(rt)
		} else {
			params = rt.NewObject()
			args = append(args, params)
		}
	case 0:
		params = rt.NewObject()
		args = []goja.Value{goja.Null(), params}
	}

	return args, params
}

// normalizeHTTPHeaders ensures that the http method params object has
// a headers object, and returns it. This method modifies the params
// object in place, and returns the headers object.
func (t *Tracing) normalizeParamsHeaders(params *goja.Object) (*goja.Object, error) {
	rt := t.vu.Runtime()
	var headers *goja.Object

	headersValue := params.Get("headers")
	if headersValue == nil || goja.IsUndefined(headersValue) || goja.IsNull(headersValue) {
		headers = rt.NewObject()
		err := params.Set("headers", headers)
		if err != nil {
			return nil, err
		}

		return headers, nil
	}

	return headersValue.ToObject(rt), nil
}

type k6HTTPMethodName string

const (
	k6HTTPDeleteMethodName  k6HTTPMethodName = "del"
	k6HTTPGetMethodName     k6HTTPMethodName = "get"
	k6HTTPPostMethodName    k6HTTPMethodName = "post"
	k6HTTPPutMethodName     k6HTTPMethodName = "put"
	k6HTTPPatchMethodName   k6HTTPMethodName = "patch"
	k6HTTPHeadMethodName    k6HTTPMethodName = "head"
	k6HTTPOptionsMethodName k6HTTPMethodName = "options"
)
