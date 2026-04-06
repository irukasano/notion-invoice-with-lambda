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
