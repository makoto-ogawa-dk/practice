package api

type Resource struct {
	ResourceID   int64   `json:"resource_id"`
	ResourceName string  `json:"resource_name"`
	Department   *string `json:"department,omitempty"`
	Note         *string `json:"note,omitempty"`
}

type Project struct {
	ProjectID   int64   `json:"project_id"`
	ProjectName string  `json:"project_name"`
	StartMonth  *string `json:"start_month,omitempty"`
	EndMonth    *string `json:"end_month,omitempty"`
	Status      *string `json:"status,omitempty"`
	Note        *string `json:"note,omitempty"`
}

type Allocation struct {
	AllocationID int64   `json:"allocation_id"`
	TargetMonth  string  `json:"target_month"`
	ResourceID   int64   `json:"resource_id"`
	ProjectID    int64   `json:"project_id"`
	Workload     float64 `json:"workload"`
	Note         *string `json:"note,omitempty"`
}

type APIResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Errors  []APIError  `json:"errors,omitempty"`
}

type APIError struct {
	Field  string `json:"field"`
	Reason string `json:"reason"`
}
