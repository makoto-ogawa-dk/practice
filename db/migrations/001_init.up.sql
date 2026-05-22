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
        CHECK (start_month IS NULL OR start_month ~ '^[0-9]{4}-(0[1-9]|1[0-2])$'),
    CONSTRAINT chk_projects_end_month_format
        CHECK (end_month IS NULL OR end_month ~ '^[0-9]{4}-(0[1-9]|1[0-2])$'),
    CONSTRAINT chk_projects_month_order
        CHECK (start_month IS NULL OR end_month IS NULL OR start_month <= end_month)
);

CREATE TABLE allocations (
    allocation_id BIGSERIAL PRIMARY KEY,
    target_month CHAR(7) NOT NULL,
    resource_id BIGINT NOT NULL REFERENCES resources(resource_id),
    project_id BIGINT NOT NULL REFERENCES projects(project_id),
    workload NUMERIC(5,2) NOT NULL,
    note VARCHAR(500),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT uq_allocations_month_resource_project UNIQUE (target_month, resource_id, project_id),
    CONSTRAINT chk_allocations_target_month_format
        CHECK (target_month ~ '^[0-9]{4}-(0[1-9]|1[0-2])$'),
    CONSTRAINT chk_allocations_workload_non_negative CHECK (workload >= 0)
);

CREATE INDEX idx_allocations_target_month ON allocations (target_month);
CREATE INDEX idx_allocations_resource_id ON allocations (resource_id);
CREATE INDEX idx_allocations_project_id ON allocations (project_id);
