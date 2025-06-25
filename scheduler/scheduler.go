package scheduler

import (
	"errors"
	"log"
	"net/smtp"
	"os"
	"time"

	"github.com/robfig/cron/v3"
)

var sleep = time.Sleep // override in tests

func StartCronJobs() error {
	c := cron.New()

	_, err := c.AddFunc("0 0 * * *", func() {
		log.Println("🧹 Running cleanup job...")
		if err := cleanupOldFiles(); err != nil {
			log.Printf(" Cleanup failed: %v\n", err)
		} else {
			log.Println(" Cleanup completed successfully.")
		}
	})
	if err != nil {
		return errors.New("failed to schedule cleanup job")
	}

	_, err = c.AddFunc("0 9 * * MON", func() {
		log.Println("📧 Sending weekly report email...")
		if err := sendWeeklyEmail(); err != nil {
			log.Printf(" Email sending failed: %v\n", err)
		} else {
			log.Println(" Email sent successfully.")
		}
	})
	if err != nil {
		return errors.New("failed to schedule email job")
	}

	c.Start()
	log.Println(" Cron jobs started.")
	return nil
}

func cleanupOldFiles() error {
	sleep(1 * time.Second)
	return nil
}

func sendWeeklyEmail() error {
	// Skip actual email sending if environment is "TEST"
	if os.Getenv("ENV") == "TEST" {
		log.Println("Skipping real email sending in TEST environment")
		sleep(1 * time.Second) // simulate delay
		return nil
	}

	// SMTP server configuration.
	smtpHost := "smtp.example.com"
	smtpPort := "587"
	smtpUser := "your_email@example.com"
	smtpPass := "your_password"

	from := smtpUser
	to := []string{"recipient@example.com"}

	subject := "Subject: Weekly Report\n"
	body := "This is the weekly report email sent by your scheduler.\n"

	msg := []byte(subject + "\n" + body)

	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, msg)
	if err != nil {
		return err
	}

	return nil
}
