package usecase

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/arvians-id/go-apriori-microservice/services/payment-service/client"
	"github.com/arvians-id/go-apriori-microservice/services/payment-service/config"
	"github.com/arvians-id/go-apriori-microservice/services/payment-service/model"
	"github.com/arvians-id/go-apriori-microservice/services/payment-service/pb"
	"github.com/arvians-id/go-apriori-microservice/services/payment-service/repository"
	"github.com/arvians-id/go-apriori-microservice/services/payment-service/util"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/veritrans/go-midtrans"
	"log"
	"reflect"
	"strings"
	"time"
)

type PaymentService struct {
	MidClient          midtrans.Client
	SnapGateway        midtrans.SnapGateway
	ServerKey          string
	ClientKey          string
	PaymentRepository  repository.PaymentRepository
	UserOrderService   client.UserOrderServiceClient
	TransactionService client.TransactionServiceClient
	DB                 *sql.DB
}

func NewPaymentService(
	configuration *config.Config,
	paymentRepository repository.PaymentRepository,
	userOrderService client.UserOrderServiceClient,
	transactionService client.TransactionServiceClient,
	db *sql.DB,
) pb.PaymentServiceServer {
	midClient := midtrans.NewClient()
	midClient.ServerKey = configuration.MidtransServerKey
	midClient.ClientKey = configuration.MidtransClientKey
	midClient.APIEnvType = midtrans.Sandbox

	return &PaymentService{
		MidClient:          midClient,
		ServerKey:          midClient.ServerKey,
		ClientKey:          midClient.ClientKey,
		PaymentRepository:  paymentRepository,
		UserOrderService:   userOrderService,
		TransactionService: transactionService,
		DB:                 db,
	}
}

func (service *PaymentService) GetClient() {
	service.SnapGateway = midtrans.SnapGateway{
		Client: service.MidClient,
	}
}

