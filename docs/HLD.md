# 基本設計書（請求書自動生成 Lambda / Notion連携）

本書は、Notion データベース（Transactions / Customers / Invoices / Cashflow）を基盤とし、  
AWS Lambda により請求書 PDF を自動生成し S3 へ保存し、  
必要に応じてメール送信およびキャッシュフロー集計を行うシステムの基本設計をまとめたものである。

---

# 0. 変更履歴

* 2025-12-11: 初版

---

# 1. システム概要

## 1.1 目的

- Notion に登録された売上明細（Transactions）を集計し、顧客ごと・請求予定日ごとに請求書 PDF を自動生成する。
- 生成された請求情報（請求ヘッダ・明細）を Notion Invoices DB に反映し、ステータス管理を一元化する。
- 承認済みの請求書 PDF を SendGrid 経由で顧客へメール送信する。
- 入金／支払期日ベースでキャッシュフロー（Cashflow）を自動集計し、将来の残高推移を可視化する。

## 1.2 業務フロー（概要）

1. ユーザーが Notion DB（Transactions）に売上・支出・残高調整を登録。
2. Lambda（請求PDF作成処理）が、顧客＋請求予定日単位に集計し、請求ヘッダ（Invoices）＋ PDF を生成。
3. ユーザーが Notion 上で内容確認し、`approval_status=承認済み` へ更新。
4. Lambda（承認済請求PDF送信処理）が、承認済みかつ未送信の請求をメール送信。
5. Lambda（キャッシュフロー更新処理）が、入金／支払期日ベースで日次キャッシュフロー（Cashflow）を更新。

## 1.3 自動化範囲

- 自動化する範囲
  - 請求書 PDF の作成・S3 保存・Notion Invoices レコード作成
  - Invoices ⇔ Transactions の Relation 更新・ステータス更新
  - 承認済み請求のメール送信（SendGrid）
  - キャッシュフロー（Cashflow）の再計算・日次レコード作成

- 人手で行う範囲
  - Transactions の登録・編集
  - 顧客（Customers）の登録・編集
  - Invoices の内容確認・承認ステータスの更新（ドラフト → 承認済み）
  - 期初残高レコードの登録（Transactions）

- 対象外（本バージョン）
  - 入金消込の自動化（銀行明細との突合せなど）
  - メール送信結果の詳細トラッキング（開封／クリックなど）
  - 請求書再発行の UI（Notion 側からのトリガは想定せず）

## 1.4 構成技術

- 実行基盤：AWS Lambda（Go1.25 / ZIP デプロイ）
- PDF 生成：GoFPDF（テキストレイヤー）＋ pdftk（テンプレ PDF との合成）
  - 日本語フォント：IPAexGothic を GoFPDF に埋め込み
- API 連携：
  - Notion REST API（Go の生 HTTP 実装）
  - SendGrid API（ベース URL は環境変数から取得）
- ストレージ：S3（PDF 保管）
- ログ・監視：CloudWatch Logs / CloudWatch Alarm
- コード構造：
  - Lambda handler は薄くし、アプリロジック `App.Run()` をユニットテスト可能な構造に分離
  - 外部サービス（Notion / S3 / pdftk / SendGrid）はインターフェイス化

テンプレート PDF は Lambda バイナリに go:embed で埋め込む。

---

# 2. Notion データベース仕様

## 2.1 共通ルール

- Notion 側のプロパティ名は、**日本語名（画面表示）**と**内部名（API 用）**を分けて定義する。
- Lambda は内部名でアクセスする。
- 各 DB ごとに、Lambda が書き込むプロパティと人手が編集するプロパティの責務を明示する。

---

## 2.2 🗃️ DB1：請求／入金明細（Transactions）

請求明細＋支出＋残高調整をまとめた元データテーブル。

### 2.2.1 プロパティ一覧

