package views

import (
	"gitlab.finema.co/finema/etda/web-portal-api/consts"
	"gitlab.finema.co/finema/etda/web-portal-api/models"
)

type RequestedVPWithQRCode struct {
	models.RequestedVP
	RequestedVPSchemaType []models.RequestedVPSchemaType `json:"schema_list,omitempty"`
	QRCodeData            *RequestedVPQRCodeData         `json:"qr_code_data,omitempty"`
	SubmittedCount        int64                          `json:"submitted_count"`
}
type RequestedVPQRCodeData struct {
	Endpoint  string `json:"endpoint"`
	Operation string `json:"operation"`
}

func NewRequestedVPWithQRCode(requestedVP *models.RequestedVP, schemaType []models.RequestedVPSchemaType, endpoint string, submittedCount int64) *RequestedVPWithQRCode {
	//TODO: Need to calculate submitted_count to be return
	return &RequestedVPWithQRCode{
		RequestedVP:           *requestedVP,
		RequestedVPSchemaType: schemaType,
		QRCodeData: &RequestedVPQRCodeData{
			Endpoint:  endpoint,
			Operation: consts.OperationGetRequestVP,
		},
		SubmittedCount: submittedCount,
	}

}

type RequestedVP struct {
	models.RequestedVP
	RequestedVPSchemaType []models.RequestedVPSchemaType `json:"schema_list"`
	VerifierDID           string                         `json:"verifier_did"`
	Verifier              string                         `json:"verifier"`
	Endpoint              string                         `json:"endpoint"`
}

func NewRequestedVPList(requestedVPs []models.RequestedVP, submittedCounts []int64) []RequestedVPWithQRCode {
	requestedVPList := make([]RequestedVPWithQRCode, 0)
	for index, requestedVP := range requestedVPs {
		requestedVPView := &RequestedVPWithQRCode{
			RequestedVP:    requestedVP,
			SubmittedCount: submittedCounts[index],
		}
		requestedVPList = append(requestedVPList, *requestedVPView)
	}
	return requestedVPList
}
