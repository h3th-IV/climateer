package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/h3th-IV/climateer/pkg/utils"
	"github.com/justinas/alice"
	"github.com/rs/cors"
	"go.uber.org/zap"
)

type GracefulShutdownServer struct {
	HTTPListenAddr  string
	RegisterHandler http.Handler // register
	LoginHandler    http.Handler // login
	ProfileHandler  http.Handler // profile
	HomeHandler     http.Handler

	CommentHandler http.Handler //make comments

	httpServer     *http.Server
	WriteTimeout   time.Duration
	ReadTimeout    time.Duration
	IdleTimeout    time.Duration
	HandlerTimeout time.Duration
}

func (server *GracefulShutdownServer) getRouter() *mux.Router {
	router := mux.NewRouter()

	mux.CORSMethodMiddleware(router)
	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})
	middleWareChain := alice.New(utils.RequestLogger, cors.Handler)
	// authRoute := alice.New(middleware.AuthRoute)
	//authed routes
	router.Handle("/register", server.RegisterHandler).Methods(http.MethodPost)
	router.Handle("/login", server.LoginHandler).Methods(http.MethodPost)
	router.Handle("/", server.HomeHandler)
	// cors.Handler(router) jim directive
	router.Use(middleWareChain.Then) //request logging will be handled here
	router.SkipClean(true)
	return router
}

func (server *GracefulShutdownServer) Start() {
	router := server.getRouter()
	server.httpServer = &http.Server{
		Addr:         server.HTTPListenAddr,
		WriteTimeout: server.WriteTimeout,
		ReadTimeout:  server.ReadTimeout,
		IdleTimeout:  server.IdleTimeout,
		Handler:      router,
	}
	utils.Logger.Info(fmt.Sprintf("listening and serving on %s", server.HTTPListenAddr))
	if err := server.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		utils.Logger.Fatal("server failed to start", zap.Error(err))
	}
}
