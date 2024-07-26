package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
)

type ItemLogger interface {
	LogItemBought(context.Context, Item) error
}

type SimpleItemLogger struct{}

func (l SimpleItemLogger) LogItemBought(ctx context.Context, i Item) error {
	slog.Info("LOGGING item bought", "id", i.ID, "name", i.Name)
	return nil
}

type BetterItemLogger struct{}

func (l BetterItemLogger) LogItemBought(ctx context.Context, i Item) error {
	slog.Info("BETTER LOGGING item bought", "id", i.ID, "name", i.Name)
	return nil
}

type Item struct {
	ID   string
	Name string
}

type ItemHandler struct {
	Logger ItemLogger
	Count  int
}

func (i *ItemHandler) handleBuyItem(w http.ResponseWriter, r *http.Request) {
	var item Item

	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		slog.Error("failed to decode request body", "err", err)
		return
	}

	if err := i.Logger.LogItemBought(r.Context(), item); err != nil {
		slog.Error("failed to log item buying", "err", err)
		return
	}

	i.Count++

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}

// func logItemBought(item Item) error {
// 	time.Sleep(time.Millisecond * 500)
// 	slog.Info("item bought", "id", item.ID, "name", item.Name)
// 	return nil
// }

func main() {
	itemHandler := &ItemHandler{
		Logger: BetterItemLogger{},
		Count:  0,
	}

	req, _ := http.NewRequest("POST", "/buy", strings.NewReader(`{"ID":"123", "Name":"Test Item"}`))
	w := httptest.NewRecorder()

	slog.Info("Initial count", "count", itemHandler.Count)

	itemHandler.handleBuyItem(w, req)

	slog.Info("Count after handleBuyItem", "count", itemHandler.Count)
}
