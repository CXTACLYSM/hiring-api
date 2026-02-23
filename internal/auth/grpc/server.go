package grpc

import (
	"context"
	"fmt"

	"github.com/CXTACLYSM/hiring-api/internal/auth/tokens"
	"github.com/CXTACLYSM/hiring-api/internal/auth/user/application/queries/findOne"
	pb "github.com/CXTACLYSM/hiring-api/pkg/grpc/auth/v1"
	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthServer struct {
	pb.UnimplementedAuthServiceServer
	secretKey   []byte
	findOneUser findOne.Handler
}

func NewAuthServer(secretKey []byte, userFinder findOne.Handler) *AuthServer {
	return &AuthServer{
		secretKey:   secretKey,
		findOneUser: userFinder,
	}
}

func (s *AuthServer) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	claims := &tokens.Claims{}
	_, err := jwt.NewParser().ParseWithClaims(req.Token, claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.secretKey, nil
	})
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid or expired token")
	}

	user, err := s.findOneUser.Handle(findOne.Query{
		Id: claims.UserId,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to find user")
	}
	if user == nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	return &pb.ValidateTokenResponse{
		Valid: true,
		User: &pb.User{
			Id:       user.Id,
			Username: user.Username,
			Email:    user.Email,
		},
	}, nil
}
