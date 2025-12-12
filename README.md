# notion-invoice-with-lambda

Go 1.25 / AWS Lambda の請求書自動生成システム。

## ディレクトリ
- `cmd/invoice-batch`: 請求PDF作成バッチのエントリ
- `cmd/invoice-mail`: 承認済み請求のメール送信エントリ
- `cmd/cashflow-update`: キャッシュフロー更新バッチのエントリ
- `internal/logger`: 共通JSONロガー
- `internal/app`: Lambda 本体のビジネスロジック（スタブ）

## ロガーの使い方
```go
log := logger.New() // STAGE 環境変数がなければ local
entry := log.WithContext(ctx) // Lambda context から request_id を取得
entry.Info("started", map[string]any{"foo": "bar"})
```

出力例（JSON, 1 行）:
```json
{"stage":"staging","level":"INFO","message":"started","request_id":"...","foo":"bar"}
```

## テスト
```sh
GOCACHE=$(pwd)/.cache/go-build GOMODCACHE=$(pwd)/.cache/gomod go test ./...
```
