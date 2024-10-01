package repository

import (
	"context"

	pb "github.com/daffaromero/retries/services/common/genproto/grpc-api"
	"github.com/daffaromero/retries/services/product-service/repository/query"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CategoryRepository interface {
	CreateCategory(c context.Context, req *pb.Category) (*pb.Category, error)
	GetCategoryByID(c context.Context, req *pb.GetCategoryFilter) (*pb.GetCategoryResponse, error)
	GetCategories(c context.Context, req *pb.GetCategoryFilter) (*pb.GetCategoryResponse, error)
	UpdateCategory(c context.Context, req *pb.Category) (*pb.Category, error)
	DeleteCategory(c context.Context, req *pb.GetCategoryFilter) (*pb.DeleteCategoryResponse, error)
}

type categoryRepository struct {
	db       Store
	catQuery query.CategoryQuery
}

func NewCategoryRepository(db Store, catQuery query.CategoryQuery) CategoryRepository {
	return &categoryRepository{db: db, catQuery: catQuery}
}

func (c *categoryRepository) CreateCategory(ctx context.Context, req *pb.Category) (*pb.Category, error) {
	var category *pb.Category
	err := c.db.WithTx(ctx, func(tx pgx.Tx) error {
		cat, err := c.catQuery.CreateCategory(ctx, tx, req)
		if err != nil {
			return err
		}
		category = cat
		return err
	})
	if err != nil {
		return nil, err
	}
	return category, nil
}

func (c *categoryRepository) GetCategoryByID(ctx context.Context, req *pb.GetCategoryFilter) (*pb.GetCategoryResponse, error) {
	var category *pb.GetCategoryResponse
	err := c.db.WithoutTx(ctx, func(pool *pgxpool.Pool) error {
		cat, err := c.catQuery.GetCategoryByID(ctx, req)
		if err != nil {
			return err
		}
		category = cat
		return nil
	})
	if err != nil {
		return nil, err
	}
	return category, nil
}

func (c *categoryRepository) GetCategories(ctx context.Context, req *pb.GetCategoryFilter) (*pb.GetCategoryResponse, error) {
	var categories *pb.GetCategoryResponse

	err := c.db.WithoutTx(ctx, func(pool *pgxpool.Pool) error {
		var err error
		categories, err = c.catQuery.GetCategories(ctx, req)
		return err
	})
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func (c *categoryRepository) UpdateCategory(ctx context.Context, req *pb.Category) (*pb.Category, error) {
	var category *pb.Category
	err := c.db.WithTx(ctx, func(tx pgx.Tx) error {
		cat, err := c.catQuery.UpdateCategory(ctx, tx, req)
		if err != nil {
			return err
		}
		category = cat
		return nil
	})
	if err != nil {
		return nil, err
	}
	return category, nil
}

func (c *categoryRepository) DeleteCategory(ctx context.Context, req *pb.GetCategoryFilter) (*pb.DeleteCategoryResponse, error) {
	var res *pb.DeleteCategoryResponse
	err := c.db.WithTx(ctx, func(tx pgx.Tx) error {
		cat, err := c.catQuery.DeleteCategory(ctx, tx, req)
		if err != nil {
			return err
		}
		res = cat
		return nil
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}
