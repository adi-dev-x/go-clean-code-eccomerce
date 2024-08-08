package user

import (
	// "encoding/json"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"strings"
	"time"

	// db "myproject/pkg/database"
	services "myproject/pkg/client"
	"myproject/pkg/config"
	db "myproject/pkg/database"

	"myproject/pkg/model"

	"net/http"

	// "time"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	service   Service
	services  services.Services
	adminjw   Adminjwt
	templates *template.Template
	cnf       config.Config
}

func NewHandler(service Service, srv services.Services, adTK Adminjwt, cnf config.Config) *Handler {

	return &Handler{
		service:  service,
		services: srv,
		adminjw:  adTK,
		cnf:      cnf,
	}
}
func (h *Handler) MountRoutes(engine *echo.Echo) {
	//applicantApi := engine.Group(basePath)
	applicantApi := engine.Group("/user")
	applicantApi.POST("/register", h.Register)
	applicantApi.POST("/login", h.Login)
	applicantApi.POST("/OtpLogin", h.OtpLogin)
	engine.GET("/RazorPay/:id/", h.Gate)
	engine.GET("/RazorPaySucess/:id/", h.GateSuccess)
	engine.GET("/RazorPayFailed/:id/", h.GateFailed)

	renderer := &Handler{
		templates: template.Must(template.ParseGlob("pkg/templates/*.html")),
	}
	engine.Renderer = renderer
	applicantApi.Use(h.adminjw.AdminAuthMiddleware())
	{
		applicantApi.POST("/UpdateUser", h.UpdateUser)
		applicantApi.GET("/listing", h.Listing)
		applicantApi.GET("/Latestlisting", h.LatestListing)
		applicantApi.GET("/PhighListing", h.PhighListing)
		applicantApi.GET("/PlowListing", h.PlowListing)
		applicantApi.GET("/InAZListing", h.InAZListing)
		applicantApi.GET("/InZAListing", h.InZAListing)
		applicantApi.POST("/AddTocart", h.AddToCart)
		//applicantApi.POST("/AddTocart", h.AddToCart)
		applicantApi.POST("/AddToWish", h.AddToWish)
		applicantApi.GET("/Listcart/", h.Listcart)
		applicantApi.GET("/ListWish/:id", h.ListWish)
		applicantApi.POST("/AddToorder", h.AddToorder)
		applicantApi.POST("/AddAddress", h.AddAddress)
		applicantApi.GET("/ListAddress", h.ListAddress)
		applicantApi.POST("/AddToCheck", h.AddToCheck)

		applicantApi.GET("/listCoupon", h.ActiveListing)
		///list orders
		applicantApi.GET("/listAllOrders", h.ListAllOrders)
		applicantApi.GET("/listReturnedOrders", h.ListReturnedOrders)
		applicantApi.GET("/listFailedOrders", h.ListFailedOrders)
		applicantApi.GET("/listCompletedOrders", h.ListCompletedOrders)
		applicantApi.GET("/listPendingOrders", h.ListPendingOrders)
		applicantApi.POST("/returnItem", h.ReturnItem)

	}

	engine.GET("/RazorPay", func(c echo.Context) error {
		return c.Render(http.StatusOK, "payment.html", nil)
	})
	//applicantApi.Use(middleware.VenAuthMiddleware())
	//{

	//}
}
func (h *Handler) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return h.templates.ExecuteTemplate(w, name, data)
}

func (h *Handler) respondWithError(c echo.Context, code int, msg interface{}) error {
	resp := map[string]interface{}{
		"msg": msg,
	}

	return c.JSON(code, resp)
}

