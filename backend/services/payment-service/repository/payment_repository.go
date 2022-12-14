package repository

import (
	"context"
	"database/sql"
	"github.com/arvians-id/go-apriori-microservice/services/payment-service/model"
	"log"
)

type PaymentRepositoryImpl struct {
}

func NewPaymentRepository() PaymentRepository {
	return &PaymentRepositoryImpl{}
}

func (repository *PaymentRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx) ([]*model.Payment, error) {
	query := `SELECT payloads.*,users.name 
			  FROM payloads
			  	LEFT JOIN users ON users.id_user = payloads.user_id
			  ORDER BY payloads.settlement_time DESC, payloads.bank_type DESC`
	rows, err := tx.QueryContext(ctx, query)
	if err != nil {
		log.Println("[PaymentRepository][FindAll] problem querying to db, err: ", err.Error())
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Println("[PaymentRepository][FindAll] problem closing query from db, err: ", err.Error())
			return
		}
	}(rows)

	var payments []*model.Payment
	for rows.Next() {
		payment := model.Payment{
			User: &model.User{},
		}
		err := rows.Scan(
			&payment.IdPayload,
			&payment.UserId,
			&payment.OrderId,
			&payment.TransactionTime,
			&payment.TransactionStatus,
			&payment.TransactionId,
			&payment.StatusCode,
			&payment.SignatureKey,
			&payment.SettlementTime,
			&payment.PaymentType,
			&payment.MerchantId,
			&payment.GrossAmount,
			&payment.FraudStatus,
			&payment.BankType,
			&payment.VANumber,
			&payment.BillerCode,
			&payment.BillKey,
			&payment.ReceiptNumber,
			&payment.Address,
			&payment.Courier,
			&payment.CourierService,
			&payment.User.Name,
		)
		if err != nil {
			log.Println("[PaymentRepository][FindAll] problem with scanning db row, err: ", err.Error())
			return nil, err
		}

		payments = append(payments, &payment)
	}

	return payments, nil
}

func (repository *PaymentRepositoryImpl) FindAllByUserId(ctx context.Context, tx *sql.Tx, userId int64) ([]*model.Payment, error) {
	query := `SELECT * FROM payloads 
			  WHERE user_id = $1 
			  ORDER BY settlement_time DESC, bank_type DESC`
	rows, err := tx.QueryContext(ctx, query, userId)
	if err != nil {
		log.Println("[PaymentRepository][FindAllByUserId] problem querying to db, err: ", err.Error())
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Println("[PaymentRepository][FindAllByUserId] problem closing query from db, err: ", err.Error())
			return
		}
	}(rows)

	var payments []*model.Payment
	for rows.Next() {
		var payment model.Payment
		err := rows.Scan(
			&payment.IdPayload,
			&payment.UserId,
			&payment.OrderId,
			&payment.TransactionTime,
			&payment.TransactionStatus,
			&payment.TransactionId,
			&payment.StatusCode,
			&payment.SignatureKey,
			&payment.SettlementTime,
			&payment.PaymentType,
			&payment.MerchantId,
			&payment.GrossAmount,
			&payment.FraudStatus,
			&payment.BankType,
			&payment.VANumber,
			&payment.BillerCode,
			&payment.BillKey,
			&payment.ReceiptNumber,
			&payment.Address,
			&payment.Courier,
			&payment.CourierService,
		)
		if err != nil {
			log.Println("[PaymentRepository][FindAllByUserId] problem with scanning db row, err: ", err.Error())
			return nil, err
		}

		payments = append(payments, &payment)
	}

	return payments, nil
}

