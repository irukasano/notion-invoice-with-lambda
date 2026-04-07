---
name: plan-issue
description: Plan issue#N implementation in this repository from docs/HLD.md and the GitHub issue, decompose it into subtasks, and create feature-level integration tests on the issue branch before subtask implementation starts.
---

# Plan Issue

このスキルは issue 全体の実装計画を固める。
目的は、HLD を根拠に subtask 分解と feature-level integration test を先に確定すること。

## 必須参照

1. ユーザーの最新指示
2. `AGENTS.md`
3. 対象の GitHub issue
4. `docs/ARCHITECTURES.md`
5. `docs/HLD.md`
6. `ai/tasks/*.md`
7. 既存テストとコード

## 出力契約

- issue branch 用の実行計画
- subtask 一覧と受け入れ条件
- HLD ベースの feature-level integration test の観点
- issue worktree と subtask worktree の切り方
- worktree root と命名規約
- 未確定事項と確認要否

## ワークフロー

1. まず現在地が対象 issue の専用 worktree 配下か判定し、そうでなければ issue worktree `../notion-invoice-with-lambda-worktrees/issue-<番号>` を作成して移動する
2. issue の要求を整理する
3. `docs/HLD.md` から該当章を抜き出す
4. HLD を起点に feature-level の期待振る舞いを列挙する
5. 期待振る舞いを subtask に分解する
6. subtask ごとに境界、依存、完了条件を明文化する
7. worktree root `../notion-invoice-with-lambda-worktrees` 配下の配置を決める
8. issue worktree と subtask worktree の対応を決める
9. `create-feature-test` で固定すべき観点を明文化する
10. issue worktree 上の `ai/tasks/todo.md` に計画を記録する

## 分解ルール

- subtask は単独で red test と green 実装が閉じる粒度にする
- 複数コンポーネントに跨っても、期待振る舞いが 1 つなら 1 subtask としてよい
- `土台の整理` と `業務仕様の追加` を混ぜない
- feature-level integration test の期待値は HLD を根拠に固定し、現在実装から逆算しない
- subtask slug は短い kebab-case に固定する

## 注意

- subtask green の判定に feature-level integration test を使わない
- HLD と issue が矛盾したら、人間確認へ切り替える
- すでに対象 issue の専用 worktree 配下にいる場合は branch 作成や worktree への移動を重ねて行わない
- `ai/tasks/todo.md` の計画は issue worktree の作成または既存 worktree の確認後に書き始める
