package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/99designs/gqlgen/handler"
	"github.com/go-chi/chi"
	chiMiddleware "github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/jinzhu/gorm"
	"github.com/rs/cors"
	"go.uber.org/zap"

	"github.com/caquillo07/graphql-server-demo/conf"
	"github.com/caquillo07/graphql-server-demo/pkg/apierrors"
	gqlServer "github.com/caquillo07/graphql-server-demo/pkg/gqlgen/server"
)

// Server the server to be used in the application
type Server interface {
	Serve() error
}

type server struct {
	db           *gorm.DB
	httpServer   *http.Server
	config       conf.Config
	closeTimeout time.Duration
}

// Handler custom handler to allow for error checking
type Handler func(w http.ResponseWriter, r *http.Request) (render.Renderer, error)

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	res, err := h(w, r)
	if err != nil {
		renderError(w, r, err)
		return
	}
	renderResponse(w, r, res)
}

// NewGQLServerWithCloseTimeout returns a server with a custom timeout on closing
func NewGQLServerWithCloseTimeout(db *gorm.DB, config conf.Config, timeout time.Duration) Server {
	r := chi.NewRouter()
	srv := &server{
		db:           db,
		httpServer:   &http.Server{Addr: ":" + config.Server.Port, Handler: r},
		config:       config,
		closeTimeout: timeout,
	}

	applyCommonMiddleware(r, config)

	if config.GraphQL.Playground {
		r.Get("/playground", handler.Playground("GraphQL playground", "/graphql"))
	}

	gql := handler.GraphQL(
		gqlServer.NewExecutableSchema(gqlServer.Config{Resolvers: srv}),
		ErrorHandler(),
		GQLLogging(config),
	)

	// authenticated routes
	r.Group(func(r chi.Router) {

		// require auth on this route group
		authMiddleware := AuthMiddleware(config)
		r.Use(authMiddleware)
		r.Post("/graphql", gql)
	})

	return srv
}

// NewGQLServer creates and returns a new server instance for the application
func NewGQLServer(db *gorm.DB, config conf.Config) Server {
	return NewGQLServerWithCloseTimeout(db, config, 10*time.Second)
}

func (s *server) Serve() error {
	s.applyGracefulShutdown()

	startMsg := "listening on http://localhost:%s"
	if s.config.GraphQL.Playground {
		startMsg = "connect to http://localhost:%s/playground for GraphQL playground"
	}
	zap.L().Info(fmt.Sprintf(startMsg, s.config.Server.Port))
	return s.httpServer.ListenAndServe()
}

func applyCommonMiddleware(r *chi.Mux, config conf.Config) {
	if config.CORS.Enabled {
		r.Use(cors.New(cors.Options{
			AllowCredentials: true,
			AllowedOrigins:   config.CORS.AllowedOrigins,
			Debug:            config.CORS.Debug,
		}).Handler)
	}

	r.Use(
		chiMiddleware.RequestID,
		chiMiddleware.Recoverer,
		chiMiddleware.Logger,
	)
}

func (s *server) applyGracefulShutdown() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c

		// sig is a ^C, handle it
		// create context with timeout
		ctx, cancel := context.WithTimeout(context.Background(), s.closeTimeout)
		defer cancel()

		// start http shutdown
		fmt.Println("shutting down..")
		if err := s.httpServer.Shutdown(ctx); err != nil {
			zap.L().Error("error when shutting down server", zap.Error(err))
		}

		// verify, in worst case call cancel via defer
		select {
		case <-time.After(s.closeTimeout + (time.Second * 1)):
			fmt.Println("not all connections done")
		case <-ctx.Done():
			// done
		}
	}()
}

func renderError(w http.ResponseWriter, r *http.Request, err error) {
	renderErr := func(v render.Renderer) {
		zap.L().Info("returning response error", zap.Error(err))
		if err := render.Render(w, r, v); err != nil {
			// if an error happens here, just log it.
			zap.L().Error("error rendering api error", zap.Error(err))
		}
	}
	switch e := err.(type) {
	case render.Renderer:
		renderErr(e)
	case apierrors.PublicAPIError:
		resErr := apierrors.NewResponseError(
			e.PublicError(),
			e.ErrorCode(),
			e.HTTPStatusCode(),
			nil,
			err,
		)

		if dErr, ok := err.(apierrors.PublicErrorDetails); ok {
			resErr.Details = dErr.Details()
		}
		renderErr(resErr)
	default:

		// check if its a db not found error, if so let those through.
		if err == gorm.ErrRecordNotFound {
			renderErr(apierrors.NewResponseError(
				"record not found",
				http.StatusNotFound,
				http.StatusNotFound,
				nil,
				err,
			))
			return
		}
		zap.L().Info("masking internal error", zap.Error(err))
		renderErr(apierrors.NewInternalError(err))
	}
}

func renderResponse(w http.ResponseWriter, r *http.Request, v render.Renderer) {
	if err := render.Render(w, r, v); err != nil {

		// log the error as we cant return error here
		zap.L().Error("error sending http response", zap.Error(err))
	}
}
