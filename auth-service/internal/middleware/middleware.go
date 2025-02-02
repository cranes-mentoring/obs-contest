package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func ServerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			log.Println("No metadata found in context")
			return handler(ctx, req)
		}

		traceIDs := md.Get(TraceIDKey)
		spanIDs := md.Get(SpanIDKey)

		var traceID, spanID string

		if len(traceIDs) > 0 {
			traceID = traceIDs[0]
		} else {
			traceID = generateID()
			log.Printf("Generated new TraceID: %s", traceID)
		}

		if len(spanIDs) > 0 {
			spanID = spanIDs[0]
		} else {
			spanID = generateID()
			log.Printf("Generated new SpanID: %s", spanID)
		}

		ctx = context.WithValue(ctx, TraceIDKey, traceID)
		ctx = context.WithValue(ctx, SpanIDKey, spanID)

		log.Printf("Received request: %s, TraceID: %s, SpanID: %s", info.FullMethod, traceID, spanID)

		return handler(ctx, req)
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
