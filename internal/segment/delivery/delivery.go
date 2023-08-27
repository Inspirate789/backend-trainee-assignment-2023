package delivery

import (
	"fmt"
	"github.com/Inspirate789/backend-trainee-assignment-2023/internal/segment/usecase/dto"
	"github.com/gofiber/fiber/v2"
	"log/slog"
)

type fiberSegmentDelivery struct {
	useCase UseCase
	logger  *slog.Logger
}

func NewFiberDelivery(api fiber.Router, useCase UseCase, logger *slog.Logger) {
	handler := &fiberSegmentDelivery{
		useCase: useCase,
		logger:  logger.WithGroup("fiberSegmentDelivery"),
	}
	api.Post("/segment", handler.postSegment)
	api.Delete("/segment", handler.deleteSegment)
}

// postSegment godoc
//
//	@Summary		Add new segment.
//	@Description	add new segment
//	@Tags			Segment API
//	@Param			SegmentDTO	body	dto.SegmentDTO	true	"Segment data"
//	@Accept			json
//	@Success		200
//	@Failure		422	{object}	string
//	@Failure		500	{object}	string
//	@Router			/segment [post]
func (d *fiberSegmentDelivery) postSegment(ctx *fiber.Ctx) error {
	var body dto.SegmentDTO
	err := ctx.BodyParser(&body)
	if err != nil {
		d.logger.Error(err.Error())
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	d.logger.Debug("request body received", slog.String("body", fmt.Sprintf("%v", body)))

	err = d.useCase.AddSegment(body)
	if err != nil {
		d.logger.Error(err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.SendStatus(fiber.StatusOK)
}

// deleteSegment godoc
//
//	@Summary		Delete segment.
//	@Description	delete segment
//	@Tags			Segment API
//	@Param			SegmentDTO	body	dto.SegmentNameDTO	true	"Segment name"
//	@Accept			json
//	@Success		200
//	@Failure		422	{object}	string
//	@Failure		500	{object}	string
//	@Router			/segment [delete]
func (d *fiberSegmentDelivery) deleteSegment(ctx *fiber.Ctx) error {
	var body dto.SegmentNameDTO
	err := ctx.BodyParser(&body)
	if err != nil {
		d.logger.Error(err.Error())
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	d.logger.Debug("request body received", slog.String("body", fmt.Sprintf("%v", body)))

	err = d.useCase.RemoveSegment(body)
	if err != nil {
		d.logger.Error(err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.SendStatus(fiber.StatusOK)
}
