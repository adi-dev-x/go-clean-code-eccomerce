package vendor

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strings"
	"time"

	"myproject/pkg/model"
)

type Repository interface {
	Register(ctx context.Context, request model.VendorRegisterRequest) error

	Login(ctx context.Context, email string) (model.VendorRegisterRequest, error)
	///product
	Listing(ctx context.Context, id string) ([]model.ProductList, error)
	CategoryListing(ctx context.Context, category string, id string) ([]model.ProductList, error)
	AddProduct(ctx context.Context, request model.Product) error
	LatestListing(ctx context.Context, id string) ([]model.ProductList, error)
	PhighListing(ctx context.Context, id string) ([]model.ProductList, error)
	PlowListing(ctx context.Context, id string) ([]model.ProductList, error)
	InAZListing(ctx context.Context, id string) ([]model.ProductList, error)
	InZAListing(ctx context.Context, id string) ([]model.ProductList, error)
	Getid(ctx context.Context, username string) string
	UpdateProduct(ctx context.Context, query string, args []interface{}) error
	//orders
	ListAllOrders(ctx context.Context, id string) ([]model.ListOrdersVendor, error)
	ListReturnedOrders(ctx context.Context, id string) ([]model.ListOrdersVendor, error)
	ListFailedOrders(ctx context.Context, id string) ([]model.ListOrdersVendor, error)
	ListCompletedOrders(ctx context.Context, id string) ([]model.ListOrdersVendor, error)
	ListPendingOrders(ctx context.Context, id string) ([]model.ListOrdersVendor, error)

	GetSalesFactByDate(ctx context.Context, filterType string, startDate, endDate time.Time, vendorID string) ([]model.Salesfact, error)
	SalesReportOrdersWeekly(ctx context.Context, vendorID string) ([]model.ListOrdersVendor, error)
	SalesReportOrdersDaily(ctx context.Context, vendorID string) ([]model.ListOrdersVendor, error)

	SalesReportOrdersMonthly(ctx context.Context, vendorID string) ([]model.ListOrdersVendor, error)
	SalesReportOrdersYearly(ctx context.Context, vendorID string) ([]model.ListOrdersVendor, error)
	SalesReportOrdersCustom(ctx context.Context, startDate, endDate time.Time, vendorID string) ([]model.ListOrdersVendor, error)
	//return
	GetSingleItem(ctx context.Context, id string, oid string) (model.ListAllOrdersCheck, error)

	//stock
	IncreaseStock(ctx context.Context, id string, unit int) error

	UpdateOiStatus(ctx context.Context, id, ty string) error
	CreditWallet(ctx context.Context, id string, amt float64) (string, error)
	UpdateWalletTransaction(ctx context.Context, value interface{}) error
	ChangeOrderStatus(ctx context.Context, id string) error
	VerifyOtp(ctx context.Context, email string)
}

type repository struct {
	sql *sql.DB
}

func NewRepository(sqlDB *sql.DB) Repository {
	return &repository{
		sql: sqlDB,
	}
}
func (r *repository) VerifyOtp(ctx context.Context, email string) {
	query := `
	UPDATE vendor
	SET verification =true
	WHERE email = $1
	`

	_, err := r.sql.ExecContext(ctx, query, email)

	if err != nil {
		fmt.Errorf("failed to execute update query: %w", err)
	}

}
func (r *repository) UpdateProduct(ctx context.Context, query string, args []interface{}) error {
	queryWithParams := query
	for _, arg := range args {
		queryWithParams = strings.Replace(queryWithParams, "?", fmt.Sprintf("'%v'", arg), 1)
	}
	fmt.Println("Executing update with query:", queryWithParams)
	fmt.Println("Arguments:", args)
	fmt.Println("Executing update for email:", args[len(args)-1]) // Email is the last argument
	_, err := r.sql.ExecContext(ctx, queryWithParams)
	if err != nil {
		return fmt.Errorf("failed to execute update query: %w", err)
	}
	return nil
}
func (r *repository) UpdateWalletTransaction(ctx context.Context, value interface{}) error {
	fmt.Println("changing innnnn   !!!!!! UpdateWalletTransaction")
	values, ok := value.([]interface{})
	if !ok {
		return fmt.Errorf("invalid input")
	}
	Amt := values[0]
	id := values[1]
	Type := values[2]
	usid := values[3]
	fmt.Println("this is the UpdateWalletTransaction in repo!!@@@@@", reflect.TypeOf(Amt), "____", Amt, "!!", id, "##", Type)
	query := `
	INSERT INTO wallet_transactions (
		wallet_id,
		amount,
		transaction_type,
		user_id,
		created_at
		
	) VALUES (
		$1, $2, $3,$4, CURRENT_TIMESTAMP
	) RETURNING id;
`
	var tid string
	fmt.Println("this is the id ", tid, "user_id,", usid)

	err := r.sql.QueryRowContext(ctx, query, id, Amt, Type, usid).Scan(&tid)
	if err != nil {
		return fmt.Errorf("there is error in insertion")
	}

	return nil

}
func (r *repository) CreditWallet(ctx context.Context, id string, amt float64) (string, error) {
	query := `
	UPDATE wallet
	SET balance = balance + $1
	WHERE user_id = $2
	RETURNING id;
`
	var Wallet_id string
	err := r.sql.QueryRowContext(ctx, query, amt, id).Scan(&Wallet_id)
	fmt.Println("hey adiii CreditWallet????!!!!!", Wallet_id, "wallet_id !", "id!", id)
	if err != nil {
		return "", fmt.Errorf("failed to execute update query: %w", err)
	}

	return Wallet_id, nil
}
func (r *repository) ChangeOrderStatus(ctx context.Context, id string) error {
	fmt.Println("changing in the order status QQQQ")
	query := `
	UPDATE orders
	SET status ='Cancelled'
	WHERE id = $1
	
`

	err := r.sql.QueryRowContext(ctx, query, id)

	if err != nil {
		return fmt.Errorf("failed to execute update query: %w", err)
	}

	return nil
}

