package views

type EKYCVerifyIDCard struct {
	Status       bool    `json:"status"`
	ErrorMessage *string `json:"error_message,omitempty"`
}
