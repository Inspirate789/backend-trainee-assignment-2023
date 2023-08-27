package main

import (
	"context"
	"fmt"
	"github.com/Inspirate789/backend-trainee-assignment-2023/internal/pkg/app"
	segmentRepository "github.com/Inspirate789/backend-trainee-assignment-2023/internal/segment/repository"
	segmentUseCase "github.com/Inspirate789/backend-trainee-assignment-2023/internal/segment/usecase"
	userFsRepository "github.com/Inspirate789/backend-trainee-assignment-2023/internal/user/repository/fs"
	userSqlRepository "github.com/Inspirate789/backend-trainee-assignment-2023/internal/user/repository/sql"
	userUseCase "github.com/Inspirate789/backend-trainee-assignment-2023/internal/user/usecase"
	"github.com/Inspirate789/backend-trainee-assignment-2023/pkg/influx"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func readConfig() error {
	var configPath string
	pflag.StringVarP(&configPath, "config", "c", "", "Config file path")
	pflag.Parse()
	if configPath == "" {
		return errors.New("config file is not specified")
	}
	slog.Info(fmt.Sprintf("Config path: %s", configPath))

	viper.SetConfigFile(configPath)

	return viper.ReadInConfig()
}

func runApp(webApp app.WebApp, port string, logger *slog.Logger) {
	logger.Debug(fmt.Sprintf("web app starts at port %s with configuration: \n%v",
		port, viper.AllSettings()),
	)

	go func() {
		err := webApp.Start(port)
		if err != nil {
			panic(err)
		}
	}()
}

func shutdownApp(webApp app.WebApp, logger *slog.Logger) {
	quit := make(chan os.Signal)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Debug("shutdown web app ...")

	err := webApp.Stop(context.Background())
	if err != nil {
		panic(errors.Wrap(err, "app shutdown"))
	}
	logger.Debug("web app exited")
}

func main() {
	err := readConfig()
	if err != nil {
		panic(err)
	}
	time.Sleep(time.Duration(viper.GetInt("INIT_SLEEP_TIME")) * time.Second)

	iw := influx.NewWriter()
	err = iw.Open(
		context.Background(),
		viper.GetString("INFLUXDB_URL"),
		viper.GetString("INFLUXDB_TOKEN"),
		viper.GetString("INFLUXDB_ORG"),
		viper.GetString("INFLUXDB_BACKEND_BUCKET_NAME"),
	)
	if err != nil {
		panic(err)
	}
	defer iw.Close()

	logLevel := new(slog.LevelVar)
	logLevel.Set(slog.LevelDebug)
	logger := slog.New(slog.NewTextHandler(iw, &slog.HandlerOptions{
		AddSource:   true,
		Level:       logLevel.Level(),
		ReplaceAttr: nil,
	}))

	db, err := sqlx.Connect(viper.GetString("DB_DRIVER_NAME"), viper.GetString("DB_CONNECTION_STRING"))
	if err != nil {
		panic(err)
	}
	defer func(db *sqlx.DB) {
		err = db.Close()
		if err != nil {
			panic(err)
		}
	}(db)

	segmentRepo := segmentRepository.NewSqlxRepository(db, logger)
	segmentUC := segmentUseCase.NewUseCase(segmentRepo, logger)

	userSqlRepo := userSqlRepository.NewSqlxRepository(db, logger)
	userFsRepo := userFsRepository.NewFsRepository(viper.GetString("VOLUME_PATH"), logger)
	userUC := userUseCase.NewUseCase(userSqlRepo, userFsRepo, logger)

	useCases := app.UseCases{
		SegmentUseCase: segmentUC,
		UserUseCase:    userUC,
	}

	settings := app.ApiSettings{
		Port:      viper.GetString("PORT"),
		ApiPrefix: viper.GetString("API_PREFIX"),
	}

	webApp := app.NewFiberApp(settings, useCases, logger)

	runApp(webApp, settings.Port, logger)
	shutdownApp(webApp, logger)
}
