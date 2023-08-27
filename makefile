.PHONY: run stop shutdown swagger

run:
	docker compose --env-file env/app.env up -d --build
	docker exec -ti influxdb influx bucket create -n app_bucket -o app_org -r 0

stop:
	docker compose --env-file env/app.env down

shutdown:
	docker compose --env-file env/app.env down --volumes

swagger:
	swag init --parseDependency --parseInternal --parseDepth 1 -g cmd/app/main.go -o swagger/
	swag fmt