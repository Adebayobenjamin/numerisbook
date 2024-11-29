package helper

import (
	"fmt"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetCustomerIDFromContext(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(nil)

	t.Run("successful retrieval", func(t *testing.T) {
		expectedID := uint(123)
		c.Set("customer_id", expectedID)

		customerID, err := GetCustomerIDFromContext(c)

		assert.NoError(t, err)
		assert.Equal(t, expectedID, customerID)
	})

	t.Run("zero value when not set", func(t *testing.T) {
		c.Set("customer_id", uint(0))

		customerID, err := GetCustomerIDFromContext(c)

		assert.NoError(t, err)
		assert.Equal(t, uint(0), customerID)
	})
}

func TestGetOffset(t *testing.T) {
	tests := []struct {
		name     string
		page     int
		limit    int
		expected int
	}{
		{"first page", 1, 10, 0},
		{"second page", 2, 10, 10},
		{"third page with different limit", 3, 5, 10},
		{"zero page", 0, 10, -10},
		{"zero limit", 1, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetOffset(tt.page, tt.limit)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestJSONUnmarshalToType(t *testing.T) {
	type TestStruct struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}

	t.Run("string input", func(t *testing.T) {
		input := `{"name":"test","value":123}`
		var result TestStruct
		err := JSONUnmarshalToType(input, &result)

		assert.NoError(t, err)
		assert.Equal(t, "test", result.Name)
		assert.Equal(t, 123, result.Value)
	})

	t.Run("byte array input", func(t *testing.T) {
		input := []byte(`{"name":"test","value":123}`)
		var result TestStruct
		err := JSONUnmarshalToType(input, &result)

		assert.NoError(t, err)
		assert.Equal(t, "test", result.Name)
		assert.Equal(t, 123, result.Value)
	})

	t.Run("struct input", func(t *testing.T) {
		input := TestStruct{Name: "test", Value: 123}
		var result TestStruct
		err := JSONUnmarshalToType(input, &result)

		assert.NoError(t, err)
		assert.Equal(t, "test", result.Name)
		assert.Equal(t, 123, result.Value)
	})

	t.Run("invalid json", func(t *testing.T) {
		input := `{"name":"test",invalid}`
		var result TestStruct
		err := JSONUnmarshalToType(input, &result)

		assert.Error(t, err)
	})
}

func TestGenerateInvoiceNumber(t *testing.T) {
	t.Run("format check", func(t *testing.T) {
		now := time.Now()
		invoiceNumber := GenerateInvoiceNumber()

		assert.Contains(t, invoiceNumber, "INV-")
		expectedTimestamp := fmt.Sprintf("%d", now.Unix())
		assert.Contains(t, invoiceNumber, expectedTimestamp)
	})

	t.Run("uniqueness check", func(t *testing.T) {
		// Generate two invoice numbers with a slight delay
		first := GenerateInvoiceNumber()
		time.Sleep(time.Second)
		second := GenerateInvoiceNumber()

		assert.NotEqual(t, first, second)
	})
}

func TestReturnPointer(t *testing.T) {
	t.Run("string pointer", func(t *testing.T) {
		value := "test"
		ptr := ReturnPointer(value)

		assert.NotNil(t, ptr)
		assert.Equal(t, value, *ptr)
	})

	t.Run("int pointer", func(t *testing.T) {
		value := 123
		ptr := ReturnPointer(value)

		assert.NotNil(t, ptr)
		assert.Equal(t, value, *ptr)
	})
}
