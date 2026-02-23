package handlers

import (
	"net/http"

	"github.com/CXTACLYSM/hiring-api/pkg/shared/infrastructure/httputils"
	"github.com/CXTACLYSM/hiring-api/pkg/shared/infrastructure/responses"
)

type InfoHandler struct {
	version string
}

func NewInfoHandler(version string) *InfoHandler {
	return &InfoHandler{
		version: version,
	}
}

// ServeHTTP godoc
// @Summary     Service info
// @Description Returns current service version
// @Tags        system
// @Produce     json
// @Success     200 {object} responses.InfoResponse
// @Router      /info [get]
func (h *InfoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	httputils.ResponseOk(w, http.StatusOK, responses.NewInfoResponse(h.version))
}
