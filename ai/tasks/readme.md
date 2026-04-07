# AI Task Notes

- LLM 自身の補助メモや一時的な運用記録をここに追記する
- 他の実装タスクに共有しないと全体の統一が崩れる遵守事項があればここに追記する
- 改善の結果、共通テンプレート変更が必要な場合だけ既存実装へ波及修正する
- 2026-04-03: `.codex/skills/implement-from-test` を追加
- `implement-from-test` は failing test を満たす最小実装に限定する
- 2026-04-06: issue 全体の受け入れ条件は `create-feature-test`、subtask の red test は `create-subtask-test` に分離した
- 2026-04-06 issue #2: `internal/config` は `LoadFromLookup` を公開し、テストから環境変数依存を切り離す
- 2026-04-06 issue #2: `PDFTK_PATH` は未指定時に `pdftk` を使うデフォルトで扱い、環境差分を最小化する