| 日本語名 | 内部名 | 属性 | 説明 | 更新主体 |
|---------|--------|------|------|---------|
| 件名 | `title` | Title | レコード名 | 人手 |
| 顧客 | `customer` | Relation（Customers） | 売上明細のみ使用、支出・調整では空で可 | 人手 |
| キャッシュ日 | `cash_date` | Date | お金が動く日（入金期日／支払日／調整日） | 人手 |
| 金額 | `amount` | Number | 入金＋／支出−／残高調整＝その時点の残高 | 人手 |
| 種別 | `category` | Select | `売上` / `支出` / `残高調整` / `期初残高` | 人手 |
| キャッシュ対象 | `include_cashflow` | Checkbox | キャッシュフロー集計に含めるか | 人手 |
| 明細内容 | `description` | Rich text | 請求書 PDF にも使う説明 | 人手 |
| 数量 | `qty` | Number | 請求明細用 | 人手 |
| 単価 | `unit_price` | Number | 請求明細用 | 人手 |
| 請求金額 | `line_total` | Formula | `prop("qty") * prop("unit_price")` | Notion |
| 請求予定日 | `billing_date` | Date | 請求書生成の基準日 | 人手 |
| 入金／支払期日 | `payment_due` | Date | 売上：入金期日、支出：支払期日 | 人手 |
| 請求 | `invoice` | Relation（Invoices） | 作成された請求レコードとの紐付け | Lambda |
| ステータス | `status` | Status | `未請求 / 請求作成済 / 入金待ち / 入金済み` | Lambda（一部人手） |

- `status` は基本的に Lambda が `未請求 → 請求作成済` へ更新。
- `入金待ち → 入金済み` などは、必要に応じて人手または別処理で管理（本バージョン対象外）。

---

## 2.3 🗃️ DB2：顧客（Customers）

請求書宛名・メール送信先を管理するマスタ。

### 2.3.1 プロパティ一覧

| 日本語名 | 内部名 | 属性 | 説明 | 更新主体 |
|---------|--------|------|------|---------|
| 顧客名 | `title` | Title | 表示名（請求PDFの宛名にも使用） | 人手 |
| 請求メール宛先 | `billing_email` | Email / Text | メール送信先 | 人手 |
| 請求書宛名 | `billing_to` | Rich text | PDFに印字する宛名（御中など） | 人手 |
| 請求書送付方法 | `send_method` | Select | `メール` / `郵送` など | 人手 |
| 請求書用コード | `customer_code` | Text | 英数2文字、請求番号生成に使用 | 人手 |
| メモ | `memo` | Rich text | 任意メモ | 人手 |

---

## 2.4 🗃️ DB3：請求（Invoices）

毎月 Lambda が作る請求書ヘッダ。

### 2.4.1 プロパティ一覧

| 日本語名 | 内部名 | 属性 | 説明 | 更新主体 |
|---------|--------|------|------|---------|
| 請求番号 | `title` | Title | `S-<customer_code><YYYYMM>-NN` | Lambda |
| 顧客 | `customer` | Relation（Customers） | 顧客情報 | Lambda |
| 請求日 | `billing_date` | Date | 請求書上の発行日（＝グルーピング billing_date） | Lambda |
| 請求対象月 | `billing_period` | Text | 例：`2025-12` | Lambda |
| 請求額合計 | `total_amount` | Rollup | 明細の `line_total` の合計 | Notion |
| 明細 | `items` | Relation（Transactions） | 請求書を構成する明細 | Lambda |
| PDF URL | `pdf_url` | URL / Files | S3 URL またはファイル添付 | Lambda |
| 承認ステータス | `approval_status` | Select / Status | `ドラフト / 承認済み` | 人手 |
| メール送信ステータス | `email_status` | Select | `未送信 / 送信済み / エラー` | Lambda |
| 入金期日 | `payment_due` | Rollup | 明細の `payment_due` の最大値 | Notion |

- `approval_status` はデフォルト `ドラフト` で作成し、ユーザーが `承認済み` へ更新する。
- メール送信処理は `approval_status=承認済み` かつ `email_status=未送信` を対象とする。

---

## 2.5 🗃️ DB4：キャッシュフロー（Cashflow）

Lambda が日ごとに集計して書き込む、キャッシュ推移テーブル。

### 2.5.1 プロパティ一覧

| 日本語名 | 内部名 | 属性 | 説明 | 更新主体 |
|---------|--------|------|------|---------|
| 日付 | `date` | Date | 1日1レコード | Lambda |
| 繰越残高 | `opening_balance` | Number | その日の朝イチの残高 | Lambda |
| 入金合計 | `cash_in` | Number | 当日の入金合計（include_cashflow=true） | Lambda |
| 支出合計 | `cash_out` | Number | 当日の支出合計（include_cashflow=true） | Lambda |
| 純増減 | `net_change` | Number | `cash_in - cash_out` | Lambda |
| 本日残高 | `closing_balance` | Number | 当日の終わりの残高 | Lambda |
| 危険フラグ | `risk_flag` | Formula | しきい値未満なら⚠️など | Notion |

