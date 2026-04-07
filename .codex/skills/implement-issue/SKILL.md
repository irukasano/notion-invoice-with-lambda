---
name: implement-issue
description: Orchestrate end-to-end implementation for issue#N in this repository using docs/HLD.md as the source of truth. Use when the user asks to implement an issue with subtask decomposition, HLD-based integration tests, iterative review, and one final PR from the issue branch.
---

# Implement Issue

このスキルは issue 全体を親オーケストレータとして進める。
標準運用は `final-only` で、subtask ごとの PR は作らず、最後に issue branch から 1 本だけ PR を作る。

## 入力

- `issue#<番号>` を必須入力とする
- 明示がない限り `pr_mode=final-only` とする
- branch 命名は `issue/<番号>` と `issue/<番号>/task/<slug>` を使う
- branch 作業は `git worktree` を前提とする
- issue branch と各 subtask branch は別 worktree に分ける
- worktree root は `../notion-invoice-with-lambda-worktrees` に固定する
- issue worktree は `../notion-invoice-with-lambda-worktrees/issue-<番号>` を使う
- subtask worktree は `../notion-invoice-with-lambda-worktrees/issue-<番号>-<subtask-slug>` を使う

## 必須参照

1. ユーザーの最新指示
2. `AGENTS.md`
3. `docs/ARCHITECTURES.md`
4. `docs/HLD.md`
5. `ai/tasks/*.md`
6. GitHub issue
7. 既存コード

## 停止条件

- HLD と issue が矛盾する
- feature-level integration test の期待値が HLD から一意に定まらない
- subtask の分け方で設計差分が大きい
- セキュリティや運用権限の前提が変わる
- `stacked PR` へ切り替える必要がある

上記では、曖昧点・選択肢・推奨案を整理して人間に確認する。
推測で仕様を追加しない。

## 標準フロー

1. まず現在地が対象 issue の専用 worktree 配下か判定し、そうでなければ `git worktree` を使って `../notion-invoice-with-lambda-worktrees/issue-<番号>` の issue worktree を作る
2. issue worktree の作成または既存 worktree の確認後に `ai/tasks/todo.md` へ計画を記録する前提を整える
3. `plan-issue` を使い、issue と HLD から subtask を定義する
4. `create-feature-test` を使い、issue branch に HLD ベースの feature-level integration test を追加する
5. subtask ごとに `../notion-invoice-with-lambda-worktrees/issue-<番号>-<subtask-slug>` の worktree を用意する
6. 各 subtask で `create-subtask-test` を使って red テストを作る
7. 各 subtask で `implement-from-test` を使って green にする
8. 各 subtask で `review-subtask` を使って HLD 適合性・security・performance を確認する
9. review 指摘のうちテスト化できるものは red テストへ戻し、実装とレビューを反復する
10. 各 subtask で `finalize-subtask` を使い、branch を issue branch へ戻せる状態に整える
11. subtask 完了後は subtask worktree 上で commit し、issue worktree へ取り込む
12. 全 subtask 完了後に feature-level integration test と issue 全体レビューを実行する
13. 問題がなければ issue branch で commit し、最後に 1 本だけ PR を作る

## 判定ルール

- `subtask green` は、その subtask 用テストと関連回帰テストが green の状態を指す
- `issue green` は、feature-level integration test を含めて green の状態を指す
- feature-level integration test が最後まで red のままでも、subtask green の判定には使わない
- review 指摘は `HLD` 違反、`security`、`performance` の順で優先する
- 文法エラー、lint エラー、ビルド不能状態は review で必ず確認する
- よりよい方向性が明確で issue スコープ内なら、積極的に改善を取り込む
- 確認でエラーが出た場合は修正し、エラーが残ったまま完了にしない
- 他の実装タスクに共有しないと全体の統一が崩れる遵守事項は `ai/tasks/readme.md` に記録する
- 改善の結果、共通テンプレート変更が必要な場合だけ既存実装へ波及修正する

## Git / PR ルール

- 標準は `final-only`
- subtask ごとの PR は明示指示がある場合か、長期化・広範囲変更・早期レビュー需要がある場合だけ使う
- branch の作成と切替は `git switch` の往復ではなく `git worktree add` を優先する
- ただし、すでに対象 issue / subtask の専用 worktree 配下にいる場合は branch 作成や worktree への移動を重ねて行わない
- issue branch は専用 worktree、各 subtask branch も専用 worktree を持つ前提で進める
- worktree root は `../notion-invoice-with-lambda-worktrees` に固定する
- path 名は issue 用が `issue-<番号>`、subtask 用が `issue-<番号>-<subtask-slug>` の kebab-case に固定する
- `final-only` の完了条件は `issue/<番号> -> main` の PR 作成完了までとする
- commit 済みでも PR 未作成なら完了扱いにしない
- GitHub 認証や権限不足で PR を作れない場合は、そこで停止して不足条件を人間に確認する
- `PR はまだ作成していません` という状態で締めない

## コミットルール

- 1 行目は `<prefix> #<ISSUE_NO> <要約>` の形式に固定する
- `prefix` は `refs` または `fixes` だけを使う
- issue 全体が完了して最終 PR に載る commit は `fixes #<ISSUE_NO>` を使う
- subtask の途中 commit や参照目的の commit は `refs #<ISSUE_NO>` を使う
- `refs issue1` のような表記は使わない
- 要約は日本語 36 文字以内にする
- commit message は次の 3 段構成にする

```text
<prefix> #<ISSUE_NO> <要約>

* git diff をもとにした変更点 1
* git diff をもとにした変更点 2
* git diff をもとにした変更点 3
```

- 2 行目は空行にする
- 箇条書きは `*` で始め、1 行 60 文字前後でまとめる
- 箇条書きは実際の `git diff` を要約し、抽象的な感想にしない

## PR ルール

- PR title は `fixes #<ISSUE_NO> <要約>` の形式に固定する
- PR description は次の形式に固定する

```md
fixes #<ISSUE_NO>

## Summary
1. 変更内容の要約
2. 変更内容の要約
3. 変更内容の要約

## Changes
### <commit-hash> <commit-first-line>
<commit body の詳細>

## Comment
- 人間にレビューしてほしい観点
```

- `fixes` の場合、description 先頭行は `fixes #<ISSUE_NO>` のみを書く
- `## Summary` は `git diff` をもとに日本語で 200 文字以内にまとめる
- `## Summary` の箇条書きは `1. 2. 3.` 形式で書く
- `## Changes` には commit ごとに `### <hash> <1行目>` を並べ、その下に commit message の詳細を転記する
- `## Comment` にはレビュー時に注目してほしい点を箇条書きで書く

## エージェントの使い分け

- 分解は `issue-planner`
- feature test と subtask red テスト設計は `test-designer`
- 実装は `implementer`
- 観点レビューは `reviewer`
- subtask の仕上げと issue 全体の最終統合は `final-integrator`

## 記録

- 計画と進捗は `ai/tasks/todo.md`
- 実装補助メモは `ai/tasks/readme.md`
- ユーザーから新しい恒久ルールを受けた場合だけ `ai/tasks/lesson.md`

## 典型プロンプト

- `@implement issue#1 を実装して`
- `@implement issue#4 を final-only で進めて`
