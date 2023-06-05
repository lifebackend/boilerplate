package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/nvsco/wallet/internal/app/wallet/server/restapi/operations/health"
)

// HealthGetLivenessProbeHandler Handler for GET /_livenessProbe
func (h *Handlers) HealthGetLivenessProbeHandler(
	params *health.GetLivenessProbeParams,
	respond *health.GetLivenessProbeResponses,
) middleware.Responder {
	return middleware.NotImplemented("operation health.GetLivenessProbe has not yet been implemented")
}
