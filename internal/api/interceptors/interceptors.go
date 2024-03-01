package interceptors

import (
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	pb "github.com/rebus2015/gophkeeper/internal/proto"
)

func HashInterceptor(key string) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {

		if key == "" {
			return handler(ctx, req)
		}
		// var err error
		data, ok := req.(*pb.RegisterRequest)
		if !ok {
			return nil, status.Errorf(codes.Canceled, "%v", "hash interceptor error: corrupted data")
		}
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			values := md.Get("single")
			if len(values) == 0 {
				return handler(ctx, req)
			}
		}

		// reader := bytes.NewReader([]byte(data.Login))
		log.Println("Incoming request Updates, before decoder")

		// bodyBytes, _ := io.ReadAll(reader)
		// metrics, err := getMetrics(bodyBytes)
		// if err != nil {
		// 	log.Printf("Failed to Decode incoming metricList %v, error: %v", string(bodyBytes), err)
		// 	return nil, status.Errorf(codes.InvalidArgument, "Failed to Decode incoming metricList %v", err)
		// }
		// log.Printf("Try to update metrics: %v", metrics)
		// for i := range metrics {
		// 	if key != "" {
		// 		pass, err := checkMetric(metrics[i], key)
		// 		if err != nil || !pass {
		// 			log.Printf("check Metrics error %v, error: %v", string(bodyBytes), err)
		// 			return nil, status.Errorf(codes.InvalidArgument, "Failed to Decode incoming metricList %v", err)
		// 		}
		// 	}
		// }
		return handler(ctx, data)
	}
}
