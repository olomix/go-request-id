package grpc

import (
	"context"

	go_request_id "github.com/olomix/go-request-id"
	google_grpc "google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const ridHeader = "x-request-id"

var ridHeaders = []string{"x-request-id", "X-Request-ID"}

// UnaryRequestIdExtractor is an interceptor to extract request-id from
// incoming metadata and insert it into context
func UnaryRequestIdExtractor(ctx context.Context, req interface{},
	_ *google_grpc.UnaryServerInfo,
	handler google_grpc.UnaryHandler) (interface{}, error) {
	rid := ridFromIncomingRequest(ctx)
	if rid != "" {
		ctx = go_request_id.WithRequestID(ctx, rid)
	}
	return handler(ctx, req)
}

// StreamRequestIdExtractor is an interceptor to extract request-id from
// incoming metadata and insert it into context
func StreamRequestIdExtractor(srv interface{}, stream google_grpc.ServerStream,
	_ *google_grpc.StreamServerInfo, handler google_grpc.StreamHandler) error {
	ctx := stream.Context()
	rid := ridFromIncomingRequest(ctx)
	if rid != "" {
		ctx = go_request_id.WithRequestID(ctx, rid)
		stream = &wrappedServerStream{
			ServerStream: stream,
			ctx:          ctx,
		}
	}
	return handler(srv, stream)
}

// UnaryRequestIdInjector adds request-id header to outgoing unary requests
func UnaryRequestIdInjector(
	ctx context.Context, method string, req, reply interface{},
	cc *google_grpc.ClientConn, invoker google_grpc.UnaryInvoker,
	opts ...google_grpc.CallOption) error {
	return invoker(ctxWithRequestIdMetadata(ctx), method, req, reply, cc, opts...)
}

// StreamRequestIdInjector add request-id header to outgoing stream requests
func StreamRequestIdInjector(
	ctx context.Context, desc *google_grpc.StreamDesc,
	cc *google_grpc.ClientConn, method string, streamer google_grpc.Streamer,
	opts ...google_grpc.CallOption) (google_grpc.ClientStream, error) {
	return streamer(ctxWithRequestIdMetadata(ctx), desc, cc, method, opts...)
}

func ridFromIncomingRequest(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}

	for _, key := range ridHeaders {
		rids, ok := md[key]
		if !ok {
			continue
		}
		if len(rids) > 0 {
			return rids[0]
		}
	}

	return ""
}

// if ctx contains request-id, append it to outgoing ctx metadata
func ctxWithRequestIdMetadata(ctx context.Context) context.Context {
	rid := go_request_id.ExtractRequestID(ctx)
	if rid != "" {
		return metadata.AppendToOutgoingContext(ctx, ridHeader, rid)
	}
	return ctx
}

type wrappedServerStream struct {
	google_grpc.ServerStream
	ctx context.Context
}

func (w *wrappedServerStream) Context() context.Context {
	return w.ctx
}
