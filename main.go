package main

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"log"
	"net"
	"net/http"
	"os"
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
	if len(os.Args) != 2 {
		log.Printf("need port number\n")
		os.Exit(1)
	}
	p := os.Args[1]

	l, err := net.Listen("tcp", ":"+p)
	if err != nil {
		log.Fatalf("failed to listen to port %s: %v", p, err)
	}
	if err := run(context.Background(), l); err != nil {
		log.Printf("failed to terminate server: %v", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, l net.Listener) error {
	// *http.Server.ListenAndServeメソッドを実行してHTTPリクエストを受け付ける
	// 引数で渡されたcontext.Contextを通じて処理の中段命令を検知したとき、*http.Server.ShutdownメソッドでHTTPサーバの機能を終了する
	// run関数の戻り値として*http.Server.ListenAndServeメソッドの戻り値のエラーを返す
	s := &http.Server{
		// 引数で受け取ったnet.listenerを利用するので、Addrフィールドは指定しない
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
