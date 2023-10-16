package gateway

import "github.com/CelsoTaliatelli/ms-walletcore/internal/entity"

type TransactionGateway interface {
	Create(transaction *entity.Transaction) error
}
