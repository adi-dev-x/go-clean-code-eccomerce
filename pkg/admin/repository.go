package admin

import (
	"context"
	"database/sql"
	"fmt"
	"time"

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

	///list orders
	ListAllOrders(ctx context.Context) ([]model.ListOrdersAdmin, error)
	ListReturnedOrders(ctx context.Context) ([]model.ListOrdersAdmin, error)
	ListFailedOrders(ctx context.Context) ([]model.ListOrdersAdmin, error)
	ListCompletedOrders(ctx context.Context) ([]model.ListOrdersAdmin, error)
	ListPendingOrders(ctx context.Context) ([]model.ListOrdersAdmin, error)

	SalesReportOrdersCustom(ctx context.Context, startDate, endDate time.Time) ([]model.ListOrdersAdmin, error)
	SalesReportOrdersYearly(ctx context.Context) ([]model.ListOrdersAdmin, error)
	SalesReportOrdersMonthly(ctx context.Context) ([]model.ListOrdersAdmin, error)
	SalesReportOrdersWeekly(ctx context.Context) ([]model.ListOrdersAdmin, error)
	SalesReportOrdersDaily(ctx context.Context) ([]model.ListOrdersAdmin, error)

	GetSalesFactByDate(ctx context.Context, filterType string, startDate, endDate time.Time) ([]model.Salesfact, error)
	/// Singlevendor
	ListAllOrdersSinglevendor(ctx context.Context, id string) ([]model.ListOrdersVendor, error)
	ListReturnedOrdersSinglevendor(ctx context.Context, id string) ([]model.ListOrdersVendor, error)
	ListFailedOrdersSinglevendor(ctx context.Context, id string) ([]model.ListOrdersVendor, error)
	ListCompletedOrdersSinglevendor(ctx context.Context, id string) ([]model.ListOrdersVendor, error)
	ListPendingOrdersSinglevendor(ctx context.Context, id string) ([]model.ListOrdersVendor, error)

	GetSalesFactByDateSinglevendor(ctx context.Context, filterType string, startDate, endDate time.Time, vendorID string) ([]model.Salesfact, error)
	SalesReportOrdersWeeklySinglevendor(ctx context.Context, vendorID string) ([]model.ListOrdersVendor, error)
	SalesReportOrdersDailySinglevendor(ctx context.Context, vendorID string) ([]model.ListOrdersVendor, error)

	SalesReportOrdersMonthlySinglevendor(ctx context.Context, vendorID string) ([]model.ListOrdersVendor, error)
	SalesReportOrdersYearlySinglevendor(ctx context.Context, vendorID string) ([]model.ListOrdersVendor, error)
	SalesReportOrdersCustomSinglevendor(ctx context.Context, startDate, endDate time.Time, vendorID string) ([]model.ListOrdersVendor, error)

	PrintingUserMainOrder(ctx context.Context) ([]model.ListingMainOrders, error)
}

type repository struct {
	sql *sql.DB
}

