package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"rrieye/db"
)

type Server struct {
	router *chi.Mux
	store  db.Store
}

func NewServer(store db.Store) *Server {
	s := &Server{
		router: chi.NewRouter(),
		store:  store,
	}
	s.routes()
	return s
}

func (s *Server) routes() {
	s.router.Get("/health", s.handleHealth)
	s.router.Post("/fruits", s.handleCreateFruit)
	s.router.Get("/fruits", s.handleListFruits)
	s.router.Get("/fruits/{id}", s.handleGetFruit)
	s.router.Delete("/fruits/{id}", s.handleDeleteFruit)
	s.router.Patch("/fruits/{id}/price", s.handleUpdateFruitPrice)
	s.router.Patch("/fruits/{id}/stock", s.handleUpdateFruitStock)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	s.writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) handleCreateFruit(w http.ResponseWriter, r *http.Request) {
	var arg db.CreateFruitParams
	if err := json.NewDecoder(r.Body).Decode(&arg); err != nil {
		s.errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	fruit, err := s.store.CreateFruit(r.Context(), arg)
	if err != nil {
		s.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	s.writeJSON(w, http.StatusCreated, fruit)
}

func (s *Server) handleGetFruit(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		s.errorResponse(w, http.StatusBadRequest, "invalid id")
		return
	}

	fruit, err := s.store.GetFruit(r.Context(), id)
	if err != nil {
		s.errorResponse(w, http.StatusNotFound, "fruit not found")
		return
	}

	s.writeJSON(w, http.StatusOK, fruit)
}

func (s *Server) handleListFruits(w http.ResponseWriter, r *http.Request) {
	fruits, err := s.store.ListFruits(r.Context())
	if err != nil {
		s.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	s.writeJSON(w, http.StatusOK, fruits)
}

func (s *Server) handleDeleteFruit(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		s.errorResponse(w, http.StatusBadRequest, "invalid id")
		return
	}

	err = s.store.DeleteFruit(r.Context(), id)
	if err != nil {
		s.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

type updatePriceRequest struct {
	Price string `json:"price"`
}

func (s *Server) handleUpdateFruitPrice(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		s.errorResponse(w, http.StatusBadRequest, "invalid id")
		return
	}

	var req updatePriceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	var numericPrice pgtype.Numeric
	numericPrice.Scan(req.Price)

	arg := db.UpdateFruitPriceParams{
		ID:         id,
		PricePerKg: numericPrice,
	}

	fruit, err := s.store.UpdateFruitPrice(r.Context(), arg)
	if err != nil {
		s.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	s.writeJSON(w, http.StatusOK, fruit)
}

type updateStockRequest struct {
	Stock int32 `json:"stock"`
}

func (s *Server) handleUpdateFruitStock(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		s.errorResponse(w, http.StatusBadRequest, "invalid id")
		return
	}

	var req updateStockRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	arg := db.UpdateFruitStockParams{
		ID:      id,
		StockKg: req.Stock,
	}

	fruit, err := s.store.UpdateFruitStock(r.Context(), arg)
	if err != nil {
		s.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	s.writeJSON(w, http.StatusOK, fruit)
}

func (s *Server) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (s *Server) errorResponse(w http.ResponseWriter, status int, message string) {
	s.writeJSON(w, status, map[string]string{"error": message})
}

func (s *Server) Run(port string) {
	http.ListenAndServe(port, s.router)
}