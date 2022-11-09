package usecase

import (
	"context"
	"crypto/md5"
	"database/sql"
	"errors"
	"fmt"
	"github.com/arvians-id/go-apriori-microservice/adapter/pkg/auth/pb"
	pbuser "github.com/arvians-id/go-apriori-microservice/adapter/pkg/user/pb"
	"github.com/arvians-id/go-apriori-microservice/model"
	"github.com/arvians-id/go-apriori-microservice/services/password-reset-service/client"
	"github.com/arvians-id/go-apriori-microservice/services/password-reset-service/repository"
	"github.com/arvians-id/go-apriori-microservice/util"
	"github.com/golang/protobuf/ptypes/empty"
	"golang.org/x/crypto/bcrypt"
	"log"
	"strconv"
	"time"
)

type PasswordResetService struct {
	PasswordResetRepository repository.PasswordResetRepository
	UserService             client.UserServiceClient
	DB                      *sql.DB
}

func NewPasswordResetService(
	resetRepository repository.PasswordResetRepository,
	userService client.UserServiceClient,
	db *sql.DB,
) pb.PasswordResetServiceServer {
	return &PasswordResetService{
		PasswordResetRepository: resetRepository,
		UserService:             userService,
		DB:                      db,
	}
}

func (service *PasswordResetService) CreateOrUpdateByEmail(ctx context.Context, req *pb.GetPasswordResetByEmailRequest) (*pb.GetPasswordResetResponse, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("[PasswordResetService][CreateOrUpdateByEmail] problem in db transaction, err: ", err.Error())
		return nil, err
	}
	defer util.CommitOrRollback(tx)

	timestamp := time.Now().Add(1 * time.Hour).Unix()
	timestampString := strconv.Itoa(int(timestamp))
	token := md5.Sum([]byte(req.Email + timestampString))
	tokenString := fmt.Sprintf("%x", token)
	passwordResetRequest := model.PasswordReset{
		Email:   req.Email,
		Token:   tokenString,
		Expired: timestamp,
	}

	// Check if email is exists in table users
	user, err := service.UserService.FindByEmail(ctx, req.Email)
	if err != nil {
		log.Println("[NotificationService][CreateOrUpdateByEmail][FindByEmail] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	// Check If email is exists in table password_resets
	_, err = service.PasswordResetRepository.FindByEmail(ctx, tx, user.User.Email)
	if err != nil {
		// Create new data if not exists
		passwordReset, err := service.PasswordResetRepository.Create(ctx, tx, &passwordResetRequest)
		if err != nil {
			log.Println("[NotificationService][CreateOrUpdateByEmail][Create] problem in getting from repository, err: ", err.Error())
			return nil, err
		}

		return &pb.GetPasswordResetResponse{
			PasswordReset: passwordReset.ToProtoBuff(),
		}, nil
	}

	// Update data if exists
	passwordReset, err := service.PasswordResetRepository.Update(ctx, tx, &passwordResetRequest)
	if err != nil {
		log.Println("[NotificationService][CreateOrUpdateByEmail][Update] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	return &pb.GetPasswordResetResponse{
		PasswordReset: passwordReset.ToProtoBuff(),
	}, nil
}

func (service *PasswordResetService) Verify(ctx context.Context, req *pb.GetVerifyRequest) (*empty.Empty, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("[PasswordResetService][Verify] problem in db transaction, err: ", err.Error())
		return nil, err
	}
	defer util.CommitOrRollback(tx)

	// Check if email and token is exists in table password_resets
	passwordResetRequest := model.PasswordReset{
		Email: req.Email,
		Token: req.Token,
	}

	reset, err := service.PasswordResetRepository.FindByEmailAndToken(ctx, tx, &passwordResetRequest)
	if err != nil {
		log.Println("[NotificationService][Verify][FindByEmailAndToken] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	// Check token expired
	now := time.Now()

	// if expired
	if now.Unix() > reset.Expired {
		err := service.PasswordResetRepository.Delete(ctx, tx, reset.Email)
		if err != nil {
			log.Println("[NotificationService][Verify][Delete] problem in getting from repository, err: ", err.Error())
			return nil, err
		}

		return nil, errors.New("reset password verification is expired")
	}

	// if not
	// Check if email is exists in table users
	user, err := service.UserService.FindByEmail(ctx, req.Email)
	if err != nil {
		log.Println("[NotificationService][Verify][FindByEmail] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	// Update the password
	password, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("[NotificationService][Verify] problem in generating password hashed, err: ", err.Error())
		return nil, err
	}

	_, err = service.UserService.UpdatePassword(ctx, &pbuser.UpdateUserPasswordRequest{
		Email:    req.Email,
		Password: string(password),
	})
	if err != nil {
		log.Println("[NotificationService][Verify][UpdatePassword] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	// Delete data from table password_reset
	err = service.PasswordResetRepository.Delete(ctx, tx, user.User.Email)
	if err != nil {
		log.Println("[NotificationService][Verify][Delete] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	return nil, nil
}
