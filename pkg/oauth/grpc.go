package oauth

import (
	"context"
	"fmt"
	"strings"

	stdjwt "github.com/dgrijalva/jwt-go"
	"github.com/go-kit/kit/auth/jwt"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func AuthFromGrpcContext() grpctransport.ClientRequestFunc {
	return func(ctx context.Context, md *metadata.MD) context.Context {
		user_id, ok := ctx.Value("user_id").(string)
		if !ok {
			return ctx
		}
		(*md)["user_id"] = []string{user_id}
		return ctx
	}
}

func AuthToGrpcContext() grpctransport.ServerRequestFunc {
	return func(ctx context.Context, md metadata.MD) context.Context {
		uid_, ok := md["user_id"]
		if !ok {
			return ctx
		}
		uid := uid_[0]
		ctx = context.WithValue(ctx, "user_id", uid)
		return ctx
	}
}

func FromGrpcContext() grpctransport.ClientRequestFunc {
	return func(ctx context.Context, md *metadata.MD) context.Context {
		token, ok := ctx.Value(jwt.JWTClaimsContextKey).(stdjwt.MapClaims)
		if ok {
			(*md)["authorization"] = []string{fmt.Sprintf("Bearer %s", token)}
		}
		return ctx
	}
}

func ToGrpcContext() grpctransport.ServerRequestFunc {
	return func(ctx context.Context, md metadata.MD) context.Context {
		authHeader, ok := md["authorization"]
		if !ok {
			return ctx
		}
		token, ok := extractTokenFromAuthHeader(authHeader[0])
		if ok {
			ctx = context.WithValue(ctx, jwt.JWTTokenContextKey, token)
		}
		return ctx
	}
}

func extractTokenFromAuthHeader(val string) (token string, ok bool) {
	authHeaderParts := strings.Split(val, " ")
	if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
		return "", false
	}

	return authHeaderParts[1], true
}

func StreamIntercepter(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	return nil
}

func UnaryIntercepter(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	return nil, nil
}