---

# 3. Lambda 機能一覧

1. **請求PDF作成処理**（バッチ）
2. **承認済請求PDF送信処理**（バッチ）
3. **キャッシュフロー更新処理**（バッチ）

各処理は別 Lambda 関数として実装する想定。  
共通ロジック（Notion クライアント、S3 クライアント、SendGrid クライアント等）は共通ライブラリとして分離。

---

# 4. 請求PDF作成処理

## 4.1 実行条件（EventBridge Cron）

Lambda は毎月 **月末** および **月初 1〜3 日** に自動実行する。

- 月末：`cron(10 15 L * ? *)`
- 月初 1〜3 日：`cron(10 15 1-3 * ? *)`

※ UTC 15:10 = JST 00:10 実行。

staging 環境では Cron は Disabled とし、手動実行のみとする（詳細は 7 章）。

---

## 4.2 対象レコード抽出条件（DB1：Transactions）

抽出条件：

- `category = "売上"`
- `status = "未請求"`
- `billing_date <= 実行日 + 1日`（翌日分まで対象とする）
- `customer` が設定されている

※ 未作成の過去分は、上記条件を満たす限りすべて対象とする。

グルーピングキー：

- 顧客（`customer`）
- 請求予定日（`billing_date`）

---

## 4.3 請求番号の仕様

フォーマット：

- `S-<customer_code><YYYYMM>-NN`

例：`S-AB202512-01`

- `customer_code`：Customers.`customer_code`（英数2文字）
- `YYYYMM`：請求日の年月（`billing_date` を基準に `202512` など）
- `NN`：同一顧客・同一年月内での連番（ゼロパディング2桁）

### 4.3.1 採番ロジック

1. 対象グループ（顧客＋請求予定日）に対し、以下で既存 Invoices を検索：
   - `customer` = グループ顧客
   - `billing_date` = グループ請求予定日
   - `title` が `S-<customer_code><YYYYMM>-*` にマッチするもの
2. 既存請求がなければ `-01` を採番。
3. 既存請求があれば、既存の `NN` の最大値を取得し、`+1` した値を NN として採番（`-02`, `-03`, …）。

※ 再発行された請求書でも、元の明細との Relation, PDF URL 等はそれぞれの請求レコード単位で管理する。

---

## 4.4 処理フロー

1. **起動・コンテキスト初期化**
   - 実行日（JST）を取得。
   - `request_id`（UUID）を生成し、全ログに付与。
2. **Transactions 抽出**
   - 4.2 の条件で DB1 をクエリ。
3. **グルーピング**
   - 顧客＋請求予定日単位に明細をグループ化。
4. **顧客情報取得**
   - Customers から必要な情報（宛名・メールアドレス・コード等）を取得。
5. **請求番号採番**
   - 4.3 の仕様に従い、既存 Invoices を検索し NN を決定。
6. **textLayer.pdf 生成（GoFPDF）**
   - IPAexGothic フォントを登録（go:embed した ttf を使用）。
   - 顧客情報・請求番号・明細・金額等を PDF に描画。
   - `/tmp/text.pdf` に保存。
7. **テンプレ PDF との合成（pdftk）**
   - テンプレ PDF を go:embed から `/tmp/template.pdf` に書き出し。
   - `pdftk /tmp/template.pdf background /tmp/text.pdf output /tmp/invoice.pdf` で合成。
8. **S3 へアップロード**
   - 4.5 のキー命名ルールで S3 にアップロード。
9. **Invoices レコード作成（DB3）**
   - 4.6 の項目を設定して Invoices にレコードを作成。
10. **Transactions 更新（DB1）**
    - 対象明細に対し、`status = "請求作成済"` へ更新。
    - `invoice` Relation に作成した請求レコードを紐付け。
11. **ログ出力**
    - 処理件数・請求番号・金額などを `INFO` ログで出力。

---

## 4.5 PDF 出力形式（S3）

S3 バケット名：環境変数 `S3_BUCKET` から取得。

S3 キー命名：

