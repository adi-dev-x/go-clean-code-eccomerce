package user

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strings"

	"myproject/pkg/model"
)

// ListWish
type Repository interface {
	Login(ctx context.Context, email string) (model.UserRegisterRequest, error)
	Register(ctx context.Context, request model.UserRegisterRequest) (string, error)
	UpdateUser(ctx context.Context, query string, args []interface{}) error
	ListAddress(ctx context.Context, id string) ([]model.Address, error)
	AddAddress(ctx context.Context, request model.Address, id string) error

	//cart and wishlist
	Listcart(ctx context.Context, id string) ([]model.Usercartview, error)
	ListWish(ctx context.Context, id string) ([]model.UserWishview, error)
	AddTocart(ctx context.Context, request model.Cart) error
	UpdateToCart(ctx context.Context, request model.Cart) error
	AddToWish(ctx context.Context, request model.Wishlist) error
	GetCoupon(ctx context.Context, id string, amount int) model.CouponRes
	GetcartAmt(ctx context.Context, id string) int
	GetcartDis(ctx context.Context, id string) int
	GetCartById(ctx context.Context, cartId string) (model.Cart, error)
	GetcartRes(ctx context.Context, id string) ([]model.Cartresponse, error)
	GetSpecificCart(ctx context.Context, userid string, pis string) (model.Cart, error)
	///

	///product listing
	Listing(ctx context.Context) ([]model.ProductListingUsers, error)
	CategoryListing(ctx context.Context, category string) ([]model.ProductListingUsers, error)
	ListingSingle(ctx context.Context, id string) ([]model.ProductListDetailed, error)
	LatestListing(ctx context.Context) ([]model.ProductListingUsers, error)
	PhighListing(ctx context.Context) ([]model.ProductListingUsers, error)
	PlowListing(ctx context.Context) ([]model.ProductListingUsers, error)
	InAZListing(ctx context.Context) ([]model.ProductListingUsers, error)
	InZAListing(ctx context.Context) ([]model.ProductListingUsers, error)
	ListingByid(ctx context.Context, id string) ([]model.ProductList, error)
	GetProductIDFromCart(ctx context.Context, cartId string) (string, error)
	///

	UpdateProductUnits(ctx context.Context, productId string, newUnits float64) error
	AddToorder(ctx context.Context, request model.Order) error
	GetorderDetails(ctx context.Context, request model.Order) (model.FirstAddOrder, error)
	ActiveListing(ctx context.Context) ([]model.Coupon, error)
	Getid(ctx context.Context, username string) string

	CreateWallet(ctx context.Context, id string) error
	CreateOrder(ctx context.Context, order model.InsertOrder) (string, string, error)
	//AddToPayment(ctx context.Context, request model.Order, fiData model.FirstAddOrder, status string, username string) (string, error)
	AddOrderItems(ctx context.Context, cartData model.CartresponseData, OrderID string, id string, pid string) error
	MakePayment(ctx context.Context, paySt model.PaymentInsert) (string, error)
	UpdateStock(ctx context.Context, value interface{}) error

	///cart
	DeleteCart(ctx context.Context, id string) error
	GetCartExist(ctx context.Context, id string) (string, error)
	UpdateUsestatusCoupon(ctx context.Context, couponId string, id string) error
	UpdateOrderStatus(ctx context.Context, id string, status string) error
	UpdatePaymentStatus(ctx context.Context, id string, status string) error
	ItemExistsInCart(ctx context.Context, id string, status string) (bool, error)
	DeleteSingleCart(ctx context.Context, id string, pid string) error
	///orderss
	ListAllOrders(ctx context.Context, id string) ([]model.ListAllOrdersUsers, error)
	ListReturnedOrders(ctx context.Context, id string) ([]model.ListAllOrders, error)
	ListFailedOrders(ctx context.Context, id string) ([]model.ListAllOrders, error)
	ListCompletedOrders(ctx context.Context, id string) ([]model.ListAllOrders, error)
	ListPendingOrders(ctx context.Context, id string) ([]model.ListAllOrders, error)
	GetSingleItem(ctx context.Context, id string, oid string) (model.ListAllOrdersCheck, error)
	IncreaseStock(ctx context.Context, id string, unit int) error
	UpdateOiStatus(ctx context.Context, id string, typ string) error
	////
	CheckCouponExist(ctx context.Context, couponID string, userID string) (bool, error)
	///wallet
	CreditWallet(ctx context.Context, id string, amt float64) (string, error)
	UpdateWallet(ctx context.Context, value interface{}) (string, error)
	UpdateWalletTransaction(ctx context.Context, value interface{}) error
	GetWallAmt(ctx context.Context, id string, amount int) float32
	GetCoupnExist(ctx context.Context, id string) (string, error)
	//transaction
	ListAllTransactions(ctx context.Context, id string) ([]model.UserTransactions, error)

	ListTypeTransactions(ctx context.Context, id string, ty string) ([]model.UserTransactions, error)

	PrintingUserMainOrder(ctx context.Context, userID string) ([]model.ListingUserMainOrders, error)
	//PrintingUserSingleMainOrder(ctx context.Context, userID string, orderUid string) (model.ListingMainOrders, error)
	PrintingUserSingleMainOrder(ctx context.Context, userID string, orderUid string) ([]model.ReturnOrderPostForUser, error)

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
	UPDATE users
	SET verification =true
	WHERE email = $1
	`

	_, err := r.sql.ExecContext(ctx, query, email)

	if err != nil {
		fmt.Errorf("failed to execute update query: %w", err)
	}

}
func (r *repository) ChangeOrderStatus(ctx context.Context, id string) error {
	fmt.Println("changing in the order status QQQQ", id)
	query := `
	UPDATE orders
	SET status ='Cancelled'
	WHERE uuid = $1
	
`

	_, err := r.sql.ExecContext(ctx, query, id)

	if err != nil {
		return fmt.Errorf("failed to execute update query: %w", err)
	}

	return nil
}

// func (r *repository) PrintingUserSingleMainOrder(ctx context.Context, userID string, orderUid string) (model.ListingMainOrders, error) {
// 	query := `SELECT mo.uuid, mo.delivered, mo.payment_method, mo.status, mo.payable_amount,
// 	                 u.firstname || ' ' || u.lastname AS user,
// 	                 COALESCE(a.address1, '') || ' ' || COALESCE(a.address2, '') || ' ' ||
// 	                 COALESCE(a.address3, '') || ' ' || COALESCE(a.city, '') || ' ' ||
// 	                 COALESCE(a.state, '') || ' ' || COALESCE(a.pin, '') || ' ' ||
// 	                 COALESCE(a.country, '') AS user_ad,
// 	                 COALESCE(DATE(mo.delivery_date)::text, '') AS delivery_date
// 			  FROM orders mo
// 			  JOIN users u ON mo.user_id = u.id
// 			  JOIN address a ON mo.address_id = a.address_id
// 			  WHERE u.id=$1 AND mo.uuid=$2`
// 	var order model.ListingMainOrders

// 	err := r.sql.QueryRowContext(ctx, query, userID, orderUid).Scan(&order.OR_id, &order.Delivery_Stat, &order.D_Type, &order.O_status, &order.Amount,
// 		&order.User, &order.UserAddress, &order.Delivery_date)
// 	if err != nil {
// 		return model.ListingMainOrders{}, fmt.Errorf("can't execute query: %w", err)
// 	}

// 	return order, nil
// }

func (r *repository) PrintingUserSingleMainOrder(ctx context.Context, userID string, orderUid string) ([]model.ReturnOrderPostForUser, error) {
	query := `
		SELECT oi.id AS oid
		FROM orders mo 
		JOIN order_items oi ON mo.id = oi.order_id 
		WHERE mo.uuid = $1 AND mo.user_id = $2
	`

	fmt.Println("Executing PrintingUserSingleMainOrder with UserID:", userID, "OrderUID:", orderUid)

	// Define a slice to hold the order item IDs (which are strings)
	var orderItemIDs []model.ReturnOrderPostForUser

	// Query the database
	rows, err := r.sql.QueryContext(ctx, query, orderUid, userID)
	if err != nil {
		return nil, fmt.Errorf("can't execute query: %w", err)
	}
	defer rows.Close()
	fmt.Println(" rooowsss  ", rows)
	// Iterate over the result rows
	for rows.Next() {
		var orderItemID model.ReturnOrderPostForUser
		if err := rows.Scan(&orderItemID.Oid); err != nil {
			fmt.Println("errr in thiss")
			return nil, fmt.Errorf("can't scan result: %w", err)
		}
		orderItemID.MoReturn = true
		orderItemID.Type = "Cancelled"
		orderItemIDs = append(orderItemIDs, orderItemID)
	}

	// Check for any errors after iterating through the result rows
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return orderItemIDs, nil
}

func (r *repository) PrintingUserMainOrder(ctx context.Context, userID string) ([]model.ListingUserMainOrders, error) {
	query := `SELECT mo.uuid, mo.delivered, mo.payment_method, mo.status, mo.payable_amount, 
	                 mo.delivery_date,
			        mo.discount,mo.coupon_amount AS cmt,COALESCE(c.code , ''),mo.wallet_money AS wmt 
					 FROM orders mo
			  JOIN users u ON mo.user_id = u.id 
			  JOIN address a ON mo.address_id = a.address_id
			  LEFT JOIN coupon c ON mo.cid =c.id
			  WHERE u.id=$1`
	var orders []model.ListingUserMainOrders

	rows, err := r.sql.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("can't execute query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var order model.ListingUserMainOrders
		err := rows.Scan(&order.OR_id, &order.Delivery_Stat, &order.D_Type, &order.O_status, &order.Amount,
			&order.Delivery_date, &order.Discount, &order.Cmt, &order.Code, &order.Wmt)
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

func (r *repository) ItemExistsInCart(ctx context.Context, userID string, productID string) (bool, error) {
	fmt.Println("check exist in ,", userID, "pid !!", productID)
	query := `
		SELECT 1 
		FROM cart 
		WHERE product_id = $1 AND user_id = $2
		LIMIT 1;
	`

	var exists int
	err := r.sql.QueryRowContext(ctx, query, userID, productID).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	return exists == 1, nil
}

// /transaction
func (r *repository) ListAllTransactions(ctx context.Context, id string) ([]model.UserTransactions, error) {
	query := `select id,amount,transaction_type from wallet_transactions where user_id=$1`
	var transactions []model.UserTransactions

	rows, err := r.sql.QueryContext(ctx, query, id)
	if err != nil {

		return nil, fmt.Errorf("cant execute query", err)
	}
	defer rows.Close()
	for rows.Next() {
		var transaction model.UserTransactions
		err := rows.Scan(&transaction.Id, &transaction.Amount, &transaction.Type)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		transactions = append(transactions, transaction)

	}
	return transactions, nil

}
func (r *repository) ListTypeTransactions(ctx context.Context, id string, typ string) ([]model.UserTransactions, error) {
	query := `select id,amount,transaction_type from wallet_transactions where user_id=$1 AND transaction_type=$2`
	var transactions []model.UserTransactions

	rows, err := r.sql.QueryContext(ctx, query, id, typ)
	if err != nil {

		return nil, fmt.Errorf("cant execute query", err)
	}
	defer rows.Close()
	for rows.Next() {
		var transaction model.UserTransactions
		err := rows.Scan(&transaction.Id, &transaction.Amount, &transaction.Type)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		transactions = append(transactions, transaction)

	}
	return transactions, nil

}

func (r *repository) CreditWallet(ctx context.Context, id string, amt float64) (string, error) {
	query := `
	UPDATE wallet
	SET balance = balance + $1
	WHERE user_id = $2
	RETURNING id;
`
	fmt.Println("wallll---", id, amt)
	var Wallet_id string
	err := r.sql.QueryRowContext(ctx, query, amt, id).Scan(&Wallet_id)
	fmt.Println("hey adiii CreditWallet????!!!!!", Wallet_id, "wallet_id !", "id!", id)
	if err != nil {
		return "", fmt.Errorf("failed to execute update query wallet transaction: %w", err)
	}

	return Wallet_id, nil
}
func (r *repository) UpdateOiStatus(ctx context.Context, id string, ty string) error {
	fmt.Println("---", id, ty)
	query := `
	UPDATE order_items
	SET returned = true,
	re_cl=$1
	WHERE id = $2;
`

	_, err := r.sql.ExecContext(ctx, query, ty, id)
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
func (r *repository) GetSingleItem(ctx context.Context, id string, oid string) (model.ListAllOrdersCheck, error) {
	var order model.ListAllOrdersCheck

	query := `SELECT p.name, oi.quantity, mo.status, oi.returned, oi.price,oi.product_id AS pid, DATE(oi.created_at) AS date,mo.user_id,oi.delivered FROM order_items oi 
   JOIN product_models p ON oi.product_id = p.id 
   JOIN orders mo ON oi.order_id = mo.id 
   where oi.user_id=$1 AND oi.id=$2;
   
`
	err := r.sql.QueryRowContext(ctx, query, id, oid).Scan(&order.Name, &order.Unit, &order.Status, &order.Returned, &order.Amount, &order.Pid, &order.Date, &order.Usid, &order.Delivery)
	if err != nil {
		return model.ListAllOrdersCheck{}, fmt.Errorf("error in exequting query in  GetSingleItem")
	}
	return order, nil
}
func (r *repository) ListAllOrders(ctx context.Context, id string) ([]model.ListAllOrdersUsers, error) {
	query := `SELECT p.name, oi.quantity, mo.status, oi.returned, oi.price,oi.product_id, DATE(oi.created_at) AS date,oi.id AS oid,oi.discount FROM order_items oi 
	JOIN product_models p ON oi.product_id = p.id 
	JOIN orders mo ON oi.order_id = mo.id 
	where oi.user_id=$1 ORDER BY oi.id DESC;
`
	var orders []model.ListAllOrdersUsers

	rows, err := r.sql.QueryContext(ctx, query, id)
	if err != nil {
		return []model.ListAllOrdersUsers{}, fmt.Errorf("error in exequting query in ListAllOrders ")
	}
	defer rows.Close()
	for rows.Next() {
		var order model.ListAllOrdersUsers
		err := rows.Scan(&order.Name, &order.Unit, &order.Status, &order.Returned, &order.Amount, &order.Pid, &order.Date, &order.Oid, &order.Discount)
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
func (r *repository) ListReturnedOrders(ctx context.Context, id string) ([]model.ListAllOrders, error) {
	query := `SELECT p.name, oi.quantity, mo.status, oi.returned, oi.price,oi.product_id, DATE(oi.created_at) AS date FROM order_items oi 
	JOIN product_models p ON oi.product_id = p.id 
	JOIN orders mo ON oi.order_id = mo.id 
	where oi.user_id=$1 AND oi.returned=true;
`
	var orders []model.ListAllOrders

	rows, err := r.sql.QueryContext(ctx, query, id)
	if err != nil {
		return []model.ListAllOrders{}, fmt.Errorf("error in exequting query in ListAllOrders ")
	}
	defer rows.Close()
	for rows.Next() {
		var order model.ListAllOrders
		err := rows.Scan(&order.Name, &order.Unit, &order.Status, &order.Returned, &order.Amount, &order.Pid, &order.Date)
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
func (r *repository) ListFailedOrders(ctx context.Context, id string) ([]model.ListAllOrders, error) {
	query := `SELECT p.name, oi.quantity, mo.status, oi.returned, oi.price,oi.product_id, DATE(oi.created_at) AS date FROM order_items oi 
	JOIN product_models p ON oi.product_id = p.id 
	JOIN orders mo ON oi.order_id = mo.id 
	where oi.user_id=$1 AND mo.status='Failed';
`
	var orders []model.ListAllOrders

	rows, err := r.sql.QueryContext(ctx, query, id)
	if err != nil {
		return []model.ListAllOrders{}, fmt.Errorf("error in exequting query in ListAllOrders ")
	}
	defer rows.Close()
	for rows.Next() {
		var order model.ListAllOrders
		err := rows.Scan(&order.Name, &order.Unit, &order.Status, &order.Returned, &order.Amount, &order.Pid, &order.Date)
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
func (r *repository) ListCompletedOrders(ctx context.Context, id string) ([]model.ListAllOrders, error) {
	query := `SELECT p.name, oi.quantity, mo.status, oi.returned, oi.price,oi.product_id, DATE(oi.created_at) AS date FROM order_items oi 
	JOIN product_models p ON oi.product_id = p.id 
	JOIN orders mo ON oi.order_id = mo.id 
	where oi.user_id=$1 AND mo.status='Completed';
`
	var orders []model.ListAllOrders

	rows, err := r.sql.QueryContext(ctx, query, id)
	if err != nil {
		return []model.ListAllOrders{}, fmt.Errorf("error in exequting query in ListAllOrders ")
	}
	defer rows.Close()
	for rows.Next() {
		var order model.ListAllOrders
		err := rows.Scan(&order.Name, &order.Unit, &order.Status, &order.Returned, &order.Amount, &order.Pid, &order.Date)
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
func (r *repository) ListPendingOrders(ctx context.Context, id string) ([]model.ListAllOrders, error) {
	query := `SELECT p.name, oi.quantity, mo.status, oi.returned, oi.price ,oi.product_id , DATE(oi.created_at) AS date FROM order_items oi 
	JOIN product_models p ON oi.product_id = p.id 
	JOIN orders mo ON oi.order_id = mo.id 
	where oi.user_id=$1 AND mo.status='Pending';
`
	var orders []model.ListAllOrders

	rows, err := r.sql.QueryContext(ctx, query, id)
	if err != nil {
		return []model.ListAllOrders{}, fmt.Errorf("error in exequting query in ListAllOrders ")
	}
	defer rows.Close()
	for rows.Next() {
		var order model.ListAllOrders
		err := rows.Scan(&order.Name, &order.Unit, &order.Status, &order.Returned, &order.Amount, &order.Pid, &order.Date)
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
func (r *repository) UpdateOrderStatus(ctx context.Context, id string, status string) error {
	fmt.Println("this is in the UpdateOrderStatus ", id)
	query := `
	UPDATE orders SET status = $1 WHERE uuid = $2;
`

	_, err := r.sql.ExecContext(ctx, query, status, id)

	if err != nil {
		return fmt.Errorf("failed to execute update query: %w", err)
	}
	fmt.Println("completed 1")
	return nil

}
func (r *repository) UpdatePaymentStatus(ctx context.Context, id string, status string) error {
	fmt.Println("this is in the UpdatePaymentStatus ", id)
	query := `
	UPDATE payment SET  payment_status = $1 WHERE id = $2;
`

	_, err := r.sql.ExecContext(ctx, query, status, id)

	if err != nil {
		return fmt.Errorf("failed to execute update query: %w", err)
	}
	fmt.Println("completed 2")
	return nil

}
func (r *repository) UpdateUsestatusCoupon(ctx context.Context, couponId string, id string) error {
	query := `
	INSERT INTO coupon_usages (coupon_id, user_id)
	VALUES ($1, $2);
`
	_, err := r.sql.ExecContext(ctx, query, couponId, id)
	if err != nil {
		return fmt.Errorf("failed to insert coupon usage: %w", err)
	}
	return nil

}
func (r *repository) DeleteCart(ctx context.Context, id string) error {
	query := ` delete from cart where user_id=$1`
	result, err := r.sql.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete from cart: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no cart found with id: %s", id)
	}

	return nil
}
func (r *repository) DeleteSingleCart(ctx context.Context, id string, pid string) error {
	fmt.Println("this is in the rep DeleteSingleCart", id, pid)
	query := ` delete from cart where user_id=$1 AND product_id=$2`
	result, err := r.sql.ExecContext(ctx, query, id, pid)
	if err != nil {
		return fmt.Errorf("failed to delete from cart: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no cart found with id: %s", id)
	}

	return nil
}
func (r *repository) UpdateWalletTransaction(ctx context.Context, value interface{}) error {
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

func (r *repository) UpdateWallet(ctx context.Context, value interface{}) (string, error) {

	fmt.Println(value)
	///!!!!!!!!!!!!!!!!!!!!!!!
	values, ok := value.([]interface{})
	if !ok {
		return "", fmt.Errorf("invalid input")
	}

	if len(values) != 2 {
		return "", fmt.Errorf("invalid input: need 2 elements", len(values))
	}

	qua := values[0]
	id := values[1]

	////1//////!!!!!!!!!!!!!!!!!!!!

	query := `
	UPDATE wallet
	SET balance = balance - $1
	WHERE user_id = $2
	RETURNING id;
`
	var Wallet_id string

	err := r.sql.QueryRowContext(ctx, query, qua, id).Scan(&Wallet_id)
	fmt.Println("hey adiii UpdateWallet????!!!!!", Wallet_id, "balance !", qua, "id!", id)
	if err != nil {
		return "", fmt.Errorf("failed to execute update query: %w", err)
	}

	return Wallet_id, nil
}
func (r *repository) UpdateStock(ctx context.Context, value interface{}) error {

	fmt.Println(value)
	///!!!!!!!!!!!!!!!!!!!!!!!
	values, ok := value.([]interface{})
	if !ok {
		return fmt.Errorf("invalid input: expected []interface{}")
	}

	// Ensure the values slice has the correct number of elements
	if len(values) != 2 {
		return fmt.Errorf("invalid input: expected 2 elements, got %d", len(values))
	}

	quaInterface := values[0]
	idInterface := values[1]

	// Type assertion to extract underlying int values
	qua, _ := quaInterface.(int)
	id, _ := idInterface.(string)

	////1//////!!!!!!!!!!!!!!!!!!!!

	query := `
	UPDATE product_models
	SET units = units - $1
	WHERE id = $2
	RETURNING id;
`
	var Product_id string

	err := r.sql.QueryRowContext(ctx, query, qua, id).Scan(&Product_id)
	if err != nil {
		return fmt.Errorf("failed to execute update query: %w", err)
	}

	return nil
}
func (r *repository) MakePayment(ctx context.Context, payment model.PaymentInsert) (string, error) {
	fmt.Println("this is the MakePayment in repo!!111", payment)
	query := `
	INSERT INTO payment (
		order_id, user_id, amount, payment_method, payment_status, 
		payment_date, created_at, updated_at
	) VALUES (
		$1, $2, $3, $4, $5,  CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
	) RETURNING id;
  `
	var paymentID string

	err := r.sql.QueryRowContext(ctx, query,
		payment.OrderId,
		payment.Usid,
		payment.Amount,
		payment.Type,
		payment.Status,
	).Scan(&paymentID)
	if err != nil {
		return "", fmt.Errorf("failed to execute insert query: %w", err)
	}

	return paymentID, nil
}
func (r *repository) AddOrderItems(ctx context.Context, cartDatas model.CartresponseData, orderID string, id string, pid string) error {
	fmt.Println("this is in the repo of AddOrderItems!!!", cartDatas.Data)
	var (
		queryBuilder strings.Builder
		values       []interface{}
	)
	cartData := cartDatas.Data

	queryBuilder.WriteString(`INSERT INTO order_items (order_id, product_id, quantity, price, discount, returned,user_id,payment_id) VALUES `)

	for i, item := range cartData {
		if i > 0 {
			queryBuilder.WriteString(", ")
		}
		queryBuilder.WriteString(fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)", i*8+1, i*8+2, i*8+3, i*8+4, i*8+5, i*8+6, i*8+7, i*8+8))
		values = append(values, orderID, item.Pid, item.Unit, item.Amount, item.Discount, false, id, pid)
	}

	query := queryBuilder.String()
	fmt.Println("Generated query: ", query)

	tx, err := r.sql.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	_, err = tx.ExecContext(ctx, query, values...)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to execute insert query: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
func (r *repository) CreateOrder(ctx context.Context, order model.InsertOrder) (string, string, error) {
	if order.CouponId == "" {

		fmt.Println("yes its empltyy !!!!!!")
	}
	query := `
        INSERT INTO orders (
            user_id,total_amount,discount,coupon_amount,wallet_money,payable_amount,
            payment_method,address_id,status,cid,order_date,created_at,updated_at,uuid
        ) VALUES (
            $1, $2, $3, $4, $5, $6, $7, $8, $9, NULLIF($10, '')::integer, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP,uuid_generate_v4()
        ) RETURNING id,uuid;
    `

	var orderID string
	var uuid string
	err := r.sql.QueryRowContext(ctx, query,
		order.Usid,
		order.Amount,
		order.Discount,
		order.CouponAmt,
		order.WalletAmt,
		order.PayableAmt,
		order.PayType,
		order.Aid,
		order.Status,
		order.CouponId,
	).Scan(&orderID, &uuid)
	if err != nil {
		return "", "", fmt.Errorf("failed to execute insert query: %w", err)
	}

	return orderID, uuid, nil
}
func (r *repository) ActiveListing(ctx context.Context) ([]model.Coupon, error) {
	fmt.Println("this is lia couppp")

	query := `
		SELECT code,expiry,min_amount,amount,max_amount from coupon WHERE TO_DATE(expiry, 'YYYY-MM-DD') > CURRENT_DATE; `

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
func (r *repository) UpdateProductUnits(ctx context.Context, productId string, newUnits float64) error {
	query := `UPDATE product_models SET units = $1 WHERE id = $2`
	_, err := r.sql.ExecContext(ctx, query, newUnits, productId)
	if err != nil {
		return fmt.Errorf("failed to update product units: %w", err)
	}
	return nil
}

func (r *repository) AddToorder(ctx context.Context, request model.Order) error {
	query := `INSERT INTO orders (cart_id, ) VALUES ($1, $2, $3, NOW())`
	_, err := r.sql.ExecContext(ctx, query, request.Cartid)
	if err != nil {
		return fmt.Errorf("failed to execute insert query: %w", err)
	}
	return nil
}
func (r *repository) GetorderDetails(ctx context.Context, request model.Order) (model.FirstAddOrder, error) {
	// cartDataChan := make(chan model.CartresponseData)
	// amountChan := make(chan int)
	// //couponChan := make(chan model.Coupon)

	// go func() {
	// 	data, err := r.GetcartRes(ctx)
	// 	cartDataChan <- model.CartresponseData{Data: data, Err: err}
	// 	close(cartDataChan)

	// }()
	// go func() {
	// 	amount := r.GetcartAmt(ctx)
	// 	amountChan <- amount
	// 	close(amountChan)

	// }()

	// cartData := <-cartDataChan
	// amount := <-amountChan
	// fmt.Println("this is amttt ", amount)
	// id := request.Cartid
	// var coupon = model.CouponRes{}
	// if id != "" {
	// 	coupon = r.GetCoupon(ctx, request, amount)
	// 	fmt.Println("this is coupon!!!!!", coupon)

	// }
	// resD := model.FirstAddOrder{
	// 	Data:    cartData,
	// 	TAmount: amount,
	// 	CData:   coupon,
	// }
	// if coupon.Present {
	// 	if !coupon.Is_eligible || !coupon.Is_expired || coupon.Used {
	// 		resD.Notvalid = true
	// 	} else {
	// 		resD.TAmount = resD.TAmount - coupon.Amount
	// 	}

	// }
	// fmt.Println("this is resD 1!!!!!", resD)
	resD := model.FirstAddOrder{}
	return resD, nil
}

func (r *repository) GetcartRes(ctx context.Context, id string) ([]model.Cartresponse, error) {
	var cres []model.Cartresponse
	fmt.Println("reached inside GetcartRes")

	query := `SELECT c.id AS cid, c.user_id AS usid,c.product_id AS pid ,c.unit,p.amount,p.discount,p.units AS p_units,p.name AS p_name  FROM cart  c JOIN product_models p on c.product_id = p.id JOIN users u ON c.user_id=u.id WHERE u.id = $1`
	//rows, err := r.sql.QueryContext(ctx, query, username)
	rows, err := r.sql.Query(query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to execute select query: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var cress model.Cartresponse
		err := rows.Scan(&cress.Cid, &cress.Usid, &cress.Pid, &cress.Unit, &cress.Amount, &cress.Discount, &cress.P_Units, &cress.P_Name)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		cres = append(cres, cress)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}
	fmt.Println("this is the data ", cres)

	return cres, nil
}
func (r *repository) GetcartAmt(ctx context.Context, id string) int {
	total := 0
	query := `SELECT SUM(c.unit * p.amount - p.discount) AS total_sum 
	          FROM cart c 
	          JOIN product_models p ON c.product_id = p.id 
	          JOIN users u ON c.user_id = u.id 
	          WHERE u.id = $1`

	row := r.sql.QueryRowContext(ctx, query, id)
	err := row.Scan(&total)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("No rows found for username:", id)
			return 0
		}
		fmt.Println("Error scanning row:", err)
		return 0
	}

	fmt.Println("Total amount:", total)
	return total
}
func (r *repository) GetcartDis(ctx context.Context, id string) int {
	total := 0
	query := `SELECT SUM(p.discount) AS total_sum 
	          FROM cart c 
	          JOIN product_models p ON c.product_id = p.id 
	          JOIN users u ON c.user_id = u.id 
	          WHERE u.id = $1`

	row := r.sql.QueryRowContext(ctx, query, id)
	err := row.Scan(&total)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("No rows found for username:", id)
			return 0
		}
		fmt.Println("Error scanning row:", err)
		return 0
	}

	fmt.Println("Total amount:", total)
	return total
}
func (r *repository) GetCoupon(ctx context.Context, id string, amount int) model.CouponRes {
	fmt.Println("thissssss issssss getCoupon")
	var coupon model.CouponRes
	fmt.Println("this is the coupon check amountsss !!!", amount, "!!!!!! this is id  ", id)

	var query = `SELECT id AS cid, code, expiry, CURRENT_DATE AS current_date, CASE WHEN TO_DATE(expiry, 'YYYY-MM-DD') < CURRENT_DATE THEN true ELSE false END AS is_expired, 
	CASE WHEN $1 > min_amount THEN true ELSE false END AS is_eligible, min_amount, amount, used,max_amount FROM coupon WHERE id = $2;`

	err := r.sql.QueryRowContext(ctx, query, amount, id).Scan(
		&coupon.Cid,
		&coupon.Code,
		&coupon.Expiry,
		&coupon.CurrentDate,
		&coupon.Is_expired,
		&coupon.Is_eligible,
		&coupon.Minamount,
		&coupon.Amount,
		&coupon.Used,
		&coupon.Maxamount,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return model.CouponRes{}
		}
		return model.CouponRes{}
	}

	return coupon
}
func (r *repository) GetWallAmt(ctx context.Context, id string, amount int) float32 {
	fmt.Println("thissssss issssss getCoupon")

	var amt float32
	var query = `SELECT balance from wallet WHERE user_id = $1;`
	row := r.sql.QueryRowContext(ctx, query, id)
	err := row.Scan(&amt)

	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("No rows found for GetWallAmt:", id)
			return 0.0
		}
		fmt.Println("Error scanning row:", err)
		return 0.0
	}

	return amt
}

func (r *repository) GetCartById(ctx context.Context, cartId string) (model.Cart, error) {
	query := `SELECT product_id, user_id, unit FROM cart WHERE id = $1 `
	var cart model.Cart
	err := r.sql.QueryRowContext(ctx, query, cartId).Scan(&cart.Productid, &cart.Userid, &cart.Unit)
	if err != nil {
		return model.Cart{}, fmt.Errorf("failed to get cart: %w", err)
	}
	return cart, nil
}
func (r *repository) GetSpecificCart(ctx context.Context, userid string, pis string) (model.Cart, error) {
	query := `SELECT product_id, user_id, unit FROM cart WHERE user_id = $1 AND product_id=$2 `
	var cart model.Cart
	err := r.sql.QueryRowContext(ctx, query, userid, pis).Scan(&cart.Productid, &cart.Userid, &cart.Unit)
	if err != nil {
		return model.Cart{}, fmt.Errorf("failed to get cart: %w", err)
	}
	return cart, nil //GetCartExist
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
func (r *repository) CheckCouponExist(ctx context.Context, couponID string, userID string) (bool, error) {
	query := `
        SELECT EXISTS (
            SELECT 1
            FROM coupon_usages
            WHERE coupon_id = $1
              AND user_id = $2
        );
    `

	var exists bool
	err := r.sql.QueryRowContext(ctx, query, couponID, userID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check coupon usage existence: %w", err)
	}

	return exists, nil
}
func (r *repository) GetCartExist(ctx context.Context, id string) (string, error) {
	query := `SELECT id FROM cart WHERE user_id = $1 LIMIT 1`
	var exist string
	err := r.sql.QueryRowContext(ctx, query, id).Scan(&exist)
	fmt.Println("this is in the repo layerrr for GetCartExist!!!", exist)
	if err != nil {
		return "", fmt.Errorf("failed to get cart: %w", err)
	}
	return exist, nil
}
func (r *repository) ListAddress(ctx context.Context, id string) ([]model.Address, error) {
	fmt.Println("this is in repo of ListAddress!!####!!!", id)
	query := `SELECT address1,address2,address3,pin,country,state FROM address WHERE user_id = $1;`
	var add []model.Address

	rows, err := r.sql.QueryContext(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var singleAdd model.Address
		err := rows.Scan(&singleAdd.Address1, &singleAdd.Address2, &singleAdd.Address3, &singleAdd.PIN, &singleAdd.Country, &singleAdd.State)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		add = append(add, singleAdd)
	}
	//Scan(&add.Address1, &add.Address2, &add.Address3, &add.PIN, &add.Country, &add.State)
	fmt.Println("this is in the data side of listAddress !!!!", add)
	if err != nil {
		return []model.Address{}, fmt.Errorf("failed to get cart: %w", err)
	}
	return add, nil
}

func (r *repository) Register(ctx context.Context, request model.UserRegisterRequest) (string, error) {
	fmt.Println("this is in the repository Register")
	var id string
	query := `INSERT INTO users (firstname, lastname, email, password) VALUES ($1, $2, $3, $4) Returning id`
	err := r.sql.QueryRowContext(ctx, query, request.FirstName, request.LastName, request.Email, request.Password).Scan(&id)
	if err != nil {
		return "", fmt.Errorf("failed to execute insert query: %w", err)
	}

	return id, nil
}
func (r *repository) CreateWallet(ctx context.Context, id string) error {
	fmt.Println("this is in the repository Register")
	query := `INSERT INTO wallet (user_id,balance) VALUES ($1, 0.0)`
	_, err := r.sql.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to execute insert query: %w", err)
	}

	return nil
}
func (r *repository) AddTocart(ctx context.Context, request model.Cart) error {
	fmt.Println("this is in the repository Register !!!")
	query := `INSERT INTO cart (product_id, unit, user_id) VALUES ($1, $2, $3)`
	fmt.Println(query, request.Productid, request.Unit, request.Userid)
	_, err := r.sql.ExecContext(ctx, query, request.Productid, request.Unit, request.Userid)
	if err != nil {
		return fmt.Errorf("failed to execute insert query: %w", err)
	}

	return nil
}
func (r *repository) UpdateToCart(ctx context.Context, request model.Cart) error {
	fmt.Println("this is in the repository Register !!!")
	//query := `INSERT INTO cart (product_id, unit, user_id) VALUES ($1, $2, $3)`
	query := `UPDATE cart set unit=$1 WHERE user_id=$2 AND  product_id=$3`
	fmt.Println(query, request.Unit, request.Userid, request.Productid)
	_, err := r.sql.ExecContext(ctx, query, request.Unit, request.Userid, request.Productid)
	if err != nil {
		return fmt.Errorf("failed to execute insert query: %w", err)
	}

	return nil
}
func (r *repository) AddToWish(ctx context.Context, request model.Wishlist) error {
	fmt.Println("this is in the repository Register !!!")
	query := `INSERT INTO wishlist (product_id, user_id) VALUES ($1, $2)`
	fmt.Println(query, request.Productid, request.Userid)
	_, err := r.sql.ExecContext(ctx, query, request.Productid, request.Userid)
	if err != nil {
		return fmt.Errorf("failed to execute insert query: %w", err)
	}

	return nil
}
func (r *repository) Getid(ctx context.Context, username string) string {
	var id string
	fmt.Println("this is in the repository Register !!!")
	query := `select id from users where email=$1;`
	fmt.Println(query, username)
	row := r.sql.QueryRowContext(ctx, query, username)
	err := row.Scan(&id)
	fmt.Println(err)
	fmt.Println("this is id returning from Getid:::", id)

	return id
}
func (r *repository) AddAddress(ctx context.Context, request model.Address, id string) error {
	fmt.Println("this is in the repository Register !!!")

	query := `INSERT INTO address (address1,address2,address3,pin,country,state,user_id) VALUES ($1, $2,$3,$4,$5,$6,$7)`
	fmt.Println(query, request.Address1, request.Address2, request.Address3, request.PIN, request.Country, request.State, id)
	_, err := r.sql.ExecContext(ctx, query, request.Address1, request.Address2, request.Address3, request.PIN, request.Country, request.State, id)
	if err != nil {
		return fmt.Errorf("failed to execute insert query: %w", err)
	}

	return nil
}
func (r *repository) Login(ctx context.Context, email string) (model.UserRegisterRequest, error) {
	fmt.Println("theee !!!!!!!!!!!  LLLLoginnnnnn  ", email)
	query := `SELECT firstname, lastname, email, password FROM users WHERE email = $1 AND verification=true`
	fmt.Println(`SELECT firstname, lastname, email, password FROM users WHERE email = 'adithyanunni258@gmail.com' ;`)

	var user model.UserRegisterRequest
	err := r.sql.QueryRowContext(ctx, query, email).Scan(&user.FirstName, &user.LastName, &user.Email, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return model.UserRegisterRequest{}, nil
		}
		return model.UserRegisterRequest{}, fmt.Errorf("failed to find user by email: %w", err)
	}
	fmt.Println("the data !!!! ", user)

	return user, nil
}
func (r *repository) UpdateUser(ctx context.Context, query string, args []interface{}) error {
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
func (r *repository) Listcart(ctx context.Context, id string) ([]model.Usercartview, error) {
	query := `SELECT pm.name AS product_name, pm.status AS status, v.name AS vendor_name, pm.amount AS product_amount, pm.category AS product_category, c.unit FROM cart c JOIN product_models pm ON c.product_id = pm.id JOIN vendor v ON pm.vendor_id = v.id JOIN users u ON c.user_id = u.id WHERE u.email = $1;`

	rows, err := r.sql.QueryContext(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to execute select query: %w", err)
	}
	defer rows.Close()

	var products []model.Usercartview
	for rows.Next() {
		var product model.Usercartview
		err := rows.Scan(&product.Productname, &product.Productstatus, &product.Vendorname, &product.Productamount, &product.Productcategory, &product.Unit)
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
func (r *repository) ListWish(ctx context.Context, id string) ([]model.UserWishview, error) {
	query := `SELECT pm.name AS product_name, pm.status AS status, v.name AS vendor_name, pm.amount AS product_amount, pm.category AS product_category  FROM wishlist c JOIN product_models pm ON c.product_id = pm.id JOIN vendor v ON pm.vendor_id = v.id JOIN users u ON c.user_id = u.id WHERE u.email = $1;`

	rows, err := r.sql.QueryContext(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to execute select query: %w", err)
	}
	defer rows.Close()

	var products []model.UserWishview
	for rows.Next() {
		var product model.UserWishview
		err := rows.Scan(&product.Productname, &product.Productstatus, &product.Vendorname, &product.Productamount, &product.Productcategory)
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

func (r *repository) Listing(ctx context.Context) ([]model.ProductListingUsers, error) {
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
		product.Pdetail = "http://localhost:8080/user/listingSingleProduct/" + product.Pid
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
func (r *repository) GetProductIDFromCart(ctx context.Context, cartId string) (string, error) {
	query := `SELECT product_id FROM cart WHERE id = $1`
	var productID string
	err := r.sql.QueryRowContext(ctx, query, cartId).Scan(&productID)
	if err != nil {
		return "", fmt.Errorf("failed to get product ID from cart: %w", err)
	}
	return productID, nil
}
func (r *repository) ListingByid(ctx context.Context, id string) ([]model.ProductList, error) {
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
			vendor ON product_models.vendor_id = vendor.id WHERE product_models.id=$1;`

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
func (r *repository) LatestListing(ctx context.Context) ([]model.ProductListingUsers, error) {
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
			vendor ON product_models.vendor_id = vendor.id WHERE product_models.units > 0 ORDER BY product_models.id DESC;`

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
