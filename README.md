# go-workspaces-gui-app
Amazon WorkSpacesの一覧を表示して再起動/起動/停止を行うGUIアプリケーション
  - ついでにEC2インスタンスの一覧も表示して再起動/起動/停止も行うことができる
  - 現状、EC2インスタンスの一覧表示と操作しかできない

## フレームワーク
https://developer.fyne.io/

## ビルド
```bash
make build
```

## 実行
```bash
make start
```

## 使い方
1. Settings画面でAWS CLIにより設定したプロファイルを登録する
2. 使用するプロファイルを選択する
3. Instances画面に遷移するとEC2インスタンスの一覧が表示される
4. インスタンスの状態に応じて以下の操作を実行できる
  - Start： インスタンスを起動する
  - Stop： インスタンスを停止する
  - Reboot： インスタンスを再起動する