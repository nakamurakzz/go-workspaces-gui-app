# go-workspaces-gui-app
Amazon WorkSpacesの一覧を表示して再起動/起動/停止を行うGUIアプリケーション
  - ついでにEC2インスタンスの一覧も表示して再起動/起動/停止も行うことができる

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
3. EC2 Instancesタブに遷移するとEC2インスタンスの一覧が表示される
- インスタンスの状態に応じて以下の操作を実行できる
  - Start： インスタンスを起動する
  - Stop： インスタンスを停止する
  - Reboot： インスタンスを再起動する
5. Workspacesタブに遷移するとWorkspacesの一覧が表示される
- Workspaceの状態に応じて以下の操作を実行できる
  - Reboot： Workspaceを再起動する