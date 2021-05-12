# Experimental 3D model converter

Goで3Dモデルファイルを読み書きするライブラリ＆変換ツールです．

| Format | Read | Write | Comment |
| ------ | ---- | ----- | ------- |
| .mqo   |  ○  |  ○  | ボーン・モーフにも対応 |
| .vrm   |      |  ○  | https://github.com/qmuntal/gltf 用のエクステンションです |
| .glb   |  △  |  ○  | https://github.com/qmuntal/gltf を使っています |
| .pmx   |  ○  |  ○  | Physicsは未対応 |
| .pmd   |  ○  |      | Read only |

データを見ながら雰囲気で実装してるので，おかしな挙動をするかもしれません．

## Usage:

[cmd/modelconv](cmd/modelconv) : モデルデータを相互変換するサンプルプログラム．

Pure golangなのでGoがあればビルドできると思います．
[Releases](https://github.com/m-shimao/modelconv/releases/latest)からビルド済みのWindows用のバイナリがダウンロードできます．

```bash
go build ./cmd/modelconv
```

# License

MIT License
