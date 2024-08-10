package query

import (
	"context"
	"fmt"
	"log"

	pb "github.com/daffaromero/retries/services/common/genproto/grpc-api"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CategoryQuery interface {
	CreateCategory(context.Context, *pb.CreateCategoryRequest) (*pb.Category, error)
	GetCategoryById(context.Context, *pb.GetCategoryFilter, string) (*pb.GetCategoryResponse, error)
	GetCategories(context.Context, *pb.GetCategoryFilter, *pb.ProductService_GetCategoriesServer) error
}

type CategoryQueryImpl struct {
	Db *pgxpool.Pool
}

func NewCategoryQueryImpl(db *pgxpool.Pool) *CategoryQueryImpl {
	return &CategoryQueryImpl{
		Db: db,
	}
}

func (c *CategoryQueryImpl) CreateCategory(ctx context.Context, ca *pb.Category) (*pb.Category, error) {
	query := `INSERT INTO categories (id, name, description, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)`
	err := c.Db.QueryRow(ctx, query, ca.Id, ca.Name, ca.Description, ca.CreatedAt, ca.UpdatedAt).Scan(&ca.Id)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &pb.Category{Id: ca.Id, Name: ca.Name, Description: ca.Description}, nil
}

func (c *CategoryQueryImpl) GetCategoryById(ctx context.Context, ca *pb.GetCategoryFilter, catId string) (*pb.GetCategoryResponse, error) {
	query := `SELECT id, name, description FROM categories WHERE id=$1`
	rows, err := c.Db.Query(ctx, query, catId)
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
	rows, err := c.Db.Query(ctx, query, req.Pagination.Count, req.Pagination.Limit)
	if err != nil {
		return fmt.Errorf("query error: %v", err)
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
