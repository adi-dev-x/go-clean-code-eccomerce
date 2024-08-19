package admin

import (
	"context"
	"database/sql"
	"fmt"

	"myproject/pkg/model"
)

type Repository interface {
	Register(ctx context.Context, request model.AdminRegisterRequest) error
	Listing(ctx context.Context) ([]model.Coupon, error)
	Login(ctx context.Context, email string) (model.AdminRegisterRequest, error)
	Addcoupon(ctx context.Context, request model.Coupon) error
	LatestListing(ctx context.Context) ([]model.Coupon, error)
	ActiveListing(ctx context.Context) ([]model.Coupon, error)

	GetCoupnExist(ctx context.Context, id string) (string, error)
	Deletecoupon(ctx context.Context, id string) error

	//Product listing
	ListingSingle(ctx context.Context, id string) ([]model.ProductListDetailed, error)
	ProductListing(ctx context.Context) ([]model.ProductListingUsers, error)
	CategoryListing(ctx context.Context, category string) ([]model.ProductListingUsers, error)
	PhighListing(ctx context.Context) ([]model.ProductListingUsers, error)
	PlowListing(ctx context.Context) ([]model.ProductListingUsers, error)
	InAZListing(ctx context.Context) ([]model.ProductListingUsers, error)
	InZAListing(ctx context.Context) ([]model.ProductListingUsers, error)
}

type repository struct {
	sql *sql.DB
}

func NewRepository(sqlDB *sql.DB) Repository {
	return &repository{
		sql: sqlDB,
	}
}
func (r *repository) CategoryListing(ctx context.Context, category string) ([]model.ProductListingUsers, error) {
	query := `
		SELECT 
			product_models.name,
			product_models.category,
			product_models.units,
			product_models.tax,
			product_models.amount,
			product_models.status,
			product_models.discount,
			
			 product_models.id AS pid 
		FROM 
			product_models 
		INNER JOIN 
			vendor ON product_models.vendor_id = vendor.id WHERE product_models.units > 0 AND product_models.category ILIKE '%' || $1 || '%';`

	rows, err := r.sql.QueryContext(ctx, query, category)
	if err != nil {
		return nil, fmt.Errorf("failed to execute select query: %w", err)
	}
	defer rows.Close()

	var products []model.ProductListingUsers
	for rows.Next() {
		var product model.ProductListingUsers
		err := rows.Scan(&product.Name, &product.Category, &product.Unit, &product.Tax, &product.Price, &product.Status, &product.Discount, &product.Pid)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		product.Pdetail = "http://localhost:8080/user/listingSingleProduct/" + product.Pid
		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return products, nil
}
func (r *repository) ProductListing(ctx context.Context) ([]model.ProductListingUsers, error) {
	query := `
		SELECT 
			product_models.name,
			product_models.category,
			product_models.units,
			product_models.tax,
			product_models.amount,
			product_models.status,
			product_models.discount,
			
			 product_models.id AS pid 
		FROM 
			product_models 
		INNER JOIN 
			vendor ON product_models.vendor_id = vendor.id WHERE product_models.units > 0;`

	rows, err := r.sql.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute select query: %w", err)
	}
	defer rows.Close()

	var products []model.ProductListingUsers
	for rows.Next() {
		var product model.ProductListingUsers
		err := rows.Scan(&product.Name, &product.Category, &product.Unit, &product.Tax, &product.Price, &product.Status, &product.Discount, &product.Pid)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		product.Pdetail = "http://localhost:8080/admin/listingSingleProduct/" + product.Pid
		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return products, nil
}

func (r *repository) Register(ctx context.Context, request model.AdminRegisterRequest) error {
	fmt.Println("this is in the repository Register")
	query := `INSERT INTO admin (name, gst, email, password,phone) VALUES ($1, $2, $3, $4,$5)`
	_, err := r.sql.ExecContext(ctx, query, request.Name, request.GST, request.Email, request.Password, request.Phone)
	if err != nil {
		return fmt.Errorf("failed to execute insert query: %w", err)
	}

	return nil
}
func (r *repository) Addcoupon(ctx context.Context, request model.Coupon) error {
	fmt.Println("this is in the repository Register")
	query := `INSERT INTO coupon (code, expiry, min_amount,amount,max_amount) VALUES ($1, $2, $3, $4,$5)`
	_, err := r.sql.ExecContext(ctx, query, request.Code, request.Expiry, request.Minamount, request.Amount, request.Maxamount)
	if err != nil {
		return fmt.Errorf("failed to execute insert query: %w", err)
	}

	return nil
}
func (r *repository) GetCoupnExist(ctx context.Context, id string) (string, error) {
	query := `SELECT id FROM coupon WHERE code = $1 `
	var exist string
	err := r.sql.QueryRowContext(ctx, query, id).Scan(&exist)
	fmt.Println("this is in the repo layerrr for GetCoupnExist!!!", exist)
	if err != nil {
		return "", fmt.Errorf("failed to get cart: %w", err)
	}
	return exist, nil
}
func (r *repository) Deletecoupon(ctx context.Context, id string) error {
	query := `UPDATE coupon SET used = false WHERE id = $1`
	_, err := r.sql.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to update coupon: %w", err)
	}
	return nil
}

