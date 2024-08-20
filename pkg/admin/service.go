package admin

import (
	"context"
	"fmt"
	services "myproject/pkg/client"
	"myproject/pkg/model"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Service interface {
	Register(ctx context.Context, request model.AdminRegisterRequest) error
	Login(ctx context.Context, request model.AdminLoginRequest) error
	Listing(ctx context.Context) ([]model.Coupon, error)
	OtpLogin(ctx context.Context, request model.VendorOtp) error
	Addcoupon(ctx context.Context, request model.Coupon) error
	LatestListing(ctx context.Context) ([]model.Coupon, error)
	ActiveListing(ctx context.Context) ([]model.Coupon, error)

	Deletecoupon(ctx context.Context, id string) error

	///product listing
	ProductListing(ctx context.Context) ([]model.ProductListingUsers, error)
	PlowListing(ctx context.Context, id string) ([]model.ProductListingUsers, error)
	InAZListing(ctx context.Context, id string) ([]model.ProductListingUsers, error)
	InZAListing(ctx context.Context, id string) ([]model.ProductListingUsers, error)
	ListingSingle(ctx context.Context, id string) ([]model.ProductListDetailed, error)
	PhighListing(ctx context.Context, id string) ([]model.ProductListingUsers, error)
	CategoryListing(ctx context.Context, category string) ([]model.ProductListingUsers, error)

	///orders
	/// Singlevendor
	SalesReportSinglevendor(ctx context.Context, id string, request model.SalesReport) (model.SendSalesReort, error)
	ListAllOrdersSinglevendor(ctx context.Context, username string) ([]model.ListOrdersVendor, error)
	ListReturnedOrdersSinglevendor(ctx context.Context, username string) ([]model.ListOrdersVendor, error)
	ListFailedOrdersSinglevendor(ctx context.Context, username string) ([]model.ListOrdersVendor, error)
	ListCompletedOrdersSinglevendor(ctx context.Context, username string) ([]model.ListOrdersVendor, error)
	ListPendingOrdersSinglevendor(ctx context.Context, username string) ([]model.ListOrdersVendor, error)

	// all orders
	ListFailedOrders(ctx context.Context) ([]model.ListOrdersAdmin, error)
	ListAllOrders(ctx context.Context) ([]model.ListOrdersAdmin, error)
	ListReturnedOrders(ctx context.Context) ([]model.ListOrdersAdmin, error)
	ListCompletedOrders(ctx context.Context) ([]model.ListOrdersAdmin, error)
	ListPendingOrders(ctx context.Context) ([]model.ListOrdersAdmin, error)
	SalesReport(ctx context.Context, request model.SalesReport) (model.SendSalesReortAdmin, error)
}

type service struct {
	repo     Repository
	services services.Services
}

func NewService(repo Repository, services services.Services) Service {
	return &service{
		repo:     repo,
		services: services,
	}
}

// //All orders
func (s *service) SalesReport(ctx context.Context, request model.SalesReport) (model.SendSalesReortAdmin, error) {

	// Parse dates
	const dateFormat = "2006-01-02"
	var startDate, endDate time.Time
	var err error
	var orders []model.ListOrdersAdmin
	if request.Type == "Custom" {
		startDate, err = time.Parse(dateFormat, request.From)
		if err != nil {
			return model.SendSalesReortAdmin{}, fmt.Errorf("invalid From date format: %w", err)
		}
		endDate, err = time.Parse(dateFormat, request.To)
		if err != nil {
			return model.SendSalesReortAdmin{}, fmt.Errorf("invalid To date format: %w", err)
		}
		orders, err = s.repo.SalesReportOrdersCustom(ctx, startDate, endDate)
		fmt.Println("in data!!", orders)
		if err != nil {
			return model.SendSalesReortAdmin{}, fmt.Errorf("error in receiving data")
		}
	}
	if request.Type == "Yearly" {
		orders, err = s.repo.SalesReportOrdersYearly(ctx)
		if err != nil {
			return model.SendSalesReortAdmin{}, fmt.Errorf("error in receiving yearly data: %w", err)
		}

	}

	if request.Type == "Monthly" {
		orders, err = s.repo.SalesReportOrdersMonthly(ctx)
		if err != nil {
			return model.SendSalesReortAdmin{}, fmt.Errorf("error in receiving monthly data: %w", err)
		}

	}

	if request.Type == "Weekly" {
		orders, err = s.repo.SalesReportOrdersWeekly(ctx)
		if err != nil {
			return model.SendSalesReortAdmin{}, fmt.Errorf("error in receiving weekly data: %w", err)
		}

	}

	if request.Type == "Daily" {
		orders, err = s.repo.SalesReportOrdersDaily(ctx)
		if err != nil {
			return model.SendSalesReortAdmin{}, fmt.Errorf("error in receiving daily data: %w", err)
		}

	}
	fmt.Println("this is the dattaaaa!!!!!", orders)
	// Call repository function
	salesFacts, err := s.repo.GetSalesFactByDate(ctx, request.Type, startDate, endDate)
	if err != nil {
		return model.SendSalesReortAdmin{}, fmt.Errorf("failed to get sales facts: %w", err)
	}
	fmt.Println("valueeee in salesFact!!!", salesFacts)
	if salesFacts == nil {
		return model.SendSalesReortAdmin{}, nil
	}
	slFact := salesFacts[0]
	//slFact.TotalSales = slFact.TotalSales * 0.02
	name, err := s.services.GenerateDailySalesReportExcelAdmin(orders, slFact, request.Type, "")
	fmt.Println(name, "@@@@@@@", err)
	excelfurl := "http://localhost:8081/" + name
	pname, err := s.services.GenerateDailySalesReportPDFAdmin(orders, slFact, request.Type, "AdminPDF")
	fmt.Println(pname, "@@@@@@@", err)
	pdffurl := "http://localhost:8081/" + pname

	fmt.Println(excelfurl, "  --  ", pdffurl)
	var data model.SendSalesReortAdmin
	data.Data = orders
	data.FactsData = slFact
	data.ExcelUrl = excelfurl
	data.PdfUrl = pdffurl

	return data, nil

}
func (s *service) ListPendingOrders(ctx context.Context) ([]model.ListOrdersAdmin, error) {

	fmt.Println("inside the ListAllOrders ")

	orders, err := s.repo.ListPendingOrders(ctx)
	if err != nil {
		return []model.ListOrdersAdmin{}, fmt.Errorf("this is the error for listing all orders", err)
	}

	return orders, nil
}
func (s *service) ListCompletedOrders(ctx context.Context) ([]model.ListOrdersAdmin, error) {

	fmt.Println("inside the ListAllOrders ")

	orders, err := s.repo.ListCompletedOrders(ctx)
	if err != nil {
		return []model.ListOrdersAdmin{}, fmt.Errorf("this is the error for listing all orders", err)
	}

	return orders, nil
}
func (s *service) ListFailedOrders(ctx context.Context) ([]model.ListOrdersAdmin, error) {

	fmt.Println("inside the ListAllOrders ")

	orders, err := s.repo.ListFailedOrders(ctx)
	if err != nil {
		return []model.ListOrdersAdmin{}, fmt.Errorf("this is the error for listing all orders", err)
	}

	return orders, nil
}
func (s *service) ListReturnedOrders(ctx context.Context) ([]model.ListOrdersAdmin, error) {

	fmt.Println("inside the ListAllOrders ")

	orders, err := s.repo.ListReturnedOrders(ctx)
	if err != nil {
		return []model.ListOrdersAdmin{}, fmt.Errorf("this is the error for listing all orders", err)
	}

	return orders, nil
}
func (s *service) ListAllOrders(ctx context.Context) ([]model.ListOrdersAdmin, error) {

	fmt.Println("inside the ListAllOrders ")

	orders, err := s.repo.ListAllOrders(ctx)
	if err != nil {
		return []model.ListOrdersAdmin{}, fmt.Errorf("this is the error for listing all orders", err)
	}

	return orders, nil
}

// //list Singlevendor begining

func (s *service) SalesReportSinglevendor(ctx context.Context, id string, request model.SalesReport) (model.SendSalesReort, error) {

	// Parse dates
	const dateFormat = "2006-01-02"
	var startDate, endDate time.Time
	var err error
	var orders []model.ListOrdersVendor
	if request.Type == "Custom" {
		startDate, err = time.Parse(dateFormat, request.From)
		if err != nil {
			return model.SendSalesReort{}, fmt.Errorf("invalid From date format: %w", err)
		}
		endDate, err = time.Parse(dateFormat, request.To)
		if err != nil {
			return model.SendSalesReort{}, fmt.Errorf("invalid To date format: %w", err)
		}
		orders, err = s.repo.SalesReportOrdersCustomSinglevendor(ctx, startDate, endDate, id)
		fmt.Println("in data!!", orders)
		if err != nil {
			return model.SendSalesReort{}, fmt.Errorf("error in receiving data")
		}
	}
	if request.Type == "Yearly" {
		orders, err = s.repo.SalesReportOrdersYearlySinglevendor(ctx, id)
		if err != nil {
			return model.SendSalesReort{}, fmt.Errorf("error in receiving yearly data: %w", err)
		}

	}

	if request.Type == "Monthly" {
		orders, err = s.repo.SalesReportOrdersMonthlySinglevendor(ctx, id)
		if err != nil {
			return model.SendSalesReort{}, fmt.Errorf("error in receiving monthly data: %w", err)
		}

	}

	if request.Type == "Weekly" {
		orders, err = s.repo.SalesReportOrdersWeeklySinglevendor(ctx, id)
		if err != nil {
			return model.SendSalesReort{}, fmt.Errorf("error in receiving weekly data: %w", err)
		}

	}

	if request.Type == "Daily" {
		orders, err = s.repo.SalesReportOrdersDailySinglevendor(ctx, id)
		if err != nil {
			return model.SendSalesReort{}, fmt.Errorf("error in receiving daily data: %w", err)
		}

	}
	fmt.Println("this is the dattaaaa!!!!!", orders)
	// Call repository function
	salesFacts, err := s.repo.GetSalesFactByDateSinglevendor(ctx, request.Type, startDate, endDate, id)
	if err != nil {
		return model.SendSalesReort{}, fmt.Errorf("failed to get sales facts: %w", err)
	}
	fmt.Println("valueeee in salesFact!!!", salesFacts)
	if salesFacts == nil {
		return model.SendSalesReort{}, nil
	}
	slFact := salesFacts[0]
	name, err := s.services.GenerateDailySalesReportExcel(orders, slFact, request.Type, "Admin_EXCEL")
	fmt.Println(name, "@@@@@@@", err)
	excelfurl := "http://localhost:8081/" + name
	pname, err := s.services.GenerateDailySalesReportPDF(orders, slFact, request.Type, "Admin_PDF")
	fmt.Println(pname, "@@@@@@@", err)
	pdffurl := "http://localhost:8081/" + pname

	fmt.Println(excelfurl, "  --  ", pdffurl)
	var data model.SendSalesReort
	data.Data = orders
	data.FactsData = slFact
	data.ExcelUrl = excelfurl
	data.PdfUrl = pdffurl

	return data, nil

}
func (s *service) ListPendingOrdersSinglevendor(ctx context.Context, id string) ([]model.ListOrdersVendor, error) {

	fmt.Println("inside the ListAllOrders ", id)

	orders, err := s.repo.ListPendingOrdersSinglevendor(ctx, id)
	if err != nil {
		return []model.ListOrdersVendor{}, fmt.Errorf("this is the error for listing all orders", err)
	}

	return orders, nil
}
func (s *service) ListFailedOrdersSinglevendor(ctx context.Context, id string) ([]model.ListOrdersVendor, error) {

	fmt.Println("inside the ListAllOrders ", id)

	orders, err := s.repo.ListFailedOrdersSinglevendor(ctx, id)
	if err != nil {
		return []model.ListOrdersVendor{}, fmt.Errorf("this is the error for listing all orders", err)
	}

	return orders, nil
}
func (s *service) ListCompletedOrdersSinglevendor(ctx context.Context, id string) ([]model.ListOrdersVendor, error) {

	fmt.Println("inside the ListAllOrders ", id)

	orders, err := s.repo.ListCompletedOrdersSinglevendor(ctx, id)
	if err != nil {
		return []model.ListOrdersVendor{}, fmt.Errorf("this is the error for listing all orders", err)
	}

	return orders, nil
}
func (s *service) ListReturnedOrdersSinglevendor(ctx context.Context, id string) ([]model.ListOrdersVendor, error) {

	fmt.Println("inside the ListAllOrders ", id)

	orders, err := s.repo.ListReturnedOrdersSinglevendor(ctx, id)
	if err != nil {
		return []model.ListOrdersVendor{}, fmt.Errorf("this is the error for listing all orders", err)
	}

	return orders, nil
}
func (s *service) ListAllOrdersSinglevendor(ctx context.Context, id string) ([]model.ListOrdersVendor, error) {

	fmt.Println("inside the ListAllOrders ", id)

	orders, err := s.repo.ListAllOrdersSinglevendor(ctx, id)
	if err != nil {
		return []model.ListOrdersVendor{}, fmt.Errorf("this is the error for listing all orders", err)
	}

	return orders, nil
}

// ///end Singlevendor ending
func (s *service) CategoryListing(ctx context.Context, category string) ([]model.ProductListingUsers, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return s.repo.CategoryListing(ctx, category)
	}
}

