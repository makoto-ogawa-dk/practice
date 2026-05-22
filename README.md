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
- Health check (backend): `http://localhost:8080/health`

## 接続モデル（local containers / ARO）
- ブラウザは **frontend のオリジンのみ** にアクセスします。
- フロントエンドの JavaScript は `fetch('/api/...')` の相対パスで API を呼び出します。
- `frontend/nginx.conf` の `/api/` は backend サービスへ内部プロキシされます。

### ローカル（Docker Compose）
- Browser -> `frontend` (`http://localhost:3000`)
- Frontend Nginx -> `backend:8080`（Compose ネットワーク内）

### ARO/OpenShift 想定
- Browser -> Frontend Route/Service
- Frontend Pod の Nginx -> Backend Service（クラスタ内部通信）

このため、利用者（ブラウザ）からは backend の URL を直接意識せず、常に frontend 経由で `/api` を利用します。

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
