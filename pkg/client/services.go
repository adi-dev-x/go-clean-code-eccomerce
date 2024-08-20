package services

import (
	"fmt"
	"log"
	"math/rand"
	"myproject/pkg/config"
	"myproject/pkg/model"
	"net/smtp"

	"strconv"
	"time"

	"github.com/jung-kurt/gofpdf"
	"github.com/xuri/excelize/v2"
)

type Services interface {
	GenerateOtp(length int) int

	SendEmailWithOTP(email string) (string, error)
	SendOrderConfirmationEmail(orderUUID string, amount float64, recipientEmail string) error
	SendOrderReturnConfirmationEmailUser(name string, amt float64, unit int, mail string) error
	SendOrderReturnConfirmationEmailToUser(name string, amt float64, unit int, mail string)
	SendOrderReturnConfirmationEmailVendor(name string, amt float64, unit int, mail string) error
	GenerateDailySalesReportExcel(orders []model.ListOrdersVendor, facts model.Salesfact, types string, id string) (string, error)
	GenerateDailySalesReportPDF(orders []model.ListOrdersVendor, facts model.Salesfact, types string, id string) (string, error)
	GenerateDailySalesReportPDFAdmin(orders []model.ListOrdersAdmin, facts model.Salesfact, types string, id string) (string, error)
	GenerateDailySalesReportExcelAdmin(orders []model.ListOrdersAdmin, facts model.Salesfact, types string, id string) (string, error)
}
type MyService struct {
	Config config.Config
}

