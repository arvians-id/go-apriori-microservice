package usecase

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/arvians-id/go-apriori-microservice/services/user-service/model"
	"github.com/arvians-id/go-apriori-microservice/services/user-service/pb"
	"github.com/arvians-id/go-apriori-microservice/services/user-service/repository"
	redisLib "github.com/arvians-id/go-apriori-microservice/services/user-service/third-party/redis"
	"github.com/arvians-id/go-apriori-microservice/services/user-service/util"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/protobuf/types/known/emptypb"
)

type UserServiceCache struct {
	UserRepository repository.UserRepository
	Redis          redisLib.Redis
	DB             *sql.DB
}

func NewUserServiceCache(userRepository repository.UserRepository, redis *redisLib.Redis, db *sql.DB) pb.UserServiceServer {
	return &UserServiceCache{
		UserRepository: userRepository,
		Redis:          *redis,
		DB:             db,
	}
}

func (service *UserServiceCache) FindAll(ctx context.Context, empty *emptypb.Empty) (*pb.ListUserResponse, error) {
	var tx *sql.Tx
	if service.DB != nil {
		transaction, err := service.DB.Begin()
		if err != nil {
			log.Println("[UserServiceCache][FindAll] problem in db transaction, err: ", err.Error())
			return nil, err
		}
		tx = transaction
	}
	defer util.CommitOrRollback(tx)

	usersCache, err := service.Redis.Get(ctx, "users")
	if err != redis.Nil {
		var users []*pb.User
		err = json.Unmarshal(usersCache, &users)
		if err != nil {
			log.Println("[UserServiceCache][FindAll] unable to unmarshal json, err: ", err.Error())
			return nil, err
		}

		return &pb.ListUserResponse{
			User: users,
		}, nil
	}

	users, err := service.UserRepository.FindAll(ctx, tx)
	if err != nil {
		log.Println("[UserServiceCache][FindAll][FindAll] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	var userListResponse []*pb.User
	for _, user := range users {
		userListResponse = append(userListResponse, user.ToProtoBuff())
	}

	err = service.Redis.Set(ctx, "users", userListResponse)
	if err != nil {
		log.Println("[UserServiceCache][FindAll][Set] unable to set value to redis cache, err: ", err.Error())
		return nil, err
	}

	return &pb.ListUserResponse{
		User: userListResponse,
	}, nil
}

func (service *UserServiceCache) FindById(ctx context.Context, req *pb.GetUserByIdRequest) (*pb.GetUserResponse, error) {
	var tx *sql.Tx
	if service.DB != nil {
		transaction, err := service.DB.Begin()
		if err != nil {
			log.Println("[UserServiceCache][FindById] problem in db transaction, err: ", err.Error())
			return nil, err
		}
		tx = transaction
	}
	defer util.CommitOrRollback(tx)

	key := fmt.Sprintf("user:%d", req.Id)
	userCache, err := service.Redis.Get(ctx, key)
	if err != redis.Nil {
		var user model.User
		err = json.Unmarshal(userCache, &user)
		if err != nil {
			log.Println("[UserServiceCache][FindById] unable to unmarshal json, err: ", err.Error())
			return nil, err
		}

		return &pb.GetUserResponse{
			User: user.ToProtoBuff(),
		}, nil
	}

	user, err := service.UserRepository.FindById(ctx, tx, req.Id)
	if err != nil {
		log.Println("[UserServiceCache][FindById][FindById] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	err = service.Redis.Set(ctx, key, user)
	if err != nil {
		log.Println("[UserServiceCache][FindById][Set] unable to set value to redis cache, err: ", err.Error())
		return nil, err
	}

	return &pb.GetUserResponse{
		User: user.ToProtoBuff(),
	}, nil
}

func (service *UserServiceCache) FindByEmail(ctx context.Context, req *pb.GetUserByEmailRequest) (*pb.GetUserResponse, error) {
	var tx *sql.Tx
	if service.DB != nil {
		transaction, err := service.DB.Begin()
		if err != nil {
			log.Println("[UserServiceCache][FindByEmail] problem in db transaction, err: ", err.Error())
			return nil, err
		}
		tx = transaction
	}
	defer util.CommitOrRollback(tx)

	key := fmt.Sprintf("user:%s", req.Email)
	userCache, err := service.Redis.Get(ctx, key)
	if err != redis.Nil {
		var user model.User
		err = json.Unmarshal(userCache, &user)
		if err != nil {
			log.Println("[UserServiceCache][FindByEmail] unable to unmarshal json, err: ", err.Error())
			return nil, err
		}

		return &pb.GetUserResponse{
			User: user.ToProtoBuff(),
		}, nil
	}

	user, err := service.UserRepository.FindByEmail(ctx, tx, req.Email)
	if err != nil {
		log.Println("[UserServiceCache][FindByEmail][FindByEmail] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	err = service.Redis.Set(ctx, key, user)
	if err != nil {
		log.Println("[UserServiceCache][FindByEmail][Set] unable to set value to redis cache, err: ", err.Error())
		return nil, err
	}

	return &pb.GetUserResponse{
		User: user.ToProtoBuff(),
	}, nil
}

func (service *UserServiceCache) VerifyCredential(ctx context.Context, req *pb.GetVerifyCredentialRequest) (*pb.GetUserResponse, error) {
	var tx *sql.Tx
	if service.DB != nil {
		transaction, err := service.DB.Begin()
		if err != nil {
			log.Println("[UserServiceCache][FindByEmail] problem in db transaction, err: ", err.Error())
			return nil, err
		}
		tx = transaction
	}
	defer util.CommitOrRollback(tx)

	user, err := service.UserRepository.FindByEmail(ctx, tx, req.Email)
	if err != nil {
		log.Println("[UserServiceCache][FindByEmail][FindByEmail] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		log.Println("[UserServiceCache][FindByEmail] problem in comparing password, err: ", err.Error())
		return nil, errors.New("wrong password")
	}
	return &pb.GetUserResponse{
		User: user.ToProtoBuff(),
	}, nil
}

func (service *UserServiceCache) Create(ctx context.Context, req *pb.CreateUserRequest) (*pb.GetUserResponse, error) {
	var tx *sql.Tx
	if service.DB != nil {
		transaction, err := service.DB.Begin()
		if err != nil {
			log.Println("[UserServiceCache][Create] problem in db transaction, err: ", err.Error())
			return nil, err
		}
		tx = transaction
	}
	defer util.CommitOrRollback(tx)

	password, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("[UserServiceCache][Create] problem in generating password hashed, err: ", err.Error())
		return nil, err
	}

	timeNow, err := time.Parse(util.TimeFormat, time.Now().Format(util.TimeFormat))
	if err != nil {
		log.Println("[UserServiceCache][Create] problem in parsing to time, err: ", err.Error())
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
		log.Println("[UserServiceCache][Create][Create] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	err = service.Redis.Del(ctx, "users")
	if err != nil {
		log.Println("[UserServiceCache][Create][Del] unable to deleting specific key cache, err: ", err.Error())
	}

	return &pb.GetUserResponse{
		User: user.ToProtoBuff(),
	}, nil
}

func (service *UserServiceCache) Update(ctx context.Context, req *pb.UpdateUserRequest) (*pb.GetUserResponse, error) {
	var tx *sql.Tx
	if service.DB != nil {
		transaction, err := service.DB.Begin()
		if err != nil {
			log.Println("[UserServiceCache][Update] problem in db transaction, err: ", err.Error())
			return nil, err
		}
		tx = transaction
	}
	defer util.CommitOrRollback(tx)

	user, err := service.UserRepository.FindById(ctx, tx, req.IdUser)
	if err != nil {
		log.Println("[UserServiceCache][Update][FindById] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	newPassword := user.Password
	if req.Password != "" {
		password, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Println("[UserServiceCache][Update] problem in generating password hashed, err: ", err.Error())
			return nil, err
		}

		newPassword = string(password)
	}

	timeNow, err := time.Parse(util.TimeFormat, time.Now().Format(util.TimeFormat))
	if err != nil {
		log.Println("[UserServiceCache][Update] problem in parsing to time, err: ", err.Error())
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
		log.Println("[UserServiceCache][Update][Update] problem in getting from repository, err: ", err.Error())
		return nil, err
	}

	key1 := fmt.Sprintf("user:%s", req.Email)
	key2 := fmt.Sprintf("user:%d", req.IdUser)
	err = service.Redis.Del(ctx, "users", key1, key2)
	if err != nil {
		log.Println("[UserServiceCache][Update][Del] unable to deleting specific key cache, err: ", err.Error())
	}

	return &pb.GetUserResponse{
		User: user.ToProtoBuff(),
	}, nil
}

func (service *UserServiceCache) UpdatePassword(ctx context.Context, req *pb.UpdateUserPasswordRequest) (*emptypb.Empty, error) {
	var tx *sql.Tx
	if service.DB != nil {
		transaction, err := service.DB.Begin()
		if err != nil {
			log.Println("[UserServiceCache][UpdatePassword] problem in db transaction, err: ", err.Error())
			return new(emptypb.Empty), err
		}
		tx = transaction
	}
	defer util.CommitOrRollback(tx)

	timeNow, err := time.Parse(util.TimeFormat, time.Now().Format(util.TimeFormat))
	if err != nil {
		log.Println("[UserServiceCache][UpdatePassword] problem in parsing to time, err: ", err.Error())
		return new(emptypb.Empty), err
	}

	err = service.UserRepository.UpdatePassword(ctx, tx, &model.User{
		Password:  req.Password,
		Email:     req.Email,
		UpdatedAt: timeNow,
	})
	if err != nil {
		log.Println("[UserServiceCache][UpdatePassword] problem in getting from repository, err: ", err.Error())
		return new(emptypb.Empty), err
	}

	key1 := fmt.Sprintf("user:%s", req.Email)
	err = service.Redis.Del(ctx, "users", key1)
	if err != nil {
		log.Println("[UserServiceCache][Update][Del] unable to deleting specific key cache, err: ", err.Error())
	}

	return new(emptypb.Empty), nil
}

func (service *UserServiceCache) Delete(ctx context.Context, req *pb.GetUserByIdRequest) (*emptypb.Empty, error) {
	var tx *sql.Tx
	if service.DB != nil {
		transaction, err := service.DB.Begin()
		if err != nil {
			log.Println("[UserServiceCache][Delete] problem in db transaction, err: ", err.Error())
			return new(emptypb.Empty), err
		}
		tx = transaction
	}
	defer util.CommitOrRollback(tx)

	user, err := service.UserRepository.FindById(ctx, tx, req.Id)
	if err != nil {
		log.Println("[UserServiceCache][Delete][FindById] problem in getting from repository, err: ", err.Error())
		return new(emptypb.Empty), err
	}

	err = service.UserRepository.Delete(ctx, tx, user.IdUser)
	if err != nil {
		log.Println("[UserServiceCache][Delete][Delete] problem in getting from repository, err: ", err.Error())
		return new(emptypb.Empty), err
	}

	err = service.Redis.Del(ctx, "users")
	if err != nil {
		log.Println("[UserServiceCache][Delete][Del] unable to deleting specific key cache, err: ", err.Error())
	}

	return new(emptypb.Empty), nil
}
