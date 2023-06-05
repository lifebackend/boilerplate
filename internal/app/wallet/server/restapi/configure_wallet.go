// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"io"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/nvsco/wallet/internal/app/wallet/server/restapi/operations"
	"github.com/nvsco/wallet/internal/app/wallet/server/restapi/operations/general"
	"github.com/nvsco/wallet/internal/app/wallet/server/restapi/operations/health"
)

func configureFlags(api *operations.WalletAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.WalletAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.UseSwaggerUI()
	// To continue using redoc as your UI, uncomment the following line
	// api.UseRedoc()

	api.JSONConsumer = runtime.JSONConsumer()

	api.HTMLProducer = runtime.ProducerFunc(func(w io.Writer, data interface{}) error {
		return errors.NotImplemented("html producer has not yet been implemented")
	})
	api.JSONProducer = runtime.JSONProducer()

	/* default handlers import
		"github.com/nvsco/wallet/internal/app/wallet/server/client"
		"github.com/nvsco/wallet/internal/app/wallet/server/models"
		"github.com/nvsco/wallet/internal/app/wallet/server/restapi"
		"github.com/nvsco/wallet/internal/app/wallet/server/restapi/operations/general"
		"github.com/nvsco/wallet/internal/app/wallet/server/restapi/operations/health"
		"github.com/nvsco/wallet/internal/app/wallet/server/restapi/operations"
	  default handlers import */

	/* default handlers declaration
	   // default handle functions
	   api.GeneralGetAppCodesHandler = general.GetAppCodesHandlerFunc(h.GeneralGetAppCodesHandler)
	   api.HealthGetLivenessProbeHandler = health.GetLivenessProbeHandlerFunc(h.HealthGetLivenessProbeHandler)
	   api.HealthGetReadinessProbeHandler = health.GetReadinessProbeHandlerFunc(h.HealthGetReadinessProbeHandler)
	 default handlers declaration */

	// Default handlers

	/* default handler for /app-codes-GET
	   // GeneralGetAppCodesHandler Handler for GET /app-codes
	   func (h *Handlers) GeneralGetAppCodesHandler (
	       params *general.GetAppCodesParams,
	       respond *general.GetAppCodesResponses,
	   ) middleware.Responder {
	       return middleware.NotImplemented("operation general.GetAppCodes has not yet been implemented")
	   }
	   default handler */

	/* default handler for /_livenessProbe-GET
	   // HealthGetLivenessProbeHandler Handler for GET /_livenessProbe
	   func (h *Handlers) HealthGetLivenessProbeHandler (
	       params *health.GetLivenessProbeParams,
	       respond *health.GetLivenessProbeResponses,
	   ) middleware.Responder {
	       return middleware.NotImplemented("operation health.GetLivenessProbe has not yet been implemented")
	   }
	   default handler */

	/* default handler for /_readinessProbe-GET
	   // HealthGetReadinessProbeHandler Handler for GET /_readinessProbe
	   func (h *Handlers) HealthGetReadinessProbeHandler (
	       params *health.GetReadinessProbeParams,
	       respond *health.GetReadinessProbeResponses,
	   ) middleware.Responder {
	       return middleware.NotImplemented("operation health.GetReadinessProbe has not yet been implemented")
	   }
	   default handler */
	// Handler for GET /app-codes
	if api.GeneralGetAppCodesHandler == nil {
		api.GeneralGetAppCodesHandler = general.GetAppCodesHandlerFunc(func(
			params *general.GetAppCodesParams,
			respond *general.GetAppCodesResponses,
		) middleware.Responder {
			return middleware.NotImplemented("operation general.GetAppCodes has not yet been implemented")
		})
	}
	// Handler for GET /_livenessProbe
	if api.HealthGetLivenessProbeHandler == nil {
		api.HealthGetLivenessProbeHandler = health.GetLivenessProbeHandlerFunc(func(
			params *health.GetLivenessProbeParams,
			respond *health.GetLivenessProbeResponses,
		) middleware.Responder {
			return middleware.NotImplemented("operation health.GetLivenessProbe has not yet been implemented")
		})
	}
	// Handler for GET /_readinessProbe
	if api.HealthGetReadinessProbeHandler == nil {
		api.HealthGetReadinessProbeHandler = health.GetReadinessProbeHandlerFunc(func(
			params *health.GetReadinessProbeParams,
			respond *health.GetReadinessProbeResponses,
		) middleware.Responder {
			return middleware.NotImplemented("operation health.GetReadinessProbe has not yet been implemented")
		})
	}

	api.PreServerShutdown = func() {}

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix"
func configureServer(s *http.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/_metrics" {
			promhttp.Handler().ServeHTTP(w, r)
			return
		}

		if r.URL.Path == "/favicon.ico" {
			w.WriteHeader(http.StatusOK)
			return
		}

		handler.ServeHTTP(w, r)
		// cors.AllowAll().Handler(handler).ServeHTTP(w, r)
	})
}

/*
func serveVersion(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
}
*/
