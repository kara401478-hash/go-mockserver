# go-mockserver

YAMLファイルを書くだけでHTTP APIのモックサーバーが立ち上がるCLIツールです。

## 特徴

- 📄 YAMLで簡単にエンドポイントを定義
- ⏱️ レスポンス遅延のシミュレーション
- 🔧 カスタムヘッダーの設定
- 🚦 任意のHTTPステータスコードを返せる

## インストール

```bash
git clone https://github.com/yourusername/go-mockserver.git
cd go-mockserver
go mod tidy
```

## 使い方

### 1. YAMLで設定ファイルを書く

```yaml
port: 8080

routes:
  - path: /hello
    method: GET
    response:
      status: 200
      body: '{"message": "Hello, World!"}'

  - path: /users
    method: POST
    response:
      status: 201
      body: '{"id": 1, "message": "作成しました"}'
```

### 2. サーバーを起動する

```bash
go run cmd/main.go --config routes.yaml
```

### 3. リクエストを送る

```bash
curl http://localhost:8080/hello
# {"message": "Hello, World!"}

curl -X POST http://localhost:8080/users
# {"id": 1, "message": "作成しました"}
```

## 設定オプション

| フィールド | 型 | 説明 |
|---|---|---|
| `port` | int | ポート番号（デフォルト: 8080）|
| `routes[].path` | string | URLパス |
| `routes[].method` | string | HTTPメソッド（GET, POST など）|
| `routes[].response.status` | int | HTTPステータスコード |
| `routes[].response.body` | string | レスポンスボディ |
| `routes[].response.headers` | map | カスタムレスポンスヘッダー |
| `routes[].response.delay_ms` | int | 遅延時間（ミリ秒）|

## サンプル

`examples/routes.yaml` にサンプル設定があります。

```bash
go run cmd/main.go --config examples/routes.yaml
```

## ライセンス

MIT
