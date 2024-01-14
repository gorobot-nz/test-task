package middleware

import (
	"context"
	"encoding/base64"
	userv1 "github.com/gorobot-nz/test-task/gen/proto/user/v1"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"slices"
	"strings"
)

var adminOnlyMethods = []string{
	userv1.UserService_NewUser_FullMethodName,
	userv1.UserService_DeleteUser_FullMethodName,
	userv1.UserService_UpdateUser_FullMethodName,
}

const Username = "username"
const Password = "password"

func AuthMiddleware() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		if !slices.Contains(adminOnlyMethods, info.FullMethod) {
			return handler(ctx, req)
		}

		token, err := grpc_auth.AuthFromMD(ctx, "basic")
		if err != nil {
			return nil, status.Error(codes.PermissionDenied, err.Error())
		}

		decodeString, err := base64.StdEncoding.DecodeString(token)
		if err != nil {
			return nil, status.Error(codes.PermissionDenied, err.Error())
		}

		split := strings.Split(string(decodeString), ":")

		if len(split) != 2 {
			return nil, status.Error(codes.PermissionDenied, err.Error())
		}

		ctx = context.WithValue(ctx, Username, split[0])
		ctx = context.WithValue(ctx, Password, split[1])

		return handler(ctx, req)
	}
}
