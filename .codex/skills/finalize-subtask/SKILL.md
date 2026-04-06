---
name: finalize-subtask
description: Finalize a subtask under issue#N in this repository after implementation and review by applying review feedback, rerunning subtask tests, and preparing the task branch to merge back into the issue branch.
---

# Finalize Subtask

このスキルは subtask branch の仕上げを担当する。
目的は、review 指摘を反映し、subtask を issue branch に戻せる状態まで閉じること。

## 必須参照

1. ユーザーの最新指示
2. `AGENTS.md`
3. 対象の GitHub issue
4. 対象 subtask の定義
5. `docs/ARCHITECTURES.md`
6. `docs/HLD.md`
7. subtask 用テスト
8. subtask branch 上の差分

## 出力契約

- review 指摘の反映結果
- subtask テストの実行結果
- issue branch へ戻す準備

## ワークフロー

1. review 指摘を整理する
2. 必要ならテストへ還元する
3. 実装を修正する
4. subtask 用テストと関連回帰テストを再実行する
5. merge 前提の残課題がないか確認する

## ルール

- subtask 完了の判定は feature-level integration test ではなく subtask テストで行う
- issue 全体の最終統合は親 skill `implement-issue` が担う