func (s MyService) GenerateDailySalesReportPDFAdmin(orders []model.ListOrdersAdmin, facts model.Salesfact, types string, id string) (string, error) {
	// Create a new PDF document
	pdf := gofpdf.New("P", "mm", "A4", "")

	// Add a page
	pdf.AddPage()

	// Set font
	pdf.SetFont("Arial", "", 12)

	// Add Salesfact data
	pdf.Cell(0, 10, fmt.Sprintf("Revenue: %.2f", facts.Revenue))
	pdf.Ln(10) // Line break

	pdf.Cell(0, 10, fmt.Sprintf("Total Discount: %.2f", facts.TotalDiscount))
	pdf.Ln(10) // Line break

	pdf.Cell(0, 10, fmt.Sprintf("Total Sales: %.2f", facts.TotalSales))
	pdf.Ln(10) // Line break

	pdf.Cell(0, 10, fmt.Sprintf("Total Orders: %d", facts.TotalOrders))
	pdf.Ln(20) // Line break before table headers

	// Set table headers
	headers := []string{"Check Date", "Product Name", "Quantity", "Status", "Returned", "Total Amount", "Product ID", "User", "User Address", "Order Date"}
	for _, header := range headers {
		pdf.Cell(30, 10, header)
	}
	pdf.Ln(10) // Line break after headers

	// Add order data
	for _, order := range orders {
		pdf.Cell(30, 10, order.ListDate)
		pdf.Cell(30, 10, order.Name)
		pdf.Cell(20, 10, fmt.Sprintf("%d", order.Unit))
		pdf.Cell(30, 10, order.Status)
		pdf.Cell(20, 10, fmt.Sprintf("%t", order.Returned))
		pdf.Cell(30, 10, fmt.Sprintf("%.2f", order.Amount))
		pdf.Cell(20, 10, order.Pid)
		pdf.Cell(30, 10, order.User)
		pdf.Cell(40, 10, order.Add)
		pdf.Cell(30, 10, order.Date)
		pdf.Ln(10) // Line break after each row
	}

	// Save the file
	rand.Seed(time.Now().UnixNano())

	// Generate a random number between 0 and 100
	randomNumber := rand.Intn(101)
	fileName := fmt.Sprintf("%s_Sales_Report_%d_%s.pdf", types, randomNumber, id)
	if err := pdf.OutputFileAndClose(fileName); err != nil {
		return "", fmt.Errorf("failed to save PDF file: %w", err)
	}

	return fileName, nil
}
func (s MyService) GenerateDailySalesReportPDF(orders []model.ListOrdersVendor, facts model.Salesfact, types string, id string) (string, error) {
	// Create a new PDF document
	pdf := gofpdf.New("P", "mm", "A4", "")

	// Add a page
	pdf.AddPage()

	// Set font
	pdf.SetFont("Arial", "", 12)

	// Add Salesfact data
	pdf.Cell(0, 10, fmt.Sprintf("Revenue: %.2f", facts.Revenue))
	pdf.Ln(10) // Line break

	pdf.Cell(0, 10, fmt.Sprintf("Total Discount: %.2f", facts.TotalDiscount))
	pdf.Ln(10) // Line break

	pdf.Cell(0, 10, fmt.Sprintf("Total Sales: %.2f", facts.TotalSales))
	pdf.Ln(10) // Line break

	pdf.Cell(0, 10, fmt.Sprintf("Total Orders: %d", facts.TotalOrders))
	pdf.Ln(20) // Line break before table headers

	// Set table headers
	headers := []string{"Check Date", "Product Name", "Quantity", "Status", "Returned", "Total Amount", "Product ID", "User", "User Address", "Order Date"}
	for _, header := range headers {
		pdf.Cell(30, 10, header)
	}
	pdf.Ln(10) // Line break after headers

	// Add order data
	for _, order := range orders {
		pdf.Cell(30, 10, order.ListDate)
		pdf.Cell(30, 10, order.Name)
		pdf.Cell(20, 10, fmt.Sprintf("%d", order.Unit))
		pdf.Cell(30, 10, order.Status)
		pdf.Cell(20, 10, fmt.Sprintf("%t", order.Returned))
		pdf.Cell(30, 10, fmt.Sprintf("%.2f", order.Amount))
		pdf.Cell(20, 10, order.Pid)
		pdf.Cell(30, 10, order.User)
		pdf.Cell(40, 10, order.Add)
		pdf.Cell(30, 10, order.Date)
		pdf.Ln(10) // Line break after each row
	}
	if id == "" {
		id = "Admin_Report"
	}
	rand.Seed(time.Now().UnixNano())

	// Generate a random number between 0 and 100
	randomNumber := rand.Intn(101)
	// Save the file
	fileName := fmt.Sprintf("%s_Sales_Report_%d_%s.pdf", types, randomNumber, id)
	if err := pdf.OutputFileAndClose(fileName); err != nil {
		return "", fmt.Errorf("failed to save PDF file: %w", err)
	}

	return fileName, nil
}
func (s MyService) GenerateDailySalesReportExcelAdmin(orders []model.ListOrdersAdmin, facts model.Salesfact, types string, id string) (string, error) {

	file := excelize.NewFile()
	sheet := "Sales Report"
	file.SetSheetName(file.GetSheetName(0), sheet)

	// Add Salesfact data at the top
	file.SetCellValue(sheet, "A1", "Revenue")
	file.SetCellValue(sheet, "B1", facts.Revenue)

	file.SetCellValue(sheet, "A2", "Total Discount")
	file.SetCellValue(sheet, "B2", facts.TotalDiscount)

	file.SetCellValue(sheet, "A3", "Total Sales")
	file.SetCellValue(sheet, "B3", facts.TotalSales)

	file.SetCellValue(sheet, "A4", "Total Orders")
	file.SetCellValue(sheet, "B4", facts.TotalOrders)

	// Add a row of space
	startingRowForOrders := 6 // Two rows after the last fact row

	// Set headers
	headers := []string{"Check Date", "Product Name", "Quantity", "Status", "Returned", "Total Amount", "Product ID", "User", "User Address", "Order Date", "Vendor Name"}
	for i, header := range headers {
		cell := fmt.Sprintf("%s%d", string(rune('A'+i)), startingRowForOrders)
		file.SetCellValue(sheet, cell, header)
	}

	// Fill data
	for i, order := range orders {
		row := startingRowForOrders + i + 1 // Start from the row after headers
		file.SetCellValue(sheet, fmt.Sprintf("A%d", row), order.ListDate)
		file.SetCellValue(sheet, fmt.Sprintf("B%d", row), order.Name)
		file.SetCellValue(sheet, fmt.Sprintf("C%d", row), order.Unit)
		file.SetCellValue(sheet, fmt.Sprintf("D%d", row), order.Status)
		file.SetCellValue(sheet, fmt.Sprintf("E%d", row), order.Returned)
		file.SetCellValue(sheet, fmt.Sprintf("F%d", row), order.Amount)
		file.SetCellValue(sheet, fmt.Sprintf("G%d", row), order.Pid)
		file.SetCellValue(sheet, fmt.Sprintf("H%d", row), order.User)
		file.SetCellValue(sheet, fmt.Sprintf("I%d", row), order.Add)
		file.SetCellValue(sheet, fmt.Sprintf("J%d", row), order.Date)
		file.SetCellValue(sheet, fmt.Sprintf("k%d", row), order.VName)
	}
	rand.Seed(time.Now().UnixNano())

	// Generate a random number between 0 and 100
	randomNumber := rand.Intn(101)
	// Save the file
	fileName := fmt.Sprintf("%s_Sales_Report_%d_%s.xlsx", types, randomNumber, id)
	err := file.SaveAs(fileName)
	if err != nil {
		return "", fmt.Errorf("failed to save Excel file: %w", err)
	}

	return fileName, nil
}
func (s MyService) GenerateDailySalesReportExcel(orders []model.ListOrdersVendor, facts model.Salesfact, types string, id string) (string, error) {

	file := excelize.NewFile()
	sheet := "Sales Report"
	file.SetSheetName(file.GetSheetName(0), sheet)

	// Add Salesfact data at the top
	file.SetCellValue(sheet, "A1", "Revenue")
	file.SetCellValue(sheet, "B1", facts.Revenue)

	file.SetCellValue(sheet, "A2", "Total Discount")
	file.SetCellValue(sheet, "B2", facts.TotalDiscount)

	file.SetCellValue(sheet, "A3", "Total Sales")
	file.SetCellValue(sheet, "B3", facts.TotalSales)

	file.SetCellValue(sheet, "A4", "Total Orders")
	file.SetCellValue(sheet, "B4", facts.TotalOrders)

	// Add a row of space
	startingRowForOrders := 6 // Two rows after the last fact row

	// Set headers
	headers := []string{"Check Date", "Product Name", "Quantity", "Status", "Returned", "Total Amount", "Product ID", "User", "User Address", "Order Date"}
	for i, header := range headers {
		cell := fmt.Sprintf("%s%d", string(rune('A'+i)), startingRowForOrders)
		file.SetCellValue(sheet, cell, header)
	}

	// Fill data
	for i, order := range orders {
		row := startingRowForOrders + i + 1 // Start from the row after headers
		file.SetCellValue(sheet, fmt.Sprintf("A%d", row), order.ListDate)
		file.SetCellValue(sheet, fmt.Sprintf("B%d", row), order.Name)
		file.SetCellValue(sheet, fmt.Sprintf("C%d", row), order.Unit)
		file.SetCellValue(sheet, fmt.Sprintf("D%d", row), order.Status)
		file.SetCellValue(sheet, fmt.Sprintf("E%d", row), order.Returned)
		file.SetCellValue(sheet, fmt.Sprintf("F%d", row), order.Amount)
		file.SetCellValue(sheet, fmt.Sprintf("G%d", row), order.Pid)
		file.SetCellValue(sheet, fmt.Sprintf("H%d", row), order.User)
		file.SetCellValue(sheet, fmt.Sprintf("I%d", row), order.Add)
		file.SetCellValue(sheet, fmt.Sprintf("J%d", row), order.Date)
	}
	rand.Seed(time.Now().UnixNano())

	// Generate a random number between 0 and 100
	randomNumber := rand.Intn(101)
	// Save the file
	fileName := fmt.Sprintf("%s_Sales_Report_%d_%s.xlsx", types, randomNumber, id)
	err := file.SaveAs(fileName)
	if err != nil {
		return "", fmt.Errorf("failed to save Excel file: %w", err)
	}

	return fileName, nil
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
func (s MyService) SendOrderReturnConfirmationEmailToUser(name string, amt float64, unit int, recipientEmail string) {
	fmt.Println("this is in the SendOrderReturnConfirmationEmail !!!--", name, amt, recipientEmail)
	// Message.
	subject := "Vendor has Cancelled your order"
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
		fmt.Errorf("failed to send email: %w", err)
	}

}
func (s MyService) SendOrderReturnConfirmationEmailUser(name string, amt float64, unit int, recipientEmail string) error {
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
func (s MyService) SendOrderReturnConfirmationEmailVendor(name string, amt float64, unit int, recipientEmail string) error {
	fmt.Println("this is in the SendOrderReturnConfirmationEmail !!!--", name, amt, recipientEmail)
	// Message.
	subject := "Customer placed for return"
	body := fmt.Sprintf("A  order %s has been placed for returning!\nunits: %d\nAmount: RS%.2f", name, unit, amt)
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
