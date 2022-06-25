package token

import (
	"github.com/gin-gonic/gin"
	"github.com/jaswdr/faker"
	"os"
	"testing"
)

var f = faker.New()

func TestMain(m *testing.M) {

	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
