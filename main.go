package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/VTRyo/go_todo_app_handson/config"
	"golang.org/x/sync/errgroup"
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

	// *http.Server.ListenAndServeメソッドを実行してHTTPリクエストを受け付ける
	// 引数で渡されたcontext.Contextを通じて処理の中段命令を検知したとき、*http.Server.ShutdownメソッドでHTTPサーバの機能を終了する
	// run関数の戻り値として*http.Server.ListenAndServeメソッドの戻り値のエラーを返す
	s := &http.Server{
		// 引数で受け取ったnet.listenerを利用するので、Addrフィールドは指定しない
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// コマンドラインで実験する用
			time.Sleep(5 * time.Second)
			fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
		}),
	}
	eg, ctx := errgroup.WithContext(ctx)
	// 別ゴルーチンでHTTPサーバを起動する
	eg.Go(func() error {
		// ListenAndServeメソッドではなくServeメソッドに変更する
		if err := s.Serve(l); err != nil &&
			err != http.ErrServerClosed {
			log.Printf("failed to close: %+v", err)
			return err
		}
		return nil
	})

	<-ctx.Done()
	if err := s.Shutdown(context.Background()); err != nil {
		log.Printf("failed to shutdown: %+v", err)
	}
	return eg.Wait()
}
