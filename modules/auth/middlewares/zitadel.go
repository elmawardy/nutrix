// Package middlewares provides a set of middleware functions used to check
// Zitadel access token for auth and roles.
package middlewares

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/elmawardy/nutrix/common/config"

	"github.com/zitadel/zitadel-go/v3/pkg/authorization"
	"github.com/zitadel/zitadel-go/v3/pkg/authorization/oauth"
	"github.com/zitadel/zitadel-go/v3/pkg/http/middleware"
	"github.com/zitadel/zitadel-go/v3/pkg/zitadel"
)

// NewZitadelAuth creates a new ZitadelAuth object with the given configuration.
// It sets up the Zitadel SDK with the given domain and key path.
func NewZitadelAuth(conf config.Config) ZitadelAuth {

	ctx := context.Background()

	za := ZitadelAuth{
		Domain: conf.Zitadel.Domain,
		Key:    conf.Zitadel.KeyPath,
	}

	authZ, err := authorization.New(ctx, zitadel.New(za.Domain, zitadel.WithInsecure("2020")), oauth.DefaultAuthorization(za.Key))
	/******  6999d1a5-4501-4af6-9052-14fbce64d1ab  *******/
	if err != nil {
		slog.Error("zitadel sdk could not initialize", "error", err)
		os.Exit(1)
	}

	za.AuthZ = authZ

	return za
}

// ZitadelAuth holds the configuration for Zitadel and the Authorizer
type ZitadelAuth struct {
	Domain string // Zitadel instance domain
	Key    string // path to key.json
	AuthZ  *authorization.Authorizer[*oauth.IntrospectionContext]
}

// AllowAuthenticated middleware checks if the given request has a valid acess token.
func (za *ZitadelAuth) AllowAuthenticated(next http.Handler) http.Handler {

	mw := middleware.New(za.AuthZ)

	handler := mw.RequireAuthorization()

	return handler(next)
}

// AllowAnyOfRoles middleware checks if the given request has a valid access token
// and if the user has any of the given roles.
func (za *ZitadelAuth) AllowAnyOfRoles(next http.Handler, roles ...string) http.Handler {

	// Initialize the HTTP middleware by providing the authorization
	// mw := middleware.New(za.AuthZ)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authorized := false

		for _, role := range roles {

			reqToken := r.Header.Get("Authorization")
			if reqToken == "" {
				http.Error(w, "no authorization header found", http.StatusForbidden)
				return
			}

			_, err := za.AuthZ.CheckAuthorization(r.Context(), reqToken, authorization.WithRole(role))

			if err == nil {
				authorized = true
				next.ServeHTTP(w, r)
			}
		}

		if !authorized {
			w.WriteHeader(http.StatusUnauthorized)
		}

	})

	// checkpoints := []authorization.CheckOption{}

	// for _, role := range roles {
	// 	checkpoints = append(checkpoints, authorization.WithRole(role))
	// }

	// handler := mw.RequireAuthorization(checkpoints...)

	// return handler(next)
}