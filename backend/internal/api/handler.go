package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type Handler struct {
	db *sql.DB
}

func NewHandler(db *sql.DB) *Handler {
	return &Handler{db: db}
}

func (h *Handler) Router() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", h.handleHealth)
	mux.HandleFunc("/api/resources", h.handleResources)
	mux.HandleFunc("/api/resources/", h.handleResourceByID)
	mux.HandleFunc("/api/projects", h.handleProjects)
	mux.HandleFunc("/api/projects/", h.handleProjectByID)
	mux.HandleFunc("/api/allocations", h.handleAllocations)
	mux.HandleFunc("/api/allocations/", h.handleAllocationByID)
	return withCORS(mux)
}

func (h *Handler) handleHealth(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, APIResponse{Message: "success", Data: map[string]string{"status": "ok"}})
}

func (h *Handler) handleResources(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.listResources(w)
	case http.MethodPost:
		h.createResource(w, r)
	default:
		writeMethodNotAllowed(w)
	}
}

func (h *Handler) handleResourceByID(w http.ResponseWriter, r *http.Request) {
	id, action, err := parseIDAndAction(r.URL.Path, "/api/resources/")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid path", "resource_id", "invalid")
		return
	}

	if action == "copy" {
		if r.Method != http.MethodPost {
			writeMethodNotAllowed(w)
			return
		}
		h.copyResource(w, id)
		return
	}

	switch r.Method {
	case http.MethodPut:
		h.updateResource(w, r, id)
	case http.MethodDelete:
		h.deleteResource(w, id)
	default:
		writeMethodNotAllowed(w)
	}
}

func (h *Handler) handleProjects(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.listProjects(w)
	case http.MethodPost:
		h.createProject(w, r)
	default:
		writeMethodNotAllowed(w)
	}
}

func (h *Handler) handleProjectByID(w http.ResponseWriter, r *http.Request) {
	id, action, err := parseIDAndAction(r.URL.Path, "/api/projects/")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid path", "project_id", "invalid")
		return
	}

	if action == "copy" {
		if r.Method != http.MethodPost {
			writeMethodNotAllowed(w)
			return
		}
		h.copyProject(w, id)
		return
	}

	switch r.Method {
	case http.MethodPut:
		h.updateProject(w, r, id)
	case http.MethodDelete:
		h.deleteProject(w, id)
	default:
		writeMethodNotAllowed(w)
	}
}

func (h *Handler) handleAllocations(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.listAllocations(w, r)
	case http.MethodPost:
		if strings.TrimSuffix(r.URL.Path, "/") == "/api/allocations/copy" {
			h.copyAllocations(w, r)
			return
		}
		h.createAllocation(w, r)
	default:
		writeMethodNotAllowed(w)
	}
}

func (h *Handler) handleAllocationByID(w http.ResponseWriter, r *http.Request) {
	if strings.TrimSuffix(r.URL.Path, "/") == "/api/allocations/copy" {
		if r.Method != http.MethodPost {
			writeMethodNotAllowed(w)
			return
		}
		h.copyAllocations(w, r)
		return
	}

	id, _, err := parseIDAndAction(r.URL.Path, "/api/allocations/")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid path", "allocation_id", "invalid")
		return
	}

	switch r.Method {
	case http.MethodPut:
		h.updateAllocation(w, r, id)
	case http.MethodDelete:
		h.deleteAllocation(w, id)
	default:
		writeMethodNotAllowed(w)
	}
}

func (h *Handler) listResources(w http.ResponseWriter) {
	rows, err := h.db.Query(`SELECT resource_id, resource_name, department, note FROM resources ORDER BY resource_id`)
	if err != nil {
		writeInternalError(w, err)
		return
	}
	defer rows.Close()

	items := []Resource{}
	for rows.Next() {
		var it Resource
		if err := rows.Scan(&it.ResourceID, &it.ResourceName, &it.Department, &it.Note); err != nil {
			writeInternalError(w, err)
			return
		}
		items = append(items, it)
	}
	writeJSON(w, http.StatusOK, APIResponse{Message: "success", Data: items})
}