- `<YYYY>/<MM>/<billing_date>_<invoice_number>_<customer_name>_イルカシステム_<total_amount>.pdf`

例：

- `2025/12/2025-12-01_S-AB202512-01_株式会社サンプル_イルカシステム_120000.pdf`

注意事項：

- `customer_name` はファイル名に利用可能な文字に正規化する（記号などは `_` に置換）。
- PDF サイズは `/tmp` の容量（512MB 上限）を超えないよう注意。

---

## 4.6 Invoices（DB3）へ書き込む内容

Lambda が作成する請求レコード：

- `title`（請求番号）  
  - 例：`S-AB202512-01`
- `customer`（Relation）  
  - グループ顧客
- `billing_date`  
  - グルーピングされた請求予定日
- `billing_period`  
  - `billing_date` から `YYYY-MM` を生成
- `items`（Relation）  
  - 対象 Transactions の Relation
- `pdf_url`  
  - アップロードした S3 の URL
- `approval_status`  
  - 初期値 `ドラフト`
- `email_status`  
  - 初期値 `未送信`
- `payment_due`  
  - 対象 Transactions の `payment_due` の最大値

Transactions 側更新：

- `status` = `請求作成済`
- `invoice` Relation に上記 Invoices レコードを紐付け

---

## 4.7 エラー処理・Notion API Rate Limit

- Notion API から **429（Rate Limit）** が返却された場合：
  - 即座に `ERROR` ログ（request_id, エンドポイント, ペイロードの一部）を出力。
  - 処理全体を異常終了とする（CloudWatch Alarm の対象）。
  - リトライは行わない。
- その他の HTTP エラー（5xx 等）：
  - 原則として同様に `ERROR` ログ出力後、処理終了。
- エラー発生時は、部分的に作成された Invoices や更新された Transactions が残りうる。  
  必要に応じて運用で修正する（将来バージョンでの補正処理は検討事項）。

---

# 5. 承認済請求PDF送信処理

## 5.1 対象レコード抽出（DB3：Invoices）

抽出条件：

- `approval_status = "承認済み"`
- `email_status = "未送信"`
- `pdf_url` が設定済み

---

## 5.2 SendGrid 連携仕様

### 5.2.1 環境変数

| 変数名 | 説明 |
|--------|------|
| SENDGRID_BASE_URL | SendGrid API のベース URL（例：`https://api.sendgrid.com/v3`） |
| SENDGRID_API_KEY | SendGrid API キー |
| SENDGRID_FROM_EMAIL | 送信元メールアドレス |
| SENDGRID_FROM_NAME | 送信元表示名（イルカシステム 等） |

### 5.2.2 送信方式

- SendGrid の `/mail/send` エンドポイントを使用（`SENDGRID_BASE_URL` + `/mail/send`）。
- コンテンツ：
  - 件名：`[請求書] <顧客名>様 <billing_period>分 ご請求のご案内`（例）
  - 本文：テキスト／HTML いずれか（簡易なテンプレートをコード内に定義）。
  - 添付：S3 から取得した PDF を BASE64 エンコードして添付する。

（将来的に、S3 URL のみを本文に記載する運用も検討余地あり）

---

## 5.3 処理フロー

1. `approval_status=承認済み` かつ `email_status=未送信` の Invoices を抽出。
2. Customers から `billing_email` を取得。空または不正な場合はエラー扱い。
3. S3 から `pdf_url` の PDF を取得。
4. SendGrid API でメール送信。
5. 成功時：
   - Invoices.`email_status` = `送信済み` へ更新。
6. 失敗時：
   - Invoices.`email_status` = `エラー` へ更新。
   - エラー内容を `ERROR` ログに出力。

---

# 6. キャッシュフロー更新処理

## 6.1 対象レコード・期間

- 対象 DB：DB1（Transactions）
- 対象フラグ：`include_cashflow = true` のレコードのみ集計。
- 対象期間：
  - 実行日時の 1 年前の元日（`YYYY-01-01`）以降、かつ
  - `cash_date` が存在する全レコードを対象とする。

（将来的に期間をパラメータ化する余地あり）

## 6.2 期首日・期初残高の扱い

