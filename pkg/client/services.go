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
}
type MyService struct {
	Config config.Config
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
