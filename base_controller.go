package gin_crud

import (
	"github.com/gin-gonic/gin"
	"github.com/zubroide/gorm-crud"
	"net/http"
)

const StatusOk = "ok"
const StatusError = "error"

type BaseControllerInterface interface {
}

type BaseController struct {
	Logger gorm_crud.LoggerInterface
}

func NewBaseController(logger gorm_crud.LoggerInterface) *BaseController {
	return &BaseController{Logger: logger}
}

func (c BaseController) ReplySuccess(context *gin.Context, data interface{}) {
	c.Response(context, gin.H{"data": data, "status": StatusOk}, http.StatusOK)
}

func (c BaseController) ReplyError(context *gin.Context, message string, code int) {
	c.Response(context, gin.H{"message": message, "status": StatusError}, code)
}

func (c BaseController) Response(context *gin.Context, obj interface{}, code int) {
	switch context.GetHeader("Accept") {
		case "application/xml":
			context.XML(code, obj)
		default:
			context.JSON(code, obj)
	}
}
