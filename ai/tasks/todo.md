# TODO

- [x] `AGENTS.md` の責務を運用ルール中心に再構成する
- [x] アーキテクチャ関連の記述を `docs/ARCHITECTURES.md` に分離する
- [x] `ai/tasks/todo.md` `ai/tasks/lesson.md` `ai/tasks/readme.md` を作成する
- [x] 変更後の文言と配置を確認する
- [x] `implement-from-test` スキルを追加する
- [ ] スキルの実行例に合わせて issue 参照導線を整える

## Review

- `AGENTS.md` に編集禁止ルールと記録先ルールを追加した
- `docs/ARCHITECTURES.md` を新設し、技術前提・実装方針・品質方針を移した
- `sed` により各ファイルの内容を確認した
- issue 起点の TDD 用スキル 2 種を `.codex/skills/` に追加した

## 2026-04-06 skills-agents orchestration

- [x] 既存 skill と運用ルールを踏まえて新しい orchestrator 構成を設計する
- [x] 親 skill と補助 skill の `SKILL.md` を追加する
- [x] 対応する `.codex/agents/*/config.toml` を追加する
- [x] 構成と記述内容を確認し、レビュー結果を追記する

### Review

- `final-only` を標準とする親 skill `implement-issue` を追加した
- `plan-issue` `review-subtask` `create-feature-test` `create-subtask-test` `finalize-subtask` を追加し、既存 `implement-from-test` を再利用する構成にした
- `.codex/agents/` に planner / test / implement / review / final integration 用の `config.toml` を追加した
- `find` と `sed` で追加ファイルを確認した
- `config.toml` の構文検証を `python3` で試したが、環境が Python 3.9 かつ `tomli` 未導入のため TOML パーサ検証は未実施
- 用語は `issue = GitHub issue`、`subtask = issue を分解した作業単位` に統一する
- agents 側の `developer_instructions` に参照する skill を明記する
- `final-only` の完了条件を commit ではなく PR 作成完了までに修正する
- issue branch / subtask branch の作業は `git worktree` 前提に修正する
- worktree path 規約を `../notion-invoice-with-lambda-worktrees/issue-<番号>` と `../notion-invoice-with-lambda-worktrees/issue-<番号>-<subtask-slug>` に固定する
- `ai/tasks/todo.md` の計画は issue worktree 作成後に記録する順序へ修正する
- commit message を `<refs|fixes> #<ISSUE_NO> <要約>` + 空行 + `*` 箇条書きに固定する
- PR title / description のテンプレートを固定する
- review に文法エラー確認と改善ルール、`ai/tasks/readme.md` への HTML 遵守事項記録ルールを追加する

## 2026-04-06 issue #1 実装

- [x] `AGENTS.md` `docs/ARCHITECTURES.md` `docs/HLD.md` `ai/tasks/*.md` と issue #1 を確認する
- [x] `issue/1` ブランチを作成する
- [x] issue #1 の受け入れ条件を満たすテスト観点を整理する
- [x] 必要最小限の実装差分でロガーと Lambda エントリの検証性を高める
- [x] 関連テストを実行し、結果を review に記録する

### Review

- issue #1 は「Go Lambda の土台」「共通 JSON ロガー」「最小限 README」「ロガーユニットテスト」が対象
- `docs/HLD.md` のうち今回の直接対象は共通構造とログ方針であり、請求業務ロジック自体はまだ対象外
- `internal/app.Runner` を追加し、各 Lambda handler が concrete type に依存せずテスト可能な形にした
- `cmd/*/main_test.go` で各 Lambda の開始・終了・失敗ログを検証し、`request_id` と `level` がエントリ全体で保たれることを確認した
- `GOCACHE=$(pwd)/.cache/go-build GOMODCACHE=$(pwd)/.cache/gomod go test ./...` は成功した

## 2026-04-06 issue #2 実装計画

- [x] issue #2 の作業範囲を `Config` と Notion ドメイン型の最小整備に限定する
- [x] `internal/config` 相当の設定読み込み・必須値検証を先に固める
- [x] Transactions / Customers / Invoices / Cashflow の内部名に沿った最小 struct を追加する
- [x] JSON マッピングと必須項目欠落時のエラー動作をテストで固定する
- [x] issue 全体の feature-level integration test と review を実行する

### Subtasks

- [x] `config` - ENV, DB IDs, S3, SendGrid, pdftk パスを持つ `Config` を追加し、必須値欠落をエラーにする
- [x] `notion-model` - Transactions / Customers / Invoices / Cashflow の最小 struct を追加し、内部名ベースの JSON マッピングを定義する
- [x] `tests` - 有効値の JSON 逆変換と、必須項目欠落時の失敗を確認する issue-level テストを追加する

### Review

- issue #2 は後続モジュールが共通利用する土台の整備であり、業務ロジックはまだ入れない
- `internal/config.Config` と `LoadFromLookup` を追加し、必須環境変数不足を `MissingEnvironmentError` で返すようにした
- `PDFTK_PATH` は未設定時に `pdftk` を使うデフォルトにし、HLD のコマンド前提を保った
- `internal/domain` に 4 DB 用の最小 Notion record struct と property struct を追加し、内部名 tag を固定した
- `internal/acceptance/issue2_feature_test.go` と unit test により JSON round-trip と設定読み込みを固定した
- `cd /home/user/workspaces/notion-invoice-with-lambda-worktrees/issue-2 && GOCACHE=/tmp/notion-invoice-issue2-go-build GOMODCACHE=/tmp/notion-invoice-issue2-gomod go test ./...` は成功した

