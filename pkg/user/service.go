package user

import (
	"context"
	"fmt"
	"myproject/pkg/config"
	"myproject/pkg/model"
	"regexp"
	"strings"

	// "github.com/go-resty/resty/v2"
	//"github.com/go-playground/validator/v10/translations/id"
	"github.com/razorpay/razorpay-go"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Service interface {
	Register(ctx context.Context, request model.UserRegisterRequest) error
	Login(ctx context.Context, request model.UserLoginRequest) error
	Listing(ctx context.Context) ([]model.ProductList, error)
	LatestListing(ctx context.Context) ([]model.ProductList, error)
	PhighListing(ctx context.Context) ([]model.ProductList, error)
	PlowListing(ctx context.Context) ([]model.ProductList, error)
	InAZListing(ctx context.Context) ([]model.ProductList, error)
	InZAListing(ctx context.Context) ([]model.ProductList, error)
	OtpLogin(ctx context.Context, request model.UserOtp) error
	UpdateUser(ctx context.Context, updatedData model.UserRegisterRequest) error
	AddTocart(ctx context.Context, request model.Cart) error
	AddToWish(ctx context.Context, request model.Wishlist) error
	AddAddress(ctx context.Context, request model.Address, username string) error
	AddToorder(ctx context.Context, request model.Order) (model.RZpayment, error)
	Listcart(ctx context.Context, id string) ([]model.Usercartview, error)
	ListWish(ctx context.Context, id string) ([]model.UserWishview, error)
	ListAddress(ctx context.Context, username string) ([]model.Address, error)
	ActiveListing(ctx context.Context) ([]model.Coupon, error)
}
type service struct {
	repo   Repository
	Config config.Config
}

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}
func (s *service) ListAddress(ctx context.Context, username string) ([]model.Address, error) {
	id := s.repo.Getid(ctx, username)
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return s.repo.ListAddress(ctx, id)
	}

}
func (s *service) ActiveListing(ctx context.Context) ([]model.Coupon, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return s.repo.ActiveListing(ctx)
	}

}
func (s *service) AddTocart(ctx context.Context, request model.Cart) error {
	if request.Productid == "" || request.Unit == 0 || request.Userid == "" {
		fmt.Println("this is in the service error value missing")
		return fmt.Errorf("missing values")
	}

	products, err := s.repo.ListingByid(ctx, request.Productid)
	if err != nil {
		fmt.Println("failed to get product:", err)
		return fmt.Errorf("failed to get product: %w", err)
	}

	if len(products) == 0 {
		fmt.Println("product not found")
		return fmt.Errorf("product not found")
	}

	product := products[0]
	if product.Unit < float64(request.Unit) {
		fmt.Println("not enough units available")
		return fmt.Errorf("not enough units available")
	}

	// Add to cart
	return s.repo.AddTocart(ctx, request)
}
func (s *service) AddToWish(ctx context.Context, request model.Wishlist) error {
	if request.Productid == "" || request.Userid == "" {
		fmt.Println("this is in the service error value missing")
		return fmt.Errorf("missing values")
	}

	_, err := s.repo.ListingByid(ctx, request.Productid)
	if err != nil {
		fmt.Println("failed to get product:", err)
		return fmt.Errorf("failed to get product: %w", err)
	}

	return s.repo.AddToWish(ctx, request)
}
func (s *service) AddAddress(ctx context.Context, request model.Address, username string) error {
	// if request.Address == "" {
	// 	fmt.Println("this is in the service error value missing")
	// 	return fmt.Errorf("missing values")
	// }
	id := s.repo.Getid(ctx, username)
	fmt.Println("this is my idddddddd!!!!", id)

	return s.repo.AddAddress(ctx, request, id)
}
func (s *service) AddToorder(ctx context.Context, request model.Order) (model.RZpayment, error) {
	//ctx context.Context, request model.Order
	fmt.Println("inside the service Addto order ")
	username := ctx.Value("username").(string)
	fmt.Println("inside the cart list ", username)
	fiData, _ := s.repo.GetorderDetails(ctx, request)
	fmt.Println("thisss is the daaa ", fiData.Data.Data, "and this is amount ", fiData.TAmount)
	if !fiData.Notvalid {
		return model.RZpayment{}, fmt.Errorf("not a valid coupon")
	}
	var k model.RZpayment
	k = s.PayGateway(ctx, fiData.TAmount)
	fmt.Println("this is kkkkkkkk!!!!!", k)
	// Pid, err := s.repo.AddToPayment(ctx, request, fiData, status, username)
	// oid, err := s.repo.AddToOrderDetails(ctx, request, fiData, status, username, Pid)
	//fmt.Println("this is the status ", status)
	return k, nil
}
func (s *service) CheckOut(ctx context.Context, request model.Order) (model.RZpayment, error) {

	fmt.Println("inside the service corrected Addto order ")
	username := ctx.Value("username").(string)
	fmt.Println("inside the cart list ", username)
	// cartDataChan := make(chan model.CartresponseData)
	// amountChan := make(chan int)
	go func() {

	}()

	fiData, _ := s.repo.GetorderDetails(ctx, request)
	fmt.Println("thisss is the daaa ", fiData.Data.Data, "and this is amount ", fiData.TAmount)
	if !fiData.Notvalid {
		return model.RZpayment{}, fmt.Errorf("not a valid coupon")
	}
	var k model.RZpayment
	k = s.PayGateway(ctx, fiData.TAmount)
	fmt.Println("this is kkkkkkkk!!!!!", k)
	// Pid, err := s.repo.AddToPayment(ctx, request, fiData, status, username)
	// oid, err := s.repo.AddToOrderDetails(ctx, request, fiData, status, username, Pid)
	//fmt.Println("this is the status ", status)
	return k, nil
}

