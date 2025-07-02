package worker

import (
	"fmt"
	"time"
)

var EmailQueue = make(chan EmailJob, 100) // buffer 100 job

type EmailSender interface {
	SendWelcomeEmail(to, name, extra string) error
	SendVerificationEmail(to, name, resetCode string) error
}

func StartEmailWorker(emailSender EmailSender) {
	for i := 0; i < 3; i++ {
		go func(id int) {
			for job := range EmailQueue {
				switch job.Type {
				case "welcome":
					_ = emailSender.SendWelcomeEmail(job.To, job.Name, "")
				case "verification":
					_ = emailSender.SendVerificationEmail(job.To, job.Name, job.ResetCode)
				}
			}
		}(i)
	}
}

func processEmailJob(workerID int, job EmailJob) {
	fmt.Printf("[Worker %d] Sending email to %s\n", workerID, job.To)
	time.Sleep(2 * time.Second) // simulasi delay
	fmt.Printf("[Worker %d] Email sent to %s\n", workerID, job.To)
}