// /
func (r *repository) UpdateOiStatus(ctx context.Context, id string, ty string) error {

	query := `
	UPDATE order_items
	SET returned = true,
	re_cl=$1
	WHERE id = $2
	RETURNING id;
`
	var Oi_id string
	err := r.sql.QueryRowContext(ctx, query, ty, id).Scan(&Oi_id)
	if err != nil {
		return fmt.Errorf("failed to execute update query: %w", err)
	}

	return nil
}

// increse stock
func (r *repository) IncreaseStock(ctx context.Context, id string, unit int) error {
	fmt.Println("this is in the IncreaseStock!!", id, "unitss ", unit)
	query := `
	UPDATE product_models
	SET units = units + $1
	WHERE id = $2
	RETURNING id;
`
	var Product_id string
	err := r.sql.QueryRowContext(ctx, query, unit, id).Scan(&Product_id)
	if err != nil {
		return fmt.Errorf("failed to execute update query: %w", err)
	}

	return nil

}

// /return
func (r *repository) GetSingleItem(ctx context.Context, id string, oid string) (model.ListAllOrdersCheck, error) {
	var order model.ListAllOrdersCheck

	query := `SELECT p.name,  oi.quantity,   mo.status, oi.returned, 
    oi.price,oi.product_id AS pid,DATE(oi.created_at) AS date,mo.user_id ,v.id AS vid,u.email AS usmail,mo.id AS mid
     FROM order_items oi 
    JOIN  product_models p ON oi.product_id = p.id 
    JOIN  orders mo ON oi.order_id = mo.id 
    JOIN  vendor v ON p.vendor_id = v.id 
	JOIN  users  u ON oi.user_id=u.id
    WHERE  v.id = $1 AND oi.id = $2;
  
`
	err := r.sql.QueryRowContext(ctx, query, id, oid).Scan(&order.Name, &order.Unit, &order.Status, &order.Returned, &order.Amount, &order.Pid, &order.Date, &order.Usid, &order.Vid, &order.Usmail, &order.Moid)
	if err != nil {
		return model.ListAllOrdersCheck{}, fmt.Errorf("error in exequting query in  GetSingleItem")
	}
	return order, nil
}

