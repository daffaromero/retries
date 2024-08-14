package query

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "github.com/daffaromero/retries/services/common/genproto/grpc-api"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CategoryQuery interface {
	CreateCategory(context.Context, pgx.Tx, *pb.Category) (*pb.Category, error)
	GetCategoryById(context.Context, *pb.GetCategoryFilter) (*pb.GetCategoryResponse, error)
	GetCategories(context.Context, *pb.GetCategoryFilter, pb.ProductService_GetCategoriesServer) error
	UpdateCategory(context.Context, pgx.Tx, *pb.Category) (*pb.Category, error)
	DeleteCategory(context.Context, pgx.Tx, *pb.GetCategoryFilter) (*pb.DeleteCategoryResponse, error)
}

type CategoryQueryImpl struct {
	db *pgxpool.Pool
}

func NewCategoryQueryImpl(db *pgxpool.Pool) *CategoryQueryImpl {
	return &CategoryQueryImpl{
		db: db,
	}
}

func (c *CategoryQueryImpl) CreateCategory(ctx context.Context, tx pgx.Tx, req *pb.Category) (*pb.Category, error) {
	query := `INSERT INTO categories (id, name, description, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)`
	err := tx.QueryRow(ctx, query, req.Id, req.Name, req.Description, req.CreatedAt, req.UpdatedAt).Scan(&req.Id)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &pb.Category{Id: req.Id, Name: req.Name, Description: req.Description}, nil
}

func (c *CategoryQueryImpl) GetCategoryById(ctx context.Context, req *pb.GetCategoryFilter) (*pb.GetCategoryResponse, error) {
	query := `SELECT id, name, description FROM categories WHERE id=$1`
	rows, err := c.db.Query(ctx, query, req.Id)
	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("no records found")
	} else if err != nil {
		return nil, err
	}
	var categories []*pb.Category
	for rows.Next() {
		var category pb.Category
		err = rows.Scan(&category.Id, &category.Name, &category.Description)
		if err != nil {
			return nil, err
		}
		newCategory := &pb.Category{
			Id:          category.Id,
			Name:        category.Name,
			Description: category.Description,
		}
		categories = append(categories, newCategory)
	}
	return &pb.GetCategoryResponse{Categories: categories}, nil
}

func (c *CategoryQueryImpl) GetCategories(ctx context.Context, req *pb.GetCategoryFilter, stream pb.ProductService_GetCategoriesServer) error {
	query := `SELECT id, name, description FROM categories LIMIT $1 OFFSET $2`
	rows, err := c.db.Query(ctx, query, req.Pagination.Count, req.Pagination.Limit)
	if err != nil {
		return fmt.Errorf("get all query error: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		category := &pb.Category{}

		if err := rows.Scan(&category.Id, &category.Name, &category.Description); err != nil {
			return fmt.Errorf("scan error: %v", err)
		}
		response := &pb.GetCategoryResponse{
			Categories: []*pb.Category{category},
		}
		err := stream.Send(response)
		if err != nil {
			return err
		}
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("rows error: %v", err)
	}
	return nil
}

func (c *CategoryQueryImpl) UpdateCategory(ctx context.Context, tx pgx.Tx, req *pb.Category) (*pb.Category, error) {
	query := `UPDATE categories SET name $1 description $2, updated_at $3 WHERE id = $4`
	err := tx.QueryRow(ctx, query, req.Id, req.Name, req.Description, req.UpdatedAt).Scan(&req.Id)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &pb.Category{Id: req.Id, Name: req.Name, Description: req.Description}, nil
}

func (c *CategoryQueryImpl) DeleteCategory(ctx context.Context, tx pgx.Tx, req *pb.GetCategoryFilter) (*pb.DeleteCategoryResponse, error) {
	query := `UPDATE categories SET deleted_at = $1 WHERE id = $2`
	rows, err := tx.Query(ctx, query, time.Now(), req.Id)
	if err != nil {
		return nil, fmt.Errorf("delete query error: %v", err)
	}
	defer rows.Close()

	return &pb.DeleteCategoryResponse{Status: true}, nil
}
