package session

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
	"github.com/gofiber/utils/v2"
)

// Config defines the config for middleware.
type Config struct {
	// Next defines a function to skip this middleware when returned true.
	//
	// Optional. Default: nil
	Next func(c fiber.Ctx) bool

	// Store defines the session store
	//
	// Required.
	Store *Store

	// ErrorHandler defines a function which is executed for errors
	//
	// Optional. Default: nil
	ErrorHandler func(*fiber.Ctx, error)

	// Allowed session idle duration
	// Optional. Default value 24 * time.Hour
	IdleTimeout time.Duration

	// Storage interface to store the session data
	// Optional. Default value memory.New()
	Storage fiber.Storage

	// KeyLookup is a string in the form of "<source>:<name>" that is used
	// to extract session id from the request.
	// Possible values: "header:<name>", "query:<name>" or "cookie:<name>"
	// Optional. Default value "cookie:session_id".
	KeyLookup string

	// Domain of the cookie.
	// Optional. Default value "".
	CookieDomain string

	// Path of the cookie.
	// Optional. Default value "".
	CookiePath string

	// Indicates if cookie is secure.
	// Optional. Default value false.
	CookieSecure bool

	// Indicates if cookie is HTTP only.
	// Optional. Default value false.
	CookieHTTPOnly bool

	// Value of SameSite cookie.
	// Optional. Default value "Lax".
	CookieSameSite string

	// Decides whether cookie should last for only the browser sesison.
	// Ignores Expiration if set to true
	// Optional. Default value false.
	CookieSessionOnly bool

	// KeyGenerator generates the session key.
	// Optional. Default value utils.UUIDv4
	KeyGenerator func() string

	// Source defines where to obtain the session id
	source Source

	// The session name
	sessionName string
}

type Source string

const (
	SourceCookie   Source = "cookie"
	SourceHeader   Source = "header"
	SourceURLQuery Source = "query"
)

// ConfigDefault is the default config
var ConfigDefault = Config{
	IdleTimeout:  24 * time.Hour,
	KeyLookup:    "cookie:session_id",
	KeyGenerator: utils.UUIDv4,
	source:       "cookie",
	sessionName:  "session_id",
}

func DefaultErrorHandler(c *fiber.Ctx, err error) {
	log.Errorf("session: %v", err)
	if c != nil {
		if err := (*c).SendStatus(fiber.StatusInternalServerError); err != nil {
			log.Errorf("session: %v", err)
		}
	}
}

// Helper function to set default values
func configDefault(config ...Config) Config {
	// Return default config if nothing provided
	if len(config) < 1 {
		return ConfigDefault
	}

	// Override default config
	cfg := config[0]

	// Set default values
	if int(cfg.IdleTimeout.Seconds()) <= 0 {
		cfg.IdleTimeout = ConfigDefault.IdleTimeout
	}
	if cfg.KeyLookup == "" {
		cfg.KeyLookup = ConfigDefault.KeyLookup
	}
	if cfg.KeyGenerator == nil {
		cfg.KeyGenerator = ConfigDefault.KeyGenerator
	}

	selectors := strings.Split(cfg.KeyLookup, ":")
	const numSelectors = 2
	if len(selectors) != numSelectors {
		panic("[session] KeyLookup must in the form of <source>:<name>")
	}
	switch Source(selectors[0]) {
	case SourceCookie:
		cfg.source = SourceCookie
	case SourceHeader:
		cfg.source = SourceHeader
	case SourceURLQuery:
		cfg.source = SourceURLQuery
	default:
		panic("[session] source is not supported")
	}
	cfg.sessionName = selectors[1]

	return cfg
}