// //orders
func (r *repository) ListAllOrders(ctx context.Context, id string) ([]model.ListOrdersVendor, error) {
	query := `SELECT 
    p.name, oi.quantity, mo.status, oi.returned, oi.price, oi.product_id AS pid, 
    DATE(oi.created_at) AS date, u.firstname || ' ' || u.lastname AS user, 
    COALESCE(a.address1, '') || ' ' || COALESCE(a.address2, '') || ' ' || COALESCE(a.address3, '') || ' ' ||
    COALESCE(a.city, '') || ' ' || COALESCE(a.state, '') || ' ' || COALESCE(a.pin, '') || ' ' || COALESCE(a.country, '') AS user_ad, 
     oi.id AS oid
    FROM order_items oi 
    JOIN product_models p ON oi.product_id = p.id 
    JOIN vendor v ON p.vendor_id = v.id 
    JOIN orders mo ON oi.order_id = mo.id 
    JOIN users u ON mo.user_id = u.id 
    JOIN address a ON mo.address_id = a.address_id 
    WHERE v.id = $1;`

	var orders []model.ListOrdersVendor

	rows, err := r.sql.QueryContext(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("error executing query in ListOrdersVendor: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var order model.ListOrdersVendor
		err := rows.Scan(&order.Name, &order.Unit, &order.Status, &order.Returned, &order.Amount, &order.Pid, &order.Date, &order.User, &order.Add, &order.Oid)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return orders, nil
}
func (r *repository) ListReturnedOrders(ctx context.Context, id string) ([]model.ListOrdersVendor, error) {
	query := `SELECT 
    p.name, oi.quantity, mo.status, oi.returned, oi.price, oi.product_id AS pid, 
    DATE(oi.created_at) AS date, u.firstname || ' ' || u.lastname AS user, 
    COALESCE(a.address1, '') || ' ' || COALESCE(a.address2, '') || ' ' || COALESCE(a.address3, '') || ' ' ||
    COALESCE(a.city, '') || ' ' || COALESCE(a.state, '') || ' ' || COALESCE(a.pin, '') || ' ' || COALESCE(a.country, '') AS user_ad 
    FROM order_items oi 
    JOIN product_models p ON oi.product_id = p.id 
    JOIN vendor v ON p.vendor_id = v.id 
    JOIN orders mo ON oi.order_id = mo.id 
    JOIN users u ON mo.user_id = u.id 
    JOIN address a ON mo.address_id = a.address_id 
    WHERE v.id = $1 AND oi.returned=true;`

	var orders []model.ListOrdersVendor

	rows, err := r.sql.QueryContext(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("error executing query in ListOrdersVendor: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var order model.ListOrdersVendor
		err := rows.Scan(&order.Name, &order.Unit, &order.Status, &order.Returned, &order.Amount, &order.Pid, &order.Date, &order.User, &order.Add)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return orders, nil
}
func (r *repository) ListFailedOrders(ctx context.Context, id string) ([]model.ListOrdersVendor, error) {
	query := `SELECT 
    p.name, oi.quantity, mo.status, oi.returned, oi.price, oi.product_id AS pid, 
    DATE(oi.created_at) AS date, u.firstname || ' ' || u.lastname AS user, 
    COALESCE(a.address1, '') || ' ' || COALESCE(a.address2, '') || ' ' || COALESCE(a.address3, '') || ' ' ||
    COALESCE(a.city, '') || ' ' || COALESCE(a.state, '') || ' ' || COALESCE(a.pin, '') || ' ' || COALESCE(a.country, '') AS user_ad 
    FROM order_items oi 
    JOIN product_models p ON oi.product_id = p.id 
    JOIN vendor v ON p.vendor_id = v.id 
    JOIN orders mo ON oi.order_id = mo.id 
    JOIN users u ON mo.user_id = u.id 
    JOIN address a ON mo.address_id = a.address_id 
    WHERE v.id = $1 AND  mo.status='Failed';`

	var orders []model.ListOrdersVendor

	rows, err := r.sql.QueryContext(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("error executing query in ListOrdersVendor: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var order model.ListOrdersVendor
		err := rows.Scan(&order.Name, &order.Unit, &order.Status, &order.Returned, &order.Amount, &order.Pid, &order.Date, &order.User, &order.Add)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return orders, nil
}
func (r *repository) ListCompletedOrders(ctx context.Context, id string) ([]model.ListOrdersVendor, error) {
	query := `SELECT 
    p.name, oi.quantity, mo.status, oi.returned, oi.price, oi.product_id AS pid, 
    DATE(oi.created_at) AS date, u.firstname || ' ' || u.lastname AS user, 
    COALESCE(a.address1, '') || ' ' || COALESCE(a.address2, '') || ' ' || COALESCE(a.address3, '') || ' ' ||
    COALESCE(a.city, '') || ' ' || COALESCE(a.state, '') || ' ' || COALESCE(a.pin, '') || ' ' || COALESCE(a.country, '') AS user_ad 
    FROM order_items oi 
    JOIN product_models p ON oi.product_id = p.id 
    JOIN vendor v ON p.vendor_id = v.id 
    JOIN orders mo ON oi.order_id = mo.id 
    JOIN users u ON mo.user_id = u.id 
    JOIN address a ON mo.address_id = a.address_id 
    WHERE v.id = $1 AND  mo.status='Completed';`

	var orders []model.ListOrdersVendor

	rows, err := r.sql.QueryContext(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("error executing query in ListOrdersVendor: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var order model.ListOrdersVendor
		err := rows.Scan(&order.Name, &order.Unit, &order.Status, &order.Returned, &order.Amount, &order.Pid, &order.Date, &order.User, &order.Add)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return orders, nil
}
func (r *repository) SalesReportOrdersCustom(ctx context.Context, startDate, endDate time.Time, vendorID string) ([]model.ListOrdersVendor, error) {

	query := `SELECT 
    p.name, oi.quantity, mo.status, oi.returned, oi.price, oi.product_id AS pid, 
    DATE(oi.created_at) AS date, u.firstname || ' ' || u.lastname AS user, 
    COALESCE(a.address1, '') || ' ' || COALESCE(a.address2, '') || ' ' || COALESCE(a.address3, '') || ' ' ||
    COALESCE(a.city, '') || ' ' || COALESCE(a.state, '') || ' ' || COALESCE(a.pin, '') || ' ' || COALESCE(a.country, '') AS user_ad 
    FROM order_items oi 
    JOIN product_models p ON oi.product_id = p.id 
    JOIN vendor v ON p.vendor_id = v.id 
    JOIN orders mo ON oi.order_id = mo.id 
    JOIN users u ON mo.user_id = u.id 
    JOIN address a ON mo.address_id = a.address_id 
    WHERE v.id = $1 AND  mo.status='Completed' AND DATE(oi.created_at) BETWEEN $2 AND $3 ;`

	var orders []model.ListOrdersVendor

	rows, err := r.sql.QueryContext(ctx, query, vendorID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("error executing query in ListOrdersVendor: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var order model.ListOrdersVendor
		err := rows.Scan(&order.Name, &order.Unit, &order.Status, &order.Returned, &order.Amount, &order.Pid, &order.Date, &order.User, &order.Add)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return orders, nil
}
func (r *repository) SalesReportOrdersYearly(ctx context.Context, vendorID string) ([]model.ListOrdersVendor, error) {

	query := `SELECT EXTRACT(YEAR FROM oi.created_at) AS checks, p.name, oi.quantity, mo.status, oi.returned, oi.price * 
	oi.quantity AS total_price, 
	oi.product_id AS pid, u.firstname || ' ' || u.lastname AS user, 
	COALESCE(a.address1, '') || ' ' || COALESCE(a.address2, '') || ' ' || COALESCE(a.address3, '') || ' ' || COALESCE(a.city, '') || ' ' || COALESCE(a.state, '') || ' ' || COALESCE(a.pin, '') || ' ' || COALESCE(a.country, '') 
	AS user_ad ,DATE(oi.created_at) AS date FROM order_items oi JOIN product_models p ON oi.product_id = p.id JOIN vendor v ON p.vendor_id = v.id JOIN orders mo ON oi.order_id = mo.id JOIN users u ON mo.user_id = u.id JOIN address a ON mo.address_id = a.address_id 
	WHERE v.id = $1 AND mo.status = 'Completed' AND oi.returned=false ORDER BY checks;
`

	var orders []model.ListOrdersVendor

	rows, err := r.sql.QueryContext(ctx, query, vendorID)
	if err != nil {
		return nil, fmt.Errorf("error executing query in ListOrdersVendor: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var order model.ListOrdersVendor
		err := rows.Scan(&order.ListDate, &order.Name, &order.Unit, &order.Status, &order.Returned, &order.Amount, &order.Pid, &order.User, &order.Add, &order.Date)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return orders, nil
}
func (r *repository) SalesReportOrdersMonthly(ctx context.Context, vendorID string) ([]model.ListOrdersVendor, error) {

	query := `SELECT EXTRACT(MONTH FROM oi.created_at) AS checks, p.name, oi.quantity, mo.status, oi.returned, oi.price * 
	oi.quantity AS total_price, 
	oi.product_id AS pid, u.firstname || ' ' || u.lastname AS user, 
	COALESCE(a.address1, '') || ' ' || COALESCE(a.address2, '') || ' ' || COALESCE(a.address3, '') || ' ' || COALESCE(a.city, '') || ' ' || COALESCE(a.state, '') || ' ' || COALESCE(a.pin, '') || ' ' || COALESCE(a.country, '') 
	AS user_ad, DATE(oi.created_at) AS date FROM order_items oi JOIN product_models p ON oi.product_id = p.id JOIN vendor v ON p.vendor_id = v.id JOIN orders mo ON oi.order_id = mo.id JOIN users u ON mo.user_id = u.id JOIN address a ON mo.address_id = a.address_id 
	WHERE v.id = $1 AND mo.status = 'Completed' AND oi.returned=false ORDER BY checks;
`

	var orders []model.ListOrdersVendor

	rows, err := r.sql.QueryContext(ctx, query, vendorID)
	if err != nil {
		return nil, fmt.Errorf("error executing query in ListOrdersVendor: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var order model.ListOrdersVendor
		err := rows.Scan(&order.ListDate, &order.Name, &order.Unit, &order.Status, &order.Returned, &order.Amount, &order.Pid, &order.User, &order.Add, &order.Date)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return orders, nil
}
func (r *repository) SalesReportOrdersWeekly(ctx context.Context, vendorID string) ([]model.ListOrdersVendor, error) {

	query := `SELECT EXTRACT(WEEK FROM oi.created_at) AS checks, p.name, oi.quantity, mo.status, oi.returned, oi.price * 
	oi.quantity AS total_price, 
	oi.product_id AS pid, u.firstname || ' ' || u.lastname AS user, 
	COALESCE(a.address1, '') || ' ' || COALESCE(a.address2, '') || ' ' || COALESCE(a.address3, '') || ' ' || COALESCE(a.city, '') || ' ' || COALESCE(a.state, '') || ' ' || COALESCE(a.pin, '') || ' ' || COALESCE(a.country, '') 
	AS user_ad ,DATE(oi.created_at) AS date FROM order_items oi JOIN product_models p ON oi.product_id = p.id JOIN vendor v ON p.vendor_id = v.id JOIN orders mo ON oi.order_id = mo.id JOIN users u ON mo.user_id = u.id JOIN address a ON mo.address_id = a.address_id 
	WHERE v.id = $1 AND mo.status = 'Completed' AND oi.returned=false ORDER BY checks;
`

	var orders []model.ListOrdersVendor

	rows, err := r.sql.QueryContext(ctx, query, vendorID)
	if err != nil {
		return nil, fmt.Errorf("error executing query in ListOrdersVendor: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var order model.ListOrdersVendor
		err := rows.Scan(&order.ListDate, &order.Name, &order.Unit, &order.Status, &order.Returned, &order.Amount, &order.Pid, &order.User, &order.Add, &order.Date)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return orders, nil
}
func (r *repository) SalesReportOrdersDaily(ctx context.Context, vendorID string) ([]model.ListOrdersVendor, error) {

	query := `SELECT EXTRACT(DAY FROM oi.created_at) AS checks, p.name, oi.quantity, mo.status, oi.returned, oi.price * 
	oi.quantity AS total_price, 
	oi.product_id AS pid, u.firstname || ' ' || u.lastname AS user, 
	COALESCE(a.address1, '') || ' ' || COALESCE(a.address2, '') || ' ' || COALESCE(a.address3, '') || ' ' || COALESCE(a.city, '') || ' ' || COALESCE(a.state, '') || ' ' || COALESCE(a.pin, '') || ' ' || COALESCE(a.country, '') 
	AS user_ad ,DATE(oi.created_at) AS date FROM order_items oi JOIN product_models p ON oi.product_id = p.id JOIN vendor v ON p.vendor_id = v.id JOIN orders mo ON oi.order_id = mo.id JOIN users u ON mo.user_id = u.id JOIN address a ON mo.address_id = a.address_id 
	WHERE v.id = $1 AND mo.status = 'Completed' AND oi.returned=false ORDER BY checks;
`

	var orders []model.ListOrdersVendor

	rows, err := r.sql.QueryContext(ctx, query, vendorID)
	if err != nil {
		return nil, fmt.Errorf("error executing query in ListOrdersVendor: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var order model.ListOrdersVendor
		err := rows.Scan(&order.ListDate, &order.Name, &order.Unit, &order.Status, &order.Returned, &order.Amount, &order.Pid, &order.User, &order.Add, &order.Date)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return orders, nil
}
func (r *repository) GetSalesFactByDate(ctx context.Context, filterType string, startDate, endDate time.Time, vendorID string) ([]model.Salesfact, error) {
	var query string
	var args []interface{}
	// const dateFormat = "2006-01-02"
	fmt.Println("this is the filter ", filterType)

	switch filterType {
	case "Yearly":
		query = `
			SELECT EXTRACT(YEAR FROM oi.created_at) AS year, 
			       SUM(oi.price * oi.quantity) AS revenue,
			       SUM(oi.price * oi.quantity - oi.discount * oi.quantity) AS total_sales, 
			       SUM(oi.discount * oi.quantity) AS total_discount,
			       COUNT(*) AS total_orders
			FROM order_items oi
			JOIN product_models pm ON oi.product_id = pm.id
			JOIN orders o ON oi.order_id = o.id
			WHERE pm.vendor_id = $1
			AND EXTRACT(YEAR FROM oi.created_at) = EXTRACT(YEAR FROM CURRENT_DATE) 
			 AND oi.returned=false AND o.status = 'Completed'
			GROUP BY year`
		args = append(args, vendorID)
	case "Weekly":
		query = `
			SELECT EXTRACT(WEEK FROM oi.created_at) AS week, 
			       SUM(oi.price * oi.quantity) AS revenue,
			       SUM(oi.price * oi.quantity - oi.discount * oi.quantity) AS total_sales, 
			       SUM(oi.discount * oi.quantity) AS total_discount,
			       COUNT(*) AS total_orders
			FROM order_items oi
			JOIN product_models pm ON oi.product_id = pm.id
			JOIN orders o ON oi.order_id = o.id

			WHERE pm.vendor_id = $1  AND oi.returned=false AND o.status = 'Completed'
			GROUP BY week`
		args = append(args, vendorID)
	case "Monthly":
		fmt.Println("this is inside monthly")

		query = `
		SELECT 
        EXTRACT(MONTH FROM oi.created_at) AS month, 
		SUM(oi.price * oi.quantity) AS revenue,
        SUM(oi.price * oi.quantity - oi.discount * oi.quantity) AS total_sales,
        SUM(oi.discount * oi.quantity) AS total_discount,
		COUNT(*) AS total_orders
        FROM order_items oi
        JOIN  product_models pm ON oi.product_id = pm.id
		JOIN orders o ON oi.order_id = o.id
        WHERE  pm.vendor_id = $1 AND EXTRACT(MONTH FROM oi.created_at) = EXTRACT(MONTH FROM CURRENT_DATE) 
		AND EXTRACT(YEAR FROM oi.created_at) = EXTRACT(YEAR FROM CURRENT_DATE)  AND oi.returned=false
		AND o.status = 'Completed'
		GROUP BY   month ORDER BY  month;
        `
		args = append(args, vendorID)
	case "Daily":
		query = `
         SELECT EXTRACT(DAY FROM oi.created_at) AS day, 
          SUM(oi.price * oi.quantity) AS revenue,
          SUM(oi.price * oi.quantity - oi.discount * oi.quantity) AS total_sales, 
          SUM(oi.discount * oi.quantity) AS total_discount,
          COUNT(*) AS total_orders
			FROM order_items oi
		JOIN product_models pm ON oi.product_id = pm.id
		JOIN orders o ON oi.order_id = o.id
		WHERE DATE(oi.created_at) = DATE(CURRENT_DATE)
        AND pm.vendor_id = $1
       AND oi.returned = false
       AND o.status = 'Completed'
		GROUP BY day;
			
			`
		args = append(args, vendorID)
	case "Custom":
		fmt.Println("inside the custom switch", endDate, "!!!!", startDate)
		query = `
		SELECT 1 AS day,
       SUM(oi.price * oi.quantity) AS revenue,
       SUM(oi.price * oi.quantity - oi.discount * oi.quantity) AS total_sales, 
       SUM(oi.discount * oi.quantity) AS total_discount,
       COUNT(*) AS total_orders
       FROM order_items oi
       JOIN product_models pm ON oi.product_id = pm.id
       JOIN orders o ON oi.order_id = o.id
       WHERE DATE(oi.created_at) BETWEEN $1 AND $2 
       AND pm.vendor_id = $3 
       AND oi.returned = false
       AND o.status = 'Completed';
	  `

		args = append(args, startDate, endDate, vendorID)
	default:
		return nil, fmt.Errorf("invalid filter type")
	}

	rows, err := r.sql.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var salesFacts []model.Salesfact
	for rows.Next() {
		var salesFact model.Salesfact
		err := rows.Scan(&salesFact.Date, &salesFact.Revenue, &salesFact.TotalSales, &salesFact.TotalDiscount, &salesFact.TotalOrders)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		salesFacts = append(salesFacts, salesFact)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return salesFacts, nil
}

func (r *repository) ListPendingOrders(ctx context.Context, id string) ([]model.ListOrdersVendor, error) {
	query := `SELECT 
    p.name, oi.quantity, mo.status, oi.returned, oi.price, oi.product_id AS pid, 
    DATE(oi.created_at) AS date, u.firstname || ' ' || u.lastname AS user, 
    COALESCE(a.address1, '') || ' ' || COALESCE(a.address2, '') || ' ' || COALESCE(a.address3, '') || ' ' ||
    COALESCE(a.city, '') || ' ' || COALESCE(a.state, '') || ' ' || COALESCE(a.pin, '') || ' ' || COALESCE(a.country, '') AS user_ad 
    FROM order_items oi 
    JOIN product_models p ON oi.product_id = p.id 
    JOIN vendor v ON p.vendor_id = v.id 
    JOIN orders mo ON oi.order_id = mo.id 
    JOIN users u ON mo.user_id = u.id 
    JOIN address a ON mo.address_id = a.address_id 
    WHERE v.id = $1 AND  mo.status='Pending';`

	var orders []model.ListOrdersVendor

	rows, err := r.sql.QueryContext(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("error executing query in ListOrdersVendor: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var order model.ListOrdersVendor
		err := rows.Scan(&order.Name, &order.Unit, &order.Status, &order.Returned, &order.Amount, &order.Pid, &order.Date, &order.User, &order.Add)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return orders, nil
}

// //
func (r *repository) Getid(ctx context.Context, username string) string {
	var id string
	fmt.Println("this is in the repository Register !!!")
	query := `select id from vendor where email=$1;`
	fmt.Println(query, username)
	row := r.sql.QueryRowContext(ctx, query, username)
	err := row.Scan(&id)
	fmt.Println(err)
	fmt.Println("this is id returning from Getid:::", id)

	return id
}

func (r *repository) Register(ctx context.Context, request model.VendorRegisterRequest) error {
	fmt.Println("this is in the repository Register")
	query := `INSERT INTO vendor (name, gst, email, password,phone) VALUES ($1, $2, $3, $4,$5)`
	_, err := r.sql.ExecContext(ctx, query, request.Name, request.GST, request.Email, request.Password, request.Phone)
	if err != nil {
		return fmt.Errorf("failed to execute insert query: %w", err)
	}

	return nil
}
func (r *repository) AddProduct(ctx context.Context, request model.Product) error {
	fmt.Println("this is in the repository Register")
	query := `INSERT INTO product_models (name, category, status, tax,amount,units,vendor_id,discount,description) VALUES ($1, $2, $3, $4,$5,$6,$7,$8,$9)`
	_, err := r.sql.ExecContext(ctx, query, request.Name, request.Category, request.Status, request.Tax, request.Price, request.Unit, request.Vendorid, request.Discount, request.Description)
	if err != nil {
		return fmt.Errorf("failed to execute insert query: %w", err)
	}

	return nil
}
func (r *repository) Login(ctx context.Context, email string) (model.VendorRegisterRequest, error) {
	fmt.Println("Attempting to login with email:", email)

	// SQL query to fetch user details based on email
	query := `SELECT name, gst, email, password FROM vendor WHERE email = $1 AND verification=true`
	fmt.Printf("Executing query: %s\n", query)

	var user model.VendorRegisterRequest

	// Execute the query and scan the result into the user struct
	err := r.sql.QueryRowContext(ctx, query, email).Scan(&user.Name, &user.GST, &user.Email, &user.Password)
	if err != nil {
		// Check if no rows were returned
		if err == sql.ErrNoRows {
			fmt.Println("No user found with the provided email.")
			return model.VendorRegisterRequest{}, nil
		}
		// For other types of errors, wrap and return the error
		return model.VendorRegisterRequest{}, fmt.Errorf("failed to find user by email: %w", err)
	}

	fmt.Println("User data retrieved:", user)

	return user, nil
}
func (r *repository) Listing(ctx context.Context, id string) ([]model.ProductList, error) {
	query := `
		SELECT 
			product_models.name,
			product_models.category,
			product_models.units,
			product_models.tax,
			product_models.amount,
			product_models.status,
			product_models.discount,
			vendor.name AS vendorName  
		FROM 
			product_models 
		INNER JOIN 
			vendor ON product_models.vendor_id = vendor.id WHERE product_models.units > 0 AND vendor.id=$1;`

	rows, err := r.sql.QueryContext(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to execute select query: %w", err)
	}
	defer rows.Close()

	var products []model.ProductList
	for rows.Next() {
		var product model.ProductList
		err := rows.Scan(&product.Name, &product.Category, &product.Unit, &product.Tax, &product.Price, &product.Status, &product.Discount, &product.VendorName)
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
func (r *repository) CategoryListing(ctx context.Context, category string, id string) ([]model.ProductList, error) {
	query := `
		SELECT 
			product_models.name,
			product_models.category,
			product_models.units,
			product_models.tax,
			product_models.amount,
			product_models.status,
			product_models.discount,
			vendor.name AS vendorName  
		FROM 
			product_models 
		INNER JOIN 
			vendor ON product_models.vendor_id = vendor.id WHERE product_models.units > 0 AND vendor.id=$1 AND product_models.category ILIKE '%' || $2 || '%';`

	rows, err := r.sql.QueryContext(ctx, query, id, category)
	if err != nil {
		return nil, fmt.Errorf("failed to execute select query: %w", err)
	}
	defer rows.Close()

	var products []model.ProductList
	for rows.Next() {
		var product model.ProductList
		err := rows.Scan(&product.Name, &product.Category, &product.Unit, &product.Tax, &product.Price, &product.Status, &product.Discount, &product.VendorName)
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
func (r *repository) LatestListing(ctx context.Context, id string) ([]model.ProductList, error) {
	query := `
		SELECT 
			product_models.name,
			product_models.category,
			product_models.units,
			product_models.tax,
			product_models.amount,
			product_models.status,
			product_models.discount,
			vendor.name AS vendorName  
		FROM 
			product_models 
		INNER JOIN 
			vendor ON product_models.vendor_id = vendor.id WHERE product_models.units > 0 AND vendor.id=$1 ORDER BY product_models.id DESC;`

	rows, err := r.sql.QueryContext(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to execute select query: %w", err)
	}
	defer rows.Close()

	var products []model.ProductList
	for rows.Next() {
		var product model.ProductList
		err := rows.Scan(&product.Name, &product.Category, &product.Unit, &product.Tax, &product.Price, &product.Status, &product.Discount, &product.VendorName)
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
func (r *repository) PhighListing(ctx context.Context, id string) ([]model.ProductList, error) {
	query := `
		SELECT 
			product_models.name,
			product_models.category,
			product_models.units,
			product_models.tax,
			product_models.amount,
			product_models.status,
			product_models.discount,
			vendor.name AS vendorName  
		FROM 
			product_models 
		INNER JOIN 
			vendor ON product_models.vendor_id = vendor.id WHERE product_models.units > 0 AND vendor.id=$1 ORDER BY product_models.amount DESC;`

	rows, err := r.sql.QueryContext(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to execute select query: %w", err)
	}
	defer rows.Close()

	var products []model.ProductList
	for rows.Next() {
		var product model.ProductList
		err := rows.Scan(&product.Name, &product.Category, &product.Unit, &product.Tax, &product.Price, &product.Status, &product.Discount, &product.VendorName)
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
func (r *repository) PlowListing(ctx context.Context, id string) ([]model.ProductList, error) {
	query := `
		SELECT 
			product_models.name,
			product_models.category,
			product_models.units,
			product_models.tax,
			product_models.amount,
			product_models.status,
			product_models.discount,
			vendor.name AS vendorName  
		FROM 
			product_models 
		INNER JOIN 
			vendor ON product_models.vendor_id = vendor.id WHERE product_models.units > 0 AND vendor.id=$1 ORDER BY product_models.amount ASC;`

	rows, err := r.sql.QueryContext(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to execute select query: %w", err)
	}
	defer rows.Close()

	var products []model.ProductList
	for rows.Next() {
		var product model.ProductList
		err := rows.Scan(&product.Name, &product.Category, &product.Unit, &product.Tax, &product.Price, &product.Status, &product.Discount, &product.VendorName)
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
func (r *repository) InAZListing(ctx context.Context, id string) ([]model.ProductList, error) {
	query := `
		SELECT 
			product_models.name,
			product_models.category,
			product_models.units,
			product_models.tax,
			product_models.amount,
			product_models.status,
			product_models.discount,
			vendor.name AS vendorName  
		FROM 
			product_models 
		INNER JOIN 
			vendor ON product_models.vendor_id = vendor.id WHERE product_models.units > 0 AND vendor.id=$1 ORDER BY  LOWER(product_models.name);`

	rows, err := r.sql.QueryContext(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to execute select query: %w", err)
	}
	defer rows.Close()

	var products []model.ProductList
	for rows.Next() {
		var product model.ProductList
		err := rows.Scan(&product.Name, &product.Category, &product.Unit, &product.Tax, &product.Price, &product.Status, &product.Discount, &product.VendorName)
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
func (r *repository) InZAListing(ctx context.Context, id string) ([]model.ProductList, error) {
	query := `
		SELECT 
			product_models.name,
			product_models.category,
			product_models.units,
			product_models.tax,
			product_models.amount,
			product_models.status,
			product_models.discount,
			vendor.name AS vendorName  
		FROM 
			product_models 
		INNER JOIN 
			vendor ON product_models.vendor_id = vendor.id WHERE product_models.units > 0 AND vendor.id=$1 ORDER BY  LOWER(product_models.name) DESC;`

	rows, err := r.sql.QueryContext(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to execute select query: %w", err)
	}
	defer rows.Close()

	var products []model.ProductList
	for rows.Next() {
		var product model.ProductList
		err := rows.Scan(&product.Name, &product.Category, &product.Unit, &product.Tax, &product.Price, &product.Status, &product.Discount, &product.VendorName)
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
