package scheduler

import (
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

func TestSendWeeklyEmail(t *testing.T) {
	sleep = func(time.Duration) {}
	defer func() { sleep = time.Sleep }()

	err := sendWeeklyEmail()
	if err != nil {
		t.Errorf("sendWeeklyEmail failed: %v", err)
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

	t.Run("manual email cron logic", func(t *testing.T) {
		err := sendWeeklyEmail()
		if err != nil {
			t.Errorf("email job failed: %v", err)
		}
	})
}
