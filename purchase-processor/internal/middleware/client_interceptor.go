package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func ClientInterceptor() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		traceID, _ := ctx.Value(TraceIDKey).(string)
		spanID, _ := ctx.Value(SpanIDKey).(string)

		if traceID == "" {
			traceID = generateID()
			log.Printf("Generated new TraceID: %s", traceID)
		}

		if spanID == "" {
			spanID = generateID()
			log.Printf("Generated new SpanID: %s", spanID)
		}

		md := metadata.Pairs(
			TraceIDKey, traceID,
			SpanIDKey, spanID,
		)

		ctx = metadata.NewOutgoingContext(ctx, md)

		log.Printf("Sending request: %s, TraceID: %s, SpanID: %s", method, traceID, spanID)

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func generateID() string {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		log.Println("Error generating ID:", err)
		return "unknown-id"
	}

	return hex.EncodeToString(bytes)
}
