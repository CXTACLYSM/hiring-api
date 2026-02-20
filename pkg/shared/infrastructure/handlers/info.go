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

func (h *InfoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	httputils.ResponseOk(w, http.StatusOK, responses.NewInfoResponse(h.version))
}
