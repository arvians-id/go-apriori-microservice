package model

import (
	"github.com/arvians-id/go-apriori-microservice/services/user-order-service/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type Transaction struct {
	IdTransaction int64     `json:"id_transaction"`
	ProductName   string    `json:"product_name"`
	CustomerName  string    `json:"customer_name"`
	NoTransaction string    `json:"no_transaction"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (transaction *Transaction) ToProtoBuff() *pb.Transaction {
	return &pb.Transaction{
		IdTransaction: transaction.IdTransaction,
		ProductName:   transaction.ProductName,
		CustomerName:  transaction.CustomerName,
		NoTransaction: transaction.NoTransaction,
		CreatedAt:     timestamppb.New(transaction.CreatedAt),
		UpdatedAt:     timestamppb.New(transaction.UpdatedAt),
	}
}
