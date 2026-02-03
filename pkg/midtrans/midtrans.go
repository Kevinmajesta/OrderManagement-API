package midtrans

import (
	"Kevinmajesta/OrderManagementAPI/configs"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

type MidtransService struct {
	Client snap.Client
}

func NewMidtransService(cfg *configs.MidtransConfig) *MidtransService {
	var client snap.Client
	client.New(cfg.ServerKey, midtrans.Sandbox)

	// Set to production if configured
	if cfg.IsProduction == "true" {
		client.New(cfg.ServerKey, midtrans.Production)
	}

	return &MidtransService{
		Client: client,
	}
}

func (s *MidtransService) CreateTransaction(orderID string, grossAmount int64, customerName, customerEmail, customerPhone string) (*snap.Response, error) {
	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  orderID,
			GrossAmt: grossAmount,
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName: customerName,
			Email: customerEmail,
			Phone: customerPhone,
		},
	}

	snapResp, err := s.Client.CreateTransaction(req)
	if err != nil {
		return nil, err
	}

	return snapResp, nil
}
