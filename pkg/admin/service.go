package admin

import (
	"context"
	"fmt"
	services "myproject/pkg/client"
	"myproject/pkg/model"
	"reflect"
	"sync"
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
	GetMainOrders(ctx context.Context, username string, orderUid string) ([]model.ListAllOrdersUsers, error)
	Deletecoupon(ctx context.Context, id string) error

	///product listing
	ProductListing(ctx context.Context) ([]model.ProductListingUsers, error)
	PlowListing(ctx context.Context, id string) ([]model.ProductListingUsers, error)
	InAZListing(ctx context.Context, id string) ([]model.ProductListingUsers, error)
	InZAListing(ctx context.Context, id string) ([]model.ProductListingUsers, error)
	ListingSingle(ctx context.Context, id string) ([]model.ProductListDetailed, error)
	PhighListing(ctx context.Context, id string) ([]model.ProductListingUsers, error)
	CategoryListing(ctx context.Context, category string) ([]model.ProductListingUsers, error)
	//BrandListing
	BrandListing(ctx context.Context, category string) ([]model.ProductListingUsers, error)
	BestSellingListingProductCategory(ctx context.Context, category string) ([]model.ProductListingUsers, error)
	BestSellingListingProduct(ctx context.Context) ([]model.ProductListingUsers, error)

	BestSellingListingCategory(ctx context.Context) ([]string, error)
	BestSellingListingBrand(ctx context.Context) ([]string, error)

	BestSellingListingProductBrand(ctx context.Context, category string) ([]model.ProductListingUsers, error)
	///orders
	/// Singlevendor
	SalesReportSinglevendor(ctx context.Context, id string, request model.SalesReport) (model.SendSalesReortVendorinAdmin, error)
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

	ListMainOrders(ctx context.Context) ([]model.ListingMainOrders, error)
	UpdateOrderDate(ctx context.Context, request model.UpdateOrderAdmin) error
	ReturnItem(ctx context.Context, request model.ReturnOrderPost, username string) error
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
func (s *service) UpdateOrderDate(ctx context.Context, request model.UpdateOrderAdmin) error {
	fmt.Println("updating UpdateOrderDate ")
	var delivery bool
	var deliverydate string
	data, err := s.repo.GetOrderForUpdating(ctx, request.Oid)
	R_D_Pay_st, _ := data["status"].(string)
	R_D_Deli_st, _ := data["delivered"].(bool)
	R_D_Deli_date, _ := data["delivery_date"].(string)

	fmt.Println("Type of x:", reflect.TypeOf(R_D_Pay_st))

	if err != nil {
		return err
	}
	fmt.Println("d", data["payment_method"])

	if data["payment_method"] == "ONLINE" && request.Payment_status != "" {
		return fmt.Errorf("Can not change payment status of Online payment")

	}
	if request.Payment_status == "" {
		request.Payment_status = R_D_Pay_st
	}
	if request.Delivery_Stat == "" {
		delivery = R_D_Deli_st
	} else if request.Delivery_Stat == "Deliverd" {
		delivery = true
	} else {
		delivery = false

	}
	if request.Delivery_date == "" {
		deliverydate = R_D_Deli_date
	} else {
		deliverydate = request.Delivery_date

	}
	s.repo.UpdateOrderFromAdminUP(ctx, request.Oid, deliverydate, request.Payment_status, delivery)

	return nil
}
func (s *service) ReturnItem(ctx context.Context, request model.ReturnOrderPost, username string) error {

	p, err := s.repo.GetSingleItem(ctx, request.Oid)
	if err != nil {
		return fmt.Errorf("entered is wrong id", err)
	}
	fmt.Println("this is the single order ", p)
	switch p.Returned {
	case "Returned":
		return fmt.Errorf("this item is already returned")

	case "Cancelled":
		return fmt.Errorf("this item is payment Cancelled")

	}

	switch p.Status {
	case "Failed":
		return fmt.Errorf("this item is payment failed")

	case "Cancelled":
		return fmt.Errorf("this item is payment Cancelled")

	}
	var w sync.WaitGroup

	VErr0 := make(chan error, 1)
	VErr := make(chan error, 1)
	VErr2 := make(chan error, 1)
	VErr3 := make(chan error, 1)
	VErr4 := make(chan error, 1)
	w.Add(5)
	go func() {
		defer w.Done()
		err := s.services.SendOrderReturnConfirmationEmailVendor(p.Name, p.Amount, p.Unit, username)
		VErr0 <- err
	}()
	go func() {
		defer w.Done()
		err := s.services.SendOrderReturnConfirmationEmailVendor(p.Name, p.Amount, p.Unit, username)
		VErr <- err
	}()
	go func() {
		defer w.Done()
		err := s.repo.IncreaseStock(ctx, p.Pid, p.Unit)
		VErr2 <- err
	}()
	go func() {
		defer w.Done()
		err := s.repo.UpdateOiStatus(ctx, request.Oid, "Cancelled")
		VErr3 <- err
	}()
	go func() {
		defer w.Done()
		var err error
		var wallet_id string
		if p.Status == "Completed" {
			fmt.Println("in 1st if")
			//f := p.Amount * float64(p.Unit)
			// value := []interface{}{p.Amount, id, "Credit"}
			cmp_r_amt, _ := s.repo.GetcpAmtRefund(ctx, p.Moid)
			fmt.Println("---", cmp_r_amt)
			var f float64
			var upcmCheck string
			if p.Amount*float64(p.Unit) > float64(cmp_r_amt) {
				f = p.Amount*float64(p.Unit) - float64(cmp_r_amt)
				upcmCheck = "normal"
			} else {
				f = p.Amount * float64(p.Unit)
				upcmCheck = "notnormal"
			}
			wallet_id, err = s.repo.CreditWallet(ctx, p.Usid, f)
			if wallet_id != "" {
				value := []interface{}{f, wallet_id, "Credit", p.Usid, p.Moid}
				er := s.repo.UpdateWalletTransaction(ctx, value)
				if er != nil {
					fmt.Println("there is erorrrr in wallet transaction")
				} else {
					///updateing the mo statussssss
					if upcmCheck == "normal" {
						s.repo.ChangeCouponRefundStatus(ctx, p.Moid)
						fmt.Println(" normalll")
					}

				}

				fmt.Println("this is workingggg ist")
			}
		} else {
			fmt.Println("in 1st else")
			wallet_id = ""
			err = nil
		}
		VErr4 <- err

	}()
	go func() {
		w.Wait()
		close(VErr0)
		close(VErr)
		close(VErr2)
		close(VErr3)
		close(VErr4)
	}()
	if err := <-VErr0; err != nil {
		return fmt.Errorf("failed to send order  return  email: %w", err)
	}
	if err := <-VErr; err != nil {
		return fmt.Errorf("failed to send order  return  email: %w", err)
	}
	if err := <-VErr2; err != nil {
		return fmt.Errorf("failed to update unit: %w", err)
	}
	if err := <-VErr3; err != nil {
		return fmt.Errorf("failed to update to redund status: %w", err)
	}
	if err := <-VErr4; err != nil {
		return fmt.Errorf("failed to update to redund status: %w", err)
	}
	// er := s.repo.ChangeOrderStatus(ctx, p.Moid)
	// fmt.Println(er)
	// if errM != nil {
	// 	return fmt.Errorf("error in updating ")
	// }

	return nil

}

// //All orders
func (s *service) ListMainOrders(ctx context.Context) ([]model.ListingMainOrders, error) {

	orders, err := s.repo.PrintingUserMainOrder(ctx)
	if err != nil {
		return []model.ListingMainOrders{}, fmt.Errorf("error in retriving data")
	}
	return orders, nil

}
func (s *service) SalesReport(ctx context.Context, request model.SalesReport) (model.SendSalesReortAdmin, error) {

	// Parse dates
	const dateFormat = "2006-01-02"
	var startDate, endDate time.Time
	var err error
	var orders []model.ListOrdersAdmin
	var rangers string
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
		rangers = request.From + " - " + request.To + " Custom"
	}
	if request.Type == "Yearly" {
		currentTime := time.Now()
		formattedDate := fmt.Sprintf("01/01/%d - %02d/%02d/%d", currentTime.Year(), currentTime.Day(), currentTime.Month(), currentTime.Year())
		fmt.Println("Formatted Date:", formattedDate)
		rangers = formattedDate + " Yearly"
		orders, err = s.repo.SalesReportOrdersYearly(ctx)
		if err != nil {
			return model.SendSalesReortAdmin{}, fmt.Errorf("error in receiving yearly data: %w", err)
		}

	}

	if request.Type == "Monthly" {
		currentTime := time.Now()

		formattedDate := fmt.Sprintf("01/%02d/%d - %02d/%02d/%d",
			currentTime.Month(),
			currentTime.Year(),
			currentTime.Day(),
			currentTime.Month(),
			currentTime.Year(),
		)

		fmt.Println("Formatted Date:", formattedDate)
		rangers = formattedDate + " Monthly"
		orders, err = s.repo.SalesReportOrdersMonthly(ctx)
		if err != nil {
			return model.SendSalesReortAdmin{}, fmt.Errorf("error in receiving monthly data: %w", err)
		}

	}

	if request.Type == "Weekly" {
		currentTime := time.Now()

		startOfWeek := currentTime.AddDate(0, 0, -int(currentTime.Weekday()-1))

		formattedDate := fmt.Sprintf("%02d/%02d/%d - %02d/%02d/%d",
			startOfWeek.Day(), startOfWeek.Month(), startOfWeek.Year(),
			currentTime.Day(), currentTime.Month(), currentTime.Year())

		fmt.Println("Formatted Date:", formattedDate)
		rangers = formattedDate + " Weekly"
		orders, err = s.repo.SalesReportOrdersWeekly(ctx)
		if err != nil {
			return model.SendSalesReortAdmin{}, fmt.Errorf("error in receiving weekly data: %w", err)
		}

	}

	if request.Type == "Daily" {
		currentTime := time.Now()
		currentDate := currentTime.Format("2006-01-02")
		fmt.Println("Current Date:", currentDate)
		rangers = currentDate + " Daily"
		orders, err = s.repo.SalesReportOrdersDaily(ctx)
		if err != nil {
			return model.SendSalesReortAdmin{}, fmt.Errorf("error in receiving daily data: %w", err)
		}

	}
	fmt.Println("this is the dattaaaa!!!!!", orders)

	var results []model.ResultsAdminsalesReport
	for _, order := range orders {
		result := model.ResultsAdminsalesReport{
			Name:     order.Name,
			Unit:     order.Unit,
			Amount:   order.Amount,
			Date:     order.Date,
			Oid:      order.Oid,
			VName:    order.VName,
			Discount: order.Discount,
			Cmt:      order.CouponAmt,
			Code:     order.CouponCode,
		}
		results = append(results, result)
	}
	fmt.Println(results, rangers)
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
	fmt.Println("this is the slFact", slFact)
	//slFact.TotalSales = slFact.TotalSales * 0.02
	name, err := s.services.GenerateDailySalesReportExcelAdmin(orders, slFact, request.Type, "", rangers)
	fmt.Println(name, "@@@@@@@", err)
	excelfurl := "https://adiecom.gitfunswokhu.in/view/" + name
	pname, err := s.services.GenerateDailySalesReportPDFAdmin(orders, slFact, request.Type, "AdminPDF", rangers)
	fmt.Println(pname, "@@@@@@@", err)
	pdffurl := "https://adiecom.gitfunswokhu.in/view/" + pname

	fmt.Println(excelfurl, "  --  ", pdffurl)
	var data model.SendSalesReortAdmin
	data.Data = results
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
func (s *service) GetMainOrders(ctx context.Context, username string, orderUid string) ([]model.ListAllOrdersUsers, error) {

	/// check if order id is valid

	orders, err := s.repo.PrintingUserSingleMainOrderCollection(ctx, orderUid)
	fmt.Println(orders)
	if err != nil {
		return []model.ListAllOrdersUsers{}, fmt.Errorf("error in retriving data")
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

func (s *service) SalesReportSinglevendor(ctx context.Context, id string, request model.SalesReport) (model.SendSalesReortVendorinAdmin, error) {

	const dateFormat = "2006-01-02"
	var startDate, endDate time.Time
	var err error
	var orders []model.ListOrdersVendor
	var rangers string
	if request.Type == "Custom" {
		startDate, err = time.Parse(dateFormat, request.From)
		if err != nil {
			return model.SendSalesReortVendorinAdmin{}, fmt.Errorf("invalid From date format: %w", err)
		}
		endDate, err = time.Parse(dateFormat, request.To)
		if err != nil {
			return model.SendSalesReortVendorinAdmin{}, fmt.Errorf("invalid To date format: %w", err)
		}
		orders, err = s.repo.SalesReportOrdersCustomSinglevendor(ctx, startDate, endDate, id)
		fmt.Println("in data!!", orders)
		if err != nil {
			return model.SendSalesReortVendorinAdmin{}, fmt.Errorf("error in receiving data")
		}
		rangers = request.From + " - " + request.To + " Custom"
	}
	if request.Type == "Yearly" {
		currentTime := time.Now()
		formattedDate := fmt.Sprintf("01/01/%d - %02d/%02d/%d", currentTime.Year(), currentTime.Day(), currentTime.Month(), currentTime.Year())
		fmt.Println("Formatted Date:", formattedDate)
		rangers = formattedDate + " Yearly"
		orders, err = s.repo.SalesReportOrdersYearlySinglevendor(ctx, id)
		if err != nil {
			return model.SendSalesReortVendorinAdmin{}, fmt.Errorf("error in receiving yearly data: %w", err)
		}

	}

	if request.Type == "Monthly" {
		currentTime := time.Now()

		formattedDate := fmt.Sprintf("01/%02d/%d - %02d/%02d/%d",
			currentTime.Month(),
			currentTime.Year(),
			currentTime.Day(),
			currentTime.Month(),
			currentTime.Year(),
		)

		fmt.Println("Formatted Date:", formattedDate)
		rangers = formattedDate + " Monthly"
		orders, err = s.repo.SalesReportOrdersMonthlySinglevendor(ctx, id)
		if err != nil {
			return model.SendSalesReortVendorinAdmin{}, fmt.Errorf("error in receiving monthly data: %w", err)
		}

	}

	if request.Type == "Weekly" {
		currentTime := time.Now()

		startOfWeek := currentTime.AddDate(0, 0, -int(currentTime.Weekday()-1))

		formattedDate := fmt.Sprintf("%02d/%02d/%d - %02d/%02d/%d",
			startOfWeek.Day(), startOfWeek.Month(), startOfWeek.Year(),
			currentTime.Day(), currentTime.Month(), currentTime.Year())

		fmt.Println("Formatted Date:", formattedDate)
		rangers = formattedDate + " Weekly"
		orders, err = s.repo.SalesReportOrdersWeeklySinglevendor(ctx, id)
		if err != nil {
			return model.SendSalesReortVendorinAdmin{}, fmt.Errorf("error in receiving weekly data: %w", err)
		}

	}

	if request.Type == "Daily" {
		currentTime := time.Now()
		currentDate := currentTime.Format("2006-01-02")
		fmt.Println("Current Date:", currentDate)
		rangers = currentDate + " Daily"
		orders, err = s.repo.SalesReportOrdersDailySinglevendor(ctx, id)
		if err != nil {
			return model.SendSalesReortVendorinAdmin{}, fmt.Errorf("error in receiving daily data: %w", err)
		}

	}

	fmt.Println("this is the dattaaaa!!!!!", rangers)
	// Call repository function
	var results []model.ResultsVendorsales
	for _, order := range orders {
		result := model.ResultsVendorsales{
			Name:   order.Name,
			Unit:   order.Unit,
			Amount: order.Amount,
			Date:   order.Date,
			Oid:    order.Oid,

			Discount: order.Discount,
			Cmt:      order.CouponAmt,
			Code:     order.CouponCode,
		}
		results = append(results, result)
	}
	salesFacts, err := s.repo.GetSalesFactByDateSinglevendor(ctx, request.Type, startDate, endDate, id)
	if err != nil {
		return model.SendSalesReortVendorinAdmin{}, fmt.Errorf("failed to get sales facts: %w", err)
	}
	fmt.Println("valueeee in salesFact!!!", salesFacts)
	if salesFacts == nil {
		return model.SendSalesReortVendorinAdmin{}, nil
	}
	slFact := salesFacts[0]
	vdata, err := s.repo.GetVendorDetails(ctx, id)
	if err != nil {
		return model.SendSalesReortVendorinAdmin{}, fmt.Errorf("error in receiving yearly data: %w", err)
	}
	fmt.Println("listing vdataaaaa!!!", vdata)
	name, err := s.services.GenerateDailySalesReportExcelAdminside(orders, slFact, request.Type, "Admin_EXCEL", rangers, vdata[0].Name, vdata[0].Email, vdata[0].Gst)
	fmt.Println(name, "@@@@@@@", err)
	excelfurl := "https://adiecom.gitfunswokhu.in/view/" + name
	pname, err := s.services.GenerateDailySalesReportPDFAdminside(orders, slFact, request.Type, "Admin_PDF", rangers, vdata[0].Name, vdata[0].Email, vdata[0].Gst)
	fmt.Println(pname, "@@@@@@@", err)
	pdffurl := "https://adiecom.gitfunswokhu.in/view/" + pname

	fmt.Println(excelfurl, "  --  ", pdffurl)
	var data model.SendSalesReortVendorinAdmin
	data.Data = results
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

// ///end Singlevendor ending BrandListing
func (s *service) CategoryListing(ctx context.Context, category string) ([]model.ProductListingUsers, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return s.repo.CategoryListing(ctx, category)
	}
}
func (s *service) BrandListing(ctx context.Context, category string) ([]model.ProductListingUsers, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return s.repo.BrandListing(ctx, category)
	}
}
func (s *service) BestSellingListingProductBrand(ctx context.Context, category string) ([]model.ProductListingUsers, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return s.repo.BestSellingListingProductBrand(ctx, category)
	}
}
func (s *service) BestSellingListingProductCategory(ctx context.Context, category string) ([]model.ProductListingUsers, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return s.repo.BestSellingListingProductCategory(ctx, category)
	}
}
func (s *service) BestSellingListingProduct(ctx context.Context) ([]model.ProductListingUsers, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return s.repo.BestSellingListingProduct(ctx)
	}
}
func (s *service) BestSellingListingCategory(ctx context.Context) ([]string, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return s.repo.BestSellingListingCategory(ctx)
	}
}
func (s *service) BestSellingListingBrand(ctx context.Context) ([]string, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return s.repo.BestSellingListingBrand(ctx)
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
