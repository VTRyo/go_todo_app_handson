package main

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"io"
	"net/http"
	"testing"
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
func TestMainFunc(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		return run(ctx)
	})
	in := "message"
	rsp, err := http.Get("http://localhost:18080/" + in)
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