func (repository *PaymentRepositoryImpl) FindByOrderId(ctx context.Context, tx *sql.Tx, orderId string) (*model.Payment, error) {
	query := "SELECT * FROM payloads WHERE order_id = $1"
	row := tx.QueryRowContext(ctx, query, orderId)

	var payment model.Payment
	err := row.Scan(
		&payment.IdPayload,
		&payment.UserId,
		&payment.OrderId,
		&payment.TransactionTime,
		&payment.TransactionStatus,
		&payment.TransactionId,
		&payment.StatusCode,
		&payment.SignatureKey,
		&payment.SettlementTime,
		&payment.PaymentType,
		&payment.MerchantId,
		&payment.GrossAmount,
		&payment.FraudStatus,
		&payment.BankType,
		&payment.VANumber,
		&payment.BillerCode,
		&payment.BillKey,
		&payment.ReceiptNumber,
		&payment.Address,
		&payment.Courier,
		&payment.CourierService,
	)
	if err != nil {
		log.Println("[PaymentRepository][FindByOrderId] problem with scanning db row, err: ", err.Error())
		return nil, err
	}

	return &payment, nil
}

func (repository *PaymentRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, payment *model.Payment) (*model.Payment, error) {
	var id int64
	query := `INSERT INTO payloads(user_id,order_id,transaction_time,transaction_status,transaction_id,status_code,signature_key,settlement_time,payment_type,merchant_id,gross_amount,fraud_status,bank_type,va_number,biller_code,bill_key,receipt_number,address,courier,courier_service) 
			  VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20)  RETURNING id_payload`
	row := tx.QueryRowContext(
		ctx,
		query,
		payment.UserId,
		payment.OrderId,
		payment.TransactionTime,
		payment.TransactionStatus,
		payment.TransactionId,
		payment.StatusCode,
		payment.SignatureKey,
		payment.SettlementTime,
		payment.PaymentType,
		payment.MerchantId,
		payment.GrossAmount,
		payment.FraudStatus,
		payment.BankType,
		payment.VANumber,
		payment.BillerCode,
		payment.BillKey,
		payment.ReceiptNumber,
		payment.Address,
		payment.Courier,
		payment.CourierService,
	)
	err := row.Scan(&id)
	if err != nil {
		log.Println("[PaymentRepository][Create] problem with scanning db row, err: ", err.Error())
		return nil, err
	}

	payment.IdPayload = id

	return payment, nil
}

func (repository *PaymentRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, payment *model.Payment) error {
	query := `UPDATE payloads 
         	  SET user_id = $1,
         	      order_id = $2,
         	      transaction_time = $3,
         	      transaction_status = $4,
         	      transaction_id = $5,
				  status_code = $6,
				  signature_key = $7,
				  settlement_time = $8,
         	      payment_type = $9,
         	      merchant_id = $10,
         	      gross_amount = $11,
         	      fraud_status = $12,
		 	      bank_type = $13,
         	      va_number = $14,
         	      biller_code = $15,
         	      bill_key = $16
			  WHERE order_id = $17`
	_, err := tx.ExecContext(
		ctx,
		query,
		payment.UserId,
		payment.OrderId,
		payment.TransactionTime,
		payment.TransactionStatus,
		payment.TransactionId,
		payment.StatusCode,
		payment.SignatureKey,
		payment.SettlementTime,
		payment.PaymentType,
		payment.MerchantId,
		payment.GrossAmount,
		payment.FraudStatus,
		payment.BankType,
		payment.VANumber,
		payment.BillerCode,
		payment.BillKey,
		payment.OrderId,
	)
	if err != nil {
		log.Println("[PaymentRepository][Update] problem querying to db, err: ", err.Error())
		return err
	}

	return nil
}

func (repository *PaymentRepositoryImpl) UpdateReceiptNumber(ctx context.Context, tx *sql.Tx, payment *model.Payment) error {
	query := `UPDATE payloads SET receipt_number = $1 WHERE order_id = $2`
	_, err := tx.ExecContext(ctx, query, payment.ReceiptNumber, payment.OrderId)
	if err != nil {
		log.Println("[PaymentRepository][UpdateReceiptNumber] problem querying to db, err: ", err.Error())
		return err
	}

	return nil
}

func (repository *PaymentRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, orderId *string) error {
	query := "DELETE FROM payloads WHERE order_id = $1"
	_, err := tx.ExecContext(ctx, query, orderId)
	if err != nil {
		log.Println("[PaymentRepository][Delete] problem querying to db, err: ", err.Error())
		return err
	}

	return nil
}
