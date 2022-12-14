package usecase

import (
	"context"
	"crypto/md5"
	"database/sql"
	"errors"
	"fmt"
	"github.com/arvians-id/go-apriori-microservice/services/password-reset-service/client"
	"github.com/arvians-id/go-apriori-microservice/services/password-reset-service/model"
	"github.com/arvians-id/go-apriori-microservice/services/password-reset-service/pb"
	"github.com/arvians-id/go-apriori-microservice/services/password-reset-service/repository"
	"github.com/arvians-id/go-apriori-microservice/services/password-reset-service/util"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/protobuf/types/known/emptypb"
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
		log.Println("error1", err)
		log.Println("[PasswordResetService][CreateOrUpdateByEmail][FindByEmail] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	// Check If email is exists in table password_resets
	_, err = service.PasswordResetRepository.FindByEmail(ctx, tx, user.User.Email)
	if err != nil {
		// Create new data if not exists
		passwordReset, err := service.PasswordResetRepository.Create(ctx, tx, &passwordResetRequest)
		if err != nil {
			log.Println("[PasswordResetService][CreateOrUpdateByEmail][Create] problem in getting from repository, err: ", err.Error())
			return nil, err
		}

		return &pb.GetPasswordResetResponse{
			PasswordReset: passwordReset.ToProtoBuff(),
		}, nil
	}

	// Update data if exists
	passwordReset, err := service.PasswordResetRepository.Update(ctx, tx, &passwordResetRequest)
	if err != nil {
		log.Println("[PasswordResetService][CreateOrUpdateByEmail][Update] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	return &pb.GetPasswordResetResponse{
		PasswordReset: passwordReset.ToProtoBuff(),
	}, nil
}

func (service *PasswordResetService) Verify(ctx context.Context, req *pb.GetVerifyRequest) (*emptypb.Empty, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("[PasswordResetService][Verify] problem in db transaction, err: ", err.Error())
		return new(emptypb.Empty), err
	}
	defer util.CommitOrRollback(tx)

	// Check if email and token is exists in table password_resets
	passwordResetRequest := model.PasswordReset{
		Email: req.Email,
		Token: req.Token,
	}

	reset, err := service.PasswordResetRepository.FindByEmailAndToken(ctx, tx, &passwordResetRequest)
	if err != nil {
		log.Println("[PasswordResetService][Verify][FindByEmailAndToken] problem in getting from repository, err: ", err.Error())
		return new(emptypb.Empty), err
	}

	// Check token expired
	now := time.Now()

	// if expired
	if now.Unix() > reset.Expired {
		err := service.PasswordResetRepository.Delete(ctx, tx, reset.Email)
		if err != nil {
			log.Println("[PasswordResetService][Verify][Delete] problem in getting from repository, err: ", err.Error())
			return new(emptypb.Empty), err
		}

		return new(emptypb.Empty), errors.New("reset password verification is expired")
	}

	// if not
	// Check if email is exists in table users
	user, err := service.UserService.FindByEmail(ctx, req.Email)
	if err != nil {
		log.Println("[PasswordResetService][Verify][FindByEmail] problem in getting from repository, err: ", err.Error())
		return new(emptypb.Empty), err
	}

	// Update the password
	password, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("[PasswordResetService][Verify] problem in generating password hashed, err: ", err.Error())
		return new(emptypb.Empty), err
	}

	_, err = service.UserService.UpdatePassword(ctx, &pb.UpdateUserPasswordRequest{
		Email:    req.Email,
		Password: string(password),
	})
	if err != nil {
		log.Println("[PasswordResetService][Verify][UpdatePassword] problem in getting from repository, err: ", err.Error())
		return new(emptypb.Empty), err
	}

	// Delete data from table password_reset
	err = service.PasswordResetRepository.Delete(ctx, tx, user.User.Email)
	if err != nil {
		log.Println("[PasswordResetService][Verify][Delete] problem in getting from repository, err: ", err.Error())
		return new(emptypb.Empty), err
	}

	return new(emptypb.Empty), nil
}
