# 要員山積みWebシステム 基本仕様書

## 1. システム目的
要員の山積み（リソースマネジメント）を円滑に行うためのWebシステムを構築する。

## 2. 利用者
- システム利用者：1名（依頼者ご本人）

## 3. 主な機能
- 月ごと、人ごと、案件ごとに工数を1つの表で表示する
- 表示項目（列）の登録・編集機能
- 人、案件のマスタ登録機能

## 4. 制約条件
- フロントエンド：JavaScript
- バックエンド：Go言語
- 両方ともARO（Azure Red Hat OpenShift）上で稼働するコンテナとして展開

## 5. 想定画面イメージ
- 工数表（人×案件×月）の一覧表示
- マスタ登録・編集画面

## 6. セキュリティ／運用
- シングルユーザのため認証等はシンプルな仕組みを想定
- バックアップは不要
- ログ取得は最低限でOK

## 7. 確認事項・今後の課題
- AROへのデプロイに必要な設定項目・方式調査
- 一部要件の詳細化（例：マスタ項目の属性、画面レイアウトの具体化など）

## 8. 画面仕様

### 8.1 山積み一覧画面
#### 目的
月ごと・人ごと・案件ごとの工数を一覧で確認し、登録・編集を行う。

#### 主な表示項目
- 月
- 要員名
- 案件名
- 工数
- 備考

#### 主な操作
- 条件検索
  - 月指定
  - 要員指定
  - 案件指定
- 一覧表示
- 工数の新規登録
- 工数の編集
- 工数データの複製

#### 補足
- 山積みデータの複製機能では、既存月のデータを別月へ複製できるものとする。
- 複製後に工数や備考を修正できるものとする。

---

### 8.2 要員マスタ管理画面
#### 目的
要員情報を登録・編集・削除する。

#### 主な表示項目
- 要員ID
- 要員名
- 所属
- 備考

#### 主な操作
- 新規登録
- 編集
- 削除
- 複製
- 一覧表示

---

### 8.3 案件マスタ管理画面
#### 目的
案件情報を登録・編集・削除する。

#### 主な表示項目
- 案件ID
- 案件名
- 開始月
- 終了月
- ステータス
- 備考

#### 主な操作
- 新規登録
- 編集
- 削除
- 複製
- 一覧表示

---

## 9. データ項目一覧

### 9.1 要員マスタ
| 項目名 | 論理名 | 型 | 必須 | 説明 |
|---|---|---|---|---|
| 要員ID | resource_id | string/int | 必須 | 要員を一意に識別するID |
| 要員名 | resource_name | string | 必須 | 表示用の氏名 |
| 所属 | department | string | 任意 | 所属部署名 |
| 備考 | note | string | 任意 | 補足情報 |

---

### 9.2 案件マスタ
| 項目名 | 論理名 | 型 | 必須 | 説明 |
|---|---|---|---|---|
| 案件ID | project_id | string/int | 必須 | 案件を一意に識別するID |
| 案件名 | project_name | string | 必須 | 案件名 |
| 開始月 | start_month | yyyy-mm | 任意 | 案件開始予定月 |
| 終了月 | end_month | yyyy-mm | 任意 | 案件終了予定月 |
| ステータス | status | string | 任意 | 進行中、完了など |
| 備考 | note | string | 任意 | 補足情報 |

---

### 9.3 山積みデータ
| 項目名 | 論理名 | 型 | 必須 | 説明 |
|---|---|---|---|---|
| データID | allocation_id | string/int | 必須 | 山積みデータを一意に識別するID |
| 月 | target_month | yyyy-mm | 必須 | 対象月 |
| 要員ID | resource_id | string/int | 必須 | 要員マスタ参照 |
| 案件ID | project_id | string/int | 必須 | 案件マスタ参照 |
| 工数 | workload | decimal | 必須 | 月間工数 |
| 備考 | note | string | 任意 | 補足情報 |

## 10. 入力チェック仕様（詳細）

### 10.1 共通方針
- 必須項目が未入力の場合はエラーとする。
- 文字列項目は前後の空白を除去して登録する。
- ID項目はシステム内部で一意に管理する。
- エラー発生時は、対象項目とエラーメッセージを画面上に表示する。
- 登録時および更新時に入力チェックを行う。

---

### 10.2 要員マスタ管理画面

| 項目名 | チェック内容 | エラー条件 | 最大文字数 | 備考 |
|---|---|---|---|---|
| 要員名 | 必須 | 未入力の場合 | 100 | |
| 要員名 | 桁数 | 100文字を超える場合 | 100 | |
| 所属 | 桁数 | 100文字を超える場合 | 100 | 任意 |
| 備考 | 桁数 | 500文字を超える場合 | 500 | 任意 |

---

### 10.3 案件マスタ管理画面