func (h *Handler) createResource(w http.ResponseWriter, r *http.Request) {
	var in struct {
		ResourceName string  `json:"resource_name"`
		Department   *string `json:"department"`
		Note         *string `json:"note"`
	}
	if err := decodeJSON(r, &in); err != nil {
		writeError(w, http.StatusBadRequest, "validation error", "body", "invalid_json")
		return
	}
	in.ResourceName = strings.TrimSpace(in.ResourceName)
	if in.ResourceName == "" {
		writeError(w, http.StatusBadRequest, "validation error", "resource_name", "required")
		return
	}

	var id int64
	err := h.db.QueryRow(
		`INSERT INTO resources (resource_name, department, note) VALUES ($1, $2, $3) RETURNING resource_id`,
		in.ResourceName, trimPtr(in.Department), trimPtr(in.Note),
	).Scan(&id)
	if err != nil {
		writeInternalError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, APIResponse{Message: "success", Data: map[string]int64{"resource_id": id}})
}

func (h *Handler) updateResource(w http.ResponseWriter, r *http.Request, id int64) {
	var in struct {
		ResourceName string  `json:"resource_name"`
		Department   *string `json:"department"`
		Note         *string `json:"note"`
	}
	if err := decodeJSON(r, &in); err != nil {
		writeError(w, http.StatusBadRequest, "validation error", "body", "invalid_json")
		return
	}
	in.ResourceName = strings.TrimSpace(in.ResourceName)
	if in.ResourceName == "" {
		writeError(w, http.StatusBadRequest, "validation error", "resource_name", "required")
		return
	}

	res, err := h.db.Exec(
		`UPDATE resources SET resource_name=$1, department=$2, note=$3, updated_at=NOW() WHERE resource_id=$4`,
		in.ResourceName, trimPtr(in.Department), trimPtr(in.Note), id,
	)
	if err != nil {
		writeInternalError(w, err)
		return
	}
	writeRowsAffected(w, res, "resource")
}

func (h *Handler) deleteResource(w http.ResponseWriter, id int64) {
	res, err := h.db.Exec(`DELETE FROM resources WHERE resource_id=$1`, id)
	if err != nil {
		writeInternalError(w, err)
		return
	}
	writeRowsAffected(w, res, "resource")
}

func (h *Handler) copyResource(w http.ResponseWriter, id int64) {
	var newID int64
	err := h.db.QueryRow(
		`INSERT INTO resources (resource_name, department, note)
 SELECT resource_name, department, note FROM resources WHERE resource_id=$1
 RETURNING resource_id`,
		id,
	).Scan(&newID)
	if errors.Is(err, sql.ErrNoRows) {
		writeError(w, http.StatusNotFound, "not found", "resource_id", "not_found")
		return
	}
	if err != nil {
		writeInternalError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, APIResponse{Message: "success", Data: map[string]int64{"resource_id": newID}})
}

func (h *Handler) listProjects(w http.ResponseWriter) {
	rows, err := h.db.Query(`SELECT project_id, project_name, start_month, end_month, status, note FROM projects ORDER BY project_id`)
	if err != nil {
		writeInternalError(w, err)
		return
	}
	defer rows.Close()

	items := []Project{}
	for rows.Next() {
		var it Project
		if err := rows.Scan(&it.ProjectID, &it.ProjectName, &it.StartMonth, &it.EndMonth, &it.Status, &it.Note); err != nil {
			writeInternalError(w, err)
			return
		}
		items = append(items, it)
	}
	writeJSON(w, http.StatusOK, APIResponse{Message: "success", Data: items})
}

