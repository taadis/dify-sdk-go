package completion

import (
	"os"
	"testing"

	"github.com/taadis/dify-sdk-go/env"
)

var (
	testBaseUrl = ""
	testApiKey  = ""
)

func TestMain(m *testing.M) {
	testBaseUrl = env.GetDifyBaseUrl()
	testApiKey = env.GetDifyApiKey()
	os.Exit(m.Run())
}
