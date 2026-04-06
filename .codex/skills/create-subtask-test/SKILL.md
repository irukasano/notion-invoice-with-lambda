---
name: create-subtask-test
description: Create failing tests for a subtask under issue#N in this repository by deriving expected behavior from docs/HLD.md, the GitHub issue, and the approved subtask definition. Use when the user asks to prepare red tests for a specific subtask branch.
---

# Create Subtask Test

このスキルは subtask ごとの red テストを作る。
目的は、承認済みの subtask 定義をもとに、実装前に期待振る舞いを failing test として固定すること。

## 入力

- `issue#<番号>` を必須入力とする
- subtask 名または subtask の完了条件が与えられている前提で進める

## 必須参照

1. ユーザーの最新指示
2. `AGENTS.md`
3. 対象の GitHub issue
4. 承認済みの subtask 定義
5. `docs/ARCHITECTURES.md`
6. `docs/HLD.md`
7. 既存のテストコード
8. 実装コード

## 出力契約

- subtask 用の failing test
- 必要なら最小限の fixture / helper
- テスト観点の記録

## ワークフロー

1. subtask の完了条件を確認する
2. HLD と issue から該当仕様を抜き出す
3. subtask で固定すべき期待振る舞いを分解する
4. failing test を追加する
5. `ai/tasks/todo.md` に根拠と未確定事項を残す

## ルール

- feature-level integration test と混同しない
- subtask 単位で red/green が閉じる粒度にする
- 現在実装の挙動を期待値にしない
- HLD にない仕様を自然補完しない
