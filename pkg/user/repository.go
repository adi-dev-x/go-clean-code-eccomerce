package user

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"myproject/pkg/model"
)

// ListWish
type Repository interface {
	Register(ctx context.Context, request model.UserRegisterRequest) (string, error)
	AddTocart(ctx context.Context, request model.Cart) error
	AddToWish(ctx context.Context, request model.Wishlist) error
	UpdateUser(ctx context.Context, query string, args []interface{}) error
	Listing(ctx context.Context) ([]model.ProductList, error)
	ListingByid(ctx context.Context, id string) ([]model.ProductList, error)
	Listcart(ctx context.Context, id string) ([]model.Usercartview, error)
	ListWish(ctx context.Context, id string) ([]model.UserWishview, error)
	LatestListing(ctx context.Context) ([]model.ProductList, error)
	PhighListing(ctx context.Context) ([]model.ProductList, error)
	PlowListing(ctx context.Context) ([]model.ProductList, error)
	InAZListing(ctx context.Context) ([]model.ProductList, error)
	InZAListing(ctx context.Context) ([]model.ProductList, error)
	Login(ctx context.Context, email string) (model.UserRegisterRequest, error)
	GetProductIDFromCart(ctx context.Context, cartId string) (string, error)
	GetCartById(ctx context.Context, cartId string) (model.Cart, error)
	UpdateProductUnits(ctx context.Context, productId string, newUnits float64) error
	AddToorder(ctx context.Context, request model.Order) error
	GetorderDetails(ctx context.Context, request model.Order) (model.FirstAddOrder, error)
	ActiveListing(ctx context.Context) ([]model.Coupon, error)
	Getid(ctx context.Context, username string) string
	GetcartRes(ctx context.Context, id string) ([]model.Cartresponse, error)
	ListAddress(ctx context.Context, id string) ([]model.Address, error)
	AddAddress(ctx context.Context, request model.Address, id string) error
	GetcartAmt(ctx context.Context, id string) int
	GetcartDis(ctx context.Context, id string) int
	GetCoupon(ctx context.Context, id string, amount int) model.CouponRes
	GetWallAmt(ctx context.Context, id string, amount int) float32
	CreateWallet(ctx context.Context, id string) error
	CreateOrder(ctx context.Context, order model.InsertOrder) (string, error)
	//AddToPayment(ctx context.Context, request model.Order, fiData model.FirstAddOrder, status string, username string) (string, error)
	AddOrderItems(ctx context.Context, cartData model.CartresponseData, OrderID string, id string) error
	MakePayment(ctx context.Context, paySt model.PaymentInsert) (string, error)
	UpdateStock(ctx context.Context, cartData model.CartresponseData) error
}

type repository struct {
	sql *sql.DB
}

func NewRepository(sqlDB *sql.DB) Repository {
	return &repository{
		sql: sqlDB,
	}
}
func (r *repository) UpdateStock(ctx context.Context, cartData model.CartresponseData) error {

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
func (r *repository) AddOrderItems(ctx context.Context, cartDatas model.CartresponseData, orderID string, id string) error {
	fmt.Println("this is in the repo of AddOrderItems!!!", cartDatas.Data)
	var (
		queryBuilder strings.Builder
		values       []interface{}
	)
	cartData := cartDatas.Data

	queryBuilder.WriteString(`INSERT INTO order_items (order_id, product_id, quantity, price, discount, returned,user_id) VALUES `)

	for i, item := range cartData {
		if i > 0 {
			queryBuilder.WriteString(", ")
		}
		queryBuilder.WriteString(fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d)", i*7+1, i*7+2, i*7+3, i*7+4, i*7+5, i*7+6, i*7+7))
		values = append(values, orderID, item.Pid, item.Unit, item.Amount, item.Discount, false, id)
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
func (r *repository) CreateOrder(ctx context.Context, order model.InsertOrder) (string, error) {

	query := `
        INSERT INTO orders (
            user_id,total_amount,discount,coupon_amount,wallet_money,payable_amount,
            payment_method,address_id,status,cid,order_date,created_at,updated_at
        ) VALUES (
            $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
        ) RETURNING id;
    `

	var orderID string
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
	).Scan(&orderID)
	if err != nil {
		return "", fmt.Errorf("failed to execute insert query: %w", err)
	}

	return orderID, nil
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

	query := `SELECT c.id AS cid, c.user_id AS usid,c.product_id AS pid ,c.unit,p.amount,p.discount  FROM cart  c JOIN product_models p on c.product_id = p.id JOIN users u ON c.user_id=u.id WHERE u.id = $1`
	//rows, err := r.sql.QueryContext(ctx, query, username)
	rows, err := r.sql.Query(query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to execute select query: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var cress model.Cartresponse
		err := rows.Scan(&cress.Cid, &cress.Usid, &cress.Pid, &cress.Unit, &cress.Amount, &cress.Discount)
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

	var query = `SELECT id AS cid, code, expiry, CURRENT_DATE AS current_date, CASE WHEN TO_DATE(expiry, 'DD/MM/YYYY') < CURRENT_DATE THEN true ELSE false END AS is_expired, 
	CASE WHEN $1 > min_amount THEN true ELSE false END AS is_eligible, min_amount, amount, used FROM coupon WHERE id = $2;`

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
	query := `SELECT product_id, user_id, unit FROM cart WHERE id = $1`
	var cart model.Cart
	err := r.sql.QueryRowContext(ctx, query, cartId).Scan(&cart.Productid, &cart.Userid, &cart.Unit)
	if err != nil {
		return model.Cart{}, fmt.Errorf("failed to get cart: %w", err)
	}
	return cart, nil
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
	query := `SELECT firstname, lastname, email, password FROM users WHERE email = $1`
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

func (r *repository) Listing(ctx context.Context) ([]model.ProductList, error) {
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
			vendor ON product_models.vendor_id = vendor.id WHERE product_models.units > 0;`

	rows, err := r.sql.QueryContext(ctx, query)
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
func (r *repository) LatestListing(ctx context.Context) ([]model.ProductList, error) {
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
			vendor ON product_models.vendor_id = vendor.id WHERE product_models.units > 0 ORDER BY product_models.id DESC;`

	rows, err := r.sql.QueryContext(ctx, query)
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
func (r *repository) PhighListing(ctx context.Context) ([]model.ProductList, error) {
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
			vendor ON product_models.vendor_id = vendor.id WHERE product_models.units > 0 ORDER BY product_models.amount DESC;`

	rows, err := r.sql.QueryContext(ctx, query)
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
func (r *repository) PlowListing(ctx context.Context) ([]model.ProductList, error) {
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
			vendor ON product_models.vendor_id = vendor.id WHERE product_models.units > 0 ORDER BY product_models.amount ASC;`

	rows, err := r.sql.QueryContext(ctx, query)
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
func (r *repository) InAZListing(ctx context.Context) ([]model.ProductList, error) {
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
			vendor ON product_models.vendor_id = vendor.id WHERE product_models.units > 0 ORDER BY  LOWER(product_models.name);`

	rows, err := r.sql.QueryContext(ctx, query)
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
func (r *repository) InZAListing(ctx context.Context) ([]model.ProductList, error) {
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
			vendor ON product_models.vendor_id = vendor.id WHERE product_models.units > 0 ORDER BY  LOWER(product_models.name) DESC;`

	rows, err := r.sql.QueryContext(ctx, query)
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
