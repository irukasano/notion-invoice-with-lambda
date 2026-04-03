---
name: implement-from-test
description: Implement a GitHub issue in this repository by starting from existing failing tests, preserving docs/HLD.md and issue requirements, and making the smallest code change needed to pass the tests. Use when the user asks to implement issue#N from tests or continue after create-test.
---

# Implement From Test

このスキルは、このリポジトリでテストを起点に実装を進めるための手順を定義する。
目的は「すでにある failing test を満たす最小差分の実装を行い、仕様の根拠を HLD と issue に保つ」こと。

## 入力

- `issue#<番号>` を必須入力として扱う。
- 対応する failing test が存在する前提で進める。

## 必須参照

この順で確認する。

1. ユーザーの最新指示
2. `AGENTS.md`
3. 対象の GitHub issue
4. `docs/HLD.md`
5. 対応する failing test
6. 関連する実装コード

## 出力契約

このスキルで行うのは次。

- failing test を通すための実装
- 必要最小限のリファクタリング
- テスト実行と結果確認
- 実装上の補助メモを `ai/tasks/readme.md` に記録する

新しい仕様を追加したり、issue 範囲外へ拡張したりしない。

## ワークフロー

1. issue と HLD を確認する
2. failing test が何を要求しているかを整理する
3. 実装コードの変更点を絞る
4. 最小差分で実装する
5. テストを実行して通ることを確認する
6. 必要なら重複除去や命名改善などの小さな整理を行う
7. 実装判断と検証結果を `ai/tasks/readme.md` に残す

## 実装ルール

- まずテストを仕様の固定点として扱う。
- テストが HLD や issue と矛盾する場合は、テストを盲信せず根拠を見直す。
- 実装はテストを満たす最小限に留める。
- 既存コードの全面書き換えは避ける。
- 追加の分岐や構造が不自然なら、小さく整理してもよいが範囲は絞る。

## 禁止事項

- failing test を弱めて通すこと
- 実装に合わせて期待値を変えること
- issue にない別機能をついでに入れること
- HLD と矛盾する独自仕様を足すこと

## 記録

`ai/tasks/readme.md` に少なくとも以下を残す。

- 対象 issue
- 通したテスト
- 実装で選んだ方針
- 残課題や見送り事項

`ai/tasks/todo.md` には、必要なら進捗や検証結果を追記する。

## 典型プロンプト

- `$implement-from-test issue#2 を実装して`
- `$implement-from-test issue#7 の failing test を通して`
