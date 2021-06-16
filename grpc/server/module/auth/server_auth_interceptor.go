package auth

import (
	"context"
	"github.com/aleale2121/DSP_LAB/Music_Service/constant"
	"github.com/aleale2121/DSP_LAB/Music_Service/grpc/client_to_server/auth_client_service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log"
)

type AuthenticationInterceptor struct {}

func NewAuthInterceptor( ) *AuthenticationInterceptor {
	return &AuthenticationInterceptor{}
}

func (interceptor *AuthenticationInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		log.Println("--> unary interceptor: ", info.FullMethod)

		err := interceptor.authorize(ctx, info.FullMethod)
		if err != nil {
			return nil, err
		}
		return handler(ctx, req)

	}
}

func (interceptor *AuthenticationInterceptor) Stream() grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		log.Println("--> stream interceptor: ", info.FullMethod)

		err := interceptor.authorize(stream.Context(), info.FullMethod)
		if err != nil {
			return err
		}

		return handler(srv, stream)
	}
}

func (interceptor *AuthenticationInterceptor) authorize(ctx context.Context, method string) error {

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	values := md["authorization"]
	if len(values) == 0 {
		return status.Errorf(codes.Unauthenticated, "authorization token is not provided")
	}

	accessToken := values[0]
	grpcDialStr, err := constant.GetUserManagementGrpcConnectionString()
	if err != nil {
		return status.Error(codes.Unavailable, "Unable To connect with remote server")

	}
	transportOption := grpc.WithInsecure()

	cc, err := grpc.Dial(grpcDialStr, transportOption)

	if err != nil {
		return status.Error(codes.Unavailable, "Unable To connect with remote server")

	}
	authClient := *auth_client_service.NewAuthClient(cc)

	isAuthorized, err := authClient.IsAuthorized(accessToken, method)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}
	if !isAuthorized {
		return status.Error(codes.PermissionDenied, "no permission to access this RPC")
	}
	return nil

}
