package tracing

import (
	"net/http"
)

const (
	// W3CHeaderName is the name of the W3C trace context header
	W3CHeaderName = "Traceparent"

	// B3HeaderName is the name of the B3 trace context header
	B3HeaderName = "b3"

	// JaegerHeaderName is the name of the Jaeger trace context header
	JaegerHeaderName = "uber-trace-id"
)

// Propagator is an interface for trace context propagation
type Propagator interface {
	Propagate(traceID string) (http.Header, error)
}

// W3CPropagator is a Propagator for the W3C trace context header
type W3CPropagator struct{}

// Propagate returns a header with a random trace ID in the W3C format
func (p *W3CPropagator) Propagate(traceID string) (http.Header, error) {
	return http.Header{
		W3CHeaderName: {"00-" + traceID + "-" + RandHexStringRunes(16) + "-01", traceID, RandHexStringRunes(16)},
	}, nil
}

// B3Propagator is a Propagator for the B3 trace context header
type B3Propagator struct{}

// Propagate returns a header with a random trace ID in the B3 format
func (p *B3Propagator) Propagate(traceID string) (http.Header, error) {
	return http.Header{
		B3HeaderName: {traceID + "-" + RandHexStringRunes(8) + "-1"},
	}, nil
}

// JaegerPropagator is a Propagator for the Jaeger trace context header
type JaegerPropagator struct{}

// Propagate returns a header with a random trace ID in the Jaeger format
func (p *JaegerPropagator) Propagate(traceID string) (http.Header, error) {
	return http.Header{
		JaegerHeaderName: {traceID + ":" + RandHexStringRunes(8) + ":0:1"},
	}, nil
}