func (h *Handler) createProject(w http.ResponseWriter, r *http.Request) {
	var in struct {
		ProjectName string  `json:"project_name"`
		StartMonth  *string `json:"start_month"`
		EndMonth    *string `json:"end_month"`
		Status      *string `json:"status"`
		Note        *string `json:"note"`
	}
	if err := decodeJSON(r, &in); err != nil {
		writeError(w, http.StatusBadRequest, "validation error", "body", "invalid_json")
		return
	}
	in.ProjectName = strings.TrimSpace(in.ProjectName)
	if in.ProjectName == "" {
		writeError(w, http.StatusBadRequest, "validation error", "project_name", "required")
		return
	}
	if err := validateProjectMonths(in.StartMonth, in.EndMonth); err != nil {
		writeError(w, http.StatusBadRequest, "validation error", "start_month/end_month", err.Error())
		return
	}

	var id int64
	err := h.db.QueryRow(`INSERT INTO projects (project_name, start_month, end_month, status, note) VALUES ($1,$2,$3,$4,$5) RETURNING project_id`,
		in.ProjectName, trimPtr(in.StartMonth), trimPtr(in.EndMonth), trimPtr(in.Status), trimPtr(in.Note)).Scan(&id)
	if err != nil {
		writeInternalError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, APIResponse{Message: "success", Data: map[string]int64{"project_id": id}})
}

func (h *Handler) updateProject(w http.ResponseWriter, r *http.Request, id int64) {
	var in struct {
		ProjectName string  `json:"project_name"`
		StartMonth  *string `json:"start_month"`
		EndMonth    *string `json:"end_month"`
		Status      *string `json:"status"`
		Note        *string `json:"note"`
	}
	if err := decodeJSON(r, &in); err != nil {
		writeError(w, http.StatusBadRequest, "validation error", "body", "invalid_json")
		return
	}
	in.ProjectName = strings.TrimSpace(in.ProjectName)
	if in.ProjectName == "" {
		writeError(w, http.StatusBadRequest, "validation error", "project_name", "required")
		return
	}
	if err := validateProjectMonths(in.StartMonth, in.EndMonth); err != nil {
		writeError(w, http.StatusBadRequest, "validation error", "start_month/end_month", err.Error())
		return
	}

	res, err := h.db.Exec(`UPDATE projects SET project_name=$1, start_month=$2, end_month=$3, status=$4, note=$5, updated_at=NOW() WHERE project_id=$6`,
		in.ProjectName, trimPtr(in.StartMonth), trimPtr(in.EndMonth), trimPtr(in.Status), trimPtr(in.Note), id)
	if err != nil {
		writeInternalError(w, err)
		return
	}
	writeRowsAffected(w, res, "project")
}

func (h *Handler) deleteProject(w http.ResponseWriter, id int64) {
	res, err := h.db.Exec(`DELETE FROM projects WHERE project_id=$1`, id)
	if err != nil {
		writeInternalError(w, err)
		return
	}
	writeRowsAffected(w, res, "project")
}

func (h *Handler) copyProject(w http.ResponseWriter, id int64) {
	var newID int64
	err := h.db.QueryRow(`INSERT INTO projects (project_name, start_month, end_month, status, note)
SELECT project_name, start_month, end_month, status, note FROM projects WHERE project_id=$1
RETURNING project_id`, id).Scan(&newID)
	if errors.Is(err, sql.ErrNoRows) {
		writeError(w, http.StatusNotFound, "not found", "project_id", "not_found")
		return
	}
	if err != nil {
		writeInternalError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, APIResponse{Message: "success", Data: map[string]int64{"project_id": newID}})
}

