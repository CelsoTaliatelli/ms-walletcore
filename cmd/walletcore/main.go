package main

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/CelsoTaliatelli/ms-walletcore/internal/database"
	"github.com/CelsoTaliatelli/ms-walletcore/internal/event"
	"github.com/CelsoTaliatelli/ms-walletcore/internal/event/handler"
	"github.com/CelsoTaliatelli/ms-walletcore/internal/usecase/create_account"
	"github.com/CelsoTaliatelli/ms-walletcore/internal/usecase/create_client"
	"github.com/CelsoTaliatelli/ms-walletcore/internal/usecase/create_transaction"
	"github.com/CelsoTaliatelli/ms-walletcore/internal/web"
	"github.com/CelsoTaliatelli/ms-walletcore/internal/web/webserver"
	"github.com/CelsoTaliatelli/ms-walletcore/pkg/events"
	"github.com/CelsoTaliatelli/ms-walletcore/pkg/kafka"
	"github.com/CelsoTaliatelli/ms-walletcore/pkg/uow"
	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", "root", "root", "mysql", "3306", "wallet"))
	if err != nil {
		println(err)
		panic(err)
	}
	defer db.Close()

	configMap := ckafka.ConfigMap{
		"bootstrap.servers": "kafka:29092",
		"group.id":          "wallet",
	}
	kafkaProducer := kafka.NewKafkaProducer(&configMap)

	eventDispatcher := events.NewEventDispatcher()
	eventDispatcher.Register("TransactionCreated", handler.NewTransactionCreatedKafkaHandler(kafkaProducer))
	transactionCreatedEvent := event.NewTransactionCreated()
	balanceUpdatedEvent := event.NewBalanceUpdated()

	clientDb := database.NewClientDB(db)
	accountDb := database.NewAccountDB(db)

	ctx := context.Background()
	uow := uow.NewUow(ctx, db)

	uow.Register("TransactionDB", func(tx *sql.Tx) interface{} {
		return database.NewTransactionDB(db)
	})

	createCientUseCase := create_client.NewCreateClientUseCase(clientDb)
	createAccountUseCase := create_account.NewCreateAccountUseCase(accountDb, clientDb)
	createTransactionUseCase := create_transaction.NewCreateTransactionUseCase(uow, eventDispatcher, transactionCreatedEvent, balanceUpdatedEvent)

	webserver := webserver.NewWebServer(":8000")

	clientHandler := web.NewWebClientHandler(*createCientUseCase)
	accountHandler := web.NewWebAccountHandler(*createAccountUseCase)
	transactionHandler := web.NewWebTransactionHandler(*createTransactionUseCase)

	webserver.AddHandler("/clients", clientHandler.CreateClient)
	webserver.AddHandler("/accounts", accountHandler.CreateAccount)
	webserver.AddHandler("/transactions", transactionHandler.CreateTransaction)

	webserver.Start()
}
