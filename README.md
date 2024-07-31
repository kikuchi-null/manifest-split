# @salesforce/cliにて取得したPackage.xmlを分割する

## 使用方法
ターミナルまたはFinderからこのツールを起動し、以下の情報を入力すえう
- -input string：分割したいpackage.xmlのパス
- -output string：出力先のパス
- -mode string：分割モード分割モード(Enterでデフォルトモード)
- -n int：1ファイルに含まれるコンポーネント数の上限(最大1万) または 分割したいファイル数 (default 1)

## モード(入力値)
- デフォルト(default)： 1ファイルに含まれるコンポーネント(members)が指定された数になるようにPackage.xmlを分割する（1〜10000）
- ファイル(files)： 指定されたファイル数になるようにPackage.xmlを分割する（1〜)
- タイプ(types): TypesごとにPackage.xmlを分割する
- サンプル(sample)： Package.xmlのサンプルを生成する 
