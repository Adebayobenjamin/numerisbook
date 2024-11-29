package helper

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func GetCustomerIDFromContext(ctx *gin.Context) (uint, error) {
	customerID := ctx.GetUint("customer_id")
	return customerID, nil
}

func GetOffset(page int, limit int) int {
	return (page - 1) * limit
}

func JSONUnmarshalToType[T any](data any, valueType *T) error {
	var (
		err error
		dd  []byte
	)

	switch data := data.(type) {
	case string:
		dd = []byte(data)
	case []byte:
		dd = data
	default:
		dd, err = json.Marshal(data)
		if err != nil {
			return err
		}
	}
	err = json.Unmarshal(dd, valueType)
	if err != nil {
		return err
	}
	return nil
}

func GenerateInvoiceNumber() string {
	// TODO: Generate invoice number format
	return fmt.Sprintf("INV-%d", time.Now().Unix())
}

func ReturnPointer[T any](value T) *T {
	return &value
}
