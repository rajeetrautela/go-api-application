package scheduler

import (
	"os"
	"testing"
	"time"
)

func TestCleanupOldFiles(t *testing.T) {
	sleep = func(time.Duration) {} // mock delay
	defer func() { sleep = time.Sleep }()

	err := cleanupOldFiles()
	if err != nil {
		t.Errorf("cleanupOldFiles failed: %v", err)
	}
}

func TestSendWeeklyEmail_SkipInTestEnv(t *testing.T) {
	sleep = func(time.Duration) {}
	defer func() { sleep = time.Sleep }()

	// Set ENV to TEST so sendWeeklyEmail skips real sending
	os.Setenv("ENV", "TEST")
	defer os.Unsetenv("ENV")

	err := sendWeeklyEmail()
	if err != nil {
		t.Errorf("sendWeeklyEmail failed in TEST env: %v", err)
	}
}

func TestSendWeeklyEmail_RealEnv(t *testing.T) {
	sleep = func(time.Duration) {}
	defer func() { sleep = time.Sleep }()

	// Ensure ENV is NOT TEST
	os.Unsetenv("ENV")

	// Since this will attempt to send real email, you may want to mock smtp.SendMail
	// For this test, just call and expect error or no panic
	err := sendWeeklyEmail()
	// We expect error here because SMTP config is dummy; test should not panic
	if err == nil {
		t.Errorf("expected error in sendWeeklyEmail with dummy SMTP config, got nil")
	}
}

func TestStartCronJobs(t *testing.T) {
	sleep = func(time.Duration) {}
	defer func() { sleep = time.Sleep }()

	err := StartCronJobs()
	if err != nil {
		t.Errorf("StartCronJobs failed: %v", err)
	}
}

func TestCronJobLogicExecution(t *testing.T) {
	sleep = func(time.Duration) {}
	defer func() { sleep = time.Sleep }()

	t.Run("manual cleanup cron logic", func(t *testing.T) {
		err := cleanupOldFiles()
		if err != nil {
			t.Errorf("cleanup job failed: %v", err)
		}
	})

	t.Run("manual email cron logic in TEST env", func(t *testing.T) {
		os.Setenv("ENV", "TEST")
		defer os.Unsetenv("ENV")

		err := sendWeeklyEmail()
		if err != nil {
			t.Errorf("email job failed in TEST env: %v", err)
		}
	})
}
