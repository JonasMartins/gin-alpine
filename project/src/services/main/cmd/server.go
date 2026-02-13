package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gin-alpine/src/services/main/internal/bootstrap"
	"gin-alpine/src/services/main/internal/handler/router"

	"github.com/gin-gonic/gin"
)

type Server struct {
	Bootstrap *bootstrap.Bootstrap
	Router    http.Handler
	Srv       *http.Server
}

func NewServer() *Server {
	b := bootstrap.MustGetBootstrapInstance()
	if b.Config.Env == "production" {
		b.Config.Env = gin.ReleaseMode
	} else {
		b.Config.Env = gin.DebugMode
	}
	b.SetInitialData()
	router := router.NewRouter(b)
	addr := fmt.Sprintf(":%s", b.Config.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}
	return &Server{
		Bootstrap: b,
		Router:    router,
		Srv:       srv,
	}
}

func Run() {
	runHTTPServer(NewServer())
}

func runHTTPServer(s *Server) {
	env := "DEVELOPMENT"
	delay := 500 * time.Millisecond
	if s.Bootstrap.Config.Env == gin.ReleaseMode {
		env = "PRODUCTION"
		delay = 3 * time.Second
	}
	s.Bootstrap.RunJobs()
	log.Printf("HTTP server running at :%s", s.Bootstrap.Config.Port)
	log.Printf("( %s MODE ) running main project", env)
	go func() {
		if err := s.Srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP Server Error %v", err)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	<-quit
	log.Println("Shutdown Server ...")
	ctx, cancel := context.WithTimeout(context.Background(), delay)
	defer cancel()
	if err := s.Srv.Shutdown(ctx); err != nil {
		log.Printf("Server Shutdown error: %v\n", err)
		defer os.Exit(1)
	} else {
		log.Printf("gracefully stopped\n")
	}
	conn := <-ctx.Done()
	log.Println("timeout of 3 seconds. ", conn)
	log.Println("Server exiting")
}

func RunLoop() {
	for {
		log.Println("Running server for ever")
		time.Sleep(1 * time.Minute)
	}
}
