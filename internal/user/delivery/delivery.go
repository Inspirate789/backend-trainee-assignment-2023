package delivery

import (
	"fmt"
	"github.com/Inspirate789/backend-trainee-assignment-2023/internal/user/usecase/dto"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/multierr"
	"log/slog"
	"os"
	"strconv"
)

type fiberUserDelivery struct {
	useCase    UseCase
	volumePath string
	logger     *slog.Logger
}

const (
	userReportRoute   = "/report"
	filenameParamName = "filename"
	reportFilename    = "report.csv"
)

func NewFiberDelivery(api fiber.Router, useCase UseCase, logger *slog.Logger) {
	handler := &fiberUserDelivery{
		useCase: useCase,
		logger:  logger.WithGroup("fiberUserDelivery"),
	}
	api.Post("/user", handler.postUser)
	api.Delete("/user", handler.deleteUser)
	api.Patch("/user/segments", handler.patchUserSegments)
	api.Get("/user/segments", handler.getUserSegments)
	api.Get("/user/history", handler.getUserHistory)
	api.Get(fmt.Sprintf("%s/:%s", userReportRoute, filenameParamName), handler.getReport)
}

func (d *fiberUserDelivery) postUser(ctx *fiber.Ctx) error {
	var body dto.UserDTO
	err := ctx.BodyParser(&body)
	if err != nil {
		d.logger.Error(err.Error())
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	d.logger.Debug("request body received", slog.String("body", fmt.Sprintf("%v", body)))

	err = d.useCase.AddUser(body)
	if err != nil {
		d.logger.Error(err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.SendStatus(fiber.StatusOK)
}

func (d *fiberUserDelivery) deleteUser(ctx *fiber.Ctx) error {
	var body dto.UserInputDTO
	err := ctx.BodyParser(&body) // TODO: generic body and params parsing (fiber_utils package)?
	if err != nil {
		d.logger.Error(err.Error())
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	d.logger.Debug("request body received", slog.String("body", fmt.Sprintf("%v", body)))

	err = d.useCase.RemoveUser(body)
	if err != nil {
		d.logger.Error(err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.SendStatus(fiber.StatusOK)
}

func (d *fiberUserDelivery) patchUserSegments(ctx *fiber.Ctx) error {
	var body dto.UserSegmentsInputDTO
	err := ctx.BodyParser(&body)
	if err != nil {
		d.logger.Error(err.Error())
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	d.logger.Debug("request body received", slog.String("body", fmt.Sprintf("%v", body)))

	err = d.useCase.ChangeUserSegments(body)
	if err != nil {
		d.logger.Error(err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.SendStatus(fiber.StatusOK)
}

func (d *fiberUserDelivery) getUserSegments(ctx *fiber.Ctx) error {
	var body dto.UserInputDTO
	err := ctx.BodyParser(&body)
	if err != nil {
		d.logger.Error(err.Error())
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	d.logger.Debug("request body received", slog.String("body", fmt.Sprintf("%v", body)))

	segments, err := d.useCase.GetUserSegments(body)
	if err != nil {
		d.logger.Error(err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(segments)
}

func (d *fiberUserDelivery) getUserHistory(ctx *fiber.Ctx) error {
	var body dto.UserHistoryInputDTO
	err := ctx.BodyParser(&body)
	if err != nil {
		d.logger.Error(err.Error())
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	d.logger.Debug("request body received", slog.String("body", fmt.Sprintf("%v", body)))

	filename, err := d.useCase.SaveUserHistory(body, strconv.FormatUint(ctx.Context().ID(), 10))
	if err != nil {
		d.logger.Error(err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"web_link": ctx.BaseURL() + userReportRoute + "/" + filename,
	})
}

func (d *fiberUserDelivery) getReport(ctx *fiber.Ctx) (err error) {
	filename := ctx.Params(filenameParamName, "")
	if filename == "" {
		d.logger.Error("empty filename received")
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "empty filename received",
		})
	}
	d.logger.Debug("path param received", slog.String("filename", filename))

	ctx.Attachment(reportFilename)
	defer func() {
		err = multierr.Combine(err, os.Remove(reportFilename))
	}()

	return ctx.Status(fiber.StatusOK).SendFile(filename)
}
