package auth_client_service

import (
	"context"
	protos "github.com/aleale2121/DSP_LAB/Music_Service/grpc/client_to_server/auth_client_service/auth"
	"google.golang.org/grpc"
	"time"
)

type AuthClient struct {
	service protos.AuthServiceClient
}

func NewAuthClient(cc *grpc.ClientConn) *AuthClient {
	service := protos.NewAuthServiceClient(cc)
	return &AuthClient{service}
}

func (client *AuthClient) Login(username string, password string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &protos.LoginRequest{
		Username: username,
		Password: password,
	}

	res, err := client.service.Login(ctx, req)
	if err != nil {
		return "", err
	}

	return res.GetAccessToken(), nil
}

func (client *AuthClient) GetUserId(token string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req := &protos.GetUserIdRequest{
		Token: token,
	}

	res, err := client.service.GetUserId(ctx, req)
	if err != nil {
		return "", err
	}

	return res.Id, nil
}

func (client *AuthClient) Register(username, password, phone, roleId string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	req := &protos.RegistrationRequest{
		Username: username,
		Password: password,
		Phone:    phone,
		RoleId:   roleId,
	}

	res, err := client.service.Register(ctx, req)
	if err != nil {
		return "", err
	}

	return res.UserId, nil
}
func (client *AuthClient) VerifyPhone(userId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	req := &protos.VerifyPhoneRequest{
		Phone: userId,
	}

	_, err := client.service.VerifyUserPhone(ctx, req)
	if err != nil {
		return err
	}

	return nil
}
func (client *AuthClient) IsAuthorized(token string, path string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	req := &protos.IsAuthorizedRequest{
		Token: token,
		Path: path,
	}

	resp, err := client.service.IsAuthorized(ctx, req)
	if err != nil {
		return false, err
	}

	return resp.IsAuthorized, nil
}
