---
name: create-feature-test
description: Create feature-level integration tests for issue#N in this repository from docs/HLD.md first, then the GitHub issue. Use when the user asks to define issue-level acceptance tests before subtask implementation starts.
---

# Create Feature Test

このスキルは GitHub issue 全体の受け入れ条件を feature-level integration test として固定する。
期待値の根拠は `docs/HLD.md` を主とし、issue は補助根拠として使う。

## 入力

- `issue#<番号>` を必須入力とする

## 必須参照

1. ユーザーの最新指示
2. `AGENTS.md`
3. 対象の GitHub issue
4. `docs/ARCHITECTURES.md`
5. `docs/HLD.md`
6. 既存の feature test
7. 既存コード

## 出力契約

- issue 全体の feature-level integration test
- 受け入れ条件とテスト観点の記録

## ワークフロー

1. issue を読み、対象範囲を整理する
2. `docs/HLD.md` から該当仕様を抜き出す
3. issue 全体として観測すべき振る舞いを列挙する
4. 期待値を feature-level integration test に落とし込む
5. `ai/tasks/todo.md` に根拠と未確定事項を残す

## ルール

- 現在実装から期待値を逆算しない
- subtask 単位の red test ではなく、issue 全体の受け入れ条件を固定する
- HLD と issue が矛盾したら止まり、人間確認へ切り替える
