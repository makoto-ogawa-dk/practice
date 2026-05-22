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
- Browser からの API 呼び出し: `http://localhost:3000/api/...`（frontend の Nginx が backend にプロキシ）
- Backend（直接確認用）: `http://localhost:8080`
- Health check（直接確認用）: `http://localhost:8080/health`

## Frontend-Backend 接続方針（local container / ARO）
- ブラウザは frontend のオリジンのみにアクセスします。
- フロントエンドの JavaScript は相対パス `/api/...` で API を呼び出します。
- frontend コンテナの Nginx が `/api/*` を backend サービスへリバースプロキシします。
- ARO/OpenShift でも同じ構成（Browser -> frontend Route/Service -> `/api` proxy -> backend Service）で運用できます。

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
