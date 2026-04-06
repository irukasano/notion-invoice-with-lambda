---
name: review-subtask
description: Review a completed issue subtask against docs/HLD.md in this repository, focusing on specification fit, security, and performance, then feed review findings back into tests and implementation before merging into the issue branch.
---

# Review Subtask

このスキルは subtask 実装のレビュー専用。
目的は、green になった実装を HLD・security・performance の観点で点検し、必要ならテストへ還元すること。

## 必須参照

1. ユーザーの最新指示
2. `AGENTS.md`
3. GitHub issue
4. `docs/ARCHITECTURES.md`
5. `docs/HLD.md`
6. subtask 用テスト
7. subtask 実装差分

## レビュー観点

1. HLD に一致しているか
2. セキュリティ上の問題がないか
3. 性能上の問題がないか

## ワークフロー

1. subtask の目的と完了条件を確認する
2. 差分とテストを読む
3. HLD とのズレを洗う
4. セキュリティ上の入力、権限、外部 I/O を点検する
5. 性能劣化や不要な処理増加を点検する
6. 再現可能な指摘はテスト観点へ戻す
7. `ai/tasks/todo.md` または `ai/tasks/readme.md` に結果を残す

## 判定ルール

- 指摘が期待値で表現できるなら、まずテストで固定する
- HLD にない拡張提案は issue スコープ外として扱う
- 指摘ゼロでも residual risk があれば残す

## 禁止事項

- 好みだけのリファクタ要求を混ぜること
- 現在実装に合わせて HLD 解釈を変えること