## 2026-04-07 agent 修正

- [x] 既存の agent / skill にある worktree 運用と PR title 規約の記述を確認する
- [x] worktree 配下での実行時は branch 作成と worktree 移動を省略する方針を関連ファイルへ反映する
- [x] PR title を `fixes #<ISSUE_NO> <要約>` に統一する
- [x] 差分確認と文言検証を行い、Review に結果を残す

### Review

- worktree 済みの実行では branch 作成と worktree 移動を重ねない方針を `.codex/agents/issue-planner/config.toml` `.codex/agents/final-integrator/config.toml` `.codex/skills/plan-issue/SKILL.md` `.codex/skills/implement-issue/SKILL.md` `.codex/skills/finalize-subtask/SKILL.md` に反映した
- PR title 規約を `fixes #<ISSUE_NO> <要約>` に統一し、`ai/tasks/lesson.md` に恒久ルールとして記録した
- `git diff` と `rg` で旧ルールの取り残しがないことを確認した

## 2026-04-07 branch naming rule

- [x] 現在の branch 命名ルールの定義箇所を確認する
- [x] issue / subtask branch の命名を `issue#<番号>` 系へ変更する
- [x] 変更後の取り残し確認と Review 追記を行う

### Review

- branch 命名規則を `.codex/skills/implement-issue/SKILL.md` で `issue#<番号>` と `issue#<番号>/task/<slug>` に更新した
- 恒久ルールとして `ai/tasks/lesson.md` に branch 命名規則を追記した
- `rg` と `git diff` で運用ルール上の旧表記が残っていないことを確認した
- `ai/tasks/todo.md` 内の `issue/1` は過去実績の記録であり、運用ルールではないため履歴として維持した

## 2026-04-07 subtask commit rule

- [x] subtask commit に関する現行ルールの定義箇所を確認する
- [x] subtask 完了時に commit を必須化するルールを関連 skill / agent / lesson に反映する
- [x] 変更後の差分確認と取り残し確認を行い、Review を追記する

### Review

- `implement-issue` の標準フローと commit ルールを更新し、各 subtask で最低 1 つの `refs #<ISSUE_NO>` commit を必須化した
- `finalize-subtask` に commit 作成ステップを追加し、未 commit のまま subtask を閉じないことを明記した
- `final-integrator` の agent 指示を更新し、subtask 完了時 commit 必須を親 orchestration と整合させた
- `ai/tasks/lesson.md` に恒久ルールとして「subtask は完了時に必ず commit」を追記した
- `rg` と `git diff` で関連ファイルの表記が一致していることを確認した

## 2026-04-07 issue #3 実装計画

- [x] `AGENTS.md` `docs/ARCHITECTURES.md` `docs/HLD.md` `ai/tasks/*.md` と issue #3 を確認する
- [x] `issue#3` 用 worktree `/home/user/workspaces/notion-invoice-with-lambda-worktrees/issue-3` を作成する
- [x] issue #3 の feature-level integration test を追加し、受け入れ条件を固定する
- [x] `http-core` subtask で Notion HTTP クライアント基盤と 429/5xx エラー分類を追加する
- [x] `database-ops` subtask で Transactions / Customers / Invoices / Cashflow の Query / Upsert API を追加する
- [x] subtask ごとの review と commit を完了し、issue branch に取り込む
- [ ] issue 全体テスト、最終 review、最終 commit、PR 作成を行う

### Subtasks

- [x] `notion-client-core` - 共通 request 生成、認証ヘッダ、レスポンス処理、429 / 5xx / 4xx の分類を固定する
- [x] `notion-query-upsert` - DB ごとの Query / Upsert を DTO から HTTP へ変換し、Notion の endpoint と payload を固定する
- [x] `notion-client-tests` - 正常系、429、5xx、DTO→HTTP リクエスト変換を feature-level test で固定する

### Acceptance

- [x] 正常系で Notion Query / Upsert が想定 endpoint、method、header、body で呼ばれる
- [x] 429 は rate limit エラーとして識別され、5xx は server error として識別される
- [x] 正常 / 異常レスポンスで API の返り値とエラーが安定する

### Worktree 方針

- issue worktree は `/home/user/workspaces/notion-invoice-with-lambda-worktrees/issue-3` を使用する
- subtask worktree は `/home/user/workspaces/notion-invoice-with-lambda-worktrees/issue-3-<subtask-slug>` を使用する
- issue worktree と subtask worktree は役割を分離し、同一 worktree の branch 往復はしない

### Review

- issue #3 は Notion API の薄い HTTP 層であり、業務ロジックや再試行戦略はこの issue に含めない
- HLD 4.7 の「429 / 5xx はログ出力して即終了、retry しない」を支えるエラー種別化を土台として実装する
- Query / Upsert は issue 記述を優先し、Transactions / Customers / Invoices / Cashflow ごとの最小 API に留める
- `internal/notion.Client` に共通 HTTP 実行、Authorization / Notion-Version / Content-Type ヘッダ、429 / 5xx / 4xx 分類を追加した
- DB ごとの `Query*` / `Upsert*` wrapper を追加し、create は `/v1/pages` + `parent.database_id`、update は `/v1/pages/{id}` に変換する
- feature-level test と unit test の `httptest` 依存を fake `HTTPDoer` へ置き換え、sandbox でも受け入れ条件を検証できるようにした
- `env GOCACHE=/tmp/notion-issue3-all-go-build GOMODCACHE=/home/user/workspaces/notion-invoice-with-lambda/.cache/gomod GOPROXY=off go test ./...` は成功した
