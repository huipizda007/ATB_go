package server

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	db "github.com/huipizda007/ATB_go/db/sqlc"
)

type mockStore struct {
	db.Store
	CreateFruitFunc      func(ctx context.Context, arg db.CreateFruitParams) (db.Fruit, error)
	GetFruitFunc         func(ctx context.Context, id int64) (db.Fruit, error)
	ListFruitsFunc       func(ctx context.Context) ([]db.Fruit, error)
	DeleteFruitFunc      func(ctx context.Context, id int64) error
	UpdateFruitPriceFunc func(ctx context.Context, arg db.UpdateFruitPriceParams) (db.Fruit, error)
	UpdateFruitStockFunc func(ctx context.Context, arg db.UpdateFruitStockParams) (db.Fruit, error)
}

func (m *mockStore) CreateFruit(ctx context.Context, arg db.CreateFruitParams) (db.Fruit, error) {
	return m.CreateFruitFunc(ctx, arg)
}
func (m *mockStore) GetFruit(ctx context.Context, id int64) (db.Fruit, error) {
	return m.GetFruitFunc(ctx, id)
}
func (m *mockStore) ListFruits(ctx context.Context) ([]db.Fruit, error) {
	return m.ListFruitsFunc(ctx)
}
func (m *mockStore) DeleteFruit(ctx context.Context, id int64) error {
	return m.DeleteFruitFunc(ctx, id)
}
func (m *mockStore) UpdateFruitPrice(ctx context.Context, arg db.UpdateFruitPriceParams) (db.Fruit, error) {
	return m.UpdateFruitPriceFunc(ctx, arg)
}
func (m *mockStore) UpdateFruitStock(ctx context.Context, arg db.UpdateFruitStockParams) (db.Fruit, error) {
	return m.UpdateFruitStockFunc(ctx, arg)
}

// ================= ЮНІТ-ТЕСТИ =================

func TestHandleCreateFruit(t *testing.T) {
	store := &mockStore{
		CreateFruitFunc: func(ctx context.Context, arg db.CreateFruitParams) (db.Fruit, error) {
			// Імітуємо повернення успішно створеного фрукта
			return db.Fruit{
				ID:   1,
				Name: arg.Name,
			}, nil
		},
	}

	server := NewServer(store)

	body := map[string]interface{}{
		"name":         "Яблуко",
		"brand":        "Голден",
		"price_per_kg": 35.50,
		"stock_kg":     100,
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/fruits", bytes.NewReader(jsonBody))
	w := httptest.NewRecorder()

	server.router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Очікувався статус %d, отримано %d. Причина: %s", http.StatusCreated, w.Code, w.Body.String())
	}
}

func TestHandleGetFruit(t *testing.T) {
	store := &mockStore{
		GetFruitFunc: func(ctx context.Context, id int64) (db.Fruit, error) {
			if id == 1 {
				return db.Fruit{ID: 1, Name: "Банан"}, nil
			}
			return db.Fruit{}, nil
		},
	}

	server := NewServer(store)

	req := httptest.NewRequest(http.MethodGet, "/fruits/1", nil)
	w := httptest.NewRecorder()

	server.router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Очікувався статус %d, отримано %d", http.StatusOK, w.Code)
	}
}

func TestHandleListFruits(t *testing.T) {
	store := &mockStore{
		ListFruitsFunc: func(ctx context.Context) ([]db.Fruit, error) {
			return []db.Fruit{
				{ID: 1, Name: "Яблуко"},
				{ID: 2, Name: "Банан"},
			}, nil
		},
	}

	server := NewServer(store)

	req := httptest.NewRequest(http.MethodGet, "/fruits", nil)
	w := httptest.NewRecorder()

	server.router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Очікувався статус %d, отримано %d", http.StatusOK, w.Code)
	}
}

func TestHandleUpdateFruitPrice(t *testing.T) {
	store := &mockStore{
		UpdateFruitPriceFunc: func(ctx context.Context, arg db.UpdateFruitPriceParams) (db.Fruit, error) {
			return db.Fruit{
				ID:         arg.ID,
				PricePerKg: arg.PricePerKg,
			}, nil
		},
	}

	server := NewServer(store)

	body := updatePriceRequest{Price: "55.00"}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPatch, "/fruits/1/price", bytes.NewReader(jsonBody))
	w := httptest.NewRecorder()

	server.router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Очікувався статус %d, отримано %d", http.StatusOK, w.Code)
	}
}

func TestHandleUpdateFruitStock(t *testing.T) {
	store := &mockStore{
		UpdateFruitStockFunc: func(ctx context.Context, arg db.UpdateFruitStockParams) (db.Fruit, error) {
			return db.Fruit{
				ID:      arg.ID,
				StockKg: arg.StockKg,
			}, nil
		},
	}

	server := NewServer(store)

	body := updateStockRequest{Stock: 150}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPatch, "/fruits/1/stock", bytes.NewReader(jsonBody))
	w := httptest.NewRecorder()

	server.router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Очікувався статус %d, отримано %d", http.StatusOK, w.Code)
	}
}

func TestHandleDeleteFruit(t *testing.T) {
	store := &mockStore{
		DeleteFruitFunc: func(ctx context.Context, id int64) error {
			return nil // Успішне видалення
		},
	}

	server := NewServer(store)

	req := httptest.NewRequest(http.MethodDelete, "/fruits/1", nil)
	w := httptest.NewRecorder()

	server.router.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("Очікувався статус %d, отримано %d", http.StatusNoContent, w.Code)
	}
}