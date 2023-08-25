package delivery

import (
	"github.com/gofiber/fiber/v2"
	"log/slog"
)

type fiberSegmentDelivery struct {
	useCase SegmentUseCase
	logger  *slog.Logger
}

func NewFiberSegmentDelivery(api fiber.Router, useCase SegmentUseCase, logger *slog.Logger) {
	handler := &fiberSegmentDelivery{
		useCase: useCase,
		logger:  logger.WithGroup("fiberSegmentDelivery"),
	}
	api.Post("/segment", handler.postSegment)
	api.Delete("/segment", handler.deleteSegment)
}

func (d *fiberSegmentDelivery) postSegment(ctx *fiber.Ctx) error {
	// TODO

	return ctx.SendStatus(fiber.StatusOK)
}

func (d *fiberSegmentDelivery) deleteSegment(ctx *fiber.Ctx) error {
	// TODO

	return ctx.SendStatus(fiber.StatusOK)
}
