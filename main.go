package main

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
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
	// simple http serverの実装では、実行できるが終了は指示できない
	// 関数外部から中断操作ができず、関数の戻り値もない
	// ポート番号も固定されているため、サーバを起動したままテストもできない
	// → run関数に分離する
	if err := run(context.Background()); err != nil {
		log.Printf("failed to terminate server: %v", err)
	}
}

func run(ctx context.Context) error {
	// *http.Server.ListenAndServeメソッドを実行してHTTPリクエストを受け付ける
	// 引数で渡されたcontext.Contextを通じて処理の中段命令を検知したとき、*http.Server.ShutdownメソッドでHTTPサーバの機能を終了する
	// run関数の戻り値として*http.Server.ListenAndServeメソッドの戻り値のエラーを返す
	s := &http.Server{
		Addr: ":18080",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
		}),
	}
	eg, ctx := errgroup.WithContext(ctx)
	// 別ゴルーチンでHTTPサーバを起動する
	eg.Go(func() error {
		// http.ErrServerClosedは
		// http.Server.Shutdown() が正常に終了したことを示すので異常ではない
		if err := s.ListenAndServe(); err != nil &&
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
