package gateway

import "github.com/CelsoTaliatelli/ms-walletcore/internal/entity"

type TransactionGateway interface {
	Create(transaction *entity.Transaction) error
	findByID(id string) (*entity.Account, error)
	UpdateBalance(entity *entity.Account) error
}