func (h *Handler) listAllocations(w http.ResponseWriter, r *http.Request) {
	query := `SELECT allocation_id, target_month, resource_id, project_id, workload, note FROM allocations WHERE 1=1`
	args := []interface{}{}

	if month := strings.TrimSpace(r.URL.Query().Get("month")); month != "" {
		query += fmt.Sprintf(" AND target_month=$%d", len(args)+1)
		args = append(args, month)
	}
	if resourceID := strings.TrimSpace(r.URL.Query().Get("resource_id")); resourceID != "" {
		query += fmt.Sprintf(" AND resource_id=$%d", len(args)+1)
		args = append(args, resourceID)
	}
	if projectID := strings.TrimSpace(r.URL.Query().Get("project_id")); projectID != "" {
		query += fmt.Sprintf(" AND project_id=$%d", len(args)+1)
		args = append(args, projectID)
	}
	query += " ORDER BY target_month, resource_id, project_id"

	rows, err := h.db.Query(query, args...)
	if err != nil {
		writeInternalError(w, err)
		return
	}
	defer rows.Close()

	items := []Allocation{}
	for rows.Next() {
		var it Allocation
		if err := rows.Scan(&it.AllocationID, &it.TargetMonth, &it.ResourceID, &it.ProjectID, &it.Workload, &it.Note); err != nil {
			writeInternalError(w, err)
			return
		}
		items = append(items, it)
	}
	writeJSON(w, http.StatusOK, APIResponse{Message: "success", Data: items})
}

func (h *Handler) createAllocation(w http.ResponseWriter, r *http.Request) {
	var in struct {
		TargetMonth string  `json:"target_month"`
		ResourceID  int64   `json:"resource_id"`
		ProjectID   int64   `json:"project_id"`
		Workload    float64 `json:"workload"`
		Note        *string `json:"note"`
	}
	if err := decodeJSON(r, &in); err != nil {
		writeError(w, http.StatusBadRequest, "validation error", "body", "invalid_json")
		return
	}
	if err := validateAllocation(in.TargetMonth, in.ResourceID, in.ProjectID, in.Workload); err != nil {
		writeError(w, http.StatusBadRequest, "validation error", "allocation", err.Error())
		return
	}

	var id int64
	err := h.db.QueryRow(`INSERT INTO allocations (target_month, resource_id, project_id, workload, note) VALUES ($1,$2,$3,$4,$5) RETURNING allocation_id`,
		strings.TrimSpace(in.TargetMonth), in.ResourceID, in.ProjectID, in.Workload, trimPtr(in.Note)).Scan(&id)
	if err != nil {
		if isConflict(err) {
			writeError(w, http.StatusConflict, "conflict", "target_month/resource_id/project_id", "duplicate")
			return
		}
		writeInternalError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, APIResponse{Message: "success", Data: map[string]int64{"allocation_id": id}})
}

func (h *Handler) updateAllocation(w http.ResponseWriter, r *http.Request, id int64) {
	var in struct {
		TargetMonth string  `json:"target_month"`
		ResourceID  int64   `json:"resource_id"`
		ProjectID   int64   `json:"project_id"`
		Workload    float64 `json:"workload"`
		Note        *string `json:"note"`
	}
	if err := decodeJSON(r, &in); err != nil {
		writeError(w, http.StatusBadRequest, "validation error", "body", "invalid_json")
		return
	}
	if err := validateAllocation(in.TargetMonth, in.ResourceID, in.ProjectID, in.Workload); err != nil {
		writeError(w, http.StatusBadRequest, "validation error", "allocation", err.Error())
		return
	}

	res, err := h.db.Exec(`UPDATE allocations SET target_month=$1, resource_id=$2, project_id=$3, workload=$4, note=$5, updated_at=NOW() WHERE allocation_id=$6`,
		strings.TrimSpace(in.TargetMonth), in.ResourceID, in.ProjectID, in.Workload, trimPtr(in.Note), id)
	if err != nil {
		if isConflict(err) {
			writeError(w, http.StatusConflict, "conflict", "target_month/resource_id/project_id", "duplicate")
			return
		}
		writeInternalError(w, err)
		return
	}
	writeRowsAffected(w, res, "allocation")
}

func (h *Handler) deleteAllocation(w http.ResponseWriter, id int64) {
	res, err := h.db.Exec(`DELETE FROM allocations WHERE allocation_id=$1`, id)
	if err != nil {
		writeInternalError(w, err)
		return
	}
	writeRowsAffected(w, res, "allocation")
}

