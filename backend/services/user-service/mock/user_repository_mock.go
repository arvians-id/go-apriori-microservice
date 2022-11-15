package mock

import (
	"context"
	"database/sql"
	"errors"
	"github.com/arvians-id/go-apriori-microservice/services/user-service/model"
	"github.com/arvians-id/go-apriori-microservice/services/user-service/util"
	"github.com/stretchr/testify/mock"
)

type UserRepositoryMock struct {
	Mock mock.Mock
}

func (repository *UserRepositoryMock) FindAll(ctx context.Context, tx *sql.Tx) ([]*model.User, error) {
	arguments := repository.Mock.Called(ctx)
	if arguments.Get(0) == nil {
		return nil, errors.New(util.ErrorNotFound)
	}

	return arguments.Get(0).([]*model.User), nil
}
func (repository *UserRepositoryMock) FindById(ctx context.Context, tx *sql.Tx, id int64) (*model.User, error) {
	arguments := repository.Mock.Called(ctx, id)
	if arguments.Get(0) == nil {
		return nil, errors.New(util.ErrorNotFound)
	}

	return arguments.Get(0).(*model.User), nil

}
func (repository *UserRepositoryMock) FindByEmail(ctx context.Context, tx *sql.Tx, email string) (*model.User, error) {
	arguments := repository.Mock.Called(ctx, email)
	if arguments.Get(0) == nil {
		return nil, errors.New(util.ErrorNotFound)
	}

	return arguments.Get(0).(*model.User), nil

}
func (repository *UserRepositoryMock) Create(ctx context.Context, tx *sql.Tx, user *model.User) (*model.User, error) {
	arguments := repository.Mock.Called(ctx)
	if arguments.Get(0) == nil {
		return nil, errors.New(util.ErrorNotFound)
	}

	return arguments.Get(0).(*model.User), nil
}
func (repository *UserRepositoryMock) Update(ctx context.Context, tx *sql.Tx, user *model.User) (*model.User, error) {
	arguments := repository.Mock.Called(ctx)
	if arguments.Get(0) == nil {
		return nil, errors.New(util.ErrorNotFound)
	}

	return arguments.Get(0).(*model.User), nil

}
func (repository *UserRepositoryMock) UpdatePassword(ctx context.Context, tx *sql.Tx, user *model.User) error {
	arguments := repository.Mock.Called(ctx)
	if arguments.Get(0) == nil {
		return errors.New(util.ErrorNotFound)
	}

	return nil

}
func (repository *UserRepositoryMock) Delete(ctx context.Context, tx *sql.Tx, id int64) error {
	arguments := repository.Mock.Called(ctx, id)
	if arguments.Get(0) == nil {
		return nil
	}

	return errors.New("something went wrong")
}
