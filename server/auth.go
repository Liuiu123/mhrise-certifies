package main

import (
	"context"
	"fmt"
	"time"

	auth "npln.nintendo.net/npln-practice/proto/auth/v1"
	"google.golang.org/protobuf/types/known/durationpb"
)

type AuthServer struct {
	auth.UnimplementedAuthServer
}

func NewAuthServer() *AuthServer { return &AuthServer{} }

func (s *AuthServer) CreateUser(ctx context.Context, req *auth.CreateUserRequest) (*auth.User, error) {
	user := req.GetUser()
	if user == nil {
		user = &auth.User{}
	}
	if user.GetName() == "" {
		if req.GetParent() != "" {
			user.Name = req.GetParent() + "/users/lab-user-001"
		} else {
			user.Name = "users/lab-user-001"
		}
	}
	if user.GetAccount() == "" {
		user.Account = "accounts/lab-account"
	}
	return user, nil
}

func (s *AuthServer) IssueToken(ctx context.Context, req *auth.IssueTokenRequest) (*auth.IssueTokenResponse, error) {
	user := req.GetUser()
	if user == "" {
		return nil, fmt.Errorf("user is required")
	}
	return &auth.IssueTokenResponse{Token: &auth.Token{
		User: user,
		AccessToken: "lab-access-token",
		RefreshToken: "lab-refresh-token",
		Ttl: durationpb.New(24 * time.Hour),
	}}, nil
}

func (s *AuthServer) RefreshToken(ctx context.Context, req *auth.RefreshTokenRequest) (*auth.RefreshTokenResponse, error) {
	if req.GetUser() == "" || req.GetRefreshToken() == "" {
		return nil, fmt.Errorf("user and refresh_token are required")
	}
	return &auth.RefreshTokenResponse{Token: &auth.Token{
		User: req.GetUser(),
		AccessToken: "lab-access-token",
		RefreshToken: req.GetRefreshToken(),
		Ttl: durationpb.New(24 * time.Hour),
	}}, nil
}

func (s *AuthServer) IssuePrearrangedUserToken(ctx context.Context, req *auth.IssuePrearrangedUserTokenRequest) (*auth.IssuePrearrangedUserTokenResponse, error) {
	user := &auth.User{Name: fmt.Sprintf("tenants/%s/users/%d", req.GetTenant(), req.GetUserIndex())}
	return &auth.IssuePrearrangedUserTokenResponse{
		User: user,
		Token: &auth.Token{User: user.GetName(), AccessToken: "lab-access-token", RefreshToken: "lab-refresh-token", Ttl: durationpb.New(24 * time.Hour)},
	}, nil
}
