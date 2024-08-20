package admin

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
	adminjw  AdminJWT
}

func NewHandler(service Service, srv services.Services, adTK AdminJWT) *Handler {

	return &Handler{
		service:  service,
		services: srv,
		adminjw:  adTK,
	}
}
func (h *Handler) MountRoutes(engine *echo.Echo) {
	//applicantApi := engine.Group(basePath)
	applicantApi := engine.Group("/admin")
	applicantApi.POST("/register", h.Register)
	applicantApi.POST("/login", h.Login)
	applicantApi.POST("/OtpLogin", h.OtpLogin)
	applicantApi.Use(h.adminjw.AdminAuthMiddleware())
	{
		applicantApi.POST("/Deletecoupon", h.Deletecoupon)
		applicantApi.POST("/Addcoupon", h.Addcoupon)
		applicantApi.GET("/listing", h.Listing)
		applicantApi.GET("/Latestlisting", h.LatestListing)
		applicantApi.GET("/ActiveListing", h.ActiveListing)

		//list vendors

		applicantApi.GET("/Vendorlisting", h.VendorListing)

		//Product listing
		applicantApi.GET("/Productlisting", h.ProductListing)

		applicantApi.GET("/InAZListing", h.InAZListing)
		applicantApi.GET("/InZAListing", h.InZAListing)
		applicantApi.GET("/PhighListing", h.PhighListing)
		applicantApi.GET("/PlowListing", h.PlowListing)
		applicantApi.GET("/listingSingleProduct/:id", h.ListingSingle)
		applicantApi.POST("/Categorylisting", h.CategoryListing)

		//applicantApi.GET("/ProductActiveListing", h.ProductActiveListing)

		///list orders
		//of a particular vendor   Singlevendor
		applicantApi.POST("/listAllOrdersSinglevendor", h.ListAllOrdersSinglevendor)
		applicantApi.POST("/listReturnedOrdersSinglevendor", h.ListReturnedOrdersSinglevendor)
		applicantApi.POST("/listFailedOrdersSinglevendor", h.ListFailedOrdersSinglevendor)
		applicantApi.POST("/listCompletedOrdersSinglevendor", h.ListCompletedOrdersSinglevendor)
		applicantApi.POST("/listPendingOrdersSinglevendor", h.ListPendingOrdersSinglevendor)
		applicantApi.POST("/SalesReportSinglevendor", h.SalesReportSinglevendor)

		//// ending particular vendor
		applicantApi.GET("/listAllOrders", h.ListAllOrders)
		applicantApi.GET("/listReturnedOrders", h.ListReturnedOrders)
		applicantApi.GET("/listFailedOrders", h.ListFailedOrders)
		applicantApi.GET("/listCompletedOrders", h.ListCompletedOrders)
		applicantApi.GET("/listPendingOrders", h.ListPendingOrders)
		applicantApi.POST("/SalesReport", h.SalesReport)

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
func (h *Handler) VendorListing(c echo.Context) error {

	return nil
}

// / orders
func (h *Handler) SalesReport(c echo.Context) error {
	fmt.Println("this is in the handler ListAllOrders")

	var request model.SalesReport
	if err := c.Bind(&request); err != nil {
		return h.respondWithError(c, http.StatusBadRequest, map[string]string{"request-parse": err.Error()})
	}
	ErrVal := request.Valid()
	if len(ErrVal) > 0 {
		return h.respondWithError(c, http.StatusInternalServerError, map[string]interface{}{"invalid-request": ErrVal})
	}

	ctx := c.Request().Context()

	orders, err := h.service.SalesReport(ctx, request)
	if err != nil {
		return h.respondWithError(c, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch orders", "details": err.Error()})
	}
	fmt.Println("this is the data ", orders)

	return h.respondWithData(c, http.StatusOK, "success", orders)
}
func (h *Handler) ListPendingOrders(c echo.Context) error {
	fmt.Println("this is in the handler ListAllOrders")

	ctx := c.Request().Context()
	orders, err := h.service.ListPendingOrders(ctx)
	if err != nil {
		return h.respondWithError(c, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch orders", "details": err.Error()})
	}
	fmt.Println("this is the data ", orders)

	return h.respondWithData(c, http.StatusOK, "success", orders)
}
func (h *Handler) ListCompletedOrders(c echo.Context) error {
	fmt.Println("this is in the handler ListAllOrders")

	ctx := c.Request().Context()
	orders, err := h.service.ListCompletedOrders(ctx)
	if err != nil {
		return h.respondWithError(c, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch orders", "details": err.Error()})
	}
	fmt.Println("this is the data ", orders)

	return h.respondWithData(c, http.StatusOK, "success", orders)
}
func (h *Handler) ListFailedOrders(c echo.Context) error {
	fmt.Println("this is in the handler ListAllOrders")

	ctx := c.Request().Context()
	orders, err := h.service.ListFailedOrders(ctx)
	if err != nil {
		return h.respondWithError(c, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch orders", "details": err.Error()})
	}
	fmt.Println("this is the data ", orders)

	return h.respondWithData(c, http.StatusOK, "success", orders)
}
func (h *Handler) ListReturnedOrders(c echo.Context) error {
	fmt.Println("this is in the handler ListAllOrders")

	ctx := c.Request().Context()
	orders, err := h.service.ListReturnedOrders(ctx)
	if err != nil {
		return h.respondWithError(c, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch orders", "details": err.Error()})
	}
	fmt.Println("this is the data ", orders)

	return h.respondWithData(c, http.StatusOK, "success", orders)
}
func (h *Handler) ListAllOrders(c echo.Context) error {
	fmt.Println("this is in the handler ListAllOrders")

	ctx := c.Request().Context()
	orders, err := h.service.ListAllOrders(ctx)
	if err != nil {
		return h.respondWithError(c, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch orders", "details": err.Error()})
	}
	fmt.Println("this is the data ", orders)

	return h.respondWithData(c, http.StatusOK, "success", orders)
}

// /////  Singlevendor listing orders of particular vendor  begining
func (h *Handler) SalesReportSinglevendor(c echo.Context) error {
	fmt.Println("this is in the handler ListAllOrders")

	var request model.SalesReport
	if err := c.Bind(&request); err != nil {
		return h.respondWithError(c, http.StatusBadRequest, map[string]string{"request-parse": err.Error()})
	}
	ErrVal := request.Valid()
	if len(ErrVal) > 0 {
		return h.respondWithError(c, http.StatusInternalServerError, map[string]interface{}{"invalid-request": ErrVal})
	}

	ctx := c.Request().Context()
	if request.Vid == "" {
		return h.respondWithError(c, http.StatusInternalServerError, map[string]interface{}{"invalid-request": "Please enter Vendor id"})

	}

	orders, err := h.service.SalesReportSinglevendor(ctx, request.Vid, request)
	if err != nil {
		return h.respondWithError(c, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch orders", "details": err.Error()})
	}
	fmt.Println("this is the data ", orders)

	return h.respondWithData(c, http.StatusOK, "success", orders)
}

func (h *Handler) ListPendingOrdersSinglevendor(c echo.Context) error {
	fmt.Println("this is in the handler ListAllOrders")

	type getId struct {
		Id string `json:"v_id"`
	}
	var bb getId
	if err := c.Bind(&bb); err != nil {
		return h.respondWithError(c, http.StatusBadRequest, map[string]string{"request-parseerr": err.Error()})

	}
	username := bb.Id
	ctx := c.Request().Context()
	orders, err := h.service.ListPendingOrdersSinglevendor(ctx, username)
	if err != nil {
		return h.respondWithError(c, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch orders", "details": err.Error()})
	}
	fmt.Println("this is the data ", orders)

	return h.respondWithData(c, http.StatusOK, "success", orders)
}
func (h *Handler) ListCompletedOrdersSinglevendor(c echo.Context) error {
	fmt.Println("this is in the handler ListAllOrders")
	type getId struct {
		Id string `json:"v_id"`
	}
	var bb getId
	if err := c.Bind(&bb); err != nil {
		return h.respondWithError(c, http.StatusBadRequest, map[string]string{"request-parseerr": err.Error()})

	}
	username := bb.Id
	ctx := c.Request().Context()
	orders, err := h.service.ListCompletedOrdersSinglevendor(ctx, username)
	if err != nil {
		return h.respondWithError(c, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch orders", "details": err.Error()})
	}
	fmt.Println("this is the data ", orders)

	return h.respondWithData(c, http.StatusOK, "success", orders)
}
func (h *Handler) ListFailedOrdersSinglevendor(c echo.Context) error {
	fmt.Println("this is in the handler ListAllOrders")
	type getId struct {
		Id string `json:"v_id"`
	}
	var bb getId
	if err := c.Bind(&bb); err != nil {
		return h.respondWithError(c, http.StatusBadRequest, map[string]string{"request-parseerr": err.Error()})

	}
	username := bb.Id
	ctx := c.Request().Context()
	orders, err := h.service.ListFailedOrdersSinglevendor(ctx, username)
	if err != nil {
		return h.respondWithError(c, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch orders", "details": err.Error()})
	}
	fmt.Println("this is the data ", orders)

	return h.respondWithData(c, http.StatusOK, "success", orders)
}
func (h *Handler) ListReturnedOrdersSinglevendor(c echo.Context) error {
	fmt.Println("this is in the handler ListAllOrders")
	type getId struct {
		Id string `json:"v_id"`
	}
	var bb getId
	if err := c.Bind(&bb); err != nil {
		return h.respondWithError(c, http.StatusBadRequest, map[string]string{"request-parseerr": err.Error()})

	}
	username := bb.Id
	ctx := c.Request().Context()
	orders, err := h.service.ListReturnedOrdersSinglevendor(ctx, username)
	if err != nil {
		return h.respondWithError(c, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch orders", "details": err.Error()})
	}
	fmt.Println("this is the data ", orders)

	return h.respondWithData(c, http.StatusOK, "success", orders)
}
func (h *Handler) ListAllOrdersSinglevendor(c echo.Context) error {
	fmt.Println("this is in the handler ListAllOrders")
	type getId struct {
		Id string `json:"v_id"`
	}
	var bb getId
	if err := c.Bind(&bb); err != nil {
		return h.respondWithError(c, http.StatusBadRequest, map[string]string{"request-parseerr": err.Error()})

	}
	username := bb.Id
	ctx := c.Request().Context()
	orders, err := h.service.ListAllOrdersSinglevendor(ctx, username)
	if err != nil {
		return h.respondWithError(c, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch orders", "details": err.Error()})
	}
	fmt.Println("this is the data ", orders)

	return h.respondWithData(c, http.StatusOK, "success", orders)
}

///////////// listing orders of particular vendor ending

func (h *Handler) Addcoupon(c echo.Context) error {
	// Parse request body into VendorRegisterRequest
	fmt.Println("this is in the handler AddProduct")
	var request model.Coupon
	if err := c.Bind(&request); err != nil {
		return h.respondWithError(c, http.StatusBadRequest, map[string]string{"request-parse": err.Error()})
	}
	errVal := request.Valid()
	if len(errVal) > 0 {
		return h.respondWithError(c, http.StatusBadRequest, map[string]interface{}{"invalid-request": errVal})
	}
	// Validate request fields

	ctx := c.Request().Context()
	if err := h.service.Addcoupon(ctx, request); err != nil {
		fmt.Println("this is the error !!!!!", err.Error())

		return h.respondWithError(c, http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	fmt.Println("this is in the handler Register")

	return h.respondWithData(c, http.StatusOK, "success", nil)
}
func (h *Handler) Deletecoupon(c echo.Context) error {
	// Parse request body into VendorRegisterRequest
	fmt.Println("this is in the handler AddProduct")
	var request struct {
		id string `json:"id"`
	}
	if err := c.Bind(&request); err != nil {
		return h.respondWithError(c, http.StatusBadRequest, map[string]string{"request-parse": err.Error()})
	}

	if request.id == "" {
		return h.respondWithError(c, http.StatusBadRequest, map[string]interface{}{"invalid-request": "enter id"})
	}
	// Validate request fields

	ctx := c.Request().Context()
	if err := h.service.Deletecoupon(ctx, request.id); err != nil {
		fmt.Println("this is the error !!!!!", err.Error())

		return h.respondWithError(c, http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	fmt.Println("this is in the handler Register")

	return h.respondWithData(c, http.StatusOK, "success", nil)
}
func (h *Handler) Register(c echo.Context) error {
	// Parse request body into VendorRegisterRequest
	fmt.Println("this is in the handler Register")
	var request model.AdminRegisterRequest
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
	var request model.AdminLoginRequest
	if err := c.Bind(&request); err != nil {
		return h.respondWithError(c, http.StatusBadRequest, map[string]string{"request-parse": err.Error()})
	}

	ctx := c.Request().Context()
	if err := h.service.Login(ctx, request); err != nil {
		return h.respondWithError(c, http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	fmt.Println("this is in the handler Register")
	fmt.Println("this is in the handler Register")
	token, err := h.adminjw.GenerateAdminToken(request.Email)
	if err != nil {
		return h.respondWithError(c, http.StatusInternalServerError, map[string]string{"token-generation": err.Error()})
	}

	fmt.Println("User logged in successfully")
	return h.respondWithData(c, http.StatusOK, "success", map[string]string{"token": token})
}
func (h *Handler) OtpLogin(c echo.Context) error {
	// Parse request body into VendorRegisterRequest
	fmt.Println("this is in the handler OtpLogin")
	var request model.AdminOtp

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

	products, err := h.service.Listing(ctx)
	if err != nil {
		return h.respondWithError(c, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch products", "details": err.Error()})
	}
	fmt.Println("this is the data ", products)
	return h.respondWithData(c, http.StatusOK, "success", products)
}
func (h *Handler) LatestListing(c echo.Context) error {
	ctx := c.Request().Context()

	products, err := h.service.LatestListing(ctx)
	if err != nil {
		return h.respondWithError(c, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch products", "details": err.Error()})
	}
	fmt.Println("this is the data ", products)
	return h.respondWithData(c, http.StatusOK, "success", products)
}
func (h *Handler) ActiveListing(c echo.Context) error {
	fmt.Println("in activeeee")
	ctx := c.Request().Context()

	products, err := h.service.ActiveListing(ctx)
	if err != nil {
		return h.respondWithError(c, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch products", "details": err.Error()})
	}
	fmt.Println("this is the data ", products)
	return h.respondWithData(c, http.StatusOK, "success", products)
}
func (h *Handler) PlowListing(c echo.Context) error {
	ctx := c.Request().Context()
	id := ""
	products, err := h.service.PlowListing(ctx, id)
	if err != nil {
		return h.respondWithError(c, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch products", "details": err.Error()})
	}
	fmt.Println("this is the data ", products)
	return h.respondWithData(c, http.StatusOK, "success", products)
}
func (h *Handler) PhighListing(c echo.Context) error {
	ctx := c.Request().Context()
	id := ""
	products, err := h.service.PhighListing(ctx, id)
	if err != nil {
		return h.respondWithError(c, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch products", "details": err.Error()})
	}
	fmt.Println("this is the data ", products)
	return h.respondWithData(c, http.StatusOK, "success", products)
}
func (h *Handler) InAZListing(c echo.Context) error {
	ctx := c.Request().Context()
	id := ""
	products, err := h.service.InAZListing(ctx, id)
	if err != nil {
		return h.respondWithError(c, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch products", "details": err.Error()})
	}
	fmt.Println("this is the data ", products)
	return h.respondWithData(c, http.StatusOK, "success", products)
}
func (h *Handler) InZAListing(c echo.Context) error {
	ctx := c.Request().Context()
	id := ""
	products, err := h.service.InZAListing(ctx, id)
	if err != nil {
		return h.respondWithError(c, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch products", "details": err.Error()})
	}
	fmt.Println("this is the data ", products)
	return h.respondWithData(c, http.StatusOK, "success", products)
}

// / Product listing
// func (h *Handler) ProductActiveListing(c echo.Context) error {
// 	fmt.Println("in activeeee")
// 	ctx := c.Request().Context()

//		products, err := h.service.ProductActiveListing(ctx)
//		if err != nil {
//			return h.respondWithError(c, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch products", "details": err.Error()})
//		}
//		fmt.Println("this is the data ", products)
//		return h.respondWithData(c, http.StatusOK, "success", products)
//	}
func (h *Handler) ProductListing(c echo.Context) error {
	ctx := c.Request().Context()

	products, err := h.service.ProductListing(ctx)
	if err != nil {
		return h.respondWithError(c, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch products", "details": err.Error()})
	}
	fmt.Println("this is the data ", products)
	return h.respondWithData(c, http.StatusOK, "success", products)
}
func (h *Handler) ListingSingle(c echo.Context) error {
	ctx := c.Request().Context()
	id := c.Param("id")
	products, err := h.service.ListingSingle(ctx, id)
	if err != nil {
		return h.respondWithError(c, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch products", "details": err.Error()})
	}
	fmt.Println("this is the data ", products)
	return h.respondWithData(c, http.StatusOK, "success", products)
}

func (h *Handler) CategoryListing(c echo.Context) error {
	ctx := c.Request().Context()
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
	products, err := h.service.CategoryListing(ctx, request.Category)
	if err != nil {
		return h.respondWithError(c, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch products", "details": err.Error()})
	}
	fmt.Println("this is the data ", products)
	return h.respondWithData(c, http.StatusOK, "success", products)
}