| 項目名 | チェック内容 | エラー条件 | 最大文字数 | 備考 |
|---|---|---|---|---|
| 案件名 | 必須 | 未入力の場合 | 100 | |
| 案件名 | 桁数 | 100文字を超える場合 | 100 | |
| 開始月 | 形式 | yyyy-mm 形式でない場合 | - | 任意 |
| 終了月 | 形式 | yyyy-mm 形式でない場合 | - | 任意 |
| 開始月/終了月 | 相関 | 開始月が終了月より後の場合 | - | 両方入力時のみチェック |
| ステータス | 桁数 | 50文字を超える場合 | 50 | 任意 |
| 備考 | 桁数 | 500文字を超える場合 | 500 | 任意 |

---

### 10.4 山積み一覧画面

| 項目名 | チェック内容 | エラー条件 | 最大文字数 | 備考 |
|---|---|---|---|---|
| 月 | 必須 | 未入力の場合 | 7 | yyyy-mm形式 |
| 月 | 形式 | yyyy-mm 形式でない場合 | 7 | yyyy-mm形式 |
| 要員ID | 必須 | 未選択の場合 | - | 要員マスタに存在すること |
| 案件ID | 必須 | 未選択の場合 | - | 案件マスタに存在すること |
| 工数 | 必須 | 未入力の場合 | - | 数値 |
| 工数 | 数値 | 数値でない場合 | - | |
| 工数 | 範囲 | 0未満の場合 | - | 上限は組織運用にあわせて別途可 |
| 備考 | 桁数 | 500文字を超える場合 | 500 | 任意 |
| 月×要員ID×案件ID | 一意性 | 同一組み合わせがすでに存在する場合 | - | 重複登録不可 |

---

### 10.5 エラーメッセージ例

| 項目名 | 条件 | メッセージ例 |
|---|---|---|
| 要員名 | 未入力 | 要員名を入力してください。 |
| 要員名 | 文字数超過 | 要員名は100文字以内で入力してください。 |
| 所属 | 文字数超過 | 所属は100文字以内で入力してください。 |
| 案件名 | 未入力 | 案件名を入力してください。 |
| 月 | 形式不正 | 月は yyyy-mm 形式で入力してください。 |
| 開始月/終了月 | 前後関係不正 | 開始月は終了月以前を入力してください。 |
| 工数 | 未入力 | 工数を入力してください。 |
| 工数 | 数値不正 | 工数は数値で入力してください。 |
| 工数 | 範囲不正 | 工数は0以上で入力してください。 |
| 備考 | 文字数超過 | 備考は500文字以内で入力してください。 |
| 月×要員×案件 | 重複 | 同一の月・要員・案件の組み合わせは重複登録できません。 |

---

### 補足事項
- 複製機能では、複製元の内容を初期表示し、必要に応じて編集して登録する運用とする。
- 「月×要員ID×案件ID」の一意性制約があるため、複製時にも重複確認を行う。

## 11. API基本仕様

### 11.1 共通仕様
- フロントエンドとバックエンドはHTTPベースのJSON APIで連携する。
- リクエストボディおよびレスポンスボディはJSON形式とする。
- 文字コードはUTF-8とする。
- 日付・月の表現は `yyyy-mm` 形式を使用する。
- 正常終了時は HTTP 200 系ステータスコードを返却する。
- 入力チェックエラー時は HTTP 400 を返却する。
- 対象データ未存在時は HTTP 404 を返却する。
- サーバ内部エラー時は HTTP 500 を返却する。

#### 共通レスポンス例
```json
{
  "message": "success",
  "data": {}
}
```

#### 共通エラーレスポンス例
```json
{
  "message": "validation error",
  "errors": [
    {
      "field": "project_name",
      "reason": "required"
    }
  ]
}
```

---

### 11.2 要員マスタAPI

#### 要員一覧取得
- **Method**: `GET`
- **Path**: `/api/resources`

#### 要員登録
- **Method**: `POST`
- **Path**: `/api/resources`

#### 要員更新
- **Method**: `PUT`
- **Path**: `/api/resources/{resource_id}`

#### 要員削除
- **Method**: `DELETE`
- **Path**: `/api/resources/{resource_id}`

#### 要員複製
- **Method**: `POST`
- **Path**: `/api/resources/{resource_id}/copy`

---

### 11.3 案件マスタAPI

#### 案件一覧取得
- **Method**: `GET`
- **Path**: `/api/projects`

#### 案件登録
- **Method**: `POST`
- **Path**: `/api/projects`

#### 案件更新
- **Method**: `PUT`
- **Path**: `/api/projects/{project_id}`

#### 案件削除
- **Method**: `DELETE`
- **Path**: `/api/projects/{project_id}`

#### 案件複製
- **Method**: `POST`
- **Path**: `/api/projects/{project_id}/copy`

---

### 11.4 山積みデータAPI

#### 山積み一覧取得
- **Method**: `GET`
- **Path**: `/api/allocations`

#### 山積み登録
- **Method**: `POST`
- **Path**: `/api/allocations`

#### 山積み更新
- **Method**: `PUT`
- **Path**: `/api/allocations/{allocation_id}`

#### 山積み削除
- **Method**: `DELETE`
- **Path**: `/api/allocations/{allocation_id}`

