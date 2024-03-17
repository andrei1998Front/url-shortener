package save

import (
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/andrei1998Front/url-shortener/internal/lib/api/response"
	sl "github.com/andrei1998Front/url-shortener/internal/lib/logger/slog"
	"github.com/andrei1998Front/url-shortener/internal/lib/random"
	"github.com/andrei1998Front/url-shortener/internal/storage"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	response.Response
	Alias string `json:"alias,omitempty`
}

// TODO: move to config
const aliasLength = 6

//go:generate go run github.com/vektra/mockery/v2@v2.42.1 --name=URLSaver
type URLSaver interface {
	SaveURL(urlToSave string, alias string) error
}

func New(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			log.Error("request body is empty")

			render.JSON(w, r, response.Error("empty request"))

			return
		}
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			render.JSON(w, r, response.Error("failed to decode request"))

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request", sl.Err(err))

			render.JSON(w, r, response.ValidationError(validateErr))
			return
		}

		var alias string

		if req.Alias != "" {
			alias = req.Alias
		} else {
			rndAlias, err := random.NewRandomString(aliasLength)

			if err != nil {
				log.Error("invalid alias size", sl.Err(err))
				render.JSON(w, r, response.Error("failed to set alias"))

				return
			}

			alias = rndAlias
		}

		err = urlSaver.SaveURL(req.URL, alias)

		if errors.Is(err, storage.ErrUrlExists) {
			log.Info("url already exists", slog.String("url", req.URL))
			render.JSON(w, r, response.Error("url already exists"))

			return
		}

		if err != nil {
			log.Error("failed to add url", sl.Err(err))

			render.JSON(w, r, response.Error("failed to add url"))

			return
		}

		log.Info("url added")

		render.JSON(w, r, Response{
			Response: response.OK(),
			Alias:    alias,
		})
	}
}
