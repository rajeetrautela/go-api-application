package scheduler

import (
	"log"
	"time"

	"github.com/robfig/cron/v3"
)

func StartCronJobs() {
	c := cron.New()

	// Cleanup job every day at midnight
	_, err := c.AddFunc("0 0 * * *", func() {
		log.Println("🧹 Running cleanup job...")
		err := cleanupOldFiles()
		if err != nil {
			log.Printf(" Cleanup failed: %v\n", err)
		} else {
			log.Println(" Cleanup completed successfully.")
		}
	})
	if err != nil {
		log.Fatalf("Failed to schedule cleanup job: %v", err)
	}

	// Email job every Monday at 9 AM
	_, err = c.AddFunc("0 9 * * MON", func() {
		log.Println("📧 Sending weekly report email...")
		err := sendWeeklyEmail()
		if err != nil {
			log.Printf(" Email sending failed: %v\n", err)
		} else {
			log.Println(" Email sent successfully.")
		}
	})
	if err != nil {
		log.Fatalf("Failed to schedule email job: %v", err)
	}

	c.Start()
	log.Println(" Cron jobs started.")
}

func cleanupOldFiles() error {
	time.Sleep(1 * time.Second)
	return nil
}

func sendWeeklyEmail() error {
	time.Sleep(2 * time.Second)
	return nil
}
