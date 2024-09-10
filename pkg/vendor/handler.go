package vendor

import (
	// "encoding/json"
	"fmt"
	"time"

	// db "myproject/pkg/database"
	services "myproject/pkg/client"
	db "myproject/pkg/database"
	"myproject/pkg/model"

	"net/http"
	"regexp"

	// "time"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	service  Service
	services services.Services
	adminjw  VendorJWT
}

func NewHandler(service Service, srv services.Services, adTK VendorJWT) *Handler {

	return &Handler{
		service:  service,
		services: srv,
		adminjw:  adTK,
	}
}
func (h *Handler) MountRoutes(engine *echo.Echo) {
	//applicantApi := engine.Group(basePath)
	applicantApi := engine.Group("/vendor")
	applicantApi.POST("/register", h.Register)
	applicantApi.POST("/login", h.Login)
	applicantApi.POST("/OtpLogin", h.OtpLogin)
	applicantApi.Use(h.adminjw.VendorAuthMiddleware())
	{
		applicantApi.POST("/UpdateProduct", h.UpdateProduct)
		applicantApi.POST("/AddProduct", h.AddProduct)
		applicantApi.GET("/listing/", h.Listing)
		applicantApi.POST("/Categorylisting", h.CategoryListing)
		applicantApi.GET("/Latestlisting/", h.LatestListing)
		applicantApi.GET("/PhighListing/", h.PhighListing)
		applicantApi.GET("/PlowListing/", h.PlowListing)
		applicantApi.GET("/InAZListing/", h.InAZListing)
		applicantApi.GET("/InZAListing/", h.InZAListing)
		/// list orders
		applicantApi.POST("/SalesReport", h.SalesReport)
		applicantApi.GET("/listAllOrders", h.ListAllOrders)
		applicantApi.GET("/listReturnedOrders", h.ListReturnedOrders)
		applicantApi.GET("/listFailedOrders", h.ListFailedOrders)
		applicantApi.GET("/listCompletedOrders", h.ListCompletedOrders)
		applicantApi.GET("/listPendingOrders", h.ListPendingOrders)
		applicantApi.POST("/cancelItem", h.CancelItem)

	}
}

func (h *Handler) respondWithError(c echo.Context, code int, msg interface{}) error {
	resp := map[string]interface{}{
		"msg": msg,
	}

	return c.JSON(code, resp)
}

func (h *Handler) respondWithData(c echo.Context, code int, message interface{}, data interface{}) error {
	resp := map[string]interface{}{
		"msg":  message,
		"data": data,
	}
	return c.JSON(code, resp)
}

