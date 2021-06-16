package auth_client_service

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"log"
)

type AuthInterceptor struct {
	authMethods map[string]bool
	accessToken string
}

func NewAuthInterceptor(
	authMethods map[string]bool,
	accessToken string,
) (*AuthInterceptor, error) {
	interceptor := &AuthInterceptor{
		authMethods: authMethods,
		accessToken: accessToken,
	}

	return interceptor, nil
}

func (interceptor *AuthInterceptor) Unary() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		log.Printf("--> unary interceptor: %s", method)

		if interceptor.authMethods[method] {
			err := invoker(interceptor.attachToken(ctx), method, req, reply, cc, opts...)
			if err != nil {
				fmt.Println(err.Error())
			}
			return err
		}

		err := invoker(ctx, method, req, reply, cc, opts...)
		if err != nil {
			fmt.Println(err.Error())
		}
		return err
	}
}

func (interceptor *AuthInterceptor) Stream() grpc.StreamClientInterceptor {
	return func(
		ctx context.Context,
		desc *grpc.StreamDesc,
		cc *grpc.ClientConn,
		method string,
		streamer grpc.Streamer,
		opts ...grpc.CallOption,
	) (grpc.ClientStream, error) {
		log.Printf("--> stream interceptor: %s", method)

		if interceptor.authMethods[method] {
			return streamer(interceptor.attachToken(ctx), desc, cc, method, opts...)
		}

		return streamer(ctx, desc, cc, method, opts...)
	}
}

func (interceptor *AuthInterceptor) attachToken(ctx context.Context) context.Context {
	return metadata.AppendToOutgoingContext(ctx, "authorization", interceptor.accessToken)
}

