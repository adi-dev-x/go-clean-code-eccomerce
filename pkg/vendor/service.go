package vendor

import (
	"context"
	"fmt"
	services "myproject/pkg/client"
	"myproject/pkg/config"
	"myproject/pkg/model"
	"sync"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Service interface {
	Register(ctx context.Context, request model.VendorRegisterRequest) error
	Login(ctx context.Context, request model.VendorLoginRequest) error
	Listing(ctx context.Context, id string) ([]model.ProductList, error)
	OtpLogin(ctx context.Context, request model.VendorOtp) error
	AddProduct(ctx context.Context, request model.Product, username string) error
	LatestListing(ctx context.Context, id string) ([]model.ProductList, error)
	PhighListing(ctx context.Context, id string) ([]model.ProductList, error)
	PlowListing(ctx context.Context, id string) ([]model.ProductList, error)
	InAZListing(ctx context.Context, id string) ([]model.ProductList, error)
	InZAListing(ctx context.Context, id string) ([]model.ProductList, error)

	///listing orders
	ListAllOrders(ctx context.Context, username string) ([]model.ListOrdersVendor, error)
	ListReturnedOrders(ctx context.Context, username string) ([]model.ListOrdersVendor, error)
	ListFailedOrders(ctx context.Context, username string) ([]model.ListOrdersVendor, error)
	ListCompletedOrders(ctx context.Context, username string) ([]model.ListOrdersVendor, error)
	ListPendingOrders(ctx context.Context, username string) ([]model.ListOrdersVendor, error)

	//returning
	ReturnItem(ctx context.Context, request model.ReturnOrderPost, username string) error
}

type service struct {
	repo     Repository
	Config   config.Config
	services services.Services
}

func NewService(repo Repository, services services.Services) Service {
	return &service{
		repo:     repo,
		services: services,
	}
}

// /returning
func (s *service) ReturnItem(ctx context.Context, request model.ReturnOrderPost, username string) error {
	id := s.repo.Getid(ctx, username)
	p, err := s.repo.GetSingleItem(ctx, id, request.Oid)
	if err != nil {
		return fmt.Errorf("entered is wrong id", err)
	}
	fmt.Println("this is the single order ", p)
	if p.Returned {
		return fmt.Errorf("this item is already returned")

	}

	if p.Status == "Failed" {
		return fmt.Errorf("this item is payment failed")

	}
	var w sync.WaitGroup

	VErr := make(chan error, 1)
	VErr2 := make(chan error, 1)
	VErr3 := make(chan error, 1)
	VErr4 := make(chan error, 1)
	w.Add(4)
	go func() {
		defer w.Done()
		err := s.services.SendOrderReturnConfirmationEmailVendor(p.Name, p.Amount, p.Unit, username)
		VErr <- err
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
		err := s.repo.UpdateOiStatus(ctx, request.Oid)
		VErr3 <- err
	}()
	go func() {
		defer w.Done()
		var err error
		var wallet_id string
		if p.Status == "Completed" {
			fmt.Println("in 1st if")
			// value := []interface{}{p.Amount, id, "Credit"}
			wallet_id, err = s.repo.CreditWallet(ctx, p.Usid, p.Amount)
			if wallet_id != "" {
				value := []interface{}{p.Amount, wallet_id, "Credit", p.Usid}
				er := s.repo.UpdateWalletTransaction(ctx, value)
				if er != nil {
					fmt.Println("there is erorrrr in wallet transaction")
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
		close(VErr)
		close(VErr2)
		close(VErr3)
		close(VErr4)
	}()
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

	return nil

}

// list alll oreders
func (s *service) ListAllOrders(ctx context.Context, username string) ([]model.ListOrdersVendor, error) {
	id := s.repo.Getid(ctx, username)
	fmt.Println("inside the ListAllOrders ", id)

	orders, err := s.repo.ListAllOrders(ctx, id)
	if err != nil {
		return []model.ListOrdersVendor{}, fmt.Errorf("this is the error for listing all orders", err)
	}

	return orders, nil
}
func (s *service) ListReturnedOrders(ctx context.Context, username string) ([]model.ListOrdersVendor, error) {
	id := s.repo.Getid(ctx, username)
	fmt.Println("inside the ListAllOrders ", id)

	orders, err := s.repo.ListReturnedOrders(ctx, id)
	if err != nil {
		return []model.ListOrdersVendor{}, fmt.Errorf("this is the error for listing all orders", err)
	}

	return orders, nil
}
func (s *service) ListFailedOrders(ctx context.Context, username string) ([]model.ListOrdersVendor, error) {
	id := s.repo.Getid(ctx, username)
	fmt.Println("inside the ListAllOrders ", id)

	orders, err := s.repo.ListFailedOrders(ctx, id)
	if err != nil {
		return []model.ListOrdersVendor{}, fmt.Errorf("this is the error for listing all orders", err)
	}

	return orders, nil
}
func (s *service) ListCompletedOrders(ctx context.Context, username string) ([]model.ListOrdersVendor, error) {
	id := s.repo.Getid(ctx, username)
	fmt.Println("inside the ListAllOrders ", id)

	orders, err := s.repo.ListCompletedOrders(ctx, id)
	if err != nil {
		return []model.ListOrdersVendor{}, fmt.Errorf("this is the error for listing all orders", err)
	}

	return orders, nil
}
func (s *service) ListPendingOrders(ctx context.Context, username string) ([]model.ListOrdersVendor, error) {
	id := s.repo.Getid(ctx, username)
	fmt.Println("inside the ListAllOrders ", id)

	orders, err := s.repo.ListPendingOrders(ctx, id)
	if err != nil {
		return []model.ListOrdersVendor{}, fmt.Errorf("this is the error for listing all orders", err)
	}

	return orders, nil
}

//

func (s *service) Register(ctx context.Context, request model.VendorRegisterRequest) error {
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
func (s *service) AddProduct(ctx context.Context, request model.Product, username string) error {

	fmt.Println("this is the product data ", request)
	id := s.repo.Getid(ctx, username)
	request.Vendorid = id

	err := s.repo.AddProduct(ctx, request)
	if err != nil {
		return fmt.Errorf("error in validating query")
	}
	return nil
}
func (s *service) Login(ctx context.Context, request model.VendorLoginRequest) error {
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
func (s *service) Listing(ctx context.Context, id string) ([]model.ProductList, error) {
	d := s.repo.Getid(ctx, id)
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return s.repo.Listing(ctx, d)
	}
}
func (s *service) LatestListing(ctx context.Context, id string) ([]model.ProductList, error) {
	d := s.repo.Getid(ctx, id)
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return s.repo.LatestListing(ctx, d)
	}
}
func (s *service) PhighListing(ctx context.Context, id string) ([]model.ProductList, error) {
	d := s.repo.Getid(ctx, id)
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return s.repo.PhighListing(ctx, d)
	}
}
func (s *service) PlowListing(ctx context.Context, id string) ([]model.ProductList, error) {
	d := s.repo.Getid(ctx, id)
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return s.repo.PlowListing(ctx, d)
	}
}
func (s *service) InAZListing(ctx context.Context, id string) ([]model.ProductList, error) {
	d := s.repo.Getid(ctx, id)
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return s.repo.InAZListing(ctx, d)
	}
}

func (s *service) InZAListing(ctx context.Context, id string) ([]model.ProductList, error) {
	d := s.repo.Getid(ctx, id)
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return s.repo.InZAListing(ctx, d)
	}
}
