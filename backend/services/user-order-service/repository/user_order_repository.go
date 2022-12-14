package repository

import (
	"context"
	"database/sql"
	"github.com/arvians-id/go-apriori-microservice/services/user-order-service/model"
	"log"
)

type UserOrderRepositoryImpl struct {
}

func NewUserOrderRepository() UserOrderRepository {
	return &UserOrderRepositoryImpl{}
}

func (repository *UserOrderRepositoryImpl) FindAllByPayloadId(ctx context.Context, tx *sql.Tx, payloadId int64) ([]*model.UserOrder, error) {
	query := "SELECT * FROM user_orders WHERE payload_id = $1"
	rows, err := tx.QueryContext(ctx, query, payloadId)
	if err != nil {
		log.Println("[UserOrderRepository][FindAllByPayloadId] problem querying to db, err: ", err.Error())
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Println("[UserOrderRepository][FindAllByPayloadId] problem closing query from db, err: ", err.Error())
			return
		}
	}(rows)

	var userOrders []*model.UserOrder
	for rows.Next() {
		var userOrder model.UserOrder
		err := rows.Scan(
			&userOrder.IdOrder,
			&userOrder.PayloadId,
			&userOrder.Code,
			&userOrder.Name,
			&userOrder.Price,
			&userOrder.Image,
			&userOrder.Quantity,
			&userOrder.TotalPriceItem,
		)
		if err != nil {
			log.Println("[UserOrderRepository][FindAllByPayloadId] problem with scanning db row, err: ", err.Error())
			return nil, err
		}

		userOrders = append(userOrders, &userOrder)
	}

	return userOrders, nil
}

func (repository *UserOrderRepositoryImpl) FindAllByUserId(ctx context.Context, tx *sql.Tx, userId int64) ([]*model.UserOrder, error) {
	query := `SELECT 
				id_order,
			    payload_id,
			    code,
      			name,
			    price,
			    image,
			    quantity,
			    total_price_item,
			    order_id,
			    transaction_status 
			  FROM user_orders uo 
			  	LEFT JOIN payloads p ON p.id_payload = uo.payload_id 
		   	  WHERE p.user_id = $1 AND p.transaction_status = 'settlement'
			  ORDER BY uo.id_order DESC`
	rows, err := tx.QueryContext(ctx, query, userId)
	if err != nil {
		log.Println("[UserOrderRepository][FindAllByUserId] problem querying to db, err: ", err.Error())
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Println("[UserOrderRepository][FindAllByUserId] problem closing query from db, err: ", err.Error())
			return
		}
	}(rows)

	var userOrders []*model.UserOrder
	for rows.Next() {
		userOrder := model.UserOrder{
			Payment: &model.Payment{},
		}
		err := rows.Scan(
			&userOrder.IdOrder,
			&userOrder.PayloadId,
			&userOrder.Code,
			&userOrder.Name,
			&userOrder.Price,
			&userOrder.Image,
			&userOrder.Quantity,
			&userOrder.TotalPriceItem,
			&userOrder.Payment.OrderId,
			&userOrder.Payment.TransactionStatus,
		)
		if err != nil {
			log.Println("[UserOrderRepository][FindAllByUserId] problem with scanning db row, err: ", err.Error())
			return nil, err
		}

		userOrders = append(userOrders, &userOrder)
	}

	return userOrders, nil
}

func (repository *UserOrderRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, id int64) (*model.UserOrder, error) {
	query := `SELECT * FROM user_orders WHERE id_order = $1`
	row := tx.QueryRowContext(ctx, query, id)

	var userOrder model.UserOrder
	err := row.Scan(
		&userOrder.IdOrder,
		&userOrder.PayloadId,
		&userOrder.Code,
		&userOrder.Name,
		&userOrder.Price,
		&userOrder.Image,
		&userOrder.Quantity,
		&userOrder.TotalPriceItem,
	)
	if err != nil {
		log.Println("[UserOrderRepository][FindById] problem with scanning db row, err: ", err.Error())
		return nil, err
	}

	return &userOrder, nil
}

func (repository *UserOrderRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, userOrder *model.UserOrder) (*model.UserOrder, error) {
	var id int64
	query := `INSERT INTO user_orders(payload_id,code,name,price,image,quantity,total_price_item) VALUES($1,$2,$3,$4,$5,$6,$7) RETURNING id_order`
	row := tx.QueryRowContext(
		ctx,
		query,
		userOrder.PayloadId,
		userOrder.Code,
		userOrder.Name,
		userOrder.Price,
		userOrder.Image,
		userOrder.Quantity,
		userOrder.TotalPriceItem,
	)
	err := row.Scan(&id)
	if err != nil {
		log.Println("[UserOrderRepository][Create] problem with scanning db row, err: ", err.Error())
		return nil, err
	}

	userOrder.IdOrder = id

	return userOrder, nil
}