func (s *service) ListingSingle(ctx context.Context, id string) ([]model.ProductListDetailed, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return s.repo.ListingSingle(ctx, id)
	}
}
func (s *service) ProductListing(ctx context.Context) ([]model.ProductListingUsers, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return s.repo.ProductListing(ctx)
	}
}
func (s *service) Register(ctx context.Context, request model.AdminRegisterRequest) error {
	var err error
	if request.Name == "" || request.Email == "" || request.Password == "" || request.Phone == "" {
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
func (s *service) Addcoupon(ctx context.Context, request model.Coupon) error {
	var err error
	fmt.Println("this is the product data ", request)

	if request.Code == "" || request.Expiry == "" {
		err = fmt.Errorf("missing values")
		return err
	}

	// Validate units, tax, and price (consider edge cases)
	if request.Minamount < 0 {
		err = fmt.Errorf("invalid unit value: must be non-negative")
		return err
		//return fmt.Errorf("invalid unit value: must be non-negative")
	}
	if request.Amount < 0 {
		err = fmt.Errorf("invalid tax value: must be between 0 and 1 (inclusive)")
		return err

	}

	return s.repo.Addcoupon(ctx, request)
}
func (s *service) Deletecoupon(ctx context.Context, id string) error {
	cid, err := s.repo.GetCoupnExist(ctx, id)

	if err != nil {
		return fmt.Errorf("there is error in finding the coupon")
	}
	errs := s.repo.Deletecoupon(ctx, cid)
	if errs != nil {
		return fmt.Errorf("error in deleting")
	}
	return nil

}
func (s *service) Login(ctx context.Context, request model.AdminLoginRequest) error {
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

	//return s.repo.Login(ctx, request)
}
func (s *service) OtpLogin(ctx context.Context, request model.VendorOtp) error {
	fmt.Println("this is in the service Login", request.Otp)
	var err error
	if request.Email == "" || request.Otp == "" {
		fmt.Println("this is in the service error value missing")
		err = fmt.Errorf("missing values")
		return err
	}
	// storedUser, err := s.repo.Login(ctx, request.Email)
	// if err != nil {
	// 	fmt.Println("this is in the service user not found")
	// 	return fmt.Errorf("user not found: %w", err)
	// }

	return nil

	//return s.repo.Login(ctx, request)
}

// func (s *service) Listing(ctx context.Context) ([]model.Product, error) {

//		return s.repo.Listing(ctx)
//	}
func (s *service) Listing(ctx context.Context) ([]model.Coupon, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return s.repo.Listing(ctx)
	}
}
func (s *service) LatestListing(ctx context.Context) ([]model.Coupon, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return s.repo.LatestListing(ctx)
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

// PhighListing
func (s *service) PhighListing(ctx context.Context, id string) ([]model.ProductListingUsers, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return s.repo.PhighListing(ctx)
	}
}
func (s *service) PlowListing(ctx context.Context, id string) ([]model.ProductListingUsers, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return s.repo.PlowListing(ctx)
	}
}
func (s *service) InAZListing(ctx context.Context, id string) ([]model.ProductListingUsers, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return s.repo.InAZListing(ctx)
	}
}

func (s *service) InZAListing(ctx context.Context, id string) ([]model.ProductListingUsers, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return s.repo.InZAListing(ctx)
	}
}
