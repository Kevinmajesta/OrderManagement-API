package email

import (
	"fmt"
	"strconv" // Import strconv untuk konversi port

	"Kevinmajesta/OrderManagementAPI/configs" // Perbaiki import path jika 'internal/entity' tidak lagi ada
	// Jika struct Config dipindahkan ke package 'configs', maka import path-nya harus disesuaikan.
	// Jika 'entity.Config' masih ada, pastikan itu mengacu pada struct Config yang benar.

	"gopkg.in/gomail.v2"
)

type EmailSender struct {
	Config *configs.Config // Menggunakan configs.Config
}

func NewEmailSender(config *configs.Config) *EmailSender {
	return &EmailSender{Config: config}
}

func (e *EmailSender) SendEmail(to []string, subject, body string) error {
	// Ambil From email dari konfigurasi SMTP_USER
	from := e.Config.SMTP.User
	password := e.Config.SMTP.Password
	smtpHost := e.Config.SMTP.Host
	smtpPortStr := e.Config.SMTP.Port // Port masih dalam string

	// Konversi port dari string ke int
	smtpPort, err := strconv.Atoi(smtpPortStr)
	if err != nil {
		return fmt.Errorf("invalid SMTP port in config: %s", smtpPortStr)
	}

	mailer := gomail.NewMessage()
	mailer.SetHeader("From", from)
	mailer.SetHeader("To", to...)
	mailer.SetHeader("Subject", subject)
	mailer.SetBody("text/plain", body)

	// Inisialisasi dialer dengan user dan password dari config
	dialer := gomail.NewDialer(smtpHost, smtpPort, from, password)
	err = dialer.DialAndSend(mailer)
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}

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