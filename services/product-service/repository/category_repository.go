package repository

import (
	"context"

	pb "github.com/daffaromero/retries/services/common/genproto/grpc-api"
	"github.com/daffaromero/retries/services/product-service/repository/query"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CategoryRepository interface {
	CreateCategory(context.Context, *pb.Category) (*pb.Category, error)
	GetCategoryById(context.Context, *pb.GetCategoryFilter) (*pb.GetCategoryResponse, error)
	GetCategories(context.Context, *pb.GetCategoryFilter, pb.ProductService_GetCategoriesServer) error
	UpdateCategory(context.Context, *pb.Category) (*pb.Category, error)
	DeleteCategory(context.Context, *pb.GetCategoryFilter) (*pb.DeleteCategoryResponse, error)
}

type CategoryRepositoryImpl struct {
	db       Store
	catQuery query.CategoryQuery
}

func NewCategoryRepository(db Store, catQuery query.CategoryQuery) CategoryRepository {
	return &CategoryRepositoryImpl{db: db, catQuery: catQuery}
}

func (c *CategoryRepositoryImpl) CreateCategory(ctx context.Context, req *pb.Category) (*pb.Category, error) {
	var res *pb.Category
	err := c.db.WithTx(ctx, func(tx pgx.Tx) error {
		cat, err := c.catQuery.CreateCategory(ctx, req)
		if err != nil {
			return err
		}
		res = cat
		return err
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *CategoryRepositoryImpl) GetCategoryById(ctx context.Context, req *pb.GetCategoryFilter) (*pb.GetCategoryResponse, error) {
	var res *pb.GetCategoryResponse
	err := c.db.WithTx(ctx, func(tx pgx.Tx) error {
		cat, err := c.catQuery.GetCategoryById(ctx, req)
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

func (c *CategoryRepositoryImpl) GetCategories(ctx context.Context, req *pb.GetCategoryFilter, sv pb.ProductService_GetCategoriesServer) error {
	if err := c.db.WithoutTx(ctx, func(pool *pgxpool.Pool) error {
		return c.catQuery.GetCategories(ctx, req, sv)
	}); err != nil {
		return err
	}
	return nil
}

func (c *CategoryRepositoryImpl) UpdateCategory(ctx context.Context, req *pb.Category) (*pb.Category, error) {
	var res *pb.Category
	err := c.db.WithoutTx(ctx, func(pool *pgxpool.Pool) error {
		cat, err := c.catQuery.UpdateCategory(ctx, req)
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

func (c *CategoryRepositoryImpl) DeleteCategory(ctx context.Context, req *pb.GetCategoryFilter) (*pb.DeleteCategoryResponse, error) {
	var res *pb.DeleteCategoryResponse
	err := c.db.WithoutTx(ctx, func(pool *pgxpool.Pool) error {
		cat, err := c.catQuery.DeleteCategory(ctx, req)
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
