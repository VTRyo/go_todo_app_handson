package main

import (
	"context"
	"fmt"
	"github.com/VTRyo/go_todo_app_handson/config"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

//func main() {
//	err := http.ListenAndServe(
//		":8080",
//		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//			fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
//		}),
//	)
//	if err != nil {
//		fmt.Printf("failed to terminate server: %v", err)
//		os.Exit(1)
//	}
//}

type Server struct {
	srv *http.Server
	l   net.Listener
}

func NewServer(l net.Listener, mux http.Handler) *Server {
	return &Server{
		srv: &http.Server{Handler: mux},
		l:   l,
	}
}

func main() {
	if err := run(context.Background()); err != nil {
		log.Printf("failed to terminate server: %v", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	// context.Context型の値を経由してシグナルの受信を検知できる
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	// configパッケージを使って起動する
	cfg, err := config.New()
	if err != nil {
		return err
	}
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		log.Fatalf("failed to listen port %d: %v", cfg.Port, err)
	}
	url := fmt.Sprintf("http://%s", l.Addr().String())
	log.Printf("start with: %v", url)

	mux := NewMux()
	s := NewServer(l, mux)
	return s.Run(ctx)
}
