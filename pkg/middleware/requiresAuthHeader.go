package middlewares

import (
	"strconv"

	"github.com/Adebayobenjamin/numerisbook/pkg/common/exceptions"
	"github.com/gin-gonic/gin"
)

func RequiresAuthHeader() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		customerID := ctx.GetHeader("x-customer-id")
		if customerID == "" {
			exceptions.ThrowForbiddenException(ctx, "customer id is required")
			return
		}

		customerIDUint, err := strconv.ParseUint(customerID, 10, 64)
		if err != nil {
			exceptions.ThrowForbiddenException(ctx, "customer id is required")
			return
		}

		ctx.Set("customer_id", uint(customerIDUint))

		ctx.Next()
	}
}
