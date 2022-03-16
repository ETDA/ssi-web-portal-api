package views

type KeySign struct {
	Signature string `json:"signature"`
	Message   string `json:"message"`
}

func NewKeySign(signature string, message string) *KeySign {
	return &KeySign{
		Signature: signature,
		Message:   message,
	}
}
