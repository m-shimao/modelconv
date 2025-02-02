
3Dモデルファイルを変換するコマンドラインツールです．

最新版はGitHubにあります．Windows(x64) 以外のバイナリが必要な場合はお手数ですがご自身でビルドしてください．
https://github.com/m-shimao/modelconv

利用方法はQiitaの記事も参考にしてください．
https://qiita.com/m-shimao/items/d29cd21b9860809f72cf

# Usage

以下の組み合わせの変換ができます．

- (.pmd | .pmx | .mqo) → (.mqo | .pmx | .glb | .vrm)
- .glb → .vrm

座標の単位については以下のように扱っています(異なる場合は変換時の-scaleオプションで調整してください)

- MQO: 1mm
- MMD: 80mm
- glTF/VRM: 1m

vrmファイルを出力するためには，設定が書かれたjsonファイルが必要です．
jsonファイルの内容はQiitaに書いた説明を参考にしてください．

## 例

### .pmx → .mqo
./modelconv -rot180 -scale 80 -autotpose 右腕,左腕 input.pmx output.mqo

-rot180 : Y軸回りに100°回転します
-scale : 変換時にスケールをかけます
-autotpose : 指定した腕のボーンがX軸に沿うように調整します(暫定実装)

入力ファイル以外は省略可能なので，ファイルを実行ファイルにドラッグ＆ドロップすれば変換されます．

MMDから変換する場合，scaleを省略すると単位をmmにするために暗黙的に"-scale 80"とみなします．

### .glb+vrmconfig.json → .vrm
./modelconv -vrmconfig input.vrmconfig.json input.glb output.vrm

glbからvrmへの変換は特別扱いしているので，スケールや回転は指定できません．

# License

https://github.com/m-shimao/modelconv
  MIT License (https://github.com/m-shimao/modelconv/blob/master/LICENSE)

glTFの読み書きに以下のライブラリを使用しています．

qmuntal/gltf https://github.com/qmuntal/gltf
  BSD 2-Clause (https://github.com/qmuntal/gltf/blob/master/LICENSE)
