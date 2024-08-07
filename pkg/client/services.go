package services

import (
	"fmt"
	"log"
	"math/rand"
	"myproject/pkg/config"
	"net/smtp"

	"strconv"
	"time"
)

type Services interface {
	GenerateOtp(length int) int

	SendEmailWithOTP(email string) (string, error)
	SendOrderConfirmationEmail(orderUUID string, amount float64, recipientEmail string) error
	SendOrderReturnConfirmationEmail(name string, amt float64, unit int, mail string) error
}
type MyService struct {
	Config config.Config
}

func (s MyService) SendOrderConfirmationEmail(orderUUID string, amount float64, recipientEmail string) error {
	fmt.Println("this is in the SendOrderConfirmationEmail !!!--", orderUUID, amount, recipientEmail)
	// Message.
	subject := "Order Confirmation"
	body := fmt.Sprintf("Your order has been placed successfully!\nOrder UUID: %s\nAmount: RS%.2f", orderUUID, amount)
	message := fmt.Sprintf("Subject: %s\r\n\r\n%s", subject, body)

	// Authentication.
	SMTPemail := s.Config.SMTPemail
	SMTPpass := s.Config.Password
	auth := smtp.PlainAuth("", SMTPemail, SMTPpass, "smtp.gmail.com")
	fmt.Println("this is my mail !_+_++_+!", SMTPemail, "!+!+!+!+", SMTPpass)
	// Sending email.
	err := smtp.SendMail("smtp.gmail.com:587", auth, SMTPemail, []string{recipientEmail}, []byte(message))
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
func (s MyService) SendOrderReturnConfirmationEmail(name string, amt float64, unit int, recipientEmail string) error {
	fmt.Println("this is in the SendOrderReturnConfirmationEmail !!!--", name, amt, recipientEmail)
	// Message.
	subject := "Order item returned"
	body := fmt.Sprintf("Your order %s has been placed for returning!\nunits: %d\nAmount: RS%.2f", name, unit, amt)
	message := fmt.Sprintf("Subject: %s\r\n\r\n%s", subject, body)

	// Authentication.
	SMTPemail := s.Config.SMTPemail
	SMTPpass := s.Config.Password
	auth := smtp.PlainAuth("", SMTPemail, SMTPpass, "smtp.gmail.com")
	fmt.Println("this is my mail !_+_++_+!", SMTPemail, "!+!+!+!+", SMTPpass)
	// Sending email.
	err := smtp.SendMail("smtp.gmail.com:587", auth, SMTPemail, []string{recipientEmail}, []byte(message))
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

// GenerateOtp generates a random OTP of the specified length.
func (s MyService) GenerateOtp(length int) int {
	rand.Seed(time.Now().UnixNano())

	// Generate a random number between 10000 and 99999
	randomNum := rand.Intn(90000) + 10000

	fmt.Println("Random 5-digit number:", randomNum)
	return randomNum
}

// SendEmailWithOTP sends an OTP to the specified email address.
func (s MyService) SendEmailWithOTP(email string) (string, error) {
	// Generate OTP
	otp := strconv.Itoa(s.GenerateOtp(6))

	// Construct email message
	message := fmt.Sprintf("Subject: OTP for Verification\n\nYour OTP is: %s", otp)
	fmt.Println("this is my email  !!!!!", s.Config.SMTPemail, "this is my email  !!!!!", s.Config.Password)

	SMTPemail := s.Config.SMTPemail
	SMTPpass := s.Config.Password
	auth := smtp.PlainAuth("", "adithyanunni258@gmail.com", SMTPpass, "smtp.gmail.com")

	// Send email using SMTP server
	err := smtp.SendMail("smtp.gmail.com:587", auth, SMTPemail, []string{email}, []byte(message))
	if err != nil {
		log.Println("Error sending email:", err)
		return "", err
	}

	return otp, nil
}
