package gin_crud

import (
	"github.com/gin-gonic/gin"
	"github.com/zubroide/gorm-crud"
	"net/http"
	"reflect"
	"strconv"
)


type ParametersHydratorInterface interface {
	Hydrate(context *gin.Context) (gorm_crud.ListParametersInterface, error)
}

type BaseParametersHydrator struct {
	Logger gorm_crud.LoggerInterface
	ParametersHydratorInterface
}

func NewBaseParametersHydrator(logger gorm_crud.LoggerInterface) *BaseParametersHydrator {
	return &BaseParametersHydrator{Logger: logger}
}

func (c BaseParametersHydrator) Hydrate(context *gin.Context) (gorm_crud.ListParametersInterface, error) {
	var parameters gorm_crud.CrudListParameters
	err := context.ShouldBindQuery(&parameters)
	return &parameters, err
}


type CrudControllerInterface interface {
	BaseControllerInterface
	Create(context *gin.Context)
	Get(context *gin.Context)
	List(context *gin.Context)
	Update(context *gin.Context)
	Delete(context *gin.Context)
}

type CrudController struct {
	CrudControllerInterface
	*BaseController
	Service            gorm_crud.CrudServiceInterface
	ParametersHydrator ParametersHydratorInterface
}

func NewCrudController(service gorm_crud.CrudServiceInterface, parametersHydrator ParametersHydratorInterface, logger gorm_crud.LoggerInterface) *CrudController {
	controller := NewBaseController(logger)
	return &CrudController{BaseController: controller, Service: service, ParametersHydrator: parametersHydrator}
}

func (c CrudController) Get(context *gin.Context) {
	recordId, err := strconv.Atoi(context.Params.ByName("id"))
	if err != nil {
		c.ReplyError(context, "Please specify record id", http.StatusBadRequest)
		return
	}

	data, err := c.Service.GetItem(uint(recordId))

	if err != nil {
		c.ReplyError(context, "Record not found", http.StatusNotFound)
		return
	}

	c.ReplySuccess(context, data)
}

func (c CrudController) List(context *gin.Context) {
	parameters, err := c.ParametersHydrator.Hydrate(context)

	if err != nil {
		c.ReplyError(context, "Cant parse request parameters", http.StatusBadRequest)
		return
	}

	data, err := c.Service.GetList(parameters)

	if err != nil {
		c.ReplyError(context, "Data not found", http.StatusBadRequest)
		return
	}

	c.ReplySuccess(context, data)
}

func (c CrudController) Create(context *gin.Context) {
	model := c.Service.GetModel()
	data := reflect.New(reflect.TypeOf(model).Elem()).Interface()
	if err := context.ShouldBindJSON(data); err != nil {
		c.ReplyError(context, "Cant parse request", http.StatusBadRequest)
		return
	}
	data = c.Service.Create(data)
	c.ReplySuccess(context, data)
}

func (c CrudController) Update(context *gin.Context) {
	recordId, err := strconv.Atoi(context.Params.ByName("id"))
	if err != nil {
		c.ReplyError(context, "Cant parse request", http.StatusBadRequest)
		return
	}

	data, err := c.Service.GetItem(uint(recordId))
	if err != nil {
		c.ReplyError(context, "Data not found", http.StatusBadRequest)
		return
	}

	if err := context.ShouldBindJSON(data); err != nil {
		c.ReplyError(context, "Cant parse request", http.StatusBadRequest)
		return
	}
	data = c.Service.Update(data)

	c.ReplySuccess(context, data)
}

func (c CrudController) Delete(context *gin.Context) {
	recordId, err := strconv.Atoi(context.Params.ByName("id"))
	if err != nil {
		c.ReplyError(context, "Please specify record id", http.StatusBadRequest)
		return
	}

	err = c.Service.Delete(uint(recordId))
	if err != nil {
		c.ReplyError(context, "Data not found", http.StatusBadRequest)
		return
	}

	c.ReplySuccess(context, nil)
}
