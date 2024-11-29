package exceptions

import (
	"net/http"

	"github.com/Adebayobenjamin/numerisbook/pkg/common"
	"github.com/gin-gonic/gin"
)

func ThrowUnProcessableEntityException(ctx *gin.Context, err string) {
	res := common.BuildErrorResponse("UnProcessable Entity", err)
	ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, res)
}

func ThrowBadRequestException(ctx *gin.Context, err string) {
	res := common.BuildErrorResponse("Bad Request", err)
	ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
}

func ThrowUnAuthorizedException(ctx *gin.Context, err string) {
	res := common.BuildErrorResponse("UnAuthorized", err)
	ctx.AbortWithStatusJSON(http.StatusUnauthorized, res)
}

func ThrowForbiddenException(ctx *gin.Context, err string) {
	res := common.BuildErrorResponse("Forbidden", err)
	ctx.AbortWithStatusJSON(http.StatusForbidden, res)
}

func ThrowInternalServerError(ctx *gin.Context, err string) {
	res := common.BuildErrorResponse("Internal Server Error", err)
	ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
}
