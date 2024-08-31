package query

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	pb "github.com/daffaromero/retries/services/common/genproto/grpc-api"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductQuery interface {
	CreateProduct(context.Context, pgx.Tx, *pb.Product) (*pb.Product, error)
	GetProductById(context.Context, *pb.GetProductFilter) (*pb.GetProductResponse, error)
	GetProducts(context.Context, *pb.GetProductFilter, *pb.ProductService_GetProductsServer) error
	UpdateProduct(context.Context, pgx.Tx, *pb.Product) (*pb.Product, error)
	ApproveProduct(context.Context, pgx.Tx, *pb.ApproveProductRequest) (*pb.ApproveProductResponse, error)
}

type ProductQueryImpl struct {
	db *pgxpool.Pool
}

func NewProductQueryImpl(db *pgxpool.Pool) *ProductQueryImpl {
	return &ProductQueryImpl{
		db: db,
	}
}

func (p *ProductQueryImpl) CreateProduct(c context.Context, tx pgx.Tx, req *pb.Product) (*pb.Product, error) {
	query := `INSERT INTO products (id, seller_id, category_id, category_name, variant_ids, name, seller_name, description, vis_time, invis_time, insider_key, voucher, voucher_discount, total_duration, variant_settings, is_reviewable, is_admin_verified, visibility, exclusion, price, pict_url, cert_url, flat_price, percentage_price, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27)`
	var variantSettings pb.VariantSettings
	varSettings, err := json.Marshal(req.VariantSettings)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	err = tx.QueryRow(c, query, req.Id, req.SellerId, req.CategoryId, req.CategoryName, req.VariantIds, req.Name, req.SellerName, req.Description, req.VisTime, req.InvisTime, req.InsiderKey, req.Voucher, req.VoucherDiscount, req.TotalDuration, string(varSettings), req.IsReviewable, req.IsAdminVerified, req.Visibility, req.Exclusion, req.Price, req.PictUrl, req.CertUrl, req.FlatPrice, req.PercentagePrice, req.CreatedAt, req.UpdatedAt).Scan(&req.Id)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	err = json.Unmarshal(varSettings, &variantSettings)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &pb.Product{Id: req.Id, SellerId: req.SellerId, CategoryId: req.CategoryId, CategoryName: req.CategoryName, VariantIds: req.VariantIds, Name: req.Name, SellerName: req.SellerName, Description: req.Description, VisTime: req.VisTime, InsiderKey: req.InsiderKey, Voucher: req.Voucher, VoucherDiscount: req.VoucherDiscount, TotalDuration: req.TotalDuration, VariantSettings: []*pb.VariantSettings{&variantSettings}, IsReviewable: req.IsReviewable, IsAdminVerified: req.IsAdminVerified, Visibility: req.Visibility, Exclusion: req.Exclusion, Price: req.Price, PictUrl: req.PictUrl, CertUrl: req.CertUrl, FlatPrice: req.FlatPrice, PercentagePrice: req.PercentagePrice, CreatedAt: req.CreatedAt, UpdatedAt: req.UpdatedAt}, nil
}

