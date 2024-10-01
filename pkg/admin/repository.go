package admin

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
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
	BestSellingListingProductCategory(ctx context.Context, category string) ([]model.ProductListingUsers, error)
	BestSellingListingProduct(ctx context.Context) ([]model.ProductListingUsers, error)

	BestSellingListingCategory(ctx context.Context) ([]string, error)
	BestSellingListingBrand(ctx context.Context) ([]string, error)

	BestSellingListingProductBrand(ctx context.Context, category string) ([]model.ProductListingUsers, error)
	//BrandListing
	BrandListing(ctx context.Context, category string) ([]model.ProductListingUsers, error)
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
	///single vendor report
	GetSalesFactByDateSinglevendor(ctx context.Context, filterType string, startDate, endDate time.Time, vendorID string) ([]model.Salesfact, error)
	SalesReportOrdersWeeklySinglevendor(ctx context.Context, vendorID string) ([]model.ListOrdersVendor, error)
	SalesReportOrdersDailySinglevendor(ctx context.Context, vendorID string) ([]model.ListOrdersVendor, error)
	SalesReportOrdersMonthlySinglevendor(ctx context.Context, vendorID string) ([]model.ListOrdersVendor, error)
	SalesReportOrdersYearlySinglevendor(ctx context.Context, vendorID string) ([]model.ListOrdersVendor, error)
	SalesReportOrdersCustomSinglevendor(ctx context.Context, startDate, endDate time.Time, vendorID string) ([]model.ListOrdersVendor, error)

	PrintingUserMainOrder(ctx context.Context) ([]model.ListingMainOrders, error)

	GetSingleItem(ctx context.Context, oid string) (model.ListAllOrdersCheck, error)

	IncreaseStock(ctx context.Context, id string, unit int) error

	UpdateOiStatus(ctx context.Context, id, ty string) error

	CreditWallet(ctx context.Context, id string, amt float64) (string, error)

	UpdateWalletTransaction(ctx context.Context, value interface{}) error

	UpdateOrderFromAdminUP(ctx context.Context, id, date, payment string, dStatus bool)

	GetVendorDetails(ctx context.Context, id string) ([]model.GetingVeDetails, error)

	GetOrderForUpdating(ctx context.Context, id string) (map[string]interface{}, error)

	GetcpAmtRefund(ctx context.Context, oid string) (float32, error)
	ChangeCouponRefundStatus(ctx context.Context, id string)
}

type repository struct {
	sql *sql.DB
}

