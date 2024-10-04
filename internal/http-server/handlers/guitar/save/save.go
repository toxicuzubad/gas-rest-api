package save

import (
	resp "gas-rest-api/internal/lib/api/response"
	"gas-rest-api/internal/lib/logger/sl"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type Request struct {
	ManufacturerName string `json:"manufacturerName" validate:"required,manufacturerName"`
	ModelName        string `json:"modelName" validate:"required,modelName"`
	Description      string `json:"description,omitempty"`
	SerialNumber     string `json:"serialNumber" validate:"required,serialNumber"`
}

type Response struct {
	resp.Response
	Id int64 `json:"id"`
}

type GuitarSaver interface {
	SaveGuitar(manufacturerName string, modelName string, description string, serialNumber string) (int64, error)
}

func New(log *slog.Logger, guitarSaver GuitarSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.guitar.save.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)

		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			validateError := err.(validator.ValidationErrors)

			log.Error("invalid request", sl.Err(err))

			render.JSON(w, r, resp.ValidationError(validateError))

			return
		}

		id, err := guitarSaver.SaveGuitar(req.ManufacturerName, req.ModelName, req.Description, req.SerialNumber)

		if err != nil {
			log.Error("Saving guitar error", sl.Err(err))

			render.JSON(w, r, resp.Error("Saving guitar error"))

			return
		}

		log.Info("guitar added", slog.Int64("id", id))

		render.JSON(w, r, Response{
			Response: resp.OK(),
			Id:       id,
		})
	}
}