- 期首日は毎年 **11/1** とする。
- 対象期間内のうち、直近の 11/1 を「期首日」とする。
- 期首日（11/1）の「繰越残高」は、以下のいずれかで決定する：
  1. `category="期初残高"` かつ `cash_date = 11/1` のレコードが存在する場合  
     → その `amount` の合計を「繰越残高」とする。
  2. 上記が存在しない場合  
     → 前日（10/31）までの `closing_balance` を流用するか、0 とする（運用ルールに依存）。  
     ※ 初期は 0 として扱う前提で実装し、運用上は 11/1 に必ず期初残高レコードを登録する想定とする。

- 11/1 に売上／支出がなくても、`cash_in=0`, `cash_out=0`, `net_change=0` となり、  
  `closing_balance = opening_balance` となる。

---

## 6.3 日次集計ロジック

1. DB1 から `include_cashflow = true` かつ対象期間内のレコードを取得。
2. `cash_date` ごとに `amount` を集計：
   - `category="売上"` → `cash_in` に加算（amount > 0 前提）
   - `category="支出"` → `cash_out` に加算（amount < 0 前提だが、集計上は絶対値や正負の扱いを明確化する）
   - `category="残高調整"` → 「売上」「支出」のどちらにも含めず、直接 `opening_balance` に反映するかは将来検討（初期実装では期初残高のみを特別扱い）。
3. 日付順（昇順）に並べ、以下を計算：
   - `net_change = cash_in - cash_out`
   - `opening_balance`：前日 `closing_balance`
   - `closing_balance = opening_balance + net_change`
4. DB4（Cashflow）に対し、日付ごとに 1 レコードを upsert：
   - 既存レコードがあれば更新
   - なければ新規作成

---

# 7. 環境（staging / production）の切り分け

## 7.1 リソース構成

| 区分 | Lambda 関数名例 | S3 バケット例 | IAM ロール |
|------|-----------------|----------------|-------------|
| staging | `invoice-batch-stg` | `invoice-pdf-stg` | `invoice-batch-stg-role` |
| production | `invoice-batch-prod` | `invoice-pdf-prod` | `invoice-batch-prod-role` |

IAM ポリシー例（staging）：

```json
{
  "Effect": "Allow",
  "Action": ["s3:PutObject", "s3:GetObject"],
  "Resource": "arn:aws:s3:::invoice-pdf-stg/*"
}
```

※ production バケットへの誤書き込みを IAM レベルで防ぐ。

## 7.2 Lambda 環境変数

| 変数 | staging | production | 用途 |
|------|---------|------------|------|
| ENV | `stg` | `prod` | 環境識別 |
| S3_BUCKET | `invoice-pdf-stg` | `invoice-pdf-prod` | PDF 保管バケット |
| NOTION_DB_TRANSACTIONS | `<stg用ID>` | `<prod用ID>` | DB1 |
| NOTION_DB_CUSTOMERS | `<stg用ID>` | `<prod用ID>` | DB2 |
| NOTION_DB_INVOICES | `<stg用ID>` | `<prod用ID>` | DB3 |
| NOTION_DB_CASHFLOW | `<stg用ID>` | `<prod用ID>` | DB4 |
| NOTION_TOKEN | stg用 | prod用 | Notion API トークン |
| SENDGRID_BASE_URL | stg用 | prod用 | API ベース URL |
| SENDGRID_API_KEY | stg用 | prod用 | API キー |
| SENDGRID_FROM_EMAIL | 共通 or 別 | 共通 or 別 | 送信元メール |
| SENDGRID_FROM_NAME | 共通 | 共通 | 送信元名 |

アプリコードは ENV による条件分岐を極力行わず、**与えられた値をそのまま使用する**。

## 7.3 スケジュール（EventBridge）

- production：Cron ルールを Enabled  
- staging：Cron ルールは Disabled（手動実行専用）

staging での手動実行例：

```bash
aws lambda invoke --function-name invoice-batch-stg output.json
```

## 7.4 Notion 側の分離

- 可能であれば、staging 用 Notion DB を独立して作成する（ID も別）。
- 単一 DB を共用する場合は、staging 用データには `STG-` 接頭辞を付けるなど明示的に区別し、  
  本番データとの衝突を避ける。

## 7.5 IaC での環境変数管理方針