func (s *service) PayGateway(ctx context.Context, amt int) model.RZpayment {

	//Razor_ID := s.Config.Razor_ID
	//Razor_Secret := s.Config.Razor_SECRET
	Razor_ID := "rzp_test_mRydipg2bgDZmQ"
	Razor_Secret := "a2oY1G5RYIQh9gH04KWATpnx"
	fmt.Println(" RZ keyyyyy ", Razor_ID, "  RZ ZSCERET ", Razor_Secret)
	razorpayClient := razorpay.NewClient(Razor_ID, Razor_Secret)
	data := map[string]interface{}{
		"amount":   amt,
		"currency": "INR",
		"receipt":  "some_receipt_id",
	}
	body, _ := razorpayClient.Order.Create(data, nil)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"Error": "Faied to create razorpay orer"})
	// }

	// Extract the order ID from the response body returned by the RazorPay API
	//fmt.Println(" in bodyyyyyyy ", body)
	value := body["id"]
	fmt.Println(" in bodyyyyyyy id!!!!!", value)
	str := value.(string)

	// homepagevariables := PageVariable{
	// 	AppointmentID: str,
	// }
	//return c.Render(http.StatusOK, "payment.html", map[string]interface{}{
	// "invoiceID":     id,
	// "totalPrice":    amountInPaisa / 100,
	// "total":         amountInPaisa,
	// "appointmentID": homepagevariables.AppointmentID,
	// })

	// return "success"
	var p = model.RZpayment{
		Id:  str,
		Amt: amt,
	}
	fmt.Println("ggg passs on RZPAy", p)

	return p
}

type PageVariable struct {
	AppointmentID string
}

func (s *service) Register(ctx context.Context, request model.UserRegisterRequest) error {
	var err error
	if request.FirstName == "" || request.Email == "" || request.Password == "" || request.Phone == "" {
		fmt.Println("this is in the service error value missing")
		err = fmt.Errorf("missing values")
		return err
	}
	if !isValidEmail(request.Email) {
		fmt.Println("this is in the service error invalid email")
		err = fmt.Errorf("invalid email")
		return err
	}
	if !isValidPhoneNumber(request.Phone) {
		fmt.Println("this is in the service error invalid phone number")
		err = fmt.Errorf("invalid phone number")
		return err
	}
	fmt.Println("this is the dataaa ", request.Email)
	existingUser, err := s.repo.Login(ctx, request.Email)
	fmt.Println("there may be a user", existingUser)
	if err != nil && err != gorm.ErrRecordNotFound {
		fmt.Println("this is in the service error checking existing user")
		err = fmt.Errorf("failed to check existing user: %w", err)
		return err
	}
	if existingUser.Email != "" {
		fmt.Println("this is in the service user already exists")
		err = fmt.Errorf("user already exists")
		return err
	}
	fmt.Println("this is in the service Register", request.Password)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("this is in the service error hashing password")
		err = fmt.Errorf("failed to hash password: %w", err)
		return err
	}
	request.Password = string(hashedPassword)
	fmt.Println("this is in the service Register", request.Password)
	return s.repo.Register(ctx, request)
}

