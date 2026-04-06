# Lesson

- `AGENTS.md` は運用ルール専用とし、ユーザーから明示依頼がない限り書き換えない
- TODO や実行計画は `ai/tasks/todo.md` に記録する
- ユーザーから指示されて覚えたことは `ai/tasks/lesson.md` に記録する
- LLM 自身の補助メモは `ai/tasks/readme.md` に記録する
- `ai/tasks/todo.md` は毎回作り直さず、issue ごとにセクションを分けて追記する
- 各 issue セクション内では、必要に応じて session 単位の小見出しを切る
- Review や検証結果は issue 単位で `ai/tasks/todo.md` に残す
- `.codex/agents/*/config.toml` は受理されるキーだけを使い、未確認の設定項目を推測で追加しない
- agent role の指示本文キーは `system_prompt` ではなく `developer_instructions` を使う
- 用語は `issue = GitHub issue`、`subtask = issue を分解した作業単位` に統一する
- `final-only` 運用では commit だけで完了扱いにせず、最終 PR 作成までを完了条件として明記する
- issue branch と subtask branch の並行作業は `git worktree` 前提で設計し、同一 worktree で branch を往復しない
- worktree root は `../notion-invoice-with-lambda-worktrees` に固定し、issue は `issue-<番号>`、subtask は `issue-<番号>-<subtask-slug>` に統一する
- commit 1 行目は `<refs|fixes> #<ISSUE_NO> <要約>` とし、issue 完了時だけ `fixes #<ISSUE_NO>` を使う
- commit body は空行のあと `*` 箇条書きで `git diff` ベースに書く
- PR title は `#<ISSUE_NO> <要約>`、PR description は `fixes #<ISSUE_NO>` + `## Summary` `## Changes` `## Comment` に固定する
