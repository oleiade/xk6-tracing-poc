package tracing

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"fmt"
)

const (
	k6Prefix               = 0o756 // Being 075 the ASCII code for 'K' :)
	k6CloudCode            = 12    // To ingest and process the related spans in k6 Cloud.
	k6LocalCode            = 33    // To not ingest and process the related spans, b/c they are part of a non-cloud run.
	metadataTraceIDKeyName = "trace_id"
)

// Tracer is the interface that wraps the TraceID method.
type Tracer interface {
	TraceID() *TraceID
}

// TraceID of 16 bytes (128 bits) are supported by w3c, b3 and jaeger
type TraceID struct {
	Prefix            int16
	Code              int8
	UnixTimestampNano uint64
}

// NewTraceID returns a new TraceID with the given prefix, code and unix timestamp in nanoseconds.
func NewTraceID(prefix int16, code int8, unixTimestampNano uint64) *TraceID {
	return &TraceID{
		Prefix:            prefix,
		Code:              code,
		UnixTimestampNano: unixTimestampNano,
	}
}

// ParseTraceIDFrom parses a TraceID from the given hex string.
func ParseTraceIDFrom(buf []byte) *TraceID {
	pre, preLen := binary.Varint(buf)
	code, codeLen := binary.Varint(buf[preLen:])
	ts, _ := binary.Uvarint(buf[preLen+codeLen:])

	return &TraceID{
		Prefix:            int16(pre),
		Code:              int8(code),
		UnixTimestampNano: ts,
	}
}

// IsValid returns true if the TraceID is valid, false otherwise.
func (t *TraceID) IsValid() bool {
	return t.Prefix == k6Prefix && (t.Code == k6CloudCode || t.Code == k6LocalCode)
}

// Encode encodes the TraceID into a hex string and a byte slice.
func (t *TraceID) Encode() (string, []byte, error) {
	var (
		isk6Prefix    = t.Prefix == k6Prefix
		isk6CloudCode = t.Code == k6CloudCode
		isk6LocalCode = t.Code == k6LocalCode
	)

	if !isk6Prefix && (!isk6CloudCode || !isk6LocalCode) {
		return "", nil, fmt.Errorf("failed to encode traceID: %v", t)
	}

	buf := make([]byte, 16)

	n := binary.PutVarint(buf, int64(t.Prefix))
	n += binary.PutVarint(buf[n:], int64(t.Code))
	n += binary.PutUvarint(buf[n:], t.UnixTimestampNano)

	randomness := make([]byte, 16-len(buf[:n]))
	err := binary.Read(rand.Reader, binary.BigEndian, randomness)
	if err != nil {
		return "", nil, err
	}

	buf = append(buf[:n], randomness...)
	hx := hex.EncodeToString(buf)
	return hx, buf, nil
}
