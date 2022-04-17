package get

import (
	"net/http"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"

	"github.com/andreyAKor/nut_client_service/internal/http/clients/nut"
)

type Handler struct {
	nutClient *nut.Client
}

func New(nutClient *nut.Client) *Handler {
	return &Handler{
		nutClient: nutClient,
	}
}

func (h *Handler) Handle() func(_ http.ResponseWriter, r *http.Request) (interface{}, error) {
	return func(w http.ResponseWriter, r *http.Request) (interface{}, error) {
		list, err := h.nutClient.GetUPSList(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Error().Err(err).Msg("get UPS list fail")

			return nil, errors.Wrap(err, "get UPS list fail")
		}

		return convertListToList(list), nil
	}
}
