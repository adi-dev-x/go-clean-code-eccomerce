package admin

import (
	"context"
	"fmt"
	"myproject/pkg/model"

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
	PhighListing(ctx context.Context) ([]model.Coupon, error)
	PlowListing(ctx context.Context, id string) ([]model.ProductList, error)
	InAZListing(ctx context.Context, id string) ([]model.ProductList, error)
	InZAListing(ctx context.Context, id string) ([]model.ProductList, error)
	Deletecoupon(ctx context.Context, id string) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
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
func (s *service) PhighListing(ctx context.Context) ([]model.Coupon, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return s.repo.PhighListing(ctx)
	}
}
func (s *service) PlowListing(ctx context.Context, id string) ([]model.ProductList, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return s.repo.PlowListing(ctx, id)
	}
}
func (s *service) InAZListing(ctx context.Context, id string) ([]model.ProductList, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return s.repo.InAZListing(ctx, id)
	}
}

func (s *service) InZAListing(ctx context.Context, id string) ([]model.ProductList, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return s.repo.InZAListing(ctx, id)
	}
}
