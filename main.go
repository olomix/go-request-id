package go_request_id

import (
	"context"
	"log"

	"github.com/google/uuid"
)

type key uint8

type logger interface {
	Errorf(fmt string, args ...interface{})
}

type defaultLogger struct{}

func (defaultLogger) Errorf(format string, args ...interface{}) {
	log.Printf(format, args...)
}

var Logger logger = defaultLogger{}

const (
	keyRequestID key = iota
)

// WithRequestID returns new context.Context with new request-id attached to it
func WithRequestID(parent context.Context, rid string) context.Context {
	if rid == "" {
		return parent
	}
	return context.WithValue(parent, keyRequestID, rid)
}

// WithNewRandomRequestID returns context.Context with new random request-id
// attached to it
func WithNewRandomRequestID(parent context.Context) context.Context {
	uuid4, err := uuid.NewRandom()
	if err != nil {
		Logger.Errorf("error getting new random UUID: %v", err)
		return WithRequestID(parent, uuid.Nil.String())
	}
	return WithRequestID(parent, uuid4.String())
}

// ExtractRequestID returns request-id from context or empty string of no
// request-id found in context.
func ExtractRequestID(ctx context.Context) string {
	v := ctx.Value(keyRequestID)
	if v == nil {
		return ""
	}
	switch value := v.(type) {
	case string:
		return value
	default:
		Logger.Errorf("Unexpected txid type found in context: %T", v)
		return ""
	}
}