func (s *service) Login(ctx context.Context, request model.UserLoginRequest) error {
	fmt.Println("this is in the service Login", request.Password)
	var err error
	if request.Email == "" || request.Password == "" {
		fmt.Println("this is in the service error value missing")
		err = fmt.Errorf("missing values")
		return err
	}
	storedUser, err := s.repo.Login(ctx, request.Email)
	fmt.Println("thisss is the dataaa ", storedUser)
	if err != nil {
		fmt.Println("this is in the service user not found")
		return fmt.Errorf("user not found: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(request.Password)); err != nil {
		fmt.Println("this is in the service incorrect password")
		return fmt.Errorf("incorrect password: %w", err)
	}

	return nil
}

func (s *service) OtpLogin(ctx context.Context, request model.UserOtp) error {
	fmt.Println("this is in the service Login", request.Otp)
	var err error
	if request.Email == "" || request.Otp == "" {
		fmt.Println("this is in the service error value missing")
		err = fmt.Errorf("missing values")
		return err
	}
	return nil
}

func (s *service) UpdateUser(ctx context.Context, updatedData model.UserRegisterRequest) error {
	var query string
	var args []interface{}

	query = "UPDATE users SET"

	if updatedData.FirstName != "" {
		query += " firstname = ?,"
		args = append(args, updatedData.FirstName)
	}
	if updatedData.LastName != "" {
		query += " lastname = ?,"
		args = append(args, updatedData.LastName)
	}
	if updatedData.Password != "" {
		// Hash the password before updating it
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(updatedData.Password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("failed to hash password: %w", err)
		}
		query += " password = ?,"
		args = append(args, string(hashedPassword))
	}
	if updatedData.Phone != "" && isValidPhoneNumber(updatedData.Phone) {
		query += " phone = ?,"
		args = append(args, updatedData.Phone)
	}

	query = strings.TrimSuffix(query, ",")

	query += " WHERE email = ?"
	args = append(args, updatedData.Email)
	fmt.Println("this is the UpdateUser ", query, " kkk ", args)

	return s.repo.UpdateUser(ctx, query, args)
}

func (s *service) Listing(ctx context.Context) ([]model.ProductList, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return s.repo.Listing(ctx)
	}
}
func (s *service) Listcart(ctx context.Context, id string) ([]model.Usercartview, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return s.repo.Listcart(ctx, id)
	}
}
func (s *service) ListWish(ctx context.Context, id string) ([]model.UserWishview, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return s.repo.ListWish(ctx, id)
	}
}
func (s *service) LatestListing(ctx context.Context) ([]model.ProductList, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return s.repo.LatestListing(ctx)
	}
}
func (s *service) PhighListing(ctx context.Context) ([]model.ProductList, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return s.repo.PhighListing(ctx)
	}
}
func (s *service) PlowListing(ctx context.Context) ([]model.ProductList, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return s.repo.PlowListing(ctx)
	}
}
func (s *service) InAZListing(ctx context.Context) ([]model.ProductList, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return s.repo.InAZListing(ctx)
	}
}

func (s *service) InZAListing(ctx context.Context) ([]model.ProductList, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return s.repo.InZAListing(ctx)
	}
}
func isValidEmail(email string) bool {
	// Simple regex pattern for basic email validation
	fmt.Println(" check email validity")
	const emailRegex = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(email)
}
func isValidPhoneNumber(phone string) bool {
	// Simple regex pattern for basic phone number validation
	fmt.Println(" check pfone validity")
	const phoneRegex = `^\+?[1-9]\d{1,14}$` // E.164 international phone number format
	re := regexp.MustCompile(phoneRegex)
	return re.MatchString(phone)
}
