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

// postUser godoc
//
//	@Summary		Add new user.
//	@Description	add new user
//	@Tags			User API
//	@Param			UserDTO	body	dto.UserDTO	true	"User data"
//	@Accept			json
//	@Success		200
//	@Failure		422	{object}	string
//	@Failure		500	{object}	string
//	@Router			/user [post]
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

// deleteUser godoc
//
//	@Summary		Delete user.
//	@Description	delete user
//	@Tags			User API
//	@Param			UserInputDTO	body	dto.UserInputDTO	true	"User ID"
//	@Accept			json
//	@Success		200
//	@Failure		422	{object}	string
//	@Failure		500	{object}	string
//	@Router			/user [delete]
func (d *fiberUserDelivery) deleteUser(ctx *fiber.Ctx) error {
	var body dto.UserInputDTO
	err := ctx.BodyParser(&body)
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

// patchUserSegments godoc
//
//	@Summary		Change user segments.
//	@Description	change user segments
//	@Tags			User API
//	@Param			UserSegmentsInputDTO	body	dto.UserSegmentsInputDTO	true	"Old and new user segments"
//	@Accept			json
//	@Success		200
//	@Failure		422	{object}	string
//	@Failure		500	{object}	string
//	@Router			/user/segments [patch]
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

// patchUserSegments godoc
//
//	@Summary		Get user segments.
//	@Description	get user segments
//	@Tags			User API
//	@Param			UserInputDTO	body	dto.UserInputDTO	true	"User ID"
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	dto.UserSegmentsOutputDTO
//	@Failure		422	{object}	string
//	@Failure		500	{object}	string
//	@Router			/user/segments [get]
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

// getUserHistory godoc
//
//	@Summary		Get the history of changing user segments.
//	@Description	get the history of changing user segments; returns the web link to csv file with report
//	@Tags			User API
//	@Param			UserHistoryInputDTO	body	dto.UserHistoryInputDTO	true	"Year and month in history"
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	string
//	@Failure		422	{object}	string
//	@Failure		500	{object}	string
//	@Router			/user/history [get]
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
