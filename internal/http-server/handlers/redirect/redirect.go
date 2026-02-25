package redirect

import (
	"accelolabs/cat-api/internal/storage"
	"accelolabs/cat-api/internal/util/ctx"
	"accelolabs/cat-api/internal/util/response"
	"errors"
	"log/slog"
	"net/http"
)

type URLGetter interface {
	GetURL(alias string) (string, error)
}

func New(log *slog.Logger, urlGetter URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.redirect.New"

		requestID := ctx.RequestID(r.Context())

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", requestID),
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.String("remote_addr", r.RemoteAddr),
		)

		alias := r.PathValue("alias")
		if alias == "" {
			log.Info("alias is empty")

			response.JSONError(w, http.StatusBadRequest, "invalid request")

			return
		}

		resURL, err := urlGetter.GetURL(alias)
		if errors.Is(err, storage.ErrURLNotFound) {
			log.Info("url not found", "alias", alias)

			response.JSONError(w, http.StatusBadRequest, "not found")

			return
		}
		if err != nil {
			log.Error("failed to get url", slog.String("error", err.Error()))

			response.JSONError(w, http.StatusBadRequest, "internal error")

			return
		}

		log.Info("got url", slog.String("url", resURL))

		http.Redirect(w, r, resURL, http.StatusFound)
	}
}
