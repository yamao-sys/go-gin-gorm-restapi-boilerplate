# go-gin-gorm-restapi-boilerplate
Go Gin × GORMによるRESTAPIのボイラープレート

## 技術構成
- go
- gorm
- gin
- go-playground/validator
- godotenv
- go-txdb
- stretchr/testify
- go-randomdata
- factory-go/factory

## 機能
### トップ
| | HTTPメソッド | URI | 権限 |
| ------------- | ------------- | ------------- | ------------- |
| 確認用 | GET | / | 権限なし |

### 認証
| | HTTPメソッド | URI | 権限 |
| ------------- | ------------- | ------------- | ------------- |
| 会員登録 | POST | /auth/sign_up | 権限なし |
| ログイン | POST | /auth/sign_in | 権限なし |

### TODOリスト
| | HTTPメソッド | URI | 権限 |
| ------------- | ------------- | ------------- | ------------- |
| 作成 | POST | /todos | 認証済み |
| 一覧取得 | GET | /todos | 認証済み |
| 詳細 | GET | /todos/:id | 認証済み |
| 更新 | PUT | /todos/:id | 認証済み |
| 削除 | DELETE | /todos/:id | 認証済み |

## 環境構築
### 1. 環境変数のファイルの作成
- appディレクトリ配下の.env.sampleをコピーし、.env.developmentとする
- .env.sampleは開発環境をサンプルとしているため、設定値の調整は不要

```
cp app/.env.sample app/.env.development
```

### 2. Dockerのbuild・立ち上げ
- ビルド
```
docker-compose build
```

- 起動
```
docker-compose up -d
```

- コンテナに入る
```
docker-compose exec api_server sh
```

### 3. 開発環境DBの作成・マイグレーション
- 2で起動・コンテナに入った上で、以下のコマンドを実行
```
cd db/migrate

go run migrate.go
```

### 4. Webサーバの起動
- 2の起動・コンテナに入った上で、以下のコマンドを実行
```
cd /app

go run main.go
```

4により、Webサーバの起動ができたら、postman等を使用し、APIアクセスができていることを確認

## 設計方針
- Controller - Service - Repositoryのレイヤードアーキテクチャ
	- DBによるCRUDはRepositoryに寄せる
	- ロジックはServiceに寄せる
	- Controllerはリクエスト・レスポンスのハンドリング
		- リクエスト・レスポンスそれぞれのデータの加工含む
	- 将来的に複数モデルの保存などを行うケースが出てきたら、Service層からTransactionsクラスに切り出したりすると良さそう

## テスト方針
- それぞれの層でテストを書く
	- DB接続があるところはモックを使わずに行う(実環境に近い形でテストする方が不具合検知に役立つため)
		- サービスが大きくなってくると、モックを使用せず結合テストした方がリグレッションテストにも寄与するため

- 正常系/異常系ともに書く
	- 可能な限りC1カバレッジで書きたいところ
	- 事故があるとまずい機能については、C2カバレッジで書いても良さそう

## テスト実行
### テスト用DBの作成・マイグレーション
- テスト用のDBをdbコンテナのホストにログインし、DB名`go_gin_gorm_restapi_boilerplate_test`で作成する
- DBを作成した上で、api_serverコンテナに入った上で、以下のコマンドを実行
```
cd db/migrate

# テスト用のDBのマイグレーション
ENV=test.local go run migrate.go
```

### テスト実行
api_serverコンテナに入った上で、以下のコマンドを実行
```
ENV=test.local go test -v ./...
```
