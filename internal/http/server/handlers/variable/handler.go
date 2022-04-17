package variable

import (
	"encoding/json"
	"io"
	"io/ioutil"
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

func (h *Handler) Handle() func(http.ResponseWriter, *http.Request) (interface{}, error) {
	return func(w http.ResponseWriter, r *http.Request) (interface{}, error) {
		req, err := prepareCommand(r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Error().Err(err).Msg("prepare command struct from request body fail")

			return nil, errors.Wrap(err, "prepare command struct from request body fail")
		}

		if err := h.nutClient.SetVariable(r.Context(), req.Name, req.VariableName, req.Value); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Error().Err(err).Msg("set variable fail")

			return nil, errors.Wrap(err, "set variable fail")
		}

		return nil, nil
	}
}

func prepareCommand(r *http.Request) (*variable, error) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, errors.Wrap(err, "reading from body fail")
	}
	if _, err := io.Copy(ioutil.Discard, r.Body); err != nil {
		return nil, errors.Wrap(err, "copying from response body fail")
	}

	req := &variable{}
	if err := json.Unmarshal(data, req); err != nil {
		return nil, errors.Wrap(err, "json unmarshal fail")
	}

	return req, err
}