func (h *Handler) UpdateProduct(c echo.Context) error {

	fmt.Println("this is in the handler UpdateProduct")
	var request model.UpdateProduct
	if err := c.Bind(&request); err != nil {
		return h.respondWithError(c, http.StatusBadRequest, map[string]string{"request-parse": err.Error()})
	}
	ErrVal := request.Valid()
	if len(ErrVal) > 0 {
		return h.respondWithError(c, http.StatusInternalServerError, map[string]interface{}{"invalid-request": ErrVal})
	}
	authHeader := c.Request().Header.Get("Authorization")
	fmt.Println("inside the cart list ", authHeader)
	username := c.Get("username").(string)
	// Validate request fields
	//errVal := request.Valid()

	ctx := c.Request().Context()
	if err := h.service.UpdateProduct(ctx, request, username); err != nil {
		return h.respondWithError(c, http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	fmt.Println("this is in the handler UpdateProduct")

	return h.respondWithData(c, http.StatusOK, "success", nil)
}

// / orders///
func (h *Handler) CancelItem(c echo.Context) error {
	fmt.Println("this is in the handler ReturnItem")
	authHeader := c.Request().Header.Get("Authorization")
	fmt.Println("inside the cart list ", authHeader)
	username := c.Get("username").(string)
	fmt.Println("inside the cart list ", username)
	ctx := c.Request().Context()

	var request model.ReturnOrderPost
	if err := c.Bind(&request); err != nil {
		return h.respondWithError(c, http.StatusBadRequest, map[string]string{"request-parse": err.Error()})
	}
	errValues := request.Valid()
	if len(errValues) > 0 {
		return h.respondWithError(c, http.StatusInternalServerError, map[string]interface{}{"invalid-request": errValues})
	}
	err := h.service.ReturnItem(ctx, request, username)
	if err != nil {
		return h.respondWithError(c, http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return h.respondWithData(c, http.StatusOK, "success", nil)
}
func (h *Handler) ListPendingOrders(c echo.Context) error {
	fmt.Println("this is in the handler ListAllOrders")
	authHeader := c.Request().Header.Get("Authorization")
	fmt.Println("inside the cart list ", authHeader)
	username := c.Get("username").(string)
	ctx := c.Request().Context()
	orders, err := h.service.ListPendingOrders(ctx, username)
	if err != nil {
		return h.respondWithError(c, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch orders", "details": err.Error()})
	}
	fmt.Println("this is the data ", orders)

	return h.respondWithData(c, http.StatusOK, "success", orders)
}
func (h *Handler) ListCompletedOrders(c echo.Context) error {
	fmt.Println("this is in the handler ListAllOrders")
	authHeader := c.Request().Header.Get("Authorization")
	fmt.Println("inside the cart list ", authHeader)
	username := c.Get("username").(string)
	ctx := c.Request().Context()
	orders, err := h.service.ListCompletedOrders(ctx, username)
	if err != nil {
		return h.respondWithError(c, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch orders", "details": err.Error()})
	}
	fmt.Println("this is the data ", orders)

	return h.respondWithData(c, http.StatusOK, "success", orders)
}
func (h *Handler) SalesReport(c echo.Context) error {
	fmt.Println("this is in the handler ListAllOrders")
	authHeader := c.Request().Header.Get("Authorization")
	fmt.Println("inside the cart list ", authHeader)
	username := c.Get("username").(string)
	var request model.SalesReport
	if err := c.Bind(&request); err != nil {
		return h.respondWithError(c, http.StatusBadRequest, map[string]string{"request-parse": err.Error()})
	}
	ErrVal := request.Valid()
	if len(ErrVal) > 0 {
		return h.respondWithError(c, http.StatusInternalServerError, map[string]interface{}{"invalid-request": ErrVal})
	}

	ctx := c.Request().Context()
	fmt.Println(ctx, username)

	orders, err := h.service.SalesReport(ctx, username, request)
	if err != nil {
		return h.respondWithError(c, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch orders", "details": err.Error()})
	}
	fmt.Println("this is the data ", orders)

	return h.respondWithData(c, http.StatusOK, "success", orders)
}
func (h *Handler) ListFailedOrders(c echo.Context) error {
	fmt.Println("this is in the handler ListAllOrders")
	authHeader := c.Request().Header.Get("Authorization")
	fmt.Println("inside the cart list ", authHeader)
	username := c.Get("username").(string)
	ctx := c.Request().Context()
	orders, err := h.service.ListFailedOrders(ctx, username)
	if err != nil {
		return h.respondWithError(c, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch orders", "details": err.Error()})
	}
	fmt.Println("this is the data ", orders)

	return h.respondWithData(c, http.StatusOK, "success", orders)
}
func (h *Handler) ListAllOrders(c echo.Context) error {
	fmt.Println("this is in the handler ListAllOrders")
	authHeader := c.Request().Header.Get("Authorization")
	fmt.Println("inside the cart list ", authHeader)
	username := c.Get("username").(string)
	ctx := c.Request().Context()
	orders, err := h.service.ListAllOrders(ctx, username)
	if err != nil {
		return h.respondWithError(c, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch orders", "details": err.Error()})
	}
	fmt.Println("this is the data ", orders)

	return h.respondWithData(c, http.StatusOK, "success", orders)
}
func (h *Handler) ListReturnedOrders(c echo.Context) error {
	fmt.Println("this is in the handler ListAllOrders")
	authHeader := c.Request().Header.Get("Authorization")
	fmt.Println("inside the cart list ", authHeader)
	username := c.Get("username").(string)
	ctx := c.Request().Context()
	orders, err := h.service.ListReturnedOrders(ctx, username)
	if err != nil {
		return h.respondWithError(c, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch orders", "details": err.Error()})
	}
	fmt.Println("this is the data ", orders)

	return h.respondWithData(c, http.StatusOK, "success", orders)
}

// ////
func (h *Handler) AddProduct(c echo.Context) error {
	// Parse request body into VendorRegisterRequest
	fmt.Println("this is in the handler AddProduct")
	var request model.Product
	if err := c.Bind(&request); err != nil {
		return h.respondWithError(c, http.StatusBadRequest, map[string]string{"request-parse": err.Error()})
	}

	// Validate request fields
	authHeader := c.Request().Header.Get("Authorization")
	fmt.Println("inside the cart list ", authHeader)
	username := c.Get("username").(string)

	ctx := c.Request().Context()
	if err := h.service.AddProduct(ctx, request, username); err != nil {
		return h.respondWithError(c, http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	errVal := request.Valid()
	if len(errVal) > 0 {
		return h.respondWithError(c, http.StatusBadRequest, map[string]interface{}{"invalid-request": errVal})
	}
	fmt.Println("this is in the handler Register")

	return h.respondWithData(c, http.StatusOK, "success", nil)
}
func (h *Handler) Register(c echo.Context) error {
	// Parse request body into VendorRegisterRequest
	fmt.Println("this is in the handler Register")
	var request model.VendorRegisterRequest
	if err := c.Bind(&request); err != nil {
		return h.respondWithError(c, http.StatusBadRequest, map[string]string{"request-parse": err.Error()})
	}

	// Validate request fields
	errVal := request.Valid()
	if len(errVal) > 0 {
		return h.respondWithError(c, http.StatusBadRequest, map[string]interface{}{"invalid-request": errVal})
	}

	ctx := c.Request().Context()
	if err := h.service.Register(ctx, request); err != nil {
		return h.respondWithError(c, http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	fmt.Println("this is in the handler Register")

	otp, err := h.services.SendEmailWithOTP(request.Email)
	if err != nil {
		return h.respondWithError(c, http.StatusInternalServerError, map[string]string{"error": "error in sending otp"})

	}
	err = db.SetRedis(request.Email, otp, time.Minute*5)
	if err != nil {
		return h.respondWithError(c, http.StatusInternalServerError, map[string]string{"error": "error in saving otp"})

	}
	storedData, err := db.GetRedis(request.Email)
	fmt.Println("this is the keyy!!!!!", storedData)

	return h.respondWithData(c, http.StatusOK, "success", nil)
}
func (h *Handler) Login(c echo.Context) error {
	// Parse request body into VendorRegisterRequest
	fmt.Println("this is in the handler Register")
	var request model.VendorLoginRequest
	if err := c.Bind(&request); err != nil {
		return h.respondWithError(c, http.StatusBadRequest, map[string]string{"request-parse": err.Error()})
	}

	ctx := c.Request().Context()
	if err := h.service.Login(ctx, request); err != nil {
		return h.respondWithError(c, http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	fmt.Println("this is in the handler Register")
	fmt.Println("this is in the handler Register")
	token, err := h.adminjw.GenerateVendorToken(request.Email)
	if err != nil {
		return h.respondWithError(c, http.StatusInternalServerError, map[string]string{"token-generation": err.Error()})
	}

	fmt.Println("User logged in successfully")
	return h.respondWithData(c, http.StatusOK, "success", map[string]string{"token": token})
}
func (h *Handler) OtpLogin(c echo.Context) error {
	// Parse request body into VendorRegisterRequest
	fmt.Println("this is in the handler OtpLogin")
	var request model.UserOtp

	if err := c.Bind(&request); err != nil {
		return h.respondWithError(c, http.StatusBadRequest, map[string]string{"request-parse": err.Error()})
	}
	fmt.Println("this is request", request)

	// Respond with success
	storedData, err := db.GetRedis(request.Email)
	fmt.Println("this is the keyy!!!!!", storedData, err)
	if storedData != request.Otp {
		return h.respondWithError(c, http.StatusInternalServerError, map[string]string{"error": "wrong otp"})

	}
	ctx := c.Request().Context()
	h.service.VerifyOtp(ctx, request.Email)
	return h.respondWithData(c, http.StatusOK, "success", nil)
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

func (h *Handler) Listing(c echo.Context) error {
	ctx := c.Request().Context()
	authHeader := c.Request().Header.Get("Authorization")
	fmt.Println("inside the cart list ", authHeader)
	id := c.Get("username").(string)
	fmt.Println("this is the id from the paramsssss !!! ", id)

	products, err := h.service.Listing(ctx, id)
	if err != nil {
		return h.respondWithError(c, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch products", "details": err.Error()})
	}
	fmt.Println("this is the data ", products)
	return h.respondWithData(c, http.StatusOK, "success", products)
}
func (h *Handler) CategoryListing(c echo.Context) error {
	ctx := c.Request().Context()
	authHeader := c.Request().Header.Get("Authorization")
	fmt.Println("inside the cart list ", authHeader)
	id := c.Get("username").(string)
	type Cat struct {
		Category string `json:"category"`
	}
	var request Cat

	if err := c.Bind(&request); err != nil {
		return h.respondWithError(c, http.StatusBadRequest, map[string]string{"request-parse": err.Error()})
	}
	if request.Category == "" {
		return h.respondWithError(c, http.StatusInternalServerError, map[string]string{"error": "Enter a valid value"})
	}
	products, err := h.service.CategoryListing(ctx, request.Category, id)
	if err != nil {
		return h.respondWithError(c, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch products", "details": err.Error()})
	}
	fmt.Println("this is the data ", products)
	return h.respondWithData(c, http.StatusOK, "success", products)
}
func (h *Handler) LatestListing(c echo.Context) error {
	ctx := c.Request().Context()
	authHeader := c.Request().Header.Get("Authorization")
	fmt.Println("inside the cart list ", authHeader)
	id := c.Get("username").(string)
	products, err := h.service.LatestListing(ctx, id)
	if err != nil {
		return h.respondWithError(c, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch products", "details": err.Error()})
	}
	fmt.Println("this is the data ", products)
	return h.respondWithData(c, http.StatusOK, "success", products)
}
func (h *Handler) PhighListing(c echo.Context) error {
	ctx := c.Request().Context()
	authHeader := c.Request().Header.Get("Authorization")
	fmt.Println("inside the cart list ", authHeader)
	id := c.Get("username").(string)
	products, err := h.service.PhighListing(ctx, id)
	if err != nil {
		return h.respondWithError(c, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch products", "details": err.Error()})
	}
	fmt.Println("this is the data ", products)
	return h.respondWithData(c, http.StatusOK, "success", products)
}
func (h *Handler) PlowListing(c echo.Context) error {
	ctx := c.Request().Context()
	authHeader := c.Request().Header.Get("Authorization")
	fmt.Println("inside the cart list ", authHeader)
	id := c.Get("username").(string)
	products, err := h.service.PlowListing(ctx, id)
	if err != nil {
		return h.respondWithError(c, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch products", "details": err.Error()})
	}
	fmt.Println("this is the data ", products)
	return h.respondWithData(c, http.StatusOK, "success", products)
}
func (h *Handler) InAZListing(c echo.Context) error {
	ctx := c.Request().Context()
	authHeader := c.Request().Header.Get("Authorization")
	fmt.Println("inside the cart list ", authHeader)
	id := c.Get("username").(string)
	products, err := h.service.InAZListing(ctx, id)
	if err != nil {
		return h.respondWithError(c, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch products", "details": err.Error()})
	}
	fmt.Println("this is the data ", products)
	return h.respondWithData(c, http.StatusOK, "success", products)
}
func (h *Handler) InZAListing(c echo.Context) error {
	ctx := c.Request().Context()
	authHeader := c.Request().Header.Get("Authorization")
	fmt.Println("inside the cart list ", authHeader)
	id := c.Get("username").(string)
	products, err := h.service.InZAListing(ctx, id)
	if err != nil {
		return h.respondWithError(c, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch products", "details": err.Error()})
	}
	fmt.Println("this is the data ", products)
	return h.respondWithData(c, http.StatusOK, "success", products)
}