func NewRepository(sqlDB *sql.DB) Repository {
	return &repository{
		sql: sqlDB,
	}
}
func (r *repository) PrintingUserMainOrder(ctx context.Context) ([]model.ListingMainOrders, error) {
	query := `SELECT mo.uuid, mo.delivered, mo.payment_method, mo.status, mo.payable_amount, 
	                 u.firstname || ' ' || u.lastname AS user, 
	                 COALESCE(a.address1, '') || ' ' || COALESCE(a.address2, '') || ' ' || 
	                 COALESCE(a.address3, '') || ' ' || COALESCE(a.city, '') || ' ' || 
	                 COALESCE(a.state, '') || ' ' || COALESCE(a.pin, '') || ' ' || 
	                 COALESCE(a.country, '') AS user_ad, 
	                 COALESCE(DATE(mo.delivery_date)::text, '') AS delivery_date 
			  FROM orders mo
			  JOIN users u ON mo.user_id = u.id 
			  JOIN address a ON mo.address_id = a.address_id
			  `
	var orders []model.ListingMainOrders

	rows, err := r.sql.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("can't execute query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var order model.ListingMainOrders
		err := rows.Scan(&order.OR_id, &order.Delivery_Stat, &order.D_Type, &order.O_status, &order.Amount,
			&order.User, &order.UserAddress, &order.Delivery_date)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		orders = append(orders, order)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return orders, nil
}

// // list orders
func (r *repository) GetSalesFactByDate(ctx context.Context, filterType string, startDate, endDate time.Time) ([]model.Salesfact, error) {
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
			WHERE  EXTRACT(YEAR FROM oi.created_at) = EXTRACT(YEAR FROM CURRENT_DATE) 
			 AND oi.returned=false AND o.status = 'Completed'
			GROUP BY year`
		args = append(args)
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

			WHERE oi.returned=false AND o.status = 'Completed'
			GROUP BY week`
		args = append(args)
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
        WHERE  EXTRACT(MONTH FROM oi.created_at) = EXTRACT(MONTH FROM CURRENT_DATE) 
		AND EXTRACT(YEAR FROM oi.created_at) = EXTRACT(YEAR FROM CURRENT_DATE)  AND oi.returned=false
		AND o.status = 'Completed'
		GROUP BY   month ORDER BY  month;
        `
		args = append(args)
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
        
       AND oi.returned = false
       AND o.status = 'Completed'
		GROUP BY day;
			
			`
		args = append(args)
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
       
       AND oi.returned = false
       AND o.status = 'Completed';
	  `

		args = append(args, startDate, endDate)
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
func (r *repository) SalesReportOrdersDaily(ctx context.Context) ([]model.ListOrdersAdmin, error) {

	query := `SELECT EXTRACT(DAY FROM oi.created_at) AS checks, p.name, oi.quantity, mo.status, oi.returned, oi.price * 
	oi.quantity AS total_price, 
	oi.product_id AS pid, u.firstname || ' ' || u.lastname AS user, 
	COALESCE(a.address1, '') || ' ' || COALESCE(a.address2, '') || ' ' || COALESCE(a.address3, '') || ' ' || COALESCE(a.city, '') || ' ' || COALESCE(a.state, '') || ' ' || COALESCE(a.pin, '') || ' ' || COALESCE(a.country, '') 
	AS user_ad ,DATE(oi.created_at) AS date,v.name AS vname FROM order_items oi JOIN product_models p ON oi.product_id = p.id JOIN vendor v ON p.vendor_id = v.id JOIN orders mo ON oi.order_id = mo.id JOIN users u ON mo.user_id = u.id JOIN address a ON mo.address_id = a.address_id 
	WHERE  mo.status = 'Completed' AND oi.returned=false ORDER BY checks;
`

	var orders []model.ListOrdersAdmin

	rows, err := r.sql.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error executing query in ListOrdersAdmin: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var order model.ListOrdersAdmin
		err := rows.Scan(&order.ListDate, &order.Name, &order.Unit, &order.Status, &order.Returned, &order.Amount, &order.Pid, &order.User, &order.Add, &order.Date, &order.VName)
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
func (r *repository) SalesReportOrdersWeekly(ctx context.Context) ([]model.ListOrdersAdmin, error) {

	query := `SELECT EXTRACT(WEEK FROM oi.created_at) AS checks, p.name, oi.quantity, mo.status, oi.returned, oi.price * 
	oi.quantity AS total_price, 
	oi.product_id AS pid, u.firstname || ' ' || u.lastname AS user, 
	COALESCE(a.address1, '') || ' ' || COALESCE(a.address2, '') || ' ' || COALESCE(a.address3, '') || ' ' || COALESCE(a.city, '') || ' ' || COALESCE(a.state, '') || ' ' || COALESCE(a.pin, '') || ' ' || COALESCE(a.country, '') 
	AS user_ad ,DATE(oi.created_at) AS date,v.name AS vname FROM order_items oi JOIN product_models p ON oi.product_id = p.id JOIN vendor v ON p.vendor_id = v.id JOIN orders mo ON oi.order_id = mo.id JOIN users u ON mo.user_id = u.id JOIN address a ON mo.address_id = a.address_id 
	WHERE mo.status = 'Completed' AND oi.returned=false ORDER BY checks;
`

	var orders []model.ListOrdersAdmin

	rows, err := r.sql.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error executing query in ListOrdersAdmin: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var order model.ListOrdersAdmin
		err := rows.Scan(&order.ListDate, &order.Name, &order.Unit, &order.Status, &order.Returned, &order.Amount, &order.Pid, &order.User, &order.Add, &order.Date, &order.VName)
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
func (r *repository) SalesReportOrdersMonthly(ctx context.Context) ([]model.ListOrdersAdmin, error) {

	query := `SELECT EXTRACT(MONTH FROM oi.created_at) AS checks, p.name, oi.quantity, mo.status, oi.returned, oi.price * 
	oi.quantity AS total_price, 
	oi.product_id AS pid, u.firstname || ' ' || u.lastname AS user, 
	COALESCE(a.address1, '') || ' ' || COALESCE(a.address2, '') || ' ' || COALESCE(a.address3, '') || ' ' || COALESCE(a.city, '') || ' ' || COALESCE(a.state, '') || ' ' || COALESCE(a.pin, '') || ' ' || COALESCE(a.country, '') 
	AS user_ad, DATE(oi.created_at) AS date,v.name AS vname FROM order_items oi JOIN product_models p ON oi.product_id = p.id JOIN vendor v ON p.vendor_id = v.id JOIN orders mo ON oi.order_id = mo.id JOIN users u ON mo.user_id = u.id JOIN address a ON mo.address_id = a.address_id 
	WHERE  mo.status = 'Completed' AND oi.returned=false ORDER BY checks;
`

	var orders []model.ListOrdersAdmin

	rows, err := r.sql.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error executing query in ListOrdersAdmin: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var order model.ListOrdersAdmin
		err := rows.Scan(&order.ListDate, &order.Name, &order.Unit, &order.Status, &order.Returned, &order.Amount, &order.Pid, &order.User, &order.Add, &order.Date, &order.VName)
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
func (r *repository) SalesReportOrdersCustom(ctx context.Context, startDate, endDate time.Time) ([]model.ListOrdersAdmin, error) {

	query := `SELECT 
    p.name, oi.quantity, mo.status, oi.returned, oi.price, oi.product_id AS pid, 
    DATE(oi.created_at) AS date, u.firstname || ' ' || u.lastname AS user, 
    COALESCE(a.address1, '') || ' ' || COALESCE(a.address2, '') || ' ' || COALESCE(a.address3, '') || ' ' ||
    COALESCE(a.city, '') || ' ' || COALESCE(a.state, '') || ' ' || COALESCE(a.pin, '') || ' ' || COALESCE(a.country, '') AS user_ad,
	v.name AS vname 
    FROM order_items oi 
    JOIN product_models p ON oi.product_id = p.id 
    JOIN vendor v ON p.vendor_id = v.id 
    JOIN orders mo ON oi.order_id = mo.id 
    JOIN users u ON mo.user_id = u.id 
    JOIN address a ON mo.address_id = a.address_id 
    WHERE  mo.status='Completed' AND DATE(oi.created_at) BETWEEN $2 AND $3 ;`

	var orders []model.ListOrdersAdmin

	rows, err := r.sql.QueryContext(ctx, query, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("error executing query in ListOrdersAdmin: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var order model.ListOrdersAdmin
		err := rows.Scan(&order.Name, &order.Unit, &order.Status, &order.Returned, &order.Amount, &order.Pid, &order.Date, &order.User, &order.Add, &order.VName)
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
func (r *repository) SalesReportOrdersYearly(ctx context.Context) ([]model.ListOrdersAdmin, error) {

	query := `SELECT EXTRACT(YEAR FROM oi.created_at) AS checks, p.name, oi.quantity, mo.status, oi.returned, oi.price * 
	oi.quantity AS total_price, 
	oi.product_id AS pid, u.firstname || ' ' || u.lastname AS user, 
	COALESCE(a.address1, '') || ' ' || COALESCE(a.address2, '') || ' ' || COALESCE(a.address3, '') || ' ' || COALESCE(a.city, '') || ' ' || COALESCE(a.state, '') || ' ' || COALESCE(a.pin, '') || ' ' || COALESCE(a.country, '') 
	AS user_ad ,DATE(oi.created_at) AS date,v.name AS vname FROM order_items oi JOIN product_models p ON oi.product_id = p.id JOIN vendor v ON p.vendor_id = v.id JOIN orders mo ON oi.order_id = mo.id JOIN users u ON mo.user_id = u.id JOIN address a ON mo.address_id = a.address_id 
	WHERE  mo.status = 'Completed' AND oi.returned=false ORDER BY checks;
`

	var orders []model.ListOrdersAdmin

	rows, err := r.sql.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error executing query in ListOrdersAdmin: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var order model.ListOrdersAdmin
		err := rows.Scan(&order.ListDate, &order.Name, &order.Unit, &order.Status, &order.Returned, &order.Amount, &order.Pid, &order.User, &order.Add, &order.Date, &order.VName)
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

// /Singlevendor  begining
func (r *repository) GetSalesFactByDateSinglevendor(ctx context.Context, filterType string, startDate, endDate time.Time, vendorID string) ([]model.Salesfact, error) {
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
func (r *repository) SalesReportOrdersCustomSinglevendor(ctx context.Context, startDate, endDate time.Time, vendorID string) ([]model.ListOrdersVendor, error) {

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
func (r *repository) SalesReportOrdersYearlySinglevendor(ctx context.Context, vendorID string) ([]model.ListOrdersVendor, error) {

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
func (r *repository) SalesReportOrdersMonthlySinglevendor(ctx context.Context, vendorID string) ([]model.ListOrdersVendor, error) {

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

func (r *repository) SalesReportOrdersDailySinglevendor(ctx context.Context, vendorID string) ([]model.ListOrdersVendor, error) {

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
func (r *repository) SalesReportOrdersWeeklySinglevendor(ctx context.Context, vendorID string) ([]model.ListOrdersVendor, error) {

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

func (r *repository) ListPendingOrdersSinglevendor(ctx context.Context, id string) ([]model.ListOrdersVendor, error) {
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
func (r *repository) ListReturnedOrdersSinglevendor(ctx context.Context, id string) ([]model.ListOrdersVendor, error) {
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
func (r *repository) ListFailedOrdersSinglevendor(ctx context.Context, id string) ([]model.ListOrdersVendor, error) {
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
func (r *repository) ListCompletedOrdersSinglevendor(ctx context.Context, id string) ([]model.ListOrdersVendor, error) {
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

func (r *repository) ListAllOrdersSinglevendor(ctx context.Context, id string) ([]model.ListOrdersVendor, error) {
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

// ////all orders
func (r *repository) ListPendingOrders(ctx context.Context) ([]model.ListOrdersAdmin, error) {
	query := `SELECT 
    p.name, oi.quantity, mo.status, oi.returned, oi.price, oi.product_id AS pid, 
    DATE(oi.created_at) AS date, u.firstname || ' ' || u.lastname AS user, 
    COALESCE(a.address1, '') || ' ' || COALESCE(a.address2, '') || ' ' || COALESCE(a.address3, '') || ' ' ||
    COALESCE(a.city, '') || ' ' || COALESCE(a.state, '') || ' ' || COALESCE(a.pin, '') || ' ' || COALESCE(a.country, '') AS user_ad,
	v.name AS vname 
    FROM order_items oi 
    JOIN product_models p ON oi.product_id = p.id 
    JOIN vendor v ON p.vendor_id = v.id 
    JOIN orders mo ON oi.order_id = mo.id 
    JOIN users u ON mo.user_id = u.id 
    JOIN address a ON mo.address_id = a.address_id 
    WHERE  mo.status='Pending';`

	var orders []model.ListOrdersAdmin

	rows, err := r.sql.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error executing query in ListOrdersAdmin: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var order model.ListOrdersAdmin
		err := rows.Scan(&order.Name, &order.Unit, &order.Status, &order.Returned, &order.Amount, &order.Pid, &order.Date, &order.User, &order.Add, &order.VName)
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
func (r *repository) ListCompletedOrders(ctx context.Context) ([]model.ListOrdersAdmin, error) {
	query := `SELECT 
    p.name, oi.quantity, mo.status, oi.returned, oi.price, oi.product_id AS pid, 
    DATE(oi.created_at) AS date, u.firstname || ' ' || u.lastname AS user, 
    COALESCE(a.address1, '') || ' ' || COALESCE(a.address2, '') || ' ' || COALESCE(a.address3, '') || ' ' ||
    COALESCE(a.city, '') || ' ' || COALESCE(a.state, '') || ' ' || COALESCE(a.pin, '') || ' ' || COALESCE(a.country, '') AS user_ad,
	v.name AS vname 
    FROM order_items oi 
    JOIN product_models p ON oi.product_id = p.id 
    JOIN vendor v ON p.vendor_id = v.id 
    JOIN orders mo ON oi.order_id = mo.id 
    JOIN users u ON mo.user_id = u.id 
    JOIN address a ON mo.address_id = a.address_id 
    WHERE  mo.status='Completed';`

	var orders []model.ListOrdersAdmin

	rows, err := r.sql.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error executing query in ListOrdersAdmin: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var order model.ListOrdersAdmin
		err := rows.Scan(&order.Name, &order.Unit, &order.Status, &order.Returned, &order.Amount, &order.Pid, &order.Date, &order.User, &order.Add, &order.VName)
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
func (r *repository) ListFailedOrders(ctx context.Context) ([]model.ListOrdersAdmin, error) {
	query := `SELECT 
    p.name, oi.quantity, mo.status, oi.returned, oi.price, oi.product_id AS pid, 
    DATE(oi.created_at) AS date, u.firstname || ' ' || u.lastname AS user, 
    COALESCE(a.address1, '') || ' ' || COALESCE(a.address2, '') || ' ' || COALESCE(a.address3, '') || ' ' ||
    COALESCE(a.city, '') || ' ' || COALESCE(a.state, '') || ' ' || COALESCE(a.pin, '') || ' ' || COALESCE(a.country, '') AS user_ad,
	v.name AS vname 
    FROM order_items oi 
    JOIN product_models p ON oi.product_id = p.id 
    JOIN vendor v ON p.vendor_id = v.id 
    JOIN orders mo ON oi.order_id = mo.id 
    JOIN users u ON mo.user_id = u.id 
    JOIN address a ON mo.address_id = a.address_id 
    WHERE   mo.status='Failed';`

	var orders []model.ListOrdersAdmin

	rows, err := r.sql.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error executing query in ListOrdersAdmin: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var order model.ListOrdersAdmin
		err := rows.Scan(&order.Name, &order.Unit, &order.Status, &order.Returned, &order.Amount, &order.Pid, &order.Date, &order.User, &order.Add, &order.VName)
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
func (r *repository) ListReturnedOrders(ctx context.Context) ([]model.ListOrdersAdmin, error) {
	query := `SELECT 
    p.name, oi.quantity, mo.status, oi.returned, oi.price, oi.product_id AS pid, 
    DATE(oi.created_at) AS date, u.firstname || ' ' || u.lastname AS user, 
    COALESCE(a.address1, '') || ' ' || COALESCE(a.address2, '') || ' ' || COALESCE(a.address3, '') || ' ' ||
    COALESCE(a.city, '') || ' ' || COALESCE(a.state, '') || ' ' || COALESCE(a.pin, '') || ' ' || COALESCE(a.country, '') AS user_ad,
	v.name AS vname 
    FROM order_items oi 
    JOIN product_models p ON oi.product_id = p.id 
    JOIN vendor v ON p.vendor_id = v.id 
    JOIN orders mo ON oi.order_id = mo.id 
    JOIN users u ON mo.user_id = u.id 
    JOIN address a ON mo.address_id = a.address_id 
    WHERE  oi.returned=true;`

	var orders []model.ListOrdersAdmin

	rows, err := r.sql.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error executing query in ListOrdersAdmin: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var order model.ListOrdersAdmin
		err := rows.Scan(&order.Name, &order.Unit, &order.Status, &order.Returned, &order.Amount, &order.Pid, &order.Date, &order.User, &order.Add, &order.VName)
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
func (r *repository) ListAllOrders(ctx context.Context) ([]model.ListOrdersAdmin, error) {
	query := `SELECT 
    p.name, oi.quantity, mo.status, oi.returned, oi.price, oi.product_id AS pid, 
    DATE(oi.created_at) AS date, u.firstname || ' ' || u.lastname AS user, 
    COALESCE(a.address1, '') || ' ' || COALESCE(a.address2, '') || ' ' || COALESCE(a.address3, '') || ' ' ||
    COALESCE(a.city, '') || ' ' || COALESCE(a.state, '') || ' ' || COALESCE(a.pin, '') || ' ' || COALESCE(a.country, '') AS user_ad, 
     oi.id AS oid,v.name AS vname
    FROM order_items oi 
    JOIN product_models p ON oi.product_id = p.id 
    JOIN vendor v ON p.vendor_id = v.id 
    JOIN orders mo ON oi.order_id = mo.id 
    JOIN users u ON mo.user_id = u.id 
    JOIN address a ON mo.address_id = a.address_id 
    `

	var orders []model.ListOrdersAdmin

	rows, err := r.sql.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error executing query in ListOrdersAdmin: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var order model.ListOrdersAdmin
		err := rows.Scan(&order.Name, &order.Unit, &order.Status, &order.Returned, &order.Amount, &order.Pid, &order.Date, &order.User, &order.Add, &order.Oid, &order.VName)
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

// // Singlevendor ending
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
