package usecase

import (
	"context"
	"database/sql"
	"errors"
	"github.com/arvians-id/go-apriori-microservice/services/user-service/model"
	"github.com/arvians-id/go-apriori-microservice/services/user-service/pb"
	"github.com/arvians-id/go-apriori-microservice/services/user-service/repository"
	"github.com/arvians-id/go-apriori-microservice/services/user-service/util"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"time"
)

type UserService struct {
	UserRepository repository.UserRepository
	DB             *sql.DB
}

func NewUserService(userRepository repository.UserRepository, db *sql.DB) pb.UserServiceServer {
	return &UserService{
		UserRepository: userRepository,
		DB:             db,
	}
}

func (service *UserService) FindAll(ctx context.Context, empty *emptypb.Empty) (*pb.ListUserResponse, error) {
	var tx *sql.Tx
	if service.DB != nil {
		transaction, err := service.DB.Begin()
		if err != nil {
			log.Println("[UserService][FindAll] problem in db transaction, err: ", err.Error())
			return nil, err
		}
		tx = transaction
	}
	defer util.CommitOrRollback(tx)

	users, err := service.UserRepository.FindAll(ctx, tx)
	if err != nil {
		log.Println("[UserService][FindAll][FindAll] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	var userListResponse []*pb.User
	for _, user := range users {
		userListResponse = append(userListResponse, user.ToProtoBuff())
	}

	return &pb.ListUserResponse{
		User: userListResponse,
	}, nil
}

func (service *UserService) FindById(ctx context.Context, req *pb.GetUserByIdRequest) (*pb.GetUserResponse, error) {
	var tx *sql.Tx
	if service.DB != nil {
		transaction, err := service.DB.Begin()
		if err != nil {
			log.Println("[UserService][FindById] problem in db transaction, err: ", err.Error())
			return nil, err
		}
		tx = transaction
	}
	defer util.CommitOrRollback(tx)

	user, err := service.UserRepository.FindById(ctx, tx, req.Id)
	if err != nil {
		log.Println("[UserService][FindById][FindById] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	return &pb.GetUserResponse{
		User: user.ToProtoBuff(),
	}, nil
}

func (service *UserService) FindByEmail(ctx context.Context, req *pb.GetUserByEmailRequest) (*pb.GetUserResponse, error) {
	var tx *sql.Tx
	if service.DB != nil {
		transaction, err := service.DB.Begin()
		if err != nil {
			log.Println("[UserService][FindByEmail] problem in db transaction, err: ", err.Error())
			return nil, err
		}
		tx = transaction
	}
	defer util.CommitOrRollback(tx)

	user, err := service.UserRepository.FindByEmail(ctx, tx, req.Email)
	if err != nil {
		log.Println("[UserService][FindByEmail][FindByEmail] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	return &pb.GetUserResponse{
		User: user.ToProtoBuff(),
	}, nil
}

func (service *UserService) VerifyCredential(ctx context.Context, req *pb.GetVerifyCredentialRequest) (*pb.GetUserResponse, error) {
	var tx *sql.Tx
	if service.DB != nil {
		transaction, err := service.DB.Begin()
		if err != nil {
			log.Println("[UserService][FindByEmail] problem in db transaction, err: ", err.Error())
			return nil, err
		}
		tx = transaction
	}
	defer util.CommitOrRollback(tx)

	user, err := service.UserRepository.FindByEmail(ctx, tx, req.Email)
	if err != nil {
		log.Println("[UserService][FindByEmail][FindByEmail] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		log.Println("[UserService][FindByEmail] problem in comparing password, err: ", err.Error())
		return nil, errors.New("wrong password")
	}
	return &pb.GetUserResponse{
		User: user.ToProtoBuff(),
	}, nil
}

func (service *UserService) Create(ctx context.Context, req *pb.CreateUserRequest) (*pb.GetUserResponse, error) {
	var tx *sql.Tx
	if service.DB != nil {
		transaction, err := service.DB.Begin()
		if err != nil {
			log.Println("[UserService][Create] problem in db transaction, err: ", err.Error())
			return nil, err
		}
		tx = transaction
	}
	defer util.CommitOrRollback(tx)

	password, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("[UserService][Create] problem in generating password hashed, err: ", err.Error())
		return nil, err
	}

	timeNow, err := time.Parse(util.TimeFormat, time.Now().Format(util.TimeFormat))
	if err != nil {
		log.Println("[UserService][Create] problem in parsing to time, err: ", err.Error())
		return nil, err
	}

	userRequest := model.User{
		Role:      2,
		Name:      req.Name,
		Email:     req.Email,
		Address:   req.Address,
		Phone:     req.Phone,
		Password:  string(password),
		CreatedAt: timeNow,
		UpdatedAt: timeNow,
	}
	user, err := service.UserRepository.Create(ctx, tx, &userRequest)
	if err != nil {
		log.Println("[UserService][Create][Create] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	return &pb.GetUserResponse{
		User: user.ToProtoBuff(),
	}, nil
}

func (service *UserService) Update(ctx context.Context, req *pb.UpdateUserRequest) (*pb.GetUserResponse, error) {
	var tx *sql.Tx
	if service.DB != nil {
		transaction, err := service.DB.Begin()
		if err != nil {
			log.Println("[UserService][Update] problem in db transaction, err: ", err.Error())
			return nil, err
		}
		tx = transaction
	}
	defer util.CommitOrRollback(tx)

	user, err := service.UserRepository.FindById(ctx, tx, req.IdUser)
	if err != nil {
		log.Println("[UserService][Update][FindById] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	newPassword := user.Password
	if req.Password != "" {
		password, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Println("[UserService][Update] problem in generating password hashed, err: ", err.Error())
			return nil, err
		}

		newPassword = string(password)
	}

	timeNow, err := time.Parse(util.TimeFormat, time.Now().Format(util.TimeFormat))
	if err != nil {
		log.Println("[UserService][Update] problem in parsing to time, err: ", err.Error())
		return nil, err
	}

	user.Name = req.Name
	user.Email = req.Email
	user.Address = req.Address
	user.Phone = req.Phone
	user.Password = newPassword
	user.UpdatedAt = timeNow

	_, err = service.UserRepository.Update(ctx, tx, user)
	if err != nil {
		log.Println("[UserService][Update][Update] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	return &pb.GetUserResponse{
		User: user.ToProtoBuff(),
	}, nil
}

func (service *UserService) UpdatePassword(ctx context.Context, req *pb.UpdateUserPasswordRequest) (*emptypb.Empty, error) {
	var tx *sql.Tx
	if service.DB != nil {
		transaction, err := service.DB.Begin()
		if err != nil {
			log.Println("[UserService][UpdatePassword] problem in db transaction, err: ", err.Error())
			return new(emptypb.Empty), err
		}
		tx = transaction
	}
	defer util.CommitOrRollback(tx)

	timeNow, err := time.Parse(util.TimeFormat, time.Now().Format(util.TimeFormat))
	if err != nil {
		log.Println("[UserService][UpdatePassword] problem in parsing to time, err: ", err.Error())
		return new(emptypb.Empty), err
	}

	err = service.UserRepository.UpdatePassword(ctx, tx, &model.User{
		Password:  req.Password,
		Email:     req.Email,
		UpdatedAt: timeNow,
	})
	if err != nil {
		log.Println("[UserService][UpdatePassword] problem in getting from repository, err: ", err.Error())
		return new(emptypb.Empty), err
	}

	return new(emptypb.Empty), nil
}

func (service *UserService) Delete(ctx context.Context, req *pb.GetUserByIdRequest) (*emptypb.Empty, error) {
	var tx *sql.Tx
	if service.DB != nil {
		transaction, err := service.DB.Begin()
		if err != nil {
			log.Println("[UserService][Delete] problem in db transaction, err: ", err.Error())
			return new(emptypb.Empty), err
		}
		tx = transaction
	}
	defer util.CommitOrRollback(tx)

	user, err := service.UserRepository.FindById(ctx, tx, req.Id)
	if err != nil {
		log.Println("[UserService][Delete][FindById] problem in getting from repository, err: ", err.Error())
		return new(emptypb.Empty), err
	}

	err = service.UserRepository.Delete(ctx, tx, user.IdUser)
	if err != nil {
		log.Println("[UserService][Delete][Delete] problem in getting from repository, err: ", err.Error())
		return new(emptypb.Empty), err
	}

	return new(emptypb.Empty), nil
}
