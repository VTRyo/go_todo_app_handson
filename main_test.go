package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"testing"

	"golang.org/x/sync/errgroup"
)

// 以下を検証すること
// 1. 期待したとおりにHTTPサーバが起動しているか
// 2. テストコードから意図通りに終了するか

// 流れ
// キャンセル可能なcontext.Contextのオブジェクトを作る
// 別ゴルーチンでテスト対象のrun関数を実行してHTTPサーバを起動する
// エンドポイントに対してGETリクエストを送信する
// cancel関数を実行する
// *errgroup.Group.Waitメソッド経由でrun関数の戻り値を検証する
// GETリクエストで取得したレスポンスボディが期待する文字列であることを検証する
func TestRun(t *testing.T) {
	t.Skip()
	l, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("failed to listen port %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		return run(ctx)
	})
	in := "message"
	url := fmt.Sprintf("http://%s/%s", l.Addr().String(), in)
	t.Logf("try request to %q", url)
	rsp, err := http.Get(url)

	if err != nil {
		t.Errorf("failed to get: %v", err)
	}
	defer rsp.Body.Close()
	got, err := io.ReadAll(rsp.Body)
	if err != nil {
		t.Fatalf("failed to read body: %v", err)
	}
	want := fmt.Sprintf("Hello, %s!", in)
	if string(got) != want {
		t.Errorf("want %q, but got %q", want, got)
	}
	cancel()
	if err := eg.Wait(); err != nil {
		t.Fatal(err)
	}
}
