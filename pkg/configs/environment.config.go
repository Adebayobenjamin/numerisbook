package configs

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
)

// Env represents environmental variable instance
type Env struct{}

// New creates a new instance of Env and returns an error if any occurs
func NewEnvironment() *Env {
	if os.Getenv("GO_ENV") == "development" {
		err := godotenv.Load()
		if err != nil {
			panic("Error loading environment variables: " + err.Error())
		}
	}
	return &Env{}
}

// NewLoadFromFile lets you load Env object from a file
func NewLoadFromFile(fileName string) (*Env, error) {
	err := godotenv.Load(fileName)
	if err != nil {
		return nil, err
	}
	return &Env{}, nil
}

// Get retrieves the string value of an environmental variable
func (e *Env) Get(key string) string {
	return os.Getenv(key)
}

// UseMock is pkg that returns true or false if the environment should use mocks when hitting 3rd party partners
func (e *Env) UseMock() bool {
	v := e.Get("APP_MOCK")
	if len(v) == 0 {
		return true
	}

	if strings.EqualFold(v, "true") {
		return true
	}

	return false
}
