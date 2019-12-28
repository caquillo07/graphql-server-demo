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
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"

	"github.com/caquillo07/graphql-server-demo/conf"
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

// NewGQLServerWithCloseTimeout returns a server with a custom timeout on closing
func NewGQLServerWithCloseTimeout(config conf.Config, timeout time.Duration) Server {
	r := chi.NewRouter()
	srv := &server{
		httpServer:   &http.Server{Addr: ":" + config.Server.Port, Handler: r},
		config:       config,
		closeTimeout: timeout,
	}

	if config.GraphQL.Playground {
		r.Get("/playground", handler.Playground("GraphQL playground", "/graphql"))
	}

	gql := handler.GraphQL(
		gqlServer.NewExecutableSchema(gqlServer.Config{Resolvers: srv}),
	)
	r.Post("/graphql", gql)

	return srv
}

// NewGQLServer creates and returns a new server instance for the application
func NewGQLServer(config conf.Config) Server {
	return NewGQLServerWithCloseTimeout(config, 10*time.Second)
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