func (h *Handler) copyAllocations(w http.ResponseWriter, r *http.Request) {
	var in struct {
		SourceMonth string `json:"source_month"`
		TargetMonth string `json:"target_month"`
	}
	if err := decodeJSON(r, &in); err != nil {
		writeError(w, http.StatusBadRequest, "validation error", "body", "invalid_json")
		return
	}
	in.SourceMonth = strings.TrimSpace(in.SourceMonth)
	in.TargetMonth = strings.TrimSpace(in.TargetMonth)
	if !isValidMonth(in.SourceMonth) || !isValidMonth(in.TargetMonth) {
		writeError(w, http.StatusBadRequest, "validation error", "source_month/target_month", "invalid_month")
		return
	}

	res, err := h.db.Exec(`INSERT INTO allocations (target_month, resource_id, project_id, workload, note)
SELECT $2, resource_id, project_id, workload, note FROM allocations WHERE target_month=$1`, in.SourceMonth, in.TargetMonth)
	if err != nil {
		if isConflict(err) {
			writeError(w, http.StatusConflict, "conflict", "target_month/resource_id/project_id", "duplicate")
			return
		}
		writeInternalError(w, err)
		return
	}

	cnt, _ := res.RowsAffected()
	writeJSON(w, http.StatusCreated, APIResponse{Message: "success", Data: map[string]int64{"copied_count": cnt}})
}

func validateProjectMonths(start, end *string) error {
	s := derefTrim(start)
	e := derefTrim(end)
	if s != "" && !isValidMonth(s) {
		return errors.New("invalid_start_month")
	}
	if e != "" && !isValidMonth(e) {
		return errors.New("invalid_end_month")
	}
	if s != "" && e != "" && s > e {
		return errors.New("start_after_end")
	}
	return nil
}

func validateAllocation(month string, resourceID, projectID int64, workload float64) error {
	if !isValidMonth(strings.TrimSpace(month)) {
		return errors.New("invalid_target_month")
	}
	if resourceID <= 0 {
		return errors.New("invalid_resource_id")
	}
	if projectID <= 0 {
		return errors.New("invalid_project_id")
	}
	if workload < 0 {
		return errors.New("invalid_workload")
	}
	return nil
}

func parseIDAndAction(path, prefix string) (int64, string, error) {
	trimmed := strings.Trim(strings.TrimPrefix(path, prefix), "/")
	parts := strings.Split(trimmed, "/")
	if len(parts) == 0 || parts[0] == "" {
		return 0, "", errors.New("missing id")
	}
	id, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return 0, "", err
	}
	action := ""
	if len(parts) > 1 {
		action = parts[1]
	}
	return id, action, nil
}

func decodeJSON(r *http.Request, dst interface{}) error {
	defer r.Body.Close()
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	return dec.Decode(dst)
}

func writeRowsAffected(w http.ResponseWriter, res sql.Result, field string) {
	count, err := res.RowsAffected()
	if err != nil {
		writeInternalError(w, err)
		return
	}
	if count == 0 {
		writeError(w, http.StatusNotFound, "not found", field+"_id", "not_found")
		return
	}
	writeJSON(w, http.StatusOK, APIResponse{Message: "success", Data: map[string]int64{"affected": count}})
}

func writeMethodNotAllowed(w http.ResponseWriter) {
	writeError(w, http.StatusMethodNotAllowed, "method not allowed", "method", "not_allowed")
}

func writeInternalError(w http.ResponseWriter, err error) {
	writeError(w, http.StatusInternalServerError, "internal error", "server", err.Error())
}

func writeError(w http.ResponseWriter, status int, message, field, reason string) {
	writeJSON(w, status, APIResponse{Message: message, Errors: []APIError{{Field: field, Reason: reason}}})
}

func writeJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func trimPtr(v *string) interface{} {
	if v == nil {
		return nil
	}
	t := strings.TrimSpace(*v)
	if t == "" {
		return nil
	}
	return t
}

func derefTrim(v *string) string {
	if v == nil {
		return ""
	}
	return strings.TrimSpace(*v)
}

func isConflict(err error) bool {
	return strings.Contains(err.Error(), "duplicate key value violates unique constraint")
}
