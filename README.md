# @salesforce/cliにて取得したPackage.xmlを分割する

## 使用方法
ターミナルまたはFinderからこのツールを起動し、以下の情報を入力する
- モード(入力せずEnterの場合はデフォルトモードで実行)
- 分割したいpackage.xmlのパス
- 出力先のパス
- 1ファイルに含まれるコンポーネント数の上限(1〜10000) または 分割したいファイル数 (デフォルト 1)

## 分割モード
カッコ内を入力して分割モードの指定が可能
- デフォルト(default)： 1ファイルに含まれるコンポーネント(members)が指定された数になるようにPackage.xmlを分割する（1〜10000）
- ファイル(files)： 指定されたファイル数になるようにPackage.xmlを分割する（1〜)
- タイプ(types): TypesごとにPackage.xmlを分割する
- サンプル(sample)： Package.xmlのサンプルを生成する 

## Credits

- https://github.com/fatih/color
Copyright (c) 2013 Fatih Arslan
Released under the MIT license
https://opensource.org/licenses/mit-license.php