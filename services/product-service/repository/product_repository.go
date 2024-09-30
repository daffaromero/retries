package repository

import (
	"context"

	pb "github.com/daffaromero/retries/services/common/genproto/grpc-api"
	"github.com/daffaromero/retries/services/product-service/repository/query"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductRepository interface {
	CreateProduct(c context.Context, product *pb.Product) (*pb.Product, error)
	GetProductById(c context.Context, filter *pb.GetProductFilter) (*pb.GetProductResponse, error)
	GetAllProducts(c context.Context, filter *pb.GetProductFilter) (*pb.GetProductResponse, error)
	UpdateProduct(c context.Context, product *pb.Product) (*pb.Product, error)
	ApproveProduct(c context.Context, approve *pb.ApproveProductRequest) (*pb.ApproveProductResponse, error)
}

type productRepository struct {
	db           Store
	productQuery query.ProductQuery
}

func NewProductRepository(db Store, productQuery query.ProductQuery) ProductRepository {
	return &productRepository{db: db, productQuery: productQuery}
}

func (p *productRepository) CreateProduct(c context.Context, product *pb.Product) (*pb.Product, error) {
	var res *pb.Product
	err := p.db.WithTx(c, func(tx pgx.Tx) error {
		prod, err := p.productQuery.CreateProduct(c, tx, product)
		if err != nil {
			return err
		}
		res = prod
		return err
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (p *productRepository) GetProductById(c context.Context, filter *pb.GetProductFilter) (*pb.GetProductResponse, error) {
	var res *pb.GetProductResponse
	err := p.db.WithoutTx(c, func(pool *pgxpool.Pool) error {
		prod, err := p.productQuery.GetProductById(c, filter)
		if err != nil {
			return err
		}
		res = prod
		return nil
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (p *productRepository) GetAllProducts(c context.Context, filter *pb.GetProductFilter) (*pb.GetProductResponse, error) {
	var res *pb.GetProductResponse
	err := p.db.WithoutTx(c, func(pool *pgxpool.Pool) error {
		prod, err := p.productQuery.GetProducts(c, filter)
		if err != nil {
			return err
		}
		res = prod
		return nil
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (p *productRepository) UpdateProduct(c context.Context, product *pb.Product) (*pb.Product, error) {
	var res *pb.Product
	err := p.db.WithTx(c, func(tx pgx.Tx) error {
		prod, err := p.productQuery.UpdateProduct(c, tx, product)
		if err != nil {
			return err
		}
		res = prod
		return nil
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (p *productRepository) ApproveProduct(c context.Context, approve *pb.ApproveProductRequest) (*pb.ApproveProductResponse, error) {
	var res *pb.ApproveProductResponse
	err := p.db.WithTx(c, func(tx pgx.Tx) error {
		prod, err := p.productQuery.ApproveProduct(c, tx, approve)
		if err != nil {
			return err
		}
		res = prod
		return nil
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}
