package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/handler"
	"github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"
	"github.com/vektah/gqlparser/gqlerror"
	"go.uber.org/zap"

	"github.com/caquillo07/graphql-server-demo/conf"
	"github.com/caquillo07/graphql-server-demo/pkg/apierrors"
)

// A private key for context that only this package can access. This is important
// to prevent collisions between different context uses
var userIDCtxKey = &contextKey{"user"}

type contextKey struct {
	name string
}

// ErrorHandler returns an error handler to be used on each API request. This
// handler will look at the message, if it adheres to the apierrors.PublicError
// interface, it will return that error's message, otherwise will return
// "internal error"
func ErrorHandler() handler.Option {
	return handler.ErrorPresenter(func(ctx context.Context, err error) *gqlerror.Error {
		zap.L().Error(err.Error())
		if err, ok := err.(apierrors.PublicError); ok {
			return gqlerror.ErrorPathf(graphql.GetResolverContext(ctx).Path(), err.PublicError())
		}

		// special case in case this error is not caught somewhere else
		if err == gorm.ErrRecordNotFound {
			return gqlerror.ErrorPathf(graphql.GetResolverContext(ctx).Path(), "record not found")
		}

		return graphql.DefaultErrorPresenter(ctx, fmt.Errorf("internal error"))
	})
}

// GQLLogging returns a logging middleware for GQL. This middleware is config
// driven, and if enabled will log the request's query, variables and
// extensions.
func GQLLogging(config conf.Config) handler.Option {
	return handler.RequestMiddleware(func(ctx context.Context, next func(ctx context.Context) []byte) []byte {
		if !config.GraphQL.LogQueries {
			return next(ctx)
		}

		reqCtx := graphql.GetRequestContext(ctx)
		logger := zap.L()

		variables, err := json.Marshal(&reqCtx.Variables)
		if err != nil {
			logger.Error("error while unmarshalling variables in GQLLogging", zap.Error(err))
			return next(ctx)
		}

		extensions, err := json.Marshal(&reqCtx.Extensions)
		if err != nil {
			logger.Error("error while unmarshalling extensions in GQLLogging", zap.Error(err))
			return next(ctx)
		}

		logger.Info(
			"gql",
			zap.String("query", reqCtx.RawQuery),
			zap.ByteString("variables", variables),
			zap.ByteString("extensions", extensions),
		)

		return next(ctx)
	})
}

// AuthMiddleware middle to check for valid JWT in request
func AuthMiddleware(config conf.Config) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {

			// Let secure process the request. If it returns an error,
			// that indicates the request should not continue.
			if config.Auth.Enabled {
				userID, err := checkJWT(w, r)
				if err != nil {
					return
				}

				// put in context
				ctx := context.WithValue(r.Context(), userIDCtxKey, userID)

				// add it to the request so next has it
				r = r.WithContext(ctx)
			}

			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

// UserIDForContext finds the user from the context. REQUIRES Middleware to have
// run.
func UserIDForContext(ctx context.Context) (uuid.UUID, error) {
	raw, ok := ctx.Value(userIDCtxKey).(string)
	if !ok {
		return uuid.Nil, errors.New("user ID not found in context")
	}

	userID, err := uuid.FromString(raw)
	if err != nil {
		return uuid.Nil, err
	}
	return userID, nil
}

func checkJWT(w http.ResponseWriter, r *http.Request) (string, error) {
	// do some jwt validation, here we just return an user ID as test
	return "123e4567-e89b-12d3-a456-426655440000", nil
}
