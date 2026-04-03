# ARCHITECTURES.md

## 概要

このリポジトリは、Notion データベース（Transactions / Customers / Invoices / Cashflow）と AWS Lambda を連携し、請求書 PDF 自動生成・メール送信・キャッシュフロー集計を行うシステムである。

詳細仕様は `docs/HLD.md` を参照すること。

---

## 技術前提

- 言語：golang 1.25
- 実行環境：AWS Lambda
- ストレージ：AWS S3（PDF 保存）
- 外部サービス：Notion API、SendGrid API
- ステージ切替（staging / production）は環境変数で管理

---

## 実装方針

### Go

- `gofmt` / `goimports` を遵守する。
- 関数はできるだけ単純な責務になるよう分割する。
- エラーは `error` 型で返し、呼び出し元で明示的にハンドリングする。
- Notion プロパティ名・外部 API パラメータなどは定数化してハードコードを避ける。
- 外部 API 呼び出し部分はモック可能な構造にする。

### API 連携

- Notion DB の schema は `docs/HLD.md` の記述に従う。
- SendGrid API の URL・API キーは環境変数から読み込む。
- HTTP 429 などの扱いは `docs/HLD.md` に従う。

### ログ・エラーハンドリング

- Lambda の標準出力に CloudWatch 用ログを出力する。
- ログには少なくともステージ、リクエスト ID、対象 billing_date / 顧客 ID / 請求 ID、外部 API 呼び出し種別と失敗理由を含める。
- リトライしても無意味なエラーは即終了し、十分なログ情報を残す。

---

## 品質方針

- ビジネスロジック層にはユニットテストを追加する。
- 外部 API はインターフェース抽象化によりモックでテスト可能な構造とする。
- 大規模関数や複雑な条件分岐は、機能追加より先にリファクタリングを検討する。
