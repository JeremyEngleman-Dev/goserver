package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
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
		h.ListEmployees(w)
	case http.MethodPost:
		h.CreateNewEmployee(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *Handler) Employee(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		idStr := strings.TrimPrefix(r.URL.Path, "/employees/")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "invalid id", http.StatusNotFound)
			return
		}

		user, err := h.repo.GetEmployeeByID(id)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("COntent-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *Handler) ListEmployees(w http.ResponseWriter) {
	users, err := h.repo.ListEmployees()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (h *Handler) CreateNewEmployee(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	log.Println(string(body))
	r.Body = io.NopCloser(bytes.NewBuffer(body))

	var u model.Employee
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	if err := h.repo.CreateEmployee(&u); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(u)
}
