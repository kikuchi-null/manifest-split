# @salesforce/cliにて取得したPackage.xmlを分割する  

## インストール
```bash
git clone https://github.com/kikuchi-null/manifest-split.git
```

## 使用方法
以下いずれかの方法でmanifest-splitを実行する
1. manifest-splitのディレクトリに移動し、 `go run .`
2. manifest-splitのディレクトリに移動し、 `./manifest-split` または `./manifest-split.exe`
3. FinderもしくはExplorerからmanifest-splitを起動する

manifest-split起動後の入力は以下の通り
- モード(入力せずEnterの場合はデフォルトモードで実行)
- 分割したいpackage.xmlのパス
- 出力先のパス
- 1ファイルに含まれるコンポーネント数の上限(1〜10000) または 分割したいファイル数 (デフォルト 1)

## 分割モード
以下の入力値で分割モードの指定が可能  

<table>
  <thead>
    <tr>
      <th>入力値</th> <th>説明</th>
    </tr>
  </thead>
  <tr>
    <td>default</td> <td>1ファイルに含まれるコンポーネント(members)が指定された数になるようにPackage.xmlを分割する（1〜10000）</td>
  </tr>
   <tr>
    <td>files</td> <td>指定されたファイル数になるようにPackage.xmlを分割する（1〜)</td>
  </tr>
   <tr>
    <td>types</td> <td>TypesごとにPackage.xmlを分割する</td>
  </tr>
   <tr>
    <td>sample</td> <td>Package.xmlのサンプルを生成する</td>
  </tr>
</table>

## Credits

- [Taiki Kikuchi](https://github.com/kikuchi-null)
- https://github.com/fatih/color  
Copyright (c) 2013 Fatih Arslan  
Released under the MIT license  
https://opensource.org/licenses/mit-license.php
