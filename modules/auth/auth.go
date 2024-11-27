// Package auth contains the authentication module which is responsible for authenticating and authorizing users on the endpoints.
package auth

import (
	"net/http"

	"github.com/elmawardy/nutrix/common/config"
	"github.com/elmawardy/nutrix/common/logger"
	"github.com/elmawardy/nutrix/common/userio"
)

// IHttpAuth is an interface for the HTTP authentication
type IHttpAuth interface {
	// AllowRoles middleware checks if the given request has a valid JWT token
	// and if the user has any of the given roles.
	AllowRoles(next http.Handler, roles ...string) http.Handler
	// AllowAuthenticated middleware checks if the given request has a valid JWT token.
	AllowAuthenticated(next http.Handler) http.Handler
}

// NewBuilder creates a new AuthModuleBuilder
func NewBuilder(config config.Config, settings config.Settings) *AuthModuleBuilder {
	mb := new(AuthModuleBuilder)
	mb.Config = config
	mb.Settings = settings

	return mb
}

// Auth is the main struct for the auth module
type Auth struct {
	Logger   logger.ILogger
	Config   config.Config
	Settings config.Settings
	Prompter userio.Prompter
}

// AuthModuleBuilder is the builder for the auth module
type AuthModuleBuilder struct {
	Logger   logger.ILogger
	Config   config.Config
	Settings config.Settings
	Prompter userio.Prompter
}