# manifest-split

`@salesforce/cli`(sf/sfdx)で取得した`package.xml`は、組織の規模が大きくなるとコンポーネント数が膨れ上がり、1回のデプロイ/取得で扱いきれなくなることがあります。**manifest-split**は、そうした巨大な`package.xml`を、用途に応じたルールで複数のファイルに分割するためのシンプルなCLIツールです。

## 目次

- [特徴](#特徴)
- [動作要件](#動作要件)
- [インストール](#インストール)
- [使い方](#使い方)
  - [対話モードで実行する](#対話モードで実行する)
  - [コマンドラインフラグで実行する(非対話)](#コマンドラインフラグで実行する非対話)
- [分割モード](#分割モード)
- [コマンドラインフラグ一覧](#コマンドラインフラグ一覧)
- [`-clean`オプションの挙動](#-cleanオプションの挙動)
- [出力ファイル名のルール](#出力ファイル名のルール)
- [入力チェック・エラーになる条件](#入力チェックエラーになる条件)
- [使用例](#使用例)
- [開発者向け情報](#開発者向け情報)
- [Credits](#credits)

## 特徴

- 📦 `package.xml`を**コンポーネント数**・**ファイル数**・**Typesの種類**のいずれかの基準で分割できる
- ⌨️ ターミナルでの対話入力に加え、**コマンドラインフラグ**による非対話実行にも対応(CI/スクリプトから呼び出し可能)
- 🧹 分割後に残る古い出力ファイルだけを安全に片付ける`-clean`オプション付き(無関係なファイルや分割元ファイルは削除しない)
- 🔁 同じ入力に対して常に同じ並び順で出力する決定的な分割結果
- 🧪 動作確認用の大量データを生成する`sample`モード付き
- 🖥️ macOS / Windows / Linuxで動作(Go製、単一バイナリ)

## 動作要件

- ソースからビルドする場合: [Go](https://go.dev/) 1.23以降
- ビルド済みの実行ファイルを使う場合: 追加の実行環境は不要

## インストール

```bash
git clone https://github.com/kikuchi-null/manifest-split.git
cd manifest-split
```

実行ファイルとしてビルドする場合は以下を実行します(ビルド済みバイナリはリポジトリに含まれていないため、初回はビルドが必要です)。

```bash
# macOS / Linux
go build -o manifest-split .

# Windows
go build -o manifest-split.exe .
```

## 使い方

### 対話モードで実行する

以下のいずれかの方法で起動すると、ターミナル上で対話形式の入力を促されます。

1. manifest-splitのディレクトリで `go run .` を実行する
2. ビルド済みの実行ファイルを実行する: `./manifest-split`(Windowsは`manifest-split.exe`)
3. FinderもしくはExplorerからビルド済みの実行ファイルをダブルクリックして起動する

起動すると、以下の項目を順番に入力します(`sample`モード選択時は「分割したいpackage.xml」の入力はスキップされます)。

| 入力順 | 項目 | 補足 |
|---|---|---|
| 1 | モード | `default` / `files` / `types` / `sample` から選択。**未入力でEnterした場合は`default`モード**になる |
| 2 | 分割したいpackage.xmlのパス | `sample`モードでは入力不要 |
| 3 | 出力先ディレクトリのパス | 存在しない場合は自動的に作成される |
| 4 | コンポーネント数またはファイル数 | `default`/`files`モード(未入力時も含む)でのみ入力。`default`は1〜10000、`files`は1以上 |

```
====== manifest-split ======
入力方法: 入力したらEnter
モード選択(default, files, types, sample): default
分割したいpackage.xml: ./package.xml
出力先: ./dist
1ファイルに含まれるコンポーネント数(1〜10000) または 分割したいファイル数: 1000
```

### コマンドラインフラグで実行する(非対話)

`-mode` `-input` `-output` `-num` の**いずれか1つでも**指定すると、対話プロンプトは一切表示されず、指定したフラグの値だけでそのまま実行されます。CIやシェルスクリプトからの呼び出しに適しています。

```bash
./manifest-split -mode=files -input=./package.xml -output=./dist -num=5
```

対話モードで入力する項目の妥当性チェックは、フラグ経由の実行でも同様に適用されます(詳細は[入力チェック・エラーになる条件](#入力チェックエラーになる条件)を参照)。

## 分割モード

| モード | 分割の基準 | Numの意味 |
|---|---|---|
| `default` | 1ファイルに含まれる**コンポーネント数(members数)**が指定した値以下になるように分割する | 1ファイルあたりの上限コンポーネント数(1〜10000) |
| `files` | 出力される**ファイル数**が指定した値になるように分割する | 出力ファイル数(1以上) |
| `types` | `<types>`(Name)ごとに1ファイルへ分割する。コンポーネント数やファイル数の指定は不要 | 使用しない |
| `sample` | 動作確認用のダミー`package.xml`を生成する(分割元の入力は不要) | 使用しない |

**`default`/`files`モードの内部動作について**: いずれのモードも、まず全コンポーネントを1メンバー単位までいったん分解し、指定の基準で複数ファイルに配分したうえで、各出力ファイル内で同じType名(例: `ApexClass`)のコンポーネントを1つの`<types>`ブロックに統合し直してから書き込みます。そのため、出力される`<types>`ブロックの並び順・統合結果は常に同じになります(実行のたびに順序が変わることはありません)。

`sample`モードで生成される`package.xml`には、以下の8種類のメタデータタイプがそれぞれ1,000件ずつ(合計8,000コンポーネント)含まれます。`default`/`files`モードの分割結果を試したいときの入力データとして利用できます。

`ApexClass`, `ApexTrigger`, `CustomApplication`, `CustomObject`, `CustomField`, `Profile`, `Workflow`, `ValidationRule`

## コマンドラインフラグ一覧

| フラグ | 説明 |
|---|---|
| `-mode` | 分割モード(`default` / `files` / `types` / `sample`)。省略時は`default`扱い |
| `-input` | 分割したいpackage.xmlのパス(`sample`モードでは不要) |
| `-output` | 出力先ディレクトリ |
| `-num` | 1ファイルに含まれるコンポーネント数、または分割したいファイル数 |
| `-clean` | 生成が成功した後、出力先ディレクトリに残る不要な生成物を削除する(対話モードでも独立して有効) |
| `-version` | バージョンを表示して終了する |

## `-clean`オプションの挙動

`-clean`は「出力先ディレクトリを、今回生成した内容だけの状態に片付ける」ためのオプションです。誤って重要なファイルを失わないよう、以下のように安全側に倒した設計になっています。

1. **削除は生成が全て成功した後にのみ行われます。** 入力の読み込みや書き込みが途中で失敗した場合、削除処理自体が実行されないため、既存の(以前の実行による)正常な出力を巻き込んで失うことはありません。
2. **削除対象は、本ツールが生成するファイル名パターン(`package.xml`、または`001_package.xml`のような数字_package.xml)に厳密に一致するファイルのみです。** それ以外のファイル(例: `backup_package.xml`など)は対象になりません。
3. **今回の実行で書き込んだファイルは削除対象から除外されます。** そのため、分割元の`package.xml`を出力先と同じディレクトリに置いていても、今回書き込まれたファイルであれば誤って削除されることはありません。
4. 出力先ディレクトリがまだ存在しない場合は何もしません(エラーにもなりません)。

```bash
# 出力先を最新の分割結果だけの状態に保ちたい場合
./manifest-split -mode=default -input=./package.xml -output=./dist -num=1000 -clean
```

## 出力ファイル名のルール

| ケース | ファイル名 |
|---|---|
| 1ファイルに収まる場合(`default`/`files`) | `package.xml` |
| 複数ファイルに分割される場合(`default`/`files`) | `001_package.xml`, `002_package.xml`, ...(3桁以上のゼロ埋め連番) |
| `types`モード | `<types>`の出現順に`001_package.xml`から連番で採番 |
| `sample`モード | `package.xml` |

## 入力チェック・エラーになる条件

以下の条件に該当する場合、エラーメッセージを表示して終了します(対話モード・フラグ指定のどちらでも共通です)。

| 条件 | エラー内容 |
|---|---|
| モードが`default`/`files`/`types`/`sample`のいずれでもない | モードの指定が不正 |
| `sample`モード以外で、分割対象または出力先が未指定 | 分割対象・出力先の指定が必要 |
| `sample`モードで出力先が未指定 | 出力先の指定が必要 |
| `default`モードでコンポーネント数が1〜10000の範囲外 | コンポーネント数の指定が不正 |
| `files`モードでファイル数が1未満 | ファイル数の指定が不正 |

さらに、`package.xml`の読み込み時には以下も検証されます。

- ファイルが存在しない、またはXMLとして不正
- ルート要素が`<Package>`でない
- `<types>`が1つも存在しない
- いずれかの`<types>`に`<members>`が1つも含まれていない

これらに該当する場合、どのファイルの読み込みに失敗したかを含むエラーメッセージが表示されます。

## 使用例

**基本の分割(対話モード)**

```bash
go run .
```

**1ファイルあたり1,000コンポーネントに分割し、出力先を最新化する**

```bash
./manifest-split -mode=default -input=./manifest/package.xml -output=./dist -num=1000 -clean
```

**5ファイルちょうどに分割する**

```bash
./manifest-split -mode=files -input=./manifest/package.xml -output=./dist -num=5
```

**Typesごとに分割する**

```bash
./manifest-split -mode=types -input=./manifest/package.xml -output=./dist
```

**動作確認用のサンプルpackage.xmlを生成してから分割する**

```bash
./manifest-split -mode=sample -output=./sample
./manifest-split -mode=default -input=./sample/package.xml -output=./dist -num=1000
```

**バージョン確認**

```bash
./manifest-split -version
```

## 開発者向け情報

```bash
# フォーマットチェック
gofmt -l .

# 静的解析
go vet ./...

# ビルド
go build ./...

# テスト実行(カバレッジ付き)
go test ./... -cover
```

主なソース構成は以下の通りです。

- `main.go` — エントリポイント。フラグ解析と各モードへの処理の振り分け
- `ms/argument.go` — 入力(対話/フラグ)の受け取りとバリデーション
- `ms/manifest.go` — `package.xml`の読み込み・分割・書き込み・クリーンアップのコアロジック
- `ms/sample.go` — `sample`モード用のダミーデータ生成
- `ms/constant.go` — モード名・ファイル名・バージョンなどの定数定義

## Credits

- [Taiki Kikuchi](https://github.com/kikuchi-null)
- https://github.com/fatih/color
Copyright (c) 2013 Fatih Arslan
Released under the MIT license
https://opensource.org/licenses/mit-license.php