func NewRepository(sqlDB *sql.DB) Repository {
	return &repository{
		sql: sqlDB,
	}
}
func (r *repository) GetOrderForUpdating(ctx context.Context, id string) (map[string]interface{}, error) {
	var delivery bool
	var payment string
	var payment_method string
	var delivery_date string

	query := `  select status,delivered,payment_method,delivery_date from
	             orders WHERE uuid=$1
	
	
	`
	err := r.sql.QueryRowContext(ctx, query, id).Scan(&payment, &delivery, &payment_method, &delivery_date)
	if err != nil {
		return nil, fmt.Errorf("there is error in insertion")
	}
	data := make(map[string]interface{})
	data["status"] = payment
	data["delivered"] = delivery
	data["payment_method"] = payment_method
	data["delivery_date"] = delivery_date

	return data, nil

}
func (r *repository) GetVendorDetails(ctx context.Context, id string) ([]model.GetingVeDetails, error) {

	query := `

	 select name,email,gst from vendor where id=$1
	`
	var vdatas []model.GetingVeDetails
	rows, err := r.sql.QueryContext(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("error executing query in ListOrdersVendor: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var vdata model.GetingVeDetails
		err := rows.Scan(&vdata.Name, &vdata.Email, &vdata.Gst)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		vdatas = append(vdatas, vdata)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return vdatas, nil

}
func (r *repository) UpdateOrderFromAdminUP(ctx context.Context, id, date, payment string, dStatus bool) {
	fmt.Println("updatingggg ---UpdateOrderDate", id, date)
	query := `
	UPDATE orders
	SET delivery_date =$1,status=$2,delivered=$3
	WHERE uuid = $4
	`

	_, err := r.sql.ExecContext(ctx, query, date, payment, dStatus, id)

	if err != nil {
		fmt.Errorf("failed to execute update query: %w", err)
	}

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
	des := values[4]
	fmt.Println("this is the UpdateWalletTransaction in repo!!@@@@@", reflect.TypeOf(Amt), "____", Amt, "!!", id, "##", Type)
	query := `
	INSERT INTO wallet_transactions (
		wallet_id,
		amount,
		transaction_type,
		user_id,
		description,
		created_at
		
	) VALUES (
		$1, $2, $3,$4,$5, CURRENT_TIMESTAMP
	) RETURNING id;
`
	var tid string
	fmt.Println("this is the id ", tid, "user_id,", usid)

	err := r.sql.QueryRowContext(ctx, query, id, Amt, Type, usid, des).Scan(&tid)
	if err != nil {
		return fmt.Errorf("there is error in insertion")
	}

	return nil

}
func (r *repository) ChangeCouponRefundStatus(ctx context.Context, id string) {
	fmt.Println("changing in the order status QQQQ", id)
	query := `
	UPDATE orders
	SET cp_amount_refund_status =true
	WHERE uuid = $1
	
`

	_, err := r.sql.ExecContext(ctx, query, id)

	if err != nil {
		fmt.Errorf("failed to execute update query: %w", err)
	}

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
func (r *repository) GetcpAmtRefund(ctx context.Context, oid string) (float32, error) {

	query := `
		SELECT COALESCE(coupon_amount, 0.0) 
		FROM orders 
		WHERE uuid = $1 AND cp_amount_refund_status=false;
	`

	var amount float32
	err := r.sql.QueryRowContext(ctx, query, oid).Scan(&amount)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0.0, nil
		}
		return 0.00, err
	}
	fmt.Println("amttt   !!", amount)
	return amount, nil
}
func (r *repository) GetSingleItem(ctx context.Context, oid string) (model.ListAllOrdersCheck, error) {
	var order model.ListAllOrdersCheck

	query := `SELECT p.name,  oi.quantity,   mo.status, oi.returned, 
    oi.price,oi.product_id AS pid,DATE(oi.created_at) AS date,mo.user_id ,v.id AS vid,u.email AS usmail,mo.uuid AS mid
     FROM order_items oi 
    JOIN  product_models p ON oi.product_id = p.id 
    JOIN  orders mo ON oi.order_id = mo.id 
    JOIN  vendor v ON p.vendor_id = v.id 
	JOIN  users  u ON oi.user_id=u.id
    WHERE  oi.id = $1;
  
`
	err := r.sql.QueryRowContext(ctx, query, oid).Scan(&order.Name, &order.Unit, &order.Status, &order.Returned, &order.Amount, &order.Pid, &order.Date, &order.Usid, &order.Vid, &order.Usmail, &order.Moid)
	if err != nil {
		return model.ListAllOrdersCheck{}, fmt.Errorf("error in exequting query in  GetSingleItem")
	}
	return order, nil
}
func (r *repository) PrintingUserMainOrder(ctx context.Context) ([]model.ListingMainOrders, error) {
	query := `SELECT mo.uuid, mo.delivered, mo.payment_method, mo.status, mo.payable_amount, 
	                 u.firstname || ' ' || u.lastname AS user, 
	                 COALESCE(a.address1, '') || ' ' || COALESCE(a.address2, '') || ' ' || 
	                 COALESCE(a.address3, '') || ' ' || COALESCE(a.city, '') || ' ' || 
	                 COALESCE(a.state, '') || ' ' || COALESCE(a.pin, '') || ' ' || 
	                 COALESCE(a.country, '') AS user_ad, 
	                 COALESCE(DATE(mo.delivery_date)::text, '') AS delivery_date 
					 , mo.discount,mo.coupon_amount AS cmt,COALESCE(c.code , ''),mo.wallet_money AS wmt 
			  FROM orders mo
			  JOIN users u ON mo.user_id = u.id 
			  LEFT JOIN coupon c ON mo.cid =c.id
			  JOIN address a ON mo.address_id = a.address_id ORDER BY mo.id DESC
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
			&order.User, &order.UserAddress, &order.Delivery_date, &order.Discount, &order.Cmt, &order.Code, &order.Wmt)
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
			 AND o.returned=false AND o.status = 'Completed'
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

			WHERE o.returned=false AND o.status = 'Completed'
			AND EXTRACT(WEEK FROM oi.created_at) = EXTRACT(WEEK FROM CURRENT_DATE) 
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
	COALESCE(a.address1, '') || ' ' || COALESCE(a.address2, '') || ' ' || COALESCE(a.address3, '') || ' ' || COALESCE(a.city, '') 
	|| ' ' || COALESCE(a.state, '') || ' ' || COALESCE(a.pin, '') || ' ' || COALESCE(a.country, '') 
	AS user_ad ,DATE(oi.created_at) AS date,v.name AS vendorname ,mo.uuid AS oid,(oi.discount * oi.quantity) AS discount,
	mo.coupon_amount AS cmt,COALESCE(c.code , '') ,mo.wallet_money AS wmt 
	FROM order_items oi JOIN product_models p ON oi.product_id = p.id JOIN vendor v ON p.vendor_id = v.id
	 JOIN orders mo ON oi.order_id = mo.id JOIN users u ON mo.user_id = u.id JOIN address a ON mo.address_id = a.address_id 
	LEFT JOIN coupon c ON mo.cid =c.id
	WHERE DATE(oi.created_at) = DATE(CURRENT_DATE) AND  mo.status = 'Completed' AND oi.returned=false ORDER BY checks;
`

	var orders []model.ListOrdersAdmin

	rows, err := r.sql.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error executing query in ListOrdersAdmin: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var order model.ListOrdersAdmin
		err := rows.Scan(&order.ListDate, &order.Name, &order.Unit, &order.Status, &order.Returned, &order.Amount, &order.Pid, &order.User, &order.Add, &order.Date, &order.VName, &order.Oid, &order.Discount, &order.CouponAmt, &order.CouponCode, &order.WalletAmt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}
	fmt.Println("this is the result in the fmtttttt", orders)
	return orders, nil
}
func (r *repository) SalesReportOrdersWeekly(ctx context.Context) ([]model.ListOrdersAdmin, error) {

	query := `SELECT EXTRACT(WEEK FROM oi.created_at) AS checks, p.name, oi.quantity, mo.status, oi.returned, oi.price * 
	oi.quantity AS total_price, 
	oi.product_id AS pid, u.firstname || ' ' || u.lastname AS user, 
	COALESCE(a.address1, '') || ' ' || COALESCE(a.address2, '') || ' ' || COALESCE(a.address3, '') || ' ' || COALESCE(a.city, '') || ' ' || COALESCE(a.state, '') || ' ' || COALESCE(a.pin, '') || ' ' || COALESCE(a.country, '') 
	AS user_ad ,DATE(oi.created_at) AS date,v.name AS vendorname,mo.uuid AS oid, (oi.discount * oi.quantity) AS discount,mo.coupon_amount AS cmt,COALESCE(c.code , '') ,mo.wallet_money AS wmt 
	 FROM order_items oi JOIN product_models p ON oi.product_id = p.id JOIN vendor v ON p.vendor_id = v.id JOIN orders mo ON oi.order_id = mo.id JOIN users u ON mo.user_id = u.id JOIN address a ON mo.address_id = a.address_id 
	LEFT  JOIN coupon c ON mo.cid =c.id
	 WHERE mo.status = 'Completed' AND oi.returned=false AND EXTRACT(WEEK FROM oi.created_at) = EXTRACT(WEEK FROM CURRENT_DATE)  ORDER BY checks;
`

	var orders []model.ListOrdersAdmin

	rows, err := r.sql.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error executing query in ListOrdersAdmin: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var order model.ListOrdersAdmin
		err := rows.Scan(&order.ListDate, &order.Name, &order.Unit, &order.Status, &order.Returned, &order.Amount, &order.Pid, &order.User, &order.Add, &order.Date, &order.VName, &order.Oid, &order.Discount, &order.CouponAmt, &order.CouponCode, &order.WalletAmt)
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
	AS user_ad, DATE(oi.created_at) AS date,v.name AS vendorname
	
	,mo.uuid AS oid  , (oi.discount * oi.quantity) AS discount,mo.coupon_amount AS cmt,COALESCE(c.code , ''),mo.wallet_money AS wmt 

	FROM order_items oi JOIN product_models p ON oi.product_id = p.id JOIN vendor v ON p.vendor_id = v.id JOIN orders mo ON oi.order_id = mo.id JOIN users u ON mo.user_id = u.id JOIN address a ON mo.address_id = a.address_id 
	LEFT JOIN coupon c ON mo.cid =c.id
	WHERE EXTRACT(MONTH FROM oi.created_at) = EXTRACT(MONTH FROM CURRENT_DATE)  AND mo.status = 'Completed' AND oi.returned=false ORDER BY checks;
`

	var orders []model.ListOrdersAdmin

	rows, err := r.sql.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error executing query in ListOrdersAdmin: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var order model.ListOrdersAdmin
		err := rows.Scan(&order.ListDate, &order.Name, &order.Unit, &order.Status, &order.Returned, &order.Amount, &order.Pid, &order.User, &order.Add, &order.Date, &order.VName, &order.Oid, &order.Discount, &order.CouponAmt, &order.CouponCode, &order.WalletAmt)
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
	v.name AS vendorname,mo.uuid AS oid , (oi.discount * oi.quantity) AS discount,mo.coupon_amount AS cmt,COALESCE(c.code , ''),mo.wallet_money AS wmt 
    FROM order_items oi 
    JOIN product_models p ON oi.product_id = p.id 
    JOIN vendor v ON p.vendor_id = v.id 
    JOIN orders mo ON oi.order_id = mo.id 
    JOIN users u ON mo.user_id = u.id 
    JOIN address a ON mo.address_id = a.address_id 
    LEFT JOIN coupon c ON mo.cid =c.id
	WHERE  mo.status='Completed' AND DATE(oi.created_at) BETWEEN $1 AND $2 ;`

	var orders []model.ListOrdersAdmin

	rows, err := r.sql.QueryContext(ctx, query, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("error executing query in ListOrdersAdmin: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var order model.ListOrdersAdmin
		err := rows.Scan(&order.Name, &order.Unit, &order.Status, &order.Returned, &order.Amount, &order.Pid, &order.Date, &order.User, &order.Add, &order.VName, &order.Oid, &order.Discount, &order.CouponAmt, &order.CouponCode, &order.WalletAmt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}
	fmt.Println("check check #############&&&&&&&&&############", orders)
	return orders, nil
}
func (r *repository) SalesReportOrdersYearly(ctx context.Context) ([]model.ListOrdersAdmin, error) {

	query := `SELECT EXTRACT(YEAR FROM oi.created_at) AS checks, p.name, oi.quantity, mo.status, oi.returned, oi.price * 
	oi.quantity AS total_price, 
	oi.product_id AS pid, u.firstname || ' ' || u.lastname AS user, 
	COALESCE(a.address1, '') || ' ' || COALESCE(a.address2, '') || ' ' || COALESCE(a.address3, '') || ' ' || COALESCE(a.city, '') || ' ' || COALESCE(a.state, '') || ' ' || COALESCE(a.pin, '') || ' ' || COALESCE(a.country, '') 
	AS user_ad ,DATE(oi.created_at) AS date,v.name AS vendorname,
	 mo.uuid AS oid , (oi.discount * oi.quantity) AS discount,mo.coupon_amount AS cmt,COALESCE(c.code , ''),mo.wallet_money AS wmt 
	FROM order_items oi JOIN product_models p ON oi.product_id = p.id JOIN vendor v ON p.vendor_id = v.id JOIN orders mo ON oi.order_id = mo.id JOIN users u ON mo.user_id = u.id JOIN address a ON mo.address_id = a.address_id 
	LEFT JOIN coupon c ON mo.cid =c.id
	WHERE  EXTRACT(YEAR FROM oi.created_at) = EXTRACT(YEAR FROM CURRENT_DATE)  AND mo.status = 'Completed' AND oi.returned=false ORDER BY checks;
`

	var orders []model.ListOrdersAdmin

	rows, err := r.sql.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error executing query in ListOrdersAdmin: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var order model.ListOrdersAdmin
		err := rows.Scan(&order.ListDate, &order.Name, &order.Unit, &order.Status, &order.Returned, &order.Amount, &order.Pid, &order.User, &order.Add, &order.Date, &order.VName, &order.Oid, &order.Discount, &order.CouponAmt, &order.CouponCode, &order.WalletAmt)
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
			AND EXTRACT(WEEK FROM oi.created_at) = EXTRACT(WEEK FROM CURRENT_DATE) 
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
    ,mo.uuid AS oid , (oi.discount * oi.quantity) AS discount,mo.coupon_amount AS cmt,COALESCE(c.code , ''),mo.wallet_money AS wmt 
	FROM order_items oi 
    JOIN product_models p ON oi.product_id = p.id 
    JOIN vendor v ON p.vendor_id = v.id 
    JOIN orders mo ON oi.order_id = mo.id 
    JOIN users u ON mo.user_id = u.id 
    JOIN address a ON mo.address_id = a.address_id 
	LEFT JOIN coupon c ON mo.cid =c.id
    WHERE v.id = $1 AND  mo.status='Completed' AND DATE(oi.created_at) BETWEEN $2 AND $3 ;`

	var orders []model.ListOrdersVendor

	rows, err := r.sql.QueryContext(ctx, query, vendorID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("error executing query in ListOrdersVendor: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var order model.ListOrdersVendor
		err := rows.Scan(&order.Name, &order.Unit, &order.Status, &order.Returned, &order.Amount, &order.Pid, &order.Date, &order.User, &order.Add, &order.Oid, &order.Discount, &order.CouponAmt, &order.CouponCode, &order.WalletAmt)
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
	AS user_ad ,DATE(oi.created_at) AS date ,mo.uuid AS oid , (oi.discount * oi.quantity) AS discount,mo.coupon_amount AS cmt,COALESCE(c.code , ''),mo.wallet_money AS wmt 
	 FROM order_items oi JOIN product_models p ON oi.product_id = p.id JOIN vendor v ON p.vendor_id = v.id JOIN orders mo ON oi.order_id = mo.id JOIN users u ON mo.user_id = u.id JOIN address a ON mo.address_id = a.address_id 
	LEFT JOIN coupon c ON mo.cid =c.id
	WHERE EXTRACT(YEAR FROM oi.created_at) = EXTRACT(YEAR FROM CURRENT_DATE)  AND v.id = $1 AND mo.status = 'Completed' AND oi.returned=false ORDER BY checks;
`

	var orders []model.ListOrdersVendor

	rows, err := r.sql.QueryContext(ctx, query, vendorID)
	if err != nil {
		return nil, fmt.Errorf("error executing query in ListOrdersVendor: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var order model.ListOrdersVendor
		err := rows.Scan(&order.ListDate, &order.Name, &order.Unit, &order.Status, &order.Returned, &order.Amount, &order.Pid, &order.User, &order.Add, &order.Date, &order.Oid, &order.Discount, &order.CouponAmt, &order.CouponCode, &order.WalletAmt)
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
	AS user_ad, DATE(oi.created_at) AS date,mo.uuid AS oid  , (oi.discount * oi.quantity) AS discount,mo.coupon_amount AS cmt,COALESCE(c.code , ''),mo.wallet_money AS wmt 
	 FROM order_items oi JOIN product_models p ON oi.product_id = p.id JOIN vendor v ON p.vendor_id = v.id JOIN orders mo ON oi.order_id = mo.id JOIN users u ON mo.user_id = u.id JOIN address a ON mo.address_id = a.address_id 
	LEFT JOIN coupon c ON mo.cid =c.id
	WHERE EXTRACT(MONTH FROM oi.created_at) = EXTRACT(MONTH FROM CURRENT_DATE)  AND v.id = $1 AND mo.status = 'Completed' AND oi.returned=false ORDER BY checks;
`

	var orders []model.ListOrdersVendor

	rows, err := r.sql.QueryContext(ctx, query, vendorID)
	if err != nil {
		return nil, fmt.Errorf("error executing query in ListOrdersVendor: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var order model.ListOrdersVendor
		err := rows.Scan(&order.ListDate, &order.Name, &order.Unit, &order.Status, &order.Returned, &order.Amount, &order.Pid, &order.User, &order.Add, &order.Date, &order.Oid, &order.Discount, &order.CouponAmt, &order.CouponCode, &order.WalletAmt)
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
	AS user_ad ,DATE(oi.created_at) AS date,mo.uuid AS oid,(oi.discount * oi.quantity) AS discount,mo.coupon_amount AS cmt,COALESCE(c.code , '') ,mo.wallet_money AS wmt  FROM order_items oi JOIN product_models p ON oi.product_id = p.id JOIN vendor v ON p.vendor_id = v.id JOIN orders mo ON oi.order_id = mo.id JOIN users u ON mo.user_id = u.id JOIN address a ON mo.address_id = a.address_id 
	LEFT JOIN coupon c ON mo.cid =c.id
	WHERE DATE(oi.created_at) = DATE(CURRENT_DATE) AND v.id = $1 AND mo.status = 'Completed' AND oi.returned=false ORDER BY checks;
`

	var orders []model.ListOrdersVendor

	rows, err := r.sql.QueryContext(ctx, query, vendorID)
	if err != nil {
		return nil, fmt.Errorf("error executing query in ListOrdersVendor: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var order model.ListOrdersVendor
		err := rows.Scan(&order.ListDate, &order.Name, &order.Unit, &order.Status, &order.Returned, &order.Amount, &order.Pid, &order.User, &order.Add, &order.Date, &order.Oid, &order.Discount, &order.CouponAmt, &order.CouponCode, &order.WalletAmt)
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
	AS user_ad ,DATE(oi.created_at) AS date,mo.uuid AS oid, (oi.discount * oi.quantity) AS discount,mo.coupon_amount AS cmt,COALESCE(c.code , '') ,mo.wallet_money AS wmt  
	FROM order_items oi JOIN product_models p ON oi.product_id = p.id JOIN vendor v ON p.vendor_id = v.id JOIN orders mo ON oi.order_id = mo.id JOIN users u ON mo.user_id = u.id JOIN address a ON mo.address_id = a.address_id LEFT  JOIN coupon c ON mo.cid =c.id
	WHERE v.id = $1 AND mo.status = 'Completed' AND oi.returned=false AND EXTRACT(WEEK FROM oi.created_at) = EXTRACT(WEEK FROM CURRENT_DATE)  ORDER BY checks;
`

	var orders []model.ListOrdersVendor

	rows, err := r.sql.QueryContext(ctx, query, vendorID)
	if err != nil {
		return nil, fmt.Errorf("error executing query in ListOrdersVendor: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var order model.ListOrdersVendor
		err := rows.Scan(&order.ListDate, &order.Name, &order.Unit, &order.Status, &order.Returned, &order.Amount, &order.Pid, &order.User, &order.Add, &order.Date, &order.Oid, &order.Discount, &order.CouponAmt, &order.CouponCode, &order.WalletAmt)
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
    COALESCE(a.city, '') || ' ' || COALESCE(a.state, '') || ' ' || COALESCE(a.pin, '') || ' ' || COALESCE(a.country, '') AS user_ad ,
    mo.uuid AS oid  , (oi.discount * oi.quantity) AS discount,mo.coupon_amount AS cmt,COALESCE(c.code , ''),mo.wallet_money AS wmt 
    FROM order_items oi 
    JOIN product_models p ON oi.product_id = p.id 
    JOIN vendor v ON p.vendor_id = v.id 
    JOIN orders mo ON oi.order_id = mo.id 
    JOIN users u ON mo.user_id = u.id 
    JOIN address a ON mo.address_id = a.address_id 
	LEFT JOIN coupon c ON mo.cid =c.id
    WHERE v.id = $1 AND  mo.status='Pending';`

	var orders []model.ListOrdersVendor

	rows, err := r.sql.QueryContext(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("error executing query in ListOrdersVendor: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var order model.ListOrdersVendor
		err := rows.Scan(&order.Name, &order.Unit, &order.Status, &order.Returned, &order.Amount, &order.Pid, &order.Date, &order.User, &order.Add, &order.Oid, &order.Discount, &order.CouponAmt, &order.CouponCode, &order.WalletAmt)
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
    COALESCE(a.city, '') || ' ' || COALESCE(a.state, '') || ' ' || COALESCE(a.pin, '') || ' ' || COALESCE(a.country, '') AS user_ad ,
    mo.uuid AS oid  , oi.discount,mo.coupon_amount AS cmt,COALESCE(c.code , ''),mo.wallet_money AS wmt 
    FROM order_items oi 
    JOIN product_models p ON oi.product_id = p.id 
    JOIN vendor v ON p.vendor_id = v.id 
    JOIN orders mo ON oi.order_id = mo.id 
    JOIN users u ON mo.user_id = u.id 
    JOIN address a ON mo.address_id = a.address_id 
	LEFT JOIN coupon c ON mo.cid =c.id
    WHERE v.id = $1 AND oi.returned=true;`

	var orders []model.ListOrdersVendor

	rows, err := r.sql.QueryContext(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("error executing query in ListOrdersVendor: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var order model.ListOrdersVendor
		err := rows.Scan(&order.Name, &order.Unit, &order.Status, &order.Returned, &order.Amount, &order.Pid, &order.Date, &order.User, &order.Add, &order.Oid, &order.Discount, &order.CouponAmt, &order.CouponCode, &order.WalletAmt)
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
    COALESCE(a.city, '') || ' ' || COALESCE(a.state, '') || ' ' || COALESCE(a.pin, '') || ' ' || COALESCE(a.country, '') AS user_ad ,
    mo.uuid AS oid  , oi.discount,mo.coupon_amount AS cmt,COALESCE(c.code , ''),mo.wallet_money AS wmt 
    FROM order_items oi 
    JOIN product_models p ON oi.product_id = p.id 
    JOIN vendor v ON p.vendor_id = v.id 
    JOIN orders mo ON oi.order_id = mo.id 
    JOIN users u ON mo.user_id = u.id 
    JOIN address a ON mo.address_id = a.address_id 
	LEFT JOIN coupon c ON mo.cid =c.id
    WHERE v.id = $1 AND  mo.status='Failed';`

	var orders []model.ListOrdersVendor

	rows, err := r.sql.QueryContext(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("error executing query in ListOrdersVendor: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var order model.ListOrdersVendor
		err := rows.Scan(&order.Name, &order.Unit, &order.Status, &order.Returned, &order.Amount, &order.Pid, &order.Date, &order.User, &order.Add, &order.Oid, &order.Discount, &order.CouponAmt, &order.CouponCode, &order.WalletAmt)
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
    COALESCE(a.city, '') || ' ' || COALESCE(a.state, '') || ' ' || COALESCE(a.pin, '') || ' ' || COALESCE(a.country, '') AS user_ad ,
    mo.uuid AS oid  , oi.discount,mo.coupon_amount AS cmt,COALESCE(c.code , ''),mo.wallet_money AS wmt 
    FROM order_items oi 
    JOIN product_models p ON oi.product_id = p.id 
    JOIN vendor v ON p.vendor_id = v.id 
    JOIN orders mo ON oi.order_id = mo.id 
    JOIN users u ON mo.user_id = u.id 
    JOIN address a ON mo.address_id = a.address_id 
	LEFT JOIN coupon c ON mo.cid =c.id
    WHERE v.id = $1 AND  mo.status='Completed';`

	var orders []model.ListOrdersVendor

	rows, err := r.sql.QueryContext(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("error executing query in ListOrdersVendor: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var order model.ListOrdersVendor
		err := rows.Scan(&order.Name, &order.Unit, &order.Status, &order.Returned, &order.Amount, &order.Pid, &order.Date, &order.User, &order.Add, &order.Oid, &order.Discount, &order.CouponAmt, &order.CouponCode, &order.WalletAmt)
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
     mo.uuid AS oid  , oi.discount,mo.coupon_amount AS cmt,COALESCE(c.code , ''),mo.wallet_money AS wmt 
    FROM order_items oi 
    JOIN product_models p ON oi.product_id = p.id 
    JOIN vendor v ON p.vendor_id = v.id 
    JOIN orders mo ON oi.order_id = mo.id 
    JOIN users u ON mo.user_id = u.id 
    JOIN address a ON mo.address_id = a.address_id 
	LEFT JOIN coupon c ON mo.cid =c.id
    WHERE v.id = $1;`

	var orders []model.ListOrdersVendor

	rows, err := r.sql.QueryContext(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("error executing query in ListOrdersVendor: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var order model.ListOrdersVendor
		err := rows.Scan(&order.Name, &order.Unit, &order.Status, &order.Returned, &order.Amount, &order.Pid, &order.Date, &order.User, &order.Add, &order.Oid, &order.Discount, &order.CouponAmt, &order.CouponCode, &order.WalletAmt)
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
	v.name AS vendorname 
	,mo.uuid AS oid  , oi.discount,mo.coupon_amount AS cmt,COALESCE(c.code , ''),mo.wallet_money AS wmt 
    FROM order_items oi 
    JOIN product_models p ON oi.product_id = p.id 
    JOIN vendor v ON p.vendor_id = v.id 
    JOIN orders mo ON oi.order_id = mo.id 
    JOIN users u ON mo.user_id = u.id 
    JOIN address a ON mo.address_id = a.address_id 
	LEFT JOIN coupon c ON mo.cid =c.id
    WHERE  mo.status='Pending';`

	var orders []model.ListOrdersAdmin

	rows, err := r.sql.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error executing query in ListOrdersAdmin: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var order model.ListOrdersAdmin
		err := rows.Scan(&order.Name, &order.Unit, &order.Status, &order.Returned, &order.Amount, &order.Pid, &order.Date, &order.User, &order.Add, &order.VName, &order.Oid, &order.Discount, &order.CouponAmt, &order.CouponCode, &order.WalletAmt)
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
	v.name AS vendorname 
	,mo.uuid AS oid  , oi.discount,mo.coupon_amount AS cmt,COALESCE(c.code , ''),mo.wallet_money AS wmt 
    FROM order_items oi 
    JOIN product_models p ON oi.product_id = p.id 
    JOIN vendor v ON p.vendor_id = v.id 
    JOIN orders mo ON oi.order_id = mo.id 
    JOIN users u ON mo.user_id = u.id 
    JOIN address a ON mo.address_id = a.address_id 
	LEFT JOIN coupon c ON mo.cid =c.id
    WHERE  mo.status='Completed';`

	var orders []model.ListOrdersAdmin

	rows, err := r.sql.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error executing query in ListOrdersAdmin: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var order model.ListOrdersAdmin
		err := rows.Scan(&order.Name, &order.Unit, &order.Status, &order.Returned, &order.Amount, &order.Pid, &order.Date, &order.User, &order.Add, &order.VName, &order.Oid, &order.Discount, &order.CouponAmt, &order.CouponCode, &order.WalletAmt)
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
	v.name AS vendorname 
		,mo.uuid AS oid  , oi.discount,mo.coupon_amount AS cmt,COALESCE(c.code , ''),mo.wallet_money AS wmt 
    FROM order_items oi 
    JOIN product_models p ON oi.product_id = p.id 
    JOIN vendor v ON p.vendor_id = v.id 
    JOIN orders mo ON oi.order_id = mo.id 
    JOIN users u ON mo.user_id = u.id 
    JOIN address a ON mo.address_id = a.address_id 
	LEFT JOIN coupon c ON mo.cid =c.id
    WHERE   mo.status='Failed';`

	var orders []model.ListOrdersAdmin

	rows, err := r.sql.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error executing query in ListOrdersAdmin: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var order model.ListOrdersAdmin
		err := rows.Scan(&order.Name, &order.Unit, &order.Status, &order.Returned, &order.Amount, &order.Pid, &order.Date, &order.User, &order.Add, &order.VName, &order.Oid, &order.Discount, &order.CouponAmt, &order.CouponCode, &order.WalletAmt)
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
	v.name AS vendorname 
	,mo.uuid AS oid  , oi.discount,mo.coupon_amount AS cmt,COALESCE(c.code , ''),mo.wallet_money AS wmt 
    FROM order_items oi 
    JOIN product_models p ON oi.product_id = p.id 
    JOIN vendor v ON p.vendor_id = v.id 
    JOIN orders mo ON oi.order_id = mo.id 
    JOIN users u ON mo.user_id = u.id 
    JOIN address a ON mo.address_id = a.address_id 
	LEFT JOIN coupon c ON mo.cid =c.id
    WHERE  oi.returned=true;`

	var orders []model.ListOrdersAdmin

	rows, err := r.sql.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error executing query in ListOrdersAdmin: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var order model.ListOrdersAdmin
		err := rows.Scan(&order.Name, &order.Unit, &order.Status, &order.Returned, &order.Amount, &order.Pid, &order.Date, &order.User, &order.Add, &order.VName, &order.Oid, &order.Discount, &order.CouponAmt, &order.CouponCode, &order.WalletAmt)
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
     oi.id AS oid,v.name AS vendorname,mo.uuid AS oid  , oi.discount,mo.coupon_amount AS cmt,COALESCE(c.code , ''),mo.wallet_money AS wmt 
    FROM order_items oi 
    JOIN product_models p ON oi.product_id = p.id 
    JOIN vendor v ON p.vendor_id = v.id 
    JOIN orders mo ON oi.order_id = mo.id 
    JOIN users u ON mo.user_id = u.id 
	LEFT JOIN coupon c ON mo.cid =c.id
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
		err := rows.Scan(&order.Name, &order.Unit, &order.Status, &order.Returned, &order.Amount, &order.Pid, &order.Date, &order.User, &order.Add, &order.Oid, &order.VName, &order.Oid, &order.Discount, &order.CouponAmt, &order.CouponCode, &order.WalletAmt)
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

// // Singlevendor ending  BrandListing
func (r *repository) BrandListing(ctx context.Context, category string) ([]model.ProductListingUsers, error) {
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
			vendor ON product_models.vendor_id = vendor.id WHERE product_models.brand ILIKE '%' || $1 || '%';`

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
		} ///  https://adiecom.gitfunswokhu.in
		product.Pdetail = "https://adiecom.gitfunswokhu.in/user/listingSingleProduct/" + product.Pid
		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return products, nil
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
			vendor ON product_models.vendor_id = vendor.id WHERE product_models.category ILIKE '%' || $1 || '%';`

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
		product.Pdetail = "https://adiecom.gitfunswokhu.in/user/listingSingleProduct/" + product.Pid
		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return products, nil
}
func (r *repository) BestSellingListingProductBrand(ctx context.Context, category string) ([]model.ProductListingUsers, error) {
	query := `
  		
	
	WITH best_selling_products AS (
    SELECT 
        pm.name,
        pm.category,
        pm.units,
        pm.tax,
        pm.amount,
        pm.status,
        pm.discount,
        COALESCE(SUM(oi.quantity), 0) AS total_sold
    FROM 
        product_models pm
    LEFT JOIN 
        order_items oi ON pm.id = oi.product_id
    LEFT JOIN 
        vendor v ON pm.vendor_id = v.id
    WHERE 
      
	pm.brand ILIKE '%' || $1 || '%'  -- Partial match with ILIKE
    GROUP BY 
        pm.name, pm.category, pm.units, pm.tax, pm.amount, pm.status, pm.discount
		HAVING 
        SUM(oi.quantity) > 0
),
ranked_products AS (
    SELECT
        *,
        ROW_NUMBER() OVER (ORDER BY total_sold DESC, name) AS rank
    FROM 
        best_selling_products
)
SELECT 
    name,
    category,
    units,
    tax,
    amount,
    status,
    discount,
    total_sold
    
FROM 
    ranked_products
WHERE 
    rank <= 10
ORDER BY
    rank;
			 
			 
			 `

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
		product.Pdetail = "https://adiecom.gitfunswokhu.in/admin/listingSingleProduct/" + product.Pid
		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return products, nil
}
func (r *repository) BestSellingListingProductCategory(ctx context.Context, category string) ([]model.ProductListingUsers, error) {
	query := `
  		
	
	WITH best_selling_products AS (
    SELECT 
        pm.name,
        pm.category,
        pm.units,
        pm.tax,
        pm.amount,
        pm.status,
        pm.discount,
        COALESCE(SUM(oi.quantity), 0) AS total_sold
    FROM 
        product_models pm
    LEFT JOIN 
        order_items oi ON pm.id = oi.product_id
    LEFT JOIN 
        vendor v ON pm.vendor_id = v.id
    WHERE 
      
	pm.category ILIKE '%' || $1 || '%'  -- Partial match with ILIKE
    GROUP BY 
        pm.name, pm.category, pm.units, pm.tax, pm.amount, pm.status, pm.discount
		HAVING 
        SUM(oi.quantity) > 0
),
ranked_products AS (
    SELECT
        *,
        ROW_NUMBER() OVER (ORDER BY total_sold DESC, name) AS rank
    FROM 
        best_selling_products
)
SELECT 
    name,
    category,
    units,
    tax,
    amount,
    status,
    discount,
    total_sold
    
FROM 
    ranked_products
WHERE 
    rank <= 10
ORDER BY
    rank;
			 
			 
			 `

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
		product.Pdetail = "https://adiecom.gitfunswokhu.in/user/listingSingleProduct/" + product.Pid
		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return products, nil
}
func (r *repository) BestSellingListingProduct(ctx context.Context) ([]model.ProductListingUsers, error) {
	query := `


WITH best_selling_products AS (
    SELECT 
        pm.name,
        pm.category,
        pm.units,
        pm.tax,
        pm.amount,
        pm.status,
        pm.discount,
        COALESCE(SUM(oi.quantity), 0) AS total_sold
    FROM 
        product_models pm
    LEFT JOIN 
        order_items oi ON pm.id = oi.product_id
    LEFT JOIN 
        vendor v ON pm.vendor_id = v.id

    GROUP BY 
        pm.name, pm.category, pm.units, pm.tax, pm.amount, pm.status, pm.discount
		HAVING 
        SUM(oi.quantity) > 0
),
ranked_products AS (
    SELECT
        *,
        ROW_NUMBER() OVER (ORDER BY total_sold DESC, name) AS rank
    FROM 
        best_selling_products
)
SELECT 
    name,
    category,
    units,
    tax,
    amount,
    status,
    discount,
    total_sold
FROM 
    ranked_products
WHERE 
    rank <= 10
ORDER BY
    rank;
			 
 `

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
		product.Pdetail = "https://adiecom.gitfunswokhu.in/user/listingSingleProduct/" + product.Pid
		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return products, nil
}

func (r *repository) BestSellingListingCategory(ctx context.Context) ([]string, error) {
	query := `
    WITH best_selling_categories AS (
        SELECT 
            pm.category,
            COALESCE(SUM(oi.quantity), 0) AS total_sold
        FROM 
            product_models pm
        LEFT JOIN 
            order_items oi ON pm.id = oi.product_id
        GROUP BY 
            pm.category
    ),
    ranked_categories AS (
        SELECT
            category,
            total_sold,
            ROW_NUMBER() OVER (ORDER BY total_sold DESC) AS rank
        FROM 
            best_selling_categories
    )
    SELECT 
        category
    FROM 
        ranked_categories
    WHERE 
        rank <= 10
    ORDER BY
        rank;
    `

	rows, err := r.sql.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute select query: %w", err)
	}
	defer rows.Close()

	var categories []string
	for rows.Next() {
		var category string
		err := rows.Scan(&category)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return categories, nil
}
func (r *repository) BestSellingListingBrand(ctx context.Context) ([]string, error) {
	query := `
    WITH best_selling_categories AS (
        SELECT 
            pm.brand,
            COALESCE(SUM(oi.quantity), 0) AS total_sold
        FROM 
            product_models pm
        LEFT JOIN 
            order_items oi ON pm.id = oi.product_id
        GROUP BY 
            pm.brand
    ),
    ranked_categories AS (
        SELECT
            brand,
            total_sold,
            ROW_NUMBER() OVER (ORDER BY total_sold DESC) AS rank
        FROM 
            best_selling_categories
    )
    SELECT 
        brand
    FROM 
        ranked_categories
    WHERE 
        rank <= 10
    ORDER BY
        rank;
    `

	rows, err := r.sql.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute select query: %w", err)
	}
	defer rows.Close()

	var categories []string
	for rows.Next() {
		var category string
		err := rows.Scan(&category)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return categories, nil
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
		product.Pdetail = "https://adiecom.gitfunswokhu.in/admin/listingSingleProduct/" + product.Pid
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
		SELECT code,expiry,min_amount,amount,max_amount from coupon;`

	rows, err := r.sql.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute select query: %w", err)
	}
	defer rows.Close()

	var products []model.Coupon
	for rows.Next() {
		var product model.Coupon
		err := rows.Scan(&product.Code, &product.Expiry, &product.Minamount, &product.Amount, &product.Maxamount)
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
		SELECT code,expiry,min_amount,amount,max_amount from coupon ORDER BY id DESC;`

	rows, err := r.sql.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute select query: %w", err)
	}
	defer rows.Close()

	var products []model.Coupon
	for rows.Next() {
		var product model.Coupon
		err := rows.Scan(&product.Code, &product.Expiry, &product.Minamount, &product.Amount, &product.Maxamount)
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
		product.Pdetail = "https://adiecom.gitfunswokhu.in/user/listingSingleProduct/" + product.Pid
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
		product.Pdetail = "https://adiecom.gitfunswokhu.in/user/listingSingleProduct/" + product.Pid
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
		product.Pdetail = "https://adiecom.gitfunswokhu.in/user/listingSingleProduct/" + product.Pid
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
		product.Pdetail = "https://adiecom.gitfunswokhu.in/user/listingSingleProduct/" + product.Pid
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
            product_models.brand,
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
		err := rows.Scan(&product.Name, &product.Category, &product.Unit, &product.Tax, &product.Price, &product.Status, &product.Discount, &product.VendorName, &product.Pid, &product.VEmail, &product.VGst, &product.VId, &product.Brand, &product.Pds)
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
