package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"example/goserver/internal/model"
	"example/goserver/internal/repository"
)

type Handler struct {
	repo *repository.Repository
}

func NewHandler(repo *repository.Repository) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) Employees(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.ListEmployees(w, r)
	case http.MethodPost:
		h.CreateNewEmployee(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *Handler) Employee(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		ctx := r.Context()

		idStr := strings.TrimPrefix(r.URL.Path, "/employees/")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "invalid id", http.StatusNotFound)
			return
		}

		user, err := h.repo.GetEmployeeByID(ctx, id)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	case http.MethodDelete:
		h.DeleteEmployee(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// Helper functions

func (h *Handler) ListEmployees(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	users, err := h.repo.ListEmployees(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (h *Handler) CreateNewEmployee(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var u model.Employee
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	if err := h.repo.CreateEmployee(ctx, &u); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(u)
}

func (h *Handler) DeleteEmployee(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := strings.TrimPrefix(r.URL.Path, "/employees/")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusNotFound)
		return
	}

	err = h.repo.DeleteEmployeeByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
