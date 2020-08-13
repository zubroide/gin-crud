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
	logger gorm_crud.LoggerInterface
}

func NewBaseController(logger gorm_crud.LoggerInterface) *BaseController {
	return &BaseController{logger: logger}
}

func (c BaseController) replySuccess(context *gin.Context, data interface{}) {
	c.response(context, gin.H{"data": data, "status": StatusOk}, http.StatusOK)
}

func (c BaseController) replyError(context *gin.Context, message string, code int) {
	c.response(context, gin.H{"message": message, "status": StatusError}, code)
}

func (c BaseController) response(context *gin.Context, obj interface{}, code int) {
	switch context.GetHeader("Accept") {
		case "application/xml":
			context.XML(code, obj)
		default:
			context.JSON(code, obj)
	}
}
