package server

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "github.com/rebus2015/gophkeeper/internal/proto"
	"github.com/rebus2015/gophkeeper/internal/storage/db"
	"google.golang.org/grpc"

	"github.com/rebus2015/gophkeeper/internal/api/interceptors"
	"github.com/rebus2015/gophkeeper/internal/logger"
	"github.com/rebus2015/gophkeeper/internal/model"
	"github.com/rebus2015/gophkeeper/internal/server/config"
)

type AuthRPCServer struct {
	srv *grpc.Server
	pb.UnimplementedUserServer
	postgreStorage db.PostgreSQLStorage
	cfg            config.Config
	log            *logger.Logger
}

func NewRPCServer(
	pgsStorage db.PostgreSQLStorage,
	conf config.Config,
	logger *logger.Logger) *AuthRPCServer {
	return &AuthRPCServer{
		postgreStorage: pgsStorage,
		cfg:            conf,
		log:            logger,
	}
}

func (s *AuthRPCServer) Run() error {
	listen, err := net.Listen("tcp", s.cfg.RunAddress)
	if err != nil {
		return fmt.Errorf("start RPC server error: %w", err)
	}

	s.srv = grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			//interceptors.SubnetCheckInterceptor(s.cfg),
			//interceptors.GzipInterceptor,
			//interceptors.RsaInterceptor(s.cfg.CryptoKey),
			interceptors.HashInterceptor(s.cfg.SecretKey),
		))

	// регистрируем сервис
	pb.RegisterUserServer(s.srv, &AuthRPCServer{
		postgreStorage: s.postgreStorage,
		cfg:            s.cfg,
	})

	fmt.Println("Сервер gRPC начал работу")
	// получаем запрос gRPC
	if err := s.srv.Serve(listen); err != nil {
		log.Fatal(err)
		return fmt.Errorf("run server err: %w", err)
	}
	return nil
}

func (s *AuthRPCServer) Register(ctx context.Context, in *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	var response pb.RegisterResponse
	log.Println("Incoming request Ping")
	user := model.User{
		Login:    "",
		Password: "",
	}
	// При успешной проверке хендлер должен вернуть HTTP-статус 200 OK, при неуспешной — 500 Internal Server Error.
	// if s.postgreStorage == nil {
	// 	//response.Status = 500
	// 	response.Error = status.Error(codes.Internal,
	// 		"Failed to ping database: nil reference exception: postgreStorage udefined").Error()
	// 	return &response, fmt.Errorf("nil reference exception: postgreStorage udefined ")
	// }
	if _, err := s.postgreStorage.UserRegister(&user); err != nil {
		log.Printf("Cannot ping database because %s", err)
		response.Token = ""
		return &response, fmt.Errorf("failed to Decode incoming metricList %w", err)
	}
	response.Token = "Token string"
	return &response, nil
}

func (s *AuthRPCServer) Shutdown() {
	s.srv.GracefulStop()
}
