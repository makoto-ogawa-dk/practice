# practice

要員山積みWebシステムの初期実装です。

## 構成
- `backend`: Go API サーバ（resources / projects / allocations のCRUD + copy）
- `frontend`: JavaScript静的フロントエンド（画面ナビゲーション + 要員一覧取得）
- `db/migrations`: PostgreSQL DDL
- `docker-compose.yml`: ローカル開発用の起動定義

## ローカル起動（Docker Compose）
```bash
docker compose up --build
```

起動後:
- Frontend: `http://localhost:3000`
- Backend (direct access for debug): `http://localhost:8080`
- Health check (backend direct): `http://localhost:8080/health`
- API access from browser/front: `http://localhost:3000/api/resources`

## フロントエンドとバックエンドの接続方式
- ブラウザは frontend のオリジン（例: `http://localhost:3000`）にのみアクセスします。
- フロントエンドの JavaScript は `/api/...` の相対パスで API を呼び出します。
- frontend コンテナ内の Nginx が `/api/*` を backend サービス（`http://backend:8080`）へリバースプロキシします。

この方式により、ローカル Docker Compose と ARO/OpenShift 配置の両方で、ブラウザ側設定を環境ごとに切り替える必要がなくなります。

## API ルート（初期スキャフォールド）
- `GET/POST /api/resources`
- `PUT/DELETE /api/resources/{resource_id}`
- `POST /api/resources/{resource_id}/copy`
- `GET/POST /api/projects`
- `PUT/DELETE /api/projects/{project_id}`
- `POST /api/projects/{project_id}/copy`
- `GET/POST /api/allocations`
- `PUT/DELETE /api/allocations/{allocation_id}`
- `POST /api/allocations/copy`

## マイグレーション
初期DDLは `db/migrations/001_init.up.sql` にあります。

## 既知の今後対応
- フロントエンドのCRUDフォーム追加
- 入力エラー表示の改善
- 認証や監査情報の追加（必要時）
