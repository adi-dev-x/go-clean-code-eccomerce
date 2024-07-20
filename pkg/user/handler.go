package user

import (
	// "encoding/json"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"time"

	// db "myproject/pkg/database"
	services "myproject/pkg/client"
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
}

func NewHandler(service Service, srv services.Services, adTK Adminjwt) *Handler {

	return &Handler{
		service:  service,
		services: srv,
		adminjw:  adTK,
	}
}
func (h *Handler) MountRoutes(engine *echo.Echo) {
	//applicantApi := engine.Group(basePath)
	applicantApi := engine.Group("/user")
	applicantApi.POST("/register", h.Register)
	applicantApi.POST("/login", h.Login)
	applicantApi.POST("/OtpLogin", h.OtpLogin)
	engine.GET("/RazorPay/:id/", h.Gate)
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
		applicantApi.GET("/Listcart/:id", h.Listcart)
		applicantApi.GET("/ListWish/:id", h.ListWish)
		applicantApi.POST("/AddToorder", h.AddToorder)
		applicantApi.POST("/AddAddress", h.AddAddress)
		applicantApi.GET("/ListAddress", h.ListAddress)
		applicantApi.POST("/AddToCheck", h.AddToCheck)

		applicantApi.GET("/listCoupon", h.ActiveListing)
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

	var request model.Cart
	if err := c.Bind(&request); err != nil {
		return h.respondWithError(c, http.StatusBadRequest, map[string]string{"request-parse": err.Error()})
	}

	ctx := c.Request().Context()
	if err := h.service.AddTocart(ctx, request); err != nil {
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

// /this is testing
func (h *Handler) AddToCheck(c echo.Context) error {
	fmt.Println("this is in the handler AddToOrder")
	authHeader := c.Request().Header.Get("Authorization")
	fmt.Println("inside the cart list ", authHeader)
	usernames := c.Get("username").(string)
	fmt.Println("inside the cart list ", usernames)

	var request model.CheckOut
	if err := c.Bind(&request); err != nil {
		return h.respondWithError(c, http.StatusBadRequest, map[string]string{"request-parse": err.Error()})
	}

	return h.respondWithData(c, http.StatusOK, "success", nil)
	//return c.Redirect(http.StatusTemporaryRedirect, redirectURL)
}
func (h *Handler) Gate(c echo.Context) error {
	fmt.Println("inside gatee RZPAy")

	id := c.Param("id")
	fmt.Println("this the id in Gate  !!!", id)
	var storedData model.RZpayment
	stored, _ := db.GetRedis(id)
	json.Unmarshal([]byte(stored), &storedData)
	fmt.Println("data !!!", storedData, "amt !!!!!", storedData.Amt*100)
	b := int64(storedData.Amt)
	data := map[string]interface{}{
		"invoiceID":     20,
		"appointmentID": 33,
		"totalPrice":    b, // Replace with actual invoice ID if available
		"totalAmount":   b,
		"UserToken":     storedData.Token,
	}
	return c.Render(http.StatusOK, "payment.html", data)

}
