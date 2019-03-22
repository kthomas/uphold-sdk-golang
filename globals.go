package uphold

import (
	"os"
	"sync"

	"github.com/kthomas/go-logger"
)

const upholdSandboxBaseURL = "https://sandbox.uphold.com"
const upholdSandboxAPIBaseURL = "https://api-sandbox.uphold.com"
const upholdSupportedScopes = "accounts:read cards:read cards:write transactions:deposit transactions:transfer:application transactions:transfer:others transactions:transfer:self transactions:withdraw transactions:read user:read contacts:read contacts:write phones:read phones:write"

var (
	log           *logger.Logger
	bootstrapOnce sync.Once

	upholdBaseURL      string
	upholdAPIBaseURL   string
	upholdClientID     string
	upholdClientSecret string
)

func init() {
	bootstrapOnce.Do(func() {
		lvl := os.Getenv("LOG_LEVEL")
		if lvl == "" {
			lvl = "INFO"
		}
		log = logger.NewLogger("uphold", lvl, true)

		if os.Getenv("UPHOLD_BASE_URL") != "" {
			upholdBaseURL = os.Getenv("UPHOLD_BASE_URL")
		} else {
			upholdBaseURL = upholdSandboxBaseURL
		}

		if os.Getenv("UPHOLD_API_BASE_URL") != "" {
			upholdAPIBaseURL = os.Getenv("UPHOLD_API_BASE_URL")
		} else {
			upholdAPIBaseURL = upholdSandboxAPIBaseURL
		}

		if os.Getenv("UPHOLD_CLIENT_ID") != "" {
			upholdClientID = os.Getenv("UPHOLD_CLIENT_ID")
		}

		if os.Getenv("UPHOLD_CLIENT_SECRET") != "" {
			upholdClientSecret = os.Getenv("UPHOLD_CLIENT_SECRET")
		}
	})
}
