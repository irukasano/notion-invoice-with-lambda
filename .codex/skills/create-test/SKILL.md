---
name: create-test
description: Create tests for a GitHub issue in this repository by deriving expected behavior from docs/HLD.md first, then the issue, and never reverse-engineering expectations from the current implementation. Use when the user asks to create tests for issue#N or wants HLD-first test design.
---

# Create Test

このスキルは、このリポジトリで issue に対応するテストを先に作るための手順を定義する。
目的は「実装から期待値を起こさず、`docs/HLD.md` と issue を根拠に failing test を作る」こと。

## 入力

- `issue#<番号>` を必須入力として扱う。
- 入力に issue 番号がなければ、このスキルだけでは進めない。

## 必須参照

この順で確認する。

1. ユーザーの最新指示
2. `AGENTS.md`
3. 対象の GitHub issue
4. `docs/HLD.md`
5. 既存のテストコード
6. 実装コード

実装コードは配置確認や API 面の整合確認には使ってよいが、期待仕様の根拠には使わない。

## 出力契約

このスキルで行うのは次まで。

- issue に対応するテストの新規作成または更新
- 必要ならテスト補助用の最小限の fixture / helper の追加
- テスト設計の記録を `ai/tasks/todo.md` に残す

このスキルでは原則として本実装を行わない。
テストを通すための製品コード変更は別スキルに委ねる。

## ワークフロー

1. issue を読む
2. `docs/HLD.md` から該当仕様を抜き出す
3. issue と HLD の差分を確認する
4. 受け入れ条件をテスト観点に分解する
5. 既存テストの流儀に合わせて failing test を作る
6. 追加したテストと根拠を `ai/tasks/todo.md` に記録する

## 判定ルール

- HLD に明記されている内容はそのまま期待値にする。
- issue に追加要件があり、HLD と矛盾しない場合は issue を補助根拠として使う。
- issue が HLD と矛盾する場合は、勝手に仕様変更せず矛盾点を明示する。
- HLD と issue の両方で根拠が弱い場合は、推測で広げず不足点を列挙する。

## テスト作成ルール

- 既存実装の現在挙動に引きずられない。
- テスト名は期待する振る舞いが読める名前にする。
- 1 テスト 1 振る舞いを基本にする。
- 境界条件、異常系、未指定時の既定動作を必要に応じて含める。
- 失敗理由が読めるアサーションを書く。
- 既存テストの書式や helper があればそれに合わせる。

## 禁止事項

- 実装コードを読んで「今こう動いているから」を期待値にすること
- テストを通すために製品コードを先に修正すること
- HLD にない仕様を自然に補完して固定すること

## 記録

`ai/tasks/todo.md` に少なくとも以下を残す。

- 対象 issue
- 参照した HLD の章
- 作成したテスト観点
- 未確定事項や仕様ギャップ

`ai/tasks/lesson.md` は、ユーザーから新しい恒久ルールを指示されたときだけ更新する。

## 典型プロンプト

- `$create-test issue#2 のテストを作成して`
- `$create-test issue#5 に対応する failing test を HLD ベースで追加して`
