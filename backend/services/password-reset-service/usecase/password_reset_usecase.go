package usecase

import (
	"context"
	"crypto/md5"
	"database/sql"
	"errors"
	"fmt"
	"github.com/arvians-id/apriori/internal/http/presenter/request"
	"github.com/arvians-id/apriori/internal/model"
	"github.com/arvians-id/apriori/internal/repository"
	"github.com/arvians-id/apriori/util"
	"golang.org/x/crypto/bcrypt"
	"log"
	"strconv"
	"time"
)

type PasswordResetServiceImpl struct {
	PasswordResetRepository repository.PasswordResetRepository
	UserRepository          repository.UserRepository
	DB                      *sql.DB
}

func NewPasswordResetService(
	resetRepository *repository.PasswordResetRepository,
	userRepository *repository.UserRepository,
	db *sql.DB,
) PasswordResetService {
	return &PasswordResetServiceImpl{
		PasswordResetRepository: *resetRepository,
		UserRepository:          *userRepository,
		DB:                      db,
	}
}

func (service *PasswordResetServiceImpl) CreateOrUpdateByEmail(ctx context.Context, email string) (*model.PasswordReset, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("[PasswordResetService][CreateOrUpdateByEmail] problem in db transaction, err: ", err.Error())
		return nil, err
	}
	defer util.CommitOrRollback(tx)

	timestamp := time.Now().Add(1 * time.Hour).Unix()
	timestampString := strconv.Itoa(int(timestamp))
	token := md5.Sum([]byte(email + timestampString))
	tokenString := fmt.Sprintf("%x", token)
	passwordResetRequest := model.PasswordReset{
		Email:   email,
		Token:   tokenString,
		Expired: timestamp,
	}

	// Check if email is exists in table users
	user, err := service.UserRepository.FindByEmail(ctx, tx, email)
	if err != nil {
		log.Println("[NotificationService][CreateOrUpdateByEmail][FindByEmail] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	// Check If email is exists in table password_resets
	_, err = service.PasswordResetRepository.FindByEmail(ctx, tx, user.Email)
	if err != nil {
		// Create new data if not exists
		passwordReset, err := service.PasswordResetRepository.Create(ctx, tx, &passwordResetRequest)
		if err != nil {
			log.Println("[NotificationService][CreateOrUpdateByEmail][Create] problem in getting from repository, err: ", err.Error())
			return nil, err
		}

		return passwordReset, nil
	}

	// Update data if exists
	passwordReset, err := service.PasswordResetRepository.Update(ctx, tx, &passwordResetRequest)
	if err != nil {
		log.Println("[NotificationService][CreateOrUpdateByEmail][Update] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	return passwordReset, nil
}

func (service *PasswordResetServiceImpl) Verify(ctx context.Context, request *request.UpdateResetPasswordUserRequest) error {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Println("[PasswordResetService][Verify] problem in db transaction, err: ", err.Error())
		return err
	}
	defer util.CommitOrRollback(tx)

	// Check if email and token is exists in table password_resets
	passwordResetRequest := model.PasswordReset{
		Email: request.Email,
		Token: request.Token,
	}

	reset, err := service.PasswordResetRepository.FindByEmailAndToken(ctx, tx, &passwordResetRequest)
	if err != nil {
		log.Println("[NotificationService][Verify][FindByEmailAndToken] problem in getting from repository, err: ", err.Error())
		return err
	}

	// Check token expired
	now := time.Now()

	// if expired
	if now.Unix() > reset.Expired {
		err := service.PasswordResetRepository.Delete(ctx, tx, reset.Email)
		if err != nil {
			log.Println("[NotificationService][Verify][Delete] problem in getting from repository, err: ", err.Error())
			return err
		}

		return errors.New("reset password verification is expired")
	}

	// if not
	// Check if email is exists in table users
	user, err := service.UserRepository.FindByEmail(ctx, tx, reset.Email)
	if err != nil {
		log.Println("[NotificationService][Verify][FindByEmail] problem in getting from repository, err: ", err.Error())
		return err
	}

	// Update the password
	timeNow, err := time.Parse(util.TimeFormat, now.Format(util.TimeFormat))
	if err != nil {
		log.Println("[NotificationService][Verify] problem in parsing to time, err: ", err.Error())
		return err
	}

	password, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("[NotificationService][Verify] problem in generating password hashed, err: ", err.Error())
		return err
	}

	user.Password = string(password)
	user.UpdatedAt = timeNow

	err = service.UserRepository.UpdatePassword(ctx, tx, user)
	if err != nil {
		log.Println("[NotificationService][Verify][UpdatePassword] problem in getting from repository, err: ", err.Error())
		return err
	}

	// Delete data from table password_reset
	err = service.PasswordResetRepository.Delete(ctx, tx, user.Email)
	if err != nil {
		log.Println("[NotificationService][Verify][Delete] problem in getting from repository, err: ", err.Error())
		return err
	}

	return nil
}