- CFn/SAM/Terraform/CDK などのテンプレートはリポジトリに含める（環境変数のキー名・参照先のみを記述）。
- 秘密値はテンプレートに直書きせず、SSM Parameter Store や Secrets Manager の参照にする（例：`/prod/notion/token`）。
- 本番値を `.tfvars` や `parameters.json`、`.env` に含めない。ローカル検証用はサンプルファイル（`example.tfvars` など）を用意し、実体は `.gitignore`。
- 公開情報（バケット名や IAM ポリシーなど）はコード化し、キーやトークンは外部ストアから読み込む運用とする。

---

# 8. デプロイ方式

## 8.1 想定コマンド

以下のような形で環境ごとのデプロイを行う：

```bash
TARGET=staging ./deploy.sh
TARGET=production ./deploy.sh
```

または Makefile を用いる場合：

```bash
TARGET=staging make deploy
TARGET=production make deploy
```

## 8.2 `deploy.sh` の役割（例）

1. `TARGET` から環境（stg/prod）を判定。
2. `GOOS=linux GOARCH=amd64` で Go バイナリをビルド。
3. 必要な依存ファイル（テンプレ PDF 等）を含めて ZIP を作成。
4. Lambda 関数ごとに `aws lambda update-function-code` を実行。
5. 成功・失敗を標準出力にログ出力。

※ 実際の実装では、環境ごとの設定（関数名・バケット名など）を `config/<env>.json` などに切り出す想定。

---

# 9. ログ出力・監視設計

## 9.1 ログ出力方針

- Lambda は CloudWatch Logs へ JSON 形式でログ出力。
- ログレベルは `INFO` / `WARN` / `ERROR` を使用。
- Lambda 実行開始時に `request_id`（UUID）を発行し、すべてのログに付与。

ログ例：

```json
{
  "level": "INFO",
  "request_id": "f3a0f4e0-...",
  "message": "invoice generated",
  "customer_code": "AB",
  "invoice_number": "S-AB202512-01",
  "count": 12,
  "total_amount": 120000
}
```

---

## 9.2 エラー出力

- Notion API エラー
- PDF 合成エラー
- S3 アップロードエラー
- SendGrid API エラー

などは、スタックトレースと入力条件（対象顧客・請求番号など）を合わせて `ERROR` として出力する。

---

## 9.3 ロググループ構成

| 環境 | CloudWatch Log Group |
|------|----------------------|
| staging | `/aws/lambda/invoice-batch-stg` |
| production | `/aws/lambda/invoice-batch-prod` |

（メール送信処理・キャッシュフロー処理はそれぞれ別の Log Group を持つ）

---

## 9.4 監視（CloudWatch Alarm）

production 環境には以下のアラームを設定：

- **ErrorCount > 0**
- **Duration > 60 秒**（無限ループ・外部サービスハングの検知）
- **PDF 出力件数が 0 件**（本来出るはずの日に 0 の場合は異常とみなす）

通知先：Slack または SNS → Email。

staging は通知オフまたは低優先度で運用。

---

# 10. テスト方針

## 10.1 ユニットテスト

- Lambda ハンドラは薄くし、`App.Run()` を直接テスト対象とする。
- 実行日（Clock）はインターフェイス経由で注入し、日付に依存しないテストを可能にする。
- Notion / S3 / pdftk / SendGrid 呼び出しはすべてインターフェイス化し、モックで検証する。

## 10.2 結合テスト（staging）

1. staging 環境へデプロイ。
2. テストデータを Notion に投入。
3. staging Lambda を手動 invoke。
4. 以下を確認：
   - S3 に PDF が生成されていること。
   - Invoices にレコードが作成されていること。
   - Transactions のステータス・Relation が更新されていること。
   - SendGrid のメール送信結果（staging では実際の送信先をテスト用に限定）。
   - Cashflow に日次レコードが作成・更新されていること。

localstack 等のモック環境は必須とはせず、staging 環境での結合テストを主とする。

---

# 11. 未決事項・今後の検討

- PDF レイアウト詳細
  - 各項目の座標・フォントサイズ・行数・改ページルール
  - `InvoiceLayout` 構造体として定義予定
- 残高調整（`category="残高調整"`）の Cashflow への反映ルール
- 請求書再発行の UI と運用（既存の `-NN` をどう扱うか）
- メール本文テンプレートの文言・多言語対応
- キャッシュフローの集計期間（現状は「1年前の元日以降」だが、運用実績に応じて調整）

---

以上。