func (h *Handler) respondWithData(c echo.Context, code int, message interface{}, data interface{}) error {
	if data == nil {
		data = "Succesfully done"
		resp := map[string]interface{}{
			"msg":     message,
			"Process": data,
		}
		return c.JSON(code, resp)

	}
	resp := map[string]interface{}{
		"msg":  message,
		"data": data,
	}
	return c.JSON(code, resp)
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

func (h *Handler) Register(c echo.Context) error {

	fmt.Println("this is in the handler Register")
	var request model.UserRegisterRequest
	if err := c.Bind(&request); err != nil {
		return h.respondWithError(c, http.StatusBadRequest, map[string]string{"request-parse": err.Error()})
	}

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
	storedData, _ := db.GetRedis(request.Email)
	fmt.Println("this is the keyy!!!!!", storedData)

	return h.respondWithData(c, http.StatusOK, "success", nil)
}
func (h *Handler) UpdateUser(c echo.Context) error {

	fmt.Println("this is in the handler UpdateUser")
	var request model.UserRegisterRequest
	if err := c.Bind(&request); err != nil {
		return h.respondWithError(c, http.StatusBadRequest, map[string]string{"request-parse": err.Error()})
	}

	// Validate request fields
	//errVal := request.Valid()

	ctx := c.Request().Context()
	if err := h.service.UpdateUser(ctx, request); err != nil {
		return h.respondWithError(c, http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	fmt.Println("this is in the handler UpdateUser")

	otp, err := h.services.SendEmailWithOTP(request.Email)
	if err != nil {
		return h.respondWithError(c, http.StatusInternalServerError, map[string]string{"error": "error in sending otp"})

	}
	err = db.SetRedis(request.Email, otp, time.Minute*5)
	if err != nil {
		return h.respondWithError(c, http.StatusInternalServerError, map[string]string{"error": "error in saving otp"})

	}
	storedData, _ := db.GetRedis(request.Email)
	fmt.Println("this is the keyy!!!!!", storedData)

	return h.respondWithData(c, http.StatusOK, "success", nil)
}
func (h *Handler) Login(c echo.Context) error {

	fmt.Println("this is in the handler Register")
	var request model.UserLoginRequest
	if err := c.Bind(&request); err != nil {
		return h.respondWithError(c, http.StatusBadRequest, map[string]string{"request-parse": err.Error()})
	}

	ctx := c.Request().Context()
	if err := h.service.Login(ctx, request); err != nil {
		return h.respondWithError(c, http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	fmt.Println("this is in the handler Register")
	token, err := h.adminjw.GenerateAdminToken(request.Email)
	if err != nil {
		return h.respondWithError(c, http.StatusInternalServerError, map[string]string{"token-generation": err.Error()})
	}

	fmt.Println("User logged in successfully")
	return h.respondWithData(c, http.StatusOK, "success", map[string]string{"token": token})
}
func (h *Handler) OtpLogin(c echo.Context) error {
	// Parse request body into UserRegisterRequest
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
	return h.respondWithData(c, http.StatusOK, "success", nil)
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
func (h *Handler) ListAddress(c echo.Context) error {
	fmt.Println("this is in the ListAddress")
	authHeader := c.Request().Header.Get("Authorization")
	fmt.Println("inside the cart list ", authHeader)
	username := c.Get("username").(string)
	ctx := c.Request().Context()

	products, err := h.service.ListAddress(ctx, username)
	if err != nil {
		return h.respondWithError(c, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch products", "details": err.Error()})
	}
	fmt.Println("this is the data ", products)
	return h.respondWithData(c, http.StatusOK, "success", products)
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
func (h *Handler) ListCompletedOrders(c echo.Context) error {
	fmt.Println("this is in the handler ListCompletedOrders")
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
func (h *Handler) ListPendingOrders(c echo.Context) error {
	fmt.Println("this is in the handler ListCompletedOrders")
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
func (h *Handler) Listcart(c echo.Context) error {
	authHeader := c.Request().Header.Get("Authorization")
	fmt.Println("inside the cart list ", authHeader)
	username := c.Get("username").(string)
	fmt.Println("inside the cart list ", username)
	// id := c.Param("id")
	ctx := c.Request().Context()

	products, err := h.service.Listcart(ctx, username)
	if err != nil {
		return h.respondWithError(c, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch products", "details": err.Error()})
	}
	fmt.Println("this is the data ", products)
	return h.respondWithData(c, http.StatusOK, "success", products)
}
func (h *Handler) ListWish(c echo.Context) error {
	authHeader := c.Request().Header.Get("Authorization")
	fmt.Println("inside the cart list ", authHeader)
	username := c.Get("username").(string)
	fmt.Println("inside the cart list ", username)
	// id := c.Param("id")
	ctx := c.Request().Context()

	products, err := h.service.ListWish(ctx, username)
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
func (h *Handler) PhighListing(c echo.Context) error {
	ctx := c.Request().Context()

	products, err := h.service.PhighListing(ctx)
	if err != nil {
		return h.respondWithError(c, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch products", "details": err.Error()})
	}
	fmt.Println("this is the data ", products)
	return h.respondWithData(c, http.StatusOK, "success", products)
}
func (h *Handler) PlowListing(c echo.Context) error {
	ctx := c.Request().Context()

	products, err := h.service.PlowListing(ctx)
	if err != nil {
		return h.respondWithError(c, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch products", "details": err.Error()})
	}
	fmt.Println("this is the data ", products)
	return h.respondWithData(c, http.StatusOK, "success", products)
}
func (h *Handler) InAZListing(c echo.Context) error {
	ctx := c.Request().Context()

	products, err := h.service.InAZListing(ctx)
	if err != nil {
		return h.respondWithError(c, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch products", "details": err.Error()})
	}
	fmt.Println("this is the data ", products)
	return h.respondWithData(c, http.StatusOK, "success", products)
}
func (h *Handler) InZAListing(c echo.Context) error {
	ctx := c.Request().Context()

	products, err := h.service.InZAListing(ctx)
	if err != nil {
		return h.respondWithError(c, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch products", "details": err.Error()})
	}
	fmt.Println("this is the data ", products)
	return h.respondWithData(c, http.StatusOK, "success", products)
}

func (h *Handler) AddToCart(c echo.Context) error {
	fmt.Println("this is in the handler AddToCart")
	authHeader := c.Request().Header.Get("Authorization")
	fmt.Println("inside the cart list ", authHeader)
	username := c.Get("username").(string)
	fmt.Println("inside the cart list ", username)

	var request model.Cart
	if err := c.Bind(&request); err != nil {
		return h.respondWithError(c, http.StatusBadRequest, map[string]string{"request-parse": err.Error()})
	}

	ctx := c.Request().Context()
	if err := h.service.AddTocart(ctx, request, username); err != nil {
		return h.respondWithError(c, http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	fmt.Println("Item added to cart successfully")

	return h.respondWithData(c, http.StatusOK, "success", nil)
}
func (h *Handler) AddToWish(c echo.Context) error {
	fmt.Println("this is in the handler AddToWish")

	var request model.Wishlist
	if err := c.Bind(&request); err != nil {
		return h.respondWithError(c, http.StatusBadRequest, map[string]string{"request-parse": err.Error()})
	}

	ctx := c.Request().Context()
	if err := h.service.AddToWish(ctx, request); err != nil {
		return h.respondWithError(c, http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	fmt.Println("Item added to cart successfully")

	return h.respondWithData(c, http.StatusOK, "success", nil)
}
func (h *Handler) AddAddress(c echo.Context) error {
	fmt.Println("this is in the handler AddToWish")
	authHeader := c.Request().Header.Get("Authorization")
	fmt.Println("inside the cart list ", authHeader)
	username := c.Get("username").(string)
	fmt.Println("inside the cart list ", username)
	var request model.Address

	if err := c.Bind(&request); err != nil {
		return h.respondWithError(c, http.StatusBadRequest, map[string]string{"request-parse": err.Error()})
	}
	errVal := request.Check()
	if len(errVal) > 0 {
		return h.respondWithError(c, http.StatusBadRequest, map[string]interface{}{"invalid-request": errVal})
	}
	ctx := c.Request().Context()
	if err := h.service.AddAddress(ctx, request, username); err != nil {
		return h.respondWithError(c, http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	fmt.Println("Item added to cart successfully")

	return h.respondWithData(c, http.StatusOK, "success", nil)
}

func (h *Handler) AddToorder(c echo.Context) error {
	fmt.Println("this is in the handler AddToOrder")
	authHeader := c.Request().Header.Get("Authorization")
	fmt.Println("inside the cart list ", authHeader)
	usernames := c.Get("username").(string)
	fmt.Println("inside the cart list ", usernames)

	var request model.Order
	if err := c.Bind(&request); err != nil {
		return h.respondWithError(c, http.StatusBadRequest, map[string]string{"request-parse": err.Error()})
	}
	errVal := request.Valid()
	if len(errVal) > 0 {
		return h.respondWithError(c, http.StatusBadRequest, map[string]interface{}{"invalid-request": errVal})
	}

	ctx := context.WithValue(c.Request().Context(), "username", usernames)
	username := ctx.Value("username").(string)
	fmt.Println("inside the cart list ", username)
	var da model.RZpayment
	da, err := h.service.AddToorder(ctx, request)
	da.Token = authHeader
	if err != nil {
		return h.respondWithError(c, http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	fmt.Println("Item added to order")
	var f string
	f = username + "RZ"
	fmt.Println("this is the data ", da)
	dam, _ := json.Marshal(&da)
	db.SetRedis(f, dam, time.Minute*5)
	redirectURL := "user/RazorPay/" + f
	fmt.Println("this is the url  ", redirectURL)

	return h.respondWithData(c, http.StatusOK, "success", nil)

}
func (h *Handler) ReturnItem(c echo.Context) error {
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

// /this is testing
func (h *Handler) AddToCheck(c echo.Context) error {
	fmt.Println("this is in the handler AddToCheck")
	authHeader := c.Request().Header.Get("Authorization")
	fmt.Println("inside the cart list ", authHeader)
	username := c.Get("username").(string)
	fmt.Println("inside the cart list ", username)
	ctx := c.Request().Context()
	var request model.CheckOut
	if err := c.Bind(&request); err != nil {
		return h.respondWithError(c, http.StatusBadRequest, map[string]string{"request-parse": err.Error()})
	}
	errValues, _ := request.Valid()
	if len(errValues) > 0 {
		return h.respondWithError(c, http.StatusInternalServerError, map[string]interface{}{"invalid-request": errValues})
	}
	///////this start

	rz, err := h.service.AddToCheck(ctx, request, username)
	if rz.Amt != 0 {
		authHeader := c.Request().Header.Get("Authorization")
		tokenString := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer"))
		url := fmt.Sprintf("http://localhost:8080/RazorPay/%s/", tokenString)
		f := username + "RZ"
		dam, _ := json.Marshal(&rz)
		db.SetRedis(f, dam, time.Minute*5)
		fmt.Println("this is the url  ", url)
		return h.respondWithData(c, http.StatusOK, "success", url)

	}
	if err != nil {
		return h.respondWithError(c, http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return h.respondWithData(c, http.StatusOK, "success", nil)
	//return c.Redirect(http.StatusTemporaryRedirect, redirectURL)
}
func (h *Handler) Gate(c echo.Context) error {
	fmt.Println("inside gatee RZPAy")

	id := c.Param("id")
	//username, err := AdminAuthentication(tokenString, s.Config.AdJWTKey)
	username, _ := AdminAuthentication(id, h.cnf.AdJWTKey)
	key := username + "RZ"
	fmt.Println("this the id in Gate  !!!", id)
	var storedData model.RZpayment
	stored, _ := db.GetRedis(key)
	json.Unmarshal([]byte(stored), &storedData)
	fmt.Println("data in razor pay!!!", storedData, "amt !!!!!", storedData.Amt*100)
	b := int64(storedData.Amt)
	data := map[string]interface{}{
		"Order_ID":    storedData.Order_ID,
		"Payment_ID":  storedData.Id,
		"totalPrice":  b, // Replace with actual invoice ID if available
		"totalAmount": b,
		"UserToken":   storedData.Token,
	}
	return c.Render(http.StatusOK, "payment.html", data)

}
func (h *Handler) GateSuccess(c echo.Context) error {
	fmt.Println("inside gatee successs !!!!!!")
	id := c.Param("id")
	username, _ := AdminAuthentication(id, h.cnf.AdJWTKey)
	key := username + "RZ"
	var storedData model.RZpayment
	stored, _ := db.GetRedis(key)
	json.Unmarshal([]byte(stored), &storedData)
	fmt.Println(storedData)
	ctx := c.Request().Context()
	err := h.service.PaymentSuccess(ctx, storedData, username)
	fmt.Println("After GateSuccess completed")
	if err != nil {
		return h.respondWithError(c, http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return h.respondWithData(c, http.StatusOK, "success", nil)

}
func (h *Handler) GateFailed(c echo.Context) error {
	fmt.Println("inside gatee failed !!!!!!")
	id := c.Param("id")
	username, _ := AdminAuthentication(id, h.cnf.AdJWTKey)
	key := username + "RZ"
	var storedData model.RZpayment
	stored, _ := db.GetRedis(key)
	json.Unmarshal([]byte(stored), &storedData)
	fmt.Println(storedData)
	ctx := c.Request().Context()
	err := h.service.PaymentFailed(ctx, storedData, username)
	fmt.Println("After Gate Failed completed")
	if err != nil {
		return h.respondWithError(c, http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return h.respondWithData(c, http.StatusOK, "success", nil)

}
