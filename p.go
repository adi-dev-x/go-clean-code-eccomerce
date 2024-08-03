package controllers

import (
	"encoding/json"
	"finmate/database"
	"finmate/models"
	"log"
	"net/http"
	"os"

	razorpay "github.com/razorpay/razorpay-go"
)

func getKey(token string) float64 {
	var ukey float64
	for key, value := range models.SaveUserToken {
		if value == token {
			ukey = key

			break
		}

	}
	return ukey
}

func CreateOrder(w http.ResponseWriter, r *http.Request) {
	var requestData struct {
		Plan   string ⁠ json:"plan" ⁠
		Coupon string ⁠ json:"coupon" ⁠
	}
	jwtToken := r.Header.Get("Authorization")
	userid := getKey(jwtToken)

	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	planName := requestData.Plan
	if planName == "" {
		http.Error(w, "Plan name is required", http.StatusBadRequest)
		return
	}
	Coupon := requestData.Coupon
	coupon:=0 
	if Coupon != "" {
		coupon, err = getPlanCoupon(Coupon)
		if err != nil {
			log.Println("Error fetching plan amount:", err)
			http.Error(w, "Failed to fetch plan amount", http.StatusInternalServerError)
			return
		}
		
	}

	amount, err := getPlanAmount(planName)
	if err != nil {
		log.Println("Error fetching plan amount:", err)
		http.Error(w, "Failed to fetch plan amount", http.StatusInternalServerError)
		return
	}
	discount, err := getPlandiscount(planName)
	if err != nil {
		log.Println("Error fetching plan amount:", err)
		http.Error(w, "Failed to fetch plan amount", http.StatusInternalServerError)
		return
	}


	client := razorpay.NewClient(os.Getenv("RAZORPAY_KEY"), os.Getenv("RAZORPAY_SECRET"))
	newamount := amount - (amount * (coupon / 100)) - discount
	data := map[string]interface{}{
		"amount":   newamount * 100, // Razorpay expects amount in paise
		"currency": "INR",
		"receipt":  "some_receipt_id",
	}

	body, err := client.Order.Create(data, nil)
	if err != nil {
		log.Println("Error creating Razorpay order:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	orderID, _ := body["id"].(string)
	amountFloat, _ := body["amount"].(float64)
	currency, _ := body["currency"].(string)
	receipt, _ := body["receipt"].(string)

	err = saveOrderToDB(orderID, amountFloat, currency, receipt, discount, Coupon, userid)
	if err != nil {
		log.Println("Error saving order to database:", err)
		http.Error(w, "Failed to save order", http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(body)
	if err != nil {
		http.Error(w, "Failed to parse response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func saveOrderToDB(orderID string, amount float64, currency string, receipt string, discount int, Coupon string, userid float64) error {
	db := database.DB
	_, err := db.Exec(`
        INSERT INTO orders (order_id, amount, currency, receipt,discount,coupon,userid) 
        VALUES ($1, $2, $3, $4,$5,$6)`,
		orderID, amount, currency, receipt, discount, Coupon, userid)
	return err
}

func getPlanAmount(planName string) (int, error) {
	db := database.DB

	var amount int
	query := "SELECT price FROM plans WHERE name=$1"
	err := db.QueryRow(query, planName).Scan(&amount)
	if err != nil {
		return 0, err
	}
	return amount, nil
}
func getPlandiscount(planName string) (int, error) {
	db := database.DB

	var amount int
	query := "SELECT discount FROM plans WHERE name=$1"
	err := db.QueryRow(query, planName).Scan(&amount)
	if err != nil {
		return 0, err
	}
	return amount, nil
}
func getPlanCoupon(Coupon string) (int, error) {
	db := database.DB

	var amount int
	query := "SELECT amount FROM copon WHERE name=$1 AND TO_DATE(expiry, 'DD/MM/YYYY') < CURRENT_DATE AND used=false"
	err := db.QueryRow(query, Coupon).Scan(&amount)
	if err != nil {
		return 0, err
	}
	return amount, nil
}