func (service *PaymentService) FindAll(ctx context.Context, empty *empty.Empty) (*pb.ListPaymentResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("[PaymentService][FindAll] problem in db transaction, err: ", err.Error())
		return nil, err
	}
	defer util.CommitOrRollback(tx)

	payments, err := service.PaymentRepository.FindAll(ctx, tx)
	if err != nil {
		log.Println("[PaymentService][FindAll][FindAll] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	var paymentListResponse []*pb.Payment
	for _, payment := range payments {
		paymentListResponse = append(paymentListResponse, payment.ToProtoBuff())
	}

	return &pb.ListPaymentResponse{
		Payment: paymentListResponse,
	}, nil
}

func (service *PaymentService) FindAllByUserId(ctx context.Context, req *pb.GetPaymentByUserIdRequest) (*pb.ListPaymentResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("[PaymentService][FindAllByUserId] problem in db transaction, err: ", err.Error())
		return nil, err
	}
	defer util.CommitOrRollback(tx)

	payments, err := service.PaymentRepository.FindAllByUserId(ctx, tx, req.UserId)
	if err != nil {
		log.Println("[PaymentService][FindAllByUserId][FindAllByUserId] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	var paymentListResponse []*pb.Payment
	for _, payment := range payments {
		paymentListResponse = append(paymentListResponse, payment.ToProtoBuff())
	}

	return &pb.ListPaymentResponse{
		Payment: paymentListResponse,
	}, nil
}

func (service *PaymentService) FindByOrderId(ctx context.Context, req *pb.GetPaymentByOrderIdRequest) (*pb.GetPaymentResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("[PaymentService][FindByOrderId] problem in db transaction, err: ", err.Error())
		return nil, err
	}
	defer util.CommitOrRollback(tx)

	payment, err := service.PaymentRepository.FindByOrderId(ctx, tx, req.OrderId)
	if err != nil {
		log.Println("[PaymentService][FindByOrderId][FindByOrderId] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	return &pb.GetPaymentResponse{
		Payment: payment.ToProtoBuff(),
	}, nil
}

func (service *PaymentService) CreateOrUpdate(ctx context.Context, req *pb.CreatePaymentRequest) (*pb.GetCreatePaymentResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("[PaymentService][CreateOrUpdate] problem in db transaction, err: ", err.Error())
		return &pb.GetCreatePaymentResponse{
			IsSuccess: false,
		}, err
	}
	defer util.CommitOrRollback(tx)

	requestPayment := make(map[string]interface{})
	err = json.Unmarshal(req.Payment, &requestPayment)
	if err != nil {
		log.Println("[PaymentController][Notification] unable to unmarshal json, err: ", err.Error())
		return &pb.GetCreatePaymentResponse{
			IsSuccess: false,
		}, err
	}

	var bankType, vaNumber, billerCode, billKey, settlementTime string

	if requestPayment["va_numbers"] != nil {
		bankType = requestPayment["va_numbers"].([]interface{})[0].(map[string]interface{})["bank"].(string)
		vaNumber = requestPayment["va_numbers"].([]interface{})[0].(map[string]interface{})["va_number"].(string)
	} else if requestPayment["permata_va_number"] != nil {
		bankType = "permata bank"
		vaNumber = requestPayment["permata_va_number"].(string)
	} else if requestPayment["biller_code"] != nil && requestPayment["bill_key"] != nil {
		billerCode = requestPayment["biller_code"].(string)
		billKey = requestPayment["bill_key"].(string)
		bankType = "mandiri"
	}

	setTime, ok := requestPayment["settlement_time"]
	if ok {
		settlementTime = setTime.(string)
	} else {
		settlementTime = ""
	}

	orderID := requestPayment["order_id"].(string)
	transactionTime := requestPayment["transaction_time"].(string)
	transactionStatus := requestPayment["transaction_status"].(string)
	transactionId := requestPayment["transaction_id"].(string)
	statusCode := requestPayment["status_code"].(string)
	signatureKey := requestPayment["signature_key"].(string)
	paymentType := requestPayment["payment_type"].(string)
	merchantId := requestPayment["merchant_id"].(string)
	grossAmount := requestPayment["gross_amount"].(string)
	fraudStatus := requestPayment["fraud_status"].(string)

	checkTransaction, err := service.PaymentRepository.FindByOrderId(ctx, tx, requestPayment["order_id"].(string))
	if err != nil {
		log.Println("[PaymentService][CreateOrUpdate][FindByOrderId] problem in getting from repository, err: ", err.Error())

		return &pb.GetCreatePaymentResponse{
			IsSuccess: false,
		}, err
	}

	checkTransaction.UserId = int64(util.StrToInt(requestPayment["custom_field1"].(string)))
	checkTransaction.OrderId = &orderID
	checkTransaction.TransactionTime = &transactionTime
	checkTransaction.TransactionStatus = &transactionStatus
	checkTransaction.TransactionId = &transactionId
	checkTransaction.StatusCode = &statusCode
	checkTransaction.SignatureKey = &signatureKey
	checkTransaction.SettlementTime = &settlementTime
	checkTransaction.PaymentType = &paymentType
	checkTransaction.MerchantId = &merchantId
	checkTransaction.GrossAmount = &grossAmount
	checkTransaction.FraudStatus = &fraudStatus
	checkTransaction.BankType = &bankType
	checkTransaction.VANumber = &vaNumber
	checkTransaction.BillerCode = &billerCode
	checkTransaction.BillKey = &billKey

	if checkTransaction.OrderId != nil {
		err := service.PaymentRepository.Update(ctx, tx, checkTransaction)
		if err != nil {
			log.Println("[PaymentService][CreateOrUpdate][Update] problem in getting from repository, err: ", err.Error())
			return &pb.GetCreatePaymentResponse{
				IsSuccess: false,
			}, err
		}

		if requestPayment["transaction_status"].(string) == "settlement" {
			payloadId := int64(util.StrToInt(requestPayment["custom_field2"].(string)))
			userOrder, err := service.UserOrderService.FindAllByPayloadId(ctx, payloadId)
			if err != nil {
				log.Println("[PaymentService][CreateOrUpdate][FindAllByPayloadId] problem in getting from repository, err: ", err.Error())
				return &pb.GetCreatePaymentResponse{
					IsSuccess: false,
				}, err
			}

			var productName []string
			for _, item := range userOrder.UserOrder {
				productName = append(productName, *item.Name)
			}

			orderId := requestPayment["order_id"].(string)
			_, err = service.TransactionService.Create(ctx, &pb.CreateTransactionRequest{
				ProductName:   strings.ToLower(strings.Join(productName, ", ")),
				CustomerName:  requestPayment["custom_field3"].(string),
				NoTransaction: &orderId,
			})
			if err != nil {
				log.Println("[PaymentService][CreateOrUpdate][TransactionCreate] problem in getting from repository, err: ", err.Error())
				return &pb.GetCreatePaymentResponse{
					IsSuccess: false,
				}, err
			}

			return &pb.GetCreatePaymentResponse{
				IsSuccess: true,
			}, nil
		}
	}

	return &pb.GetCreatePaymentResponse{
		IsSuccess: false,
	}, nil
}

func (service *PaymentService) UpdateReceiptNumber(ctx context.Context, req *pb.UpdateReceiptNumberRequest) (*pb.GetPaymentResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("[PaymentService][UpdateReceiptNumber] problem in db transaction, err: ", err.Error())
		return nil, err
	}
	defer util.CommitOrRollback(tx)

	payment, err := service.PaymentRepository.FindByOrderId(ctx, tx, req.OrderId)
	if err != nil {
		log.Println("[PaymentService][UpdateReceiptNumber][FindByOrderId] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	payment.ReceiptNumber = &req.ReceiptNumber

	err = service.PaymentRepository.UpdateReceiptNumber(ctx, tx, payment)
	if err != nil {
		log.Println("[PaymentService][UpdateReceiptNumber][UpdateReceiptNumber] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	return &pb.GetPaymentResponse{
		Payment: payment.ToProtoBuff(),
	}, nil
}

func (service *PaymentService) Delete(ctx context.Context, req *pb.GetPaymentByOrderIdRequest) (*empty.Empty, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("[PaymentService][Delete] problem in db transaction, err: ", err.Error())
		return nil, err
	}
	defer util.CommitOrRollback(tx)

	payment, err := service.PaymentRepository.FindByOrderId(ctx, tx, req.OrderId)
	if err != nil {
		log.Println("[PaymentService][Delete][FindByOrderId] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	err = service.PaymentRepository.Delete(ctx, tx, payment.OrderId)
	if err != nil {
		log.Println("[PaymentService][Delete][FindByOrderId] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	return nil, nil
}
func (service *PaymentService) GetToken(ctx context.Context, req *pb.GetPaymentTokenRequest) (*pb.GetPaymentTokenResponse, error) {
	var getClient = func() {
		service.SnapGateway = midtrans.SnapGateway{
			Client: service.MidClient,
		}
	}
	getClient()

	var items []map[string]interface{}
	for _, item := range req.Items {
		err := json.Unmarshal([]byte(item), &items)
		if err != nil {
			log.Println("[PaymentService][GetToken] unable to unmarshal json, err: ", err.Error())
			return nil, err
		}
	}

	var itemDetails []midtrans.ItemDetail
	for _, item := range items {
		code := CheckCode(item)
		itemDetails = append(itemDetails, midtrans.ItemDetail{
			ID:    code,
			Name:  item["name"].(string),
			Price: int64(item["price"].(float64)),
			Qty:   int32(item["quantity"].(float64)),
		})
	}

	itemDetails = append(itemDetails, midtrans.ItemDetail{
		ID:    itemDetails[len(itemDetails)-1].ID,
		Name:  "Pajak",
		Price: 5000,
		Qty:   1,
	})

	itemDetails = append(itemDetails, midtrans.ItemDetail{
		ID:    itemDetails[len(itemDetails)-1].ID,
		Name:  "Ongkos Kirim",
		Price: req.ShippingCost,
		Qty:   1,
	})

	orderID := util.RandomString(20)
	var snapRequest midtrans.SnapReq
	snapRequest.TransactionDetails.OrderID = orderID
	snapRequest.TransactionDetails.GrossAmt = req.GrossAmount
	snapRequest.Items = &itemDetails
	snapRequest.CustomerDetail = &midtrans.CustDetail{
		FName: req.CustomerName,
	}
	snapRequest.CustomField1 = util.IntToStr(int(req.UserId))

	// Save to database
	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("[PaymentService][GetToken] problem in db transaction, err: ", err.Error())
		return nil, err
	}
	defer util.CommitOrRollback(tx)

	canceled := "canceled"
	timeNow := time.Now().Format("2006-01-02 15:04:05")
	paymentRequest := model.Payment{
		UserId:            req.UserId,
		OrderId:           &orderID,
		TransactionStatus: &canceled,
		TransactionTime:   &timeNow,
		Address:           &req.Address,
		Courier:           &req.Courier,
		CourierService:    &req.CourierService,
	}
	payment, err := service.PaymentRepository.Create(ctx, tx, &paymentRequest)
	if err != nil {
		log.Println("[PaymentService][GetToken][Create] problem in getting from repository, err: ", err.Error())
		return nil, err
	}
	// Send id payload
	snapRequest.CustomField2 = util.IntToStr(int(payment.IdPayload))
	snapRequest.CustomField3 = req.CustomerName

	token, err := service.SnapGateway.GetToken(&snapRequest)
	if err != nil {
		log.Println("[PaymentService][GetToken][GetToken] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	var tokenResponse map[string]string
	tokenResponse["clientKey"] = service.ClientKey
	tokenResponse["token"] = token.Token

	for _, item := range items {
		code := CheckCode(item)
		price := int64(item["price"].(float64))
		quantity := int32(item["quantity"].(float64))
		totalPriceItem := int64(item["totalPricePerItem"].(float64))
		name := item["name"].(string)
		image := item["image"].(string)
		_, err := service.UserOrderService.Create(ctx, &pb.CreateUserOrderRequest{
			PayloadId:      payment.IdPayload,
			Code:           &code,
			Name:           &name,
			Price:          &price,
			Image:          &image,
			Quantity:       &quantity,
			TotalPriceItem: &totalPriceItem,
		})
		if err != nil {
			log.Println("[PaymentService][GetToken][Create] problem in getting from repository, err: ", err.Error())
			return nil, err
		}
	}

	return &pb.GetPaymentTokenResponse{
		Payment: tokenResponse,
	}, nil
}

func CheckCode(value map[string]interface{}) string {
	checkCode := reflect.ValueOf(value["code"]).Kind()
	var code string
	if checkCode == reflect.Float64 {
		code = util.IntToStr(int(value["code"].(float64)))
	} else if checkCode == reflect.String {
		code = value["code"].(string)
	}

	return code
}
