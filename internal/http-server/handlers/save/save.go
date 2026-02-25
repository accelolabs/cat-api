package save

import (
	"accelolabs/cat-api/internal/meow"
	"accelolabs/cat-api/internal/storage"
	"accelolabs/cat-api/internal/util/ctx"
	"accelolabs/cat-api/internal/util/response"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type Request struct {
	URL string `json:"url" validate:"required,url"`
}

type Response struct {
	response.Response
	Alias string `json:"alias,omitempty"`
}

type URLSaver interface {
	SaveURL(targetUrl string, alias string) (int64, error)
	GetAlias(targetUrl string) (string, error)
}

func New(log *slog.Logger, urlSaver URLSaver, aliasLength int, maxStretch int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.save.New"

		requestID := ctx.RequestID(r.Context())

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", requestID),
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.String("remote_addr", r.RemoteAddr),
		)

		w.Header().Set("Content-Type", "application/json")

		if r.Method != http.MethodPost {
			log.Warn("method not allowed")
			response.JSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		var req Request

		err := json.NewDecoder(r.Body).Decode(&req)
		if errors.Is(err, io.EOF) {
			log.Error("request body is empty")
			response.JSONError(w, http.StatusBadRequest, "empty request body")
			return
		}
		if err != nil {
			log.Error("failed to decode request body", slog.String("error", err.Error()))
			response.JSONError(w, http.StatusBadRequest, "invalid request format")
			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request", slog.String("error", err.Error()))

			response.JSONError(w, http.StatusBadRequest, response.ValidationError(validateErr))

			return
		}

		alias, err := urlSaver.GetAlias(req.URL)
		if err != nil {
			log.Info("generated new alias")
			alias = meow.Meow(aliasLength, maxStretch) // TODO: remove magic values
		}

		id, err := urlSaver.SaveURL(req.URL, alias)
		if err != nil && !errors.Is(err, storage.ErrURLExists) {
			log.Error("failed to add url", slog.String("error", err.Error()))
			response.JSONError(w, http.StatusBadRequest, "failed to add url")
			return
		}

		log.Info("url added", slog.Int64("id", id))

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(Response{
			Response: response.OK(),
			Alias:    alias,
		})
	}
}