func (r *repository) Login(ctx context.Context, email string) (model.AdminRegisterRequest, error) {
	fmt.Println("theee !!!!!!!!!!!  LLLLoginnnnnn  ", email)
	query := `SELECT name, gst, email, password FROM admin WHERE email = $1`
	fmt.Println(`SELECT name, gst, email, password FROM admin WHERE email =  = 'adithyanunni258@gmail.com' ;`)

	var user model.AdminRegisterRequest
	err := r.sql.QueryRowContext(ctx, query, email).Scan(&user.Name, &user.GST, &user.Email, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return model.AdminRegisterRequest{}, nil
		}
		return model.AdminRegisterRequest{}, fmt.Errorf("failed to find user by email: %w", err)
	}
	fmt.Println("the data !!!! ", user)

	return user, nil
}
func (r *repository) Listing(ctx context.Context) ([]model.Coupon, error) {
	query := `
		SELECT code,expiry,min_amount,amount from coupon;`

	rows, err := r.sql.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute select query: %w", err)
	}
	defer rows.Close()

	var products []model.Coupon
	for rows.Next() {
		var product model.Coupon
		err := rows.Scan(&product.Code, &product.Expiry, &product.Minamount, &product.Amount)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return products, nil
}
func (r *repository) LatestListing(ctx context.Context) ([]model.Coupon, error) {
	fmt.Println("this is lia couppp")
	query := `
		SELECT code,expiry,min_amount,amount from coupon ORDER BY id DESC;`

	rows, err := r.sql.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute select query: %w", err)
	}
	defer rows.Close()

	var products []model.Coupon
	for rows.Next() {
		var product model.Coupon
		err := rows.Scan(&product.Code, &product.Expiry, &product.Minamount, &product.Amount)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return products, nil

}
func (r *repository) ActiveListing(ctx context.Context) ([]model.Coupon, error) {
	fmt.Println("this is lia couppp")
	query := `
		SELECT code,expiry,min_amount,amount from coupon WHERE TO_DATE(expiry, 'DD/MM/YYYY') > CURRENT_DATE;`

	rows, err := r.sql.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute select query: %w", err)
	}
	defer rows.Close()

	var products []model.Coupon
	for rows.Next() {
		var product model.Coupon
		err := rows.Scan(&product.Code, &product.Expiry, &product.Minamount, &product.Amount)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return products, nil

}
func (r *repository) PhighListing(ctx context.Context) ([]model.ProductListingUsers, error) {
	query := `
		SELECT 
			product_models.name,
			product_models.category,
			product_models.units,
			product_models.tax,
			product_models.amount,
			product_models.status,
			product_models.discount,
			
			 product_models.id AS pid 
		FROM 
			product_models 
		INNER JOIN 
			vendor ON product_models.vendor_id = vendor.id WHERE product_models.units > 0 ORDER BY product_models.amount DESC;`

	rows, err := r.sql.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute select query: %w", err)
	}
	defer rows.Close()

	var products []model.ProductListingUsers
	for rows.Next() {
		var product model.ProductListingUsers
		err := rows.Scan(&product.Name, &product.Category, &product.Unit, &product.Tax, &product.Price, &product.Status, &product.Discount, &product.Pid)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		product.Pdetail = "http://localhost:8080/user/listingSingleProduct/" + product.Pid
		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return products, nil
}
func (r *repository) PlowListing(ctx context.Context) ([]model.ProductListingUsers, error) {
	query := `
		SELECT 
			product_models.name,
			product_models.category,
			product_models.units,
			product_models.tax,
			product_models.amount,
			product_models.status,
			product_models.discount,
			
			 product_models.id AS pid 
		FROM 
			product_models 
		INNER JOIN 
			vendor ON product_models.vendor_id = vendor.id WHERE product_models.units > 0 ORDER BY product_models.amount ASC;`

	rows, err := r.sql.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute select query: %w", err)
	}
	defer rows.Close()

	var products []model.ProductListingUsers
	for rows.Next() {
		var product model.ProductListingUsers
		err := rows.Scan(&product.Name, &product.Category, &product.Unit, &product.Tax, &product.Price, &product.Status, &product.Discount, &product.Pid)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		product.Pdetail = "http://localhost:8080/user/listingSingleProduct/" + product.Pid
		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return products, nil
}
func (r *repository) InAZListing(ctx context.Context) ([]model.ProductListingUsers, error) {
	query := `
		SELECT 
			product_models.name,
			product_models.category,
			product_models.units,
			product_models.tax,
			product_models.amount,
			product_models.status,
			product_models.discount,
			
			 product_models.id AS pid 
		FROM 
			product_models 
		INNER JOIN 
			vendor ON product_models.vendor_id = vendor.id WHERE product_models.units > 0 ORDER BY  LOWER(product_models.name);`

	rows, err := r.sql.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute select query: %w", err)
	}
	defer rows.Close()

	var products []model.ProductListingUsers
	for rows.Next() {
		var product model.ProductListingUsers
		err := rows.Scan(&product.Name, &product.Category, &product.Unit, &product.Tax, &product.Price, &product.Status, &product.Discount, &product.Pid)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		product.Pdetail = "http://localhost:8080/user/listingSingleProduct/" + product.Pid
		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return products, nil
}
func (r *repository) InZAListing(ctx context.Context) ([]model.ProductListingUsers, error) {
	query := `
		SELECT 
			product_models.name,
			product_models.category,
			product_models.units,
			product_models.tax,
			product_models.amount,
			product_models.status,
			product_models.discount,
			
			 product_models.id AS pid 
		FROM 
			product_models 
		INNER JOIN 
			vendor ON product_models.vendor_id = vendor.id WHERE product_models.units > 0 ORDER BY  LOWER(product_models.name) DESC;`

	rows, err := r.sql.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute select query: %w", err)
	}
	defer rows.Close()

	var products []model.ProductListingUsers
	for rows.Next() {
		var product model.ProductListingUsers
		err := rows.Scan(&product.Name, &product.Category, &product.Unit, &product.Tax, &product.Price, &product.Status, &product.Discount, &product.Pid)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		product.Pdetail = "http://localhost:8080/user/listingSingleProduct/" + product.Pid
		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return products, nil
}
func (r *repository) ListingSingle(ctx context.Context, id string) ([]model.ProductListDetailed, error) {
	query := `
		SELECT 
			product_models.name,
			product_models.category,
			product_models.units,
			product_models.tax,
			product_models.amount,
			product_models.status,
			product_models.discount,
			vendor.name AS vendorName,
            product_models.id AS pid, 
			vendor.email AS vendorEmail,
			vendor.gst AS vendorgst,
			
			vendor.id AS vendorid,
			product_models.description AS pds
		FROM 
			product_models 
		INNER JOIN 
			vendor ON product_models.vendor_id = vendor.id WHERE product_models.units > 0 AND product_models.id=$1;`

	rows, err := r.sql.QueryContext(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to execute select query: %w", err)
	}
	defer rows.Close()

	var products []model.ProductListDetailed
	for rows.Next() {
		var product model.ProductListDetailed
		err := rows.Scan(&product.Name, &product.Category, &product.Unit, &product.Tax, &product.Price, &product.Status, &product.Discount, &product.VendorName, &product.Pid, &product.VEmail, &product.VGst, &product.VId, &product.Pds)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return products, nil
}
