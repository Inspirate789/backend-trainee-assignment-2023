package main

import (
	"context"
	"fmt"
	"github.com/Inspirate789/backend-trainee-assignment-2023/internal/pkg/app"
	"github.com/Inspirate789/backend-trainee-assignment-2023/internal/segments/repository"
	"github.com/Inspirate789/backend-trainee-assignment-2023/internal/segments/usecase"
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

func runApp(webApp app.WebApp, logger *slog.Logger) {
	logger.Debug(fmt.Sprintf("web app starts at port %s with configuration: \n%v",
		viper.GetString("PARSER_PORT"), viper.AllSettings()),
	)

	go func() {
		err := webApp.Start(viper.GetString("PARSER_PORT"))
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
	logLevel.Set(slog.LevelInfo)
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

	repo := repository.NewSqlxSegmentRepository(db, logger)
	useCase := usecase.NewSegmentUseCase(repo, logger)
	webApp := app.NewFiberApp(viper.GetString("PARSER_PORT"), viper.GetString("API_PREFIX"), useCase, logger, logLevel.Level())

	runApp(webApp, logger)
	shutdownApp(webApp, logger)
}
