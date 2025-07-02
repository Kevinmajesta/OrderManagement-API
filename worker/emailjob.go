package worker

type EmailJob struct {
	Type     string 
	To       string
	Name     string
	ResetCode string
}
