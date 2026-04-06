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
