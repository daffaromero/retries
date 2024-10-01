package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/daffaromero/retries/services/common/utils/logger"
	"github.com/daffaromero/retries/services/product-service/config"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store interface {
	WithTx(ctx context.Context, fn func(pgx.Tx) error) error
	WithoutTx(ctx context.Context, fn func(*pgxpool.Pool) error) error
}

type StoreImpl struct {
	Db     *pgxpool.Pool
	logger *logger.Log
}

func NewStore(db *pgxpool.Pool) Store {
	return &StoreImpl{
		Db:     db,
		logger: logger.NewLog("database_store"),
	}
}

func (s *StoreImpl) WithTx(ctx context.Context, fn func(pgx.Tx) error) error {
	c, cancel := context.WithTimeout(context.Background(), time.Duration(config.TimeOutDuration)*time.Second)
	defer cancel()

	tx, err := s.Db.Begin(c)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
				s.logger.Error(fmt.Sprintf("Rollback error: %v, original error: %v", rollbackErr, err))

				err = fmt.Errorf("rollback error: %v (original error: %w)", rollbackErr, err)
			}
		}
	}()

	if err = fn(tx); err != nil {
		return fmt.Errorf("transaction function failed: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (s *StoreImpl) WithoutTx(ctx context.Context, fn func(*pgxpool.Pool) error) error {
	if err := fn(s.Db); err != nil {
		return err
	}
	return nil
}