#### 山積み複製
- **Method**: `POST`
- **Path**: `/api/allocations/copy`

---

### 11.5 ステータスコード一覧

| HTTPステータス | 内容 |
|---|---|
| 200 | 正常終了（取得、更新、削除成功） |
| 201 | 登録成功 |
| 400 | 入力チェックエラー |
| 404 | 対象データなし |
| 409 | 重複データあり |
| 500 | サーバ内部エラー |

## 12. テーブル定義

### 12.1 要員マスタ（resources）

| 項目名 | カラム名 | 型 | PK | NOT NULL | UK | 説明 |
|---|---|---|---|---|---|---|
| 要員ID | resource_id | bigint | ○ | ○ | ○ | 要員を一意に識別するID |
| 要員名 | resource_name | varchar(100) |  | ○ |  | 要員名 |
| 所属 | department | varchar(100) |  |  |  | 所属部署名 |
| 備考 | note | varchar(500) |  |  |  | 補足情報 |
| 作成日時 | created_at | timestamp |  | ○ |  | 作成日時 |
| 更新日時 | updated_at | timestamp |  | ○ |  | 更新日時 |

### 12.2 案件マスタ（projects）

| 項目名 | カラム名 | 型 | PK | NOT NULL | UK | 説明 |
|---|---|---|---|---|---|---|
| 案件ID | project_id | bigint | ○ | ○ | ○ | 案件を一意に識別するID |
| 案件名 | project_name | varchar(100) |  | ○ |  | 案件名 |
| 開始月 | start_month | char(7) |  |  |  | 開始月（yyyy-mm） |
| 終了月 | end_month | char(7) |  |  |  | 終了月（yyyy-mm） |
| ステータス | status | varchar(50) |  |  |  | 案件ステータス |
| 備考 | note | varchar(500) |  |  |  | 補足情報 |
| 作成日時 | created_at | timestamp |  | ○ |  | 作成日時 |
| 更新日時 | updated_at | timestamp |  | ○ |  | 更新日時 |

### 12.3 山積みデータ（allocations）

| 項目名 | カラム名 | 型 | PK | NOT NULL | UK | 説明 |
|---|---|---|---|---|---|---|
| 山積みID | allocation_id | bigint | ○ | ○ | ○ | 山積みデータを一意に識別するID |
| 対象月 | target_month | char(7) |  | ○ |  | 対象月（yyyy-mm） |
| 要員ID | resource_id | bigint |  | ○ |  | 要員ID |
| 案件ID | project_id | bigint |  | ○ |  | 案件ID |
| 工数 | workload | decimal(5,2) |  | ○ |  | 工数 |
| 備考 | note | varchar(500) |  |  |  | 補足情報 |
| 作成日時 | created_at | timestamp |  | ○ |  | 作成日時 |
| 更新日時 | updated_at | timestamp |  | ○ |  | 更新日時 |

## 13. PostgreSQL DDL案

```sql
CREATE TABLE resources (
    resource_id BIGSERIAL PRIMARY KEY,
    resource_name VARCHAR(100) NOT NULL,
    department VARCHAR(100),
    note VARCHAR(500),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE projects (
    project_id BIGSERIAL PRIMARY KEY,
    project_name VARCHAR(100) NOT NULL,
    start_month CHAR(7),
    end_month CHAR(7),
    status VARCHAR(50),
    note VARCHAR(500),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT chk_projects_start_month_format
        CHECK (start_month IS NULL OR start_month ~ '^[0-9]{4}-[0-9]{2}$'),
    CONSTRAINT chk_projects_end_month_format
        CHECK (end_month IS NULL OR end_month ~ '^[0-9]{4}-[0-9]{2}$'),
    CONSTRAINT chk_projects_month_order
        CHECK (
            start_month IS NULL
            OR end_month IS NULL
            OR start_month <= end_month
        )
);

CREATE TABLE allocations (
    allocation_id BIGSERIAL PRIMARY KEY,
    target_month CHAR(7) NOT NULL,
    resource_id BIGINT NOT NULL,
    project_id BIGINT NOT NULL,
    workload NUMERIC(5,2) NOT NULL,
    note VARCHAR(500),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_allocations_resource
        FOREIGN KEY (resource_id) REFERENCES resources(resource_id),
    CONSTRAINT fk_allocations_project
        FOREIGN KEY (project_id) REFERENCES projects(project_id),
    CONSTRAINT uq_allocations_month_resource_project
        UNIQUE (target_month, resource_id, project_id),
    CONSTRAINT chk_allocations_target_month_format
        CHECK (target_month ~ '^[0-9]{4}-[0-9]{2}$'),
    CONSTRAINT chk_allocations_workload_non_negative
        CHECK (workload >= 0)
);

CREATE INDEX idx_allocations_target_month
    ON allocations (target_month);

CREATE INDEX idx_allocations_resource_id
    ON allocations (resource_id);

CREATE INDEX idx_allocations_project_id
    ON allocations (project_id);
```
