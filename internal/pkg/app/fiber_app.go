package app

import (
	"context"
	segmentDelivery "github.com/Inspirate789/backend-trainee-assignment-2023/internal/segment/delivery"
	userDelivery "github.com/Inspirate789/backend-trainee-assignment-2023/internal/user/delivery"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"log/slog"
)

type fiberApp struct {
	fiber  *fiber.App
	logger *slog.Logger
}

type ApiSettings struct {
	Port      string
	ApiPrefix string
}

type UseCases struct {
	SegmentUseCase segmentDelivery.UseCase
	UserUseCase    userDelivery.UseCase
}

func NewFiberApp(settings ApiSettings, useCases UseCases, log *slog.Logger) WebApp {
	app := fiber.New()

	app.Use(recover.New())
	app.Use(requestid.New())
	app.Use(logger.New(logger.Config{
		Format: "${pid} ${locals:requestid} ${status} - ${latency} ${method} ${path}\n",
		Output: slog.NewLogLogger(log.Handler(), slog.LevelDebug).Writer(),
	}))

	//app.Get("/swagger/*", swagger.New(swagger.Config{ // TODO
	//	URL:          fmt.Sprintf("http://localhost:%s/swagger/doc.json", port),
	//	DeepLinking:  false,
	//	DocExpansion: "none",
	//}))

	api := app.Group(settings.ApiPrefix)
	segmentDelivery.NewFiberDelivery(api, useCases.SegmentUseCase, log)
	userDelivery.NewFiberDelivery(api, useCases.UserUseCase, log)

	return &fiberApp{
		fiber:  app,
		logger: log,
	}
}

func (f *fiberApp) Start(port string) error {
	return f.fiber.Listen(":" + port)
}

func (f *fiberApp) Stop(ctx context.Context) error {
	return f.fiber.ShutdownWithContext(ctx)
}
