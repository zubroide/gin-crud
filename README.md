# CRUD for GORM

## Features

- [x] CRUD: grom-crud
- [x] Base CRUD controller

## How to use

Repository:

```go
package repository

import (
	"github.com/jinzhu/gorm"
	"github.com/zubroide/gorm-crud"
	"go-api-boilerplate/logger"
	"go-api-boilerplate/model/entity"
)


type UserListQueryBuilder struct {
	gorm_crud.ListQueryBuilderInterface
	*gorm_crud.BaseListQueryBuilder
}

func NewUserListQueryBuilder(db *gorm.DB, logger logger.LoggerInterface) gorm_crud.ListQueryBuilderInterface {
	base := gorm_crud.NewBaseListQueryBuilder(db, logger).(*gorm_crud.BaseListQueryBuilder)
	return &UserListQueryBuilder{BaseListQueryBuilder: base}
}

func (c UserListQueryBuilder) ListQuery(parameters gorm_crud.ListParametersInterface) (*gorm.DB, error) {
	query, err := c.BaseListQueryBuilder.ListQuery(parameters)
	params := parameters.(*UserListParameters)
	if err == nil && params.Name != "" {
		query = query.Where("name LIKE ?", params.Name + "%")
	}
	return query, err
}


type UserRepositoryInterface interface {
	gorm_crud.CrudRepositoryInterface
}

type UserListParameters struct {
	*gorm_crud.CrudListParameters
	Name string
}

type UserRepository struct {
	*gorm_crud.CrudRepository
	model entity.User
}

func NewUserRepository(db *gorm.DB, logger logger.LoggerInterface) UserRepositoryInterface {
	var model entity.User
	queryBuilder := NewUserListQueryBuilder(db, logger)
	repo := gorm_crud.NewCrudRepository(db, &model, queryBuilder, logger).(*gorm_crud.CrudRepository)
	return &UserRepository{repo, model}
}
```

Service:

```go
package service

import (
	"github.com/zubroide/gorm-crud"
	"go-api-boilerplate/logger"
	"go-api-boilerplate/model/repository"
)

type UserServiceInterface interface {
	gorm_crud.CrudServiceInterface
}

type UserService struct {
	*gorm_crud.CrudService
	repository repository.UserRepositoryInterface
}

func NewUserService(repository repository.UserRepositoryInterface, logger logger.LoggerInterface) UserServiceInterface {
	crudService := gorm_crud.NewCrudService(repository, logger).(*gorm_crud.CrudService)
	service := &UserService{crudService, repository}
	return service
}
```

Controller:
```go
package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/zubroide/gin-crud"
	"github.com/zubroide/gorm-crud"
	"go-api-boilerplate/logger"
	"go-api-boilerplate/model/repository"
	"go-api-boilerplate/model/service"
)


type UserListParametersHydrator struct {
	*gin_crud.BaseParametersHydrator
}

func NewUserListParametersHydrator(logger logger.LoggerInterface) gin_crud.ParametersHydratorInterface {
	base := gin_crud.NewBaseParametersHydrator(logger).(*gin_crud.BaseParametersHydrator)
	return &UserListParametersHydrator{BaseParametersHydrator: base}
}

func (c UserListParametersHydrator) Hydrate(context *gin.Context) (gorm_crud.ListParametersInterface, error) {
	crudParams, _ := c.BaseParametersHydrator.Hydrate(context)
	parameters := &repository.UserListParameters{
		CrudListParameters: crudParams.(*gorm_crud.CrudListParameters),
	}
	if err := context.ShouldBindQuery(parameters); err != nil {
		return crudParams, err
	}

	return parameters, nil
}


type UserController struct {
	*gin_crud.CrudController
	service service.UserServiceInterface
}

func NewUserController(service service.UserServiceInterface, logger logger.LoggerInterface) *UserController {
	parametersHydrator := NewUserListParametersHydrator(logger)
	controller := gin_crud.NewCrudController(service, parametersHydrator, logger)
	return &UserController{CrudController: controller, service: service}
}
```

## Requirements
  - Go 1.12+

# License

MIT
