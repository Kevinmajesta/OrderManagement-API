// pkg/email/email.go
package email

import (
	"fmt"
	"strconv"

	"Kevinmajesta/OrderManagementAPI/configs" // Pastikan ini path modul Anda

	"gopkg.in/gomail.v2"
)

// --- MULAI PERUBAHAN ---

// EmailSenderService adalah interface yang mendefinisikan metode pengiriman email.
// Kita memberikan nama yang sedikit berbeda (EmailSenderService) agar tidak bentrok
// dengan nama struct EmailSender yang sudah ada, atau Anda bisa menamai interface EmailSender
// dan struct EmailSenderImpl.
type EmailSenderService interface {
	SendEmail(to []string, subject, body string) error
	SendWelcomeEmail(to, name, message string) error
	SendResetPasswordEmail(to, name, resetCode string) error
	SendVerificationEmail(to, name, code string) error
	SendTransactionInfo(to, Transactions_id, Cart_id, User_id, Fullname_user, Trx_date, Payment, Payment_url, Amount string) error
}

// EmailSender adalah implementasi konkret dari EmailSenderService.
type EmailSender struct { // Nama struct tetap EmailSender
	Config *configs.Config
}

// NewEmailSender adalah constructor yang mengembalikan interface EmailSenderService.
func NewEmailSender(config *configs.Config) EmailSenderService { // <--- Return type-nya sekarang interface
	return &EmailSender{Config: config} // Mengembalikan pointer ke struct yang mengimplementasikan interface
}

// --- AKHIR PERUBAHAN ---

func (e *EmailSender) SendEmail(to []string, subject, body string) error {
	from := e.Config.SMTP.User
	password := e.Config.SMTP.Password
	smtpHost := e.Config.SMTP.Host
	smtpPortStr := e.Config.SMTP.Port

	smtpPort, err := strconv.Atoi(smtpPortStr)
	if err != nil {
		return fmt.Errorf("invalid SMTP port in config: %s", smtpPortStr)
	}

	mailer := gomail.NewMessage()
	mailer.SetHeader("From", from)
	mailer.SetHeader("To", to...)
	mailer.SetHeader("Subject", subject)
	mailer.SetBody("text/plain", body)

	dialer := gomail.NewDialer(smtpHost, smtpPort, from, password)
	err = dialer.DialAndSend(mailer)
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}

// Metode lainnya tetap sama, karena mereka sudah memiliki receiver *EmailSender
// yang secara implisit akan mengimplementasikan metode-metode pada EmailSenderService.
func (e *EmailSender) SendWelcomeEmail(to, name, message string) error {
	subject := "Welcome Email | Depublic"
	body := fmt.Sprintf("Dear %s,\nThis is a welcome email message from depublic\n\nDepublic Team", name)
	return e.SendEmail([]string{to}, subject, body)
}

func (e *EmailSender) SendResetPasswordEmail(to, name, resetCode string) error {
	subject := "Reset Password | Depublic"
	body := fmt.Sprintf("Dear %s,\nPlease use the following code to reset your password: %s\n\nDepublic Team", name, resetCode)
	return e.SendEmail([]string{to}, subject, body)
}

func (e *EmailSender) SendVerificationEmail(to, name, code string) error {
	subject := "Verify Your Email | Depublic"
	body := fmt.Sprintf("Dear %s,\nPlease use the following code to verify your email: %s\n\nDepublic Team", name, code)
	return e.SendEmail([]string{to}, subject, body)
}

func (e *EmailSender) SendTransactionInfo(to, Transactions_id, Cart_id, User_id,
	Fullname_user, Trx_date, Payment, Payment_url, Amount string) error {
	subject := "Transaction Info | Depublic"
	body := fmt.Sprintf("Dear %s,\nThis is your transaction info from Depublic:\n\nTransaction ID: %s\n\nCart ID: %s\n\nUser ID: %s\n\nFullname: %s\n\nTransaction Date: %s\n\nPayment Type: %s\n\nURL Payment: %s\n\nTotal Amount: %s\n\n\nDepublic Team ",
		Fullname_user, Transactions_id, Cart_id, User_id, Fullname_user, Trx_date, Payment, Payment_url, Amount)
	return e.SendEmail([]string{to}, subject, body)
}