func (p *ProductQueryImpl) GetProductById(c context.Context, req *pb.GetProductFilter) (*pb.GetProductResponse, error) {
	query := `SELECT id, seller_id, category_id, category_name, variant_ids, name, seller_name, description, vis_time, invis_time, insider_key, voucher, voucher_discount, total_duration, variant_settings, is_reviewable, is_admin_verified, visibility, exclusion, price, pict_url, cert_url, flat_price, percentage_price, created_at, updated_at FROM products WHERE id = $1 AND deleted_at IS NULL`
	var productWithId []*pb.Product
	var variantSettings pb.VariantSettings
	rows, err := p.db.Query(c, query, req.Id)
	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("no records found")
	} else if err != nil {
		return nil, err
	}
	for rows.Next() {
		var product pb.Product
		err = rows.Scan(&product.Id, &product.SellerId, &product.CategoryId, &product.CategoryName, &product.VariantIds, &product.Name, &product.SellerName, &product.Description, &product.VisTime, &product.InvisTime, &product.InsiderKey, &product.Voucher, &product.VoucherDiscount, &product.TotalDuration, []*pb.VariantSettings{&variantSettings}, &product.IsReviewable, &product.IsAdminVerified, &product.Visibility, &product.Exclusion, &product.Price, &product.PictUrl, &product.CertUrl, &product.FlatPrice, &product.PercentagePrice, &product.CreatedAt, &product.UpdatedAt)
		if err != nil {
			return nil, err
		}
		productStruct := &pb.Product{
			Id:              product.Id,
			SellerId:        product.SellerId,
			CategoryId:      product.CategoryId,
			CategoryName:    product.CategoryName,
			VariantIds:      product.VariantIds,
			Name:            product.Name,
			SellerName:      product.SellerName,
			Description:     product.Description,
			VisTime:         product.VisTime,
			InvisTime:       product.InvisTime,
			InsiderKey:      product.InsiderKey,
			Voucher:         product.Voucher,
			VoucherDiscount: product.VoucherDiscount,
			TotalDuration:   product.TotalDuration,
			VariantSettings: product.VariantSettings,
			IsReviewable:    product.IsReviewable,
			IsAdminVerified: product.IsAdminVerified,
			Visibility:      product.Visibility,
			Exclusion:       product.Exclusion,
			Price:           product.Price,
			PictUrl:         product.PictUrl,
			CertUrl:         product.CertUrl,
			FlatPrice:       product.FlatPrice,
			PercentagePrice: product.PercentagePrice,
			CreatedAt:       product.CreatedAt,
			UpdatedAt:       product.UpdatedAt,
		}
		productWithId = append(productWithId, productStruct)
	}
	return &pb.GetProductResponse{Products: productWithId}, nil
}

func (p *ProductQueryImpl) GetProducts(c context.Context, req *pb.GetProductFilter, stream pb.ProductService_GetProductsServer) error {
	var variantSettings pb.VariantSettings
	filter, page, sort, earliest, latest, err := ProductFilters("user", req)
	if err != nil {
<<<<<<< HEAD
		return fmt.Errorf("error building product filters, %v", err)
=======
		return fmt.Errorf("error building product filters", err)
>>>>>>> 0e07b87fbf4446536d32081f8c55fe3390d3a46d
	}
	query := fmt.Sprintf(`SELECT id, seller_id, category_id, category_name, variant_ids, name, seller_name, description, vis_time, invis_time, insider_key, voucher, voucher_discount, total_duration, variant_settings, is_reviewable, is_admin_verified, visibility, exclusion, price, pict_url, cert_url, flat_price, percentage_price, created_at, updated_at FROM products WHERE deleted_at IS NULL %s %s %s`, filter, sort, page)
	rows, err := p.db.Query(c, query, earliest, latest)
	if err == pgx.ErrNoRows {
		return fmt.Errorf("no records found")
	} else if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		product := &pb.Product{}

		if err := rows.Scan(&product.Id, &product.SellerId, &product.CategoryId, &product.CategoryName, &product.VariantIds, &product.Name, &product.SellerName, &product.Description, &product.VisTime, &product.InvisTime, &product.InsiderKey, &product.Voucher, &product.VoucherDiscount, &product.TotalDuration, []*pb.VariantSettings{&variantSettings}, &product.IsReviewable, &product.IsAdminVerified, &product.Visibility, &product.Exclusion, &product.Price, &product.PictUrl, &product.CertUrl, &product.FlatPrice, &product.PercentagePrice, &product.CreatedAt, &product.UpdatedAt); err != nil {
			return fmt.Errorf("scan all products error: %v", err)
		}
		response := &pb.GetProductResponse{
			Products: []*pb.Product{product},
		}
		if err := stream.Send(response); err != nil {
			return err
		}
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("rows error: %v", err)
	}
	return nil
}
