package influx

import (
	"context"
	"errors"
	"sync"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/influxdata/influxdb-client-go/v2/domain"
)

type InfluxWriter struct {
	client   influxdb2.Client
	mx       sync.Mutex
	writeAPI api.WriteAPI
}

var (
	influxWriter     *InfluxWriter // Singleton
	influxWriterOnce sync.Once
)

func NewWriter() *InfluxWriter {
	influxWriterOnce.Do(func() {
		influxWriter = &InfluxWriter{}
	})

	return influxWriter
}

func (il *InfluxWriter) Open(ctx context.Context, url, token, org, bucketName string) error {
	client := influxdb2.NewClient(url, token)

	health, err := client.Health(ctx) // validate client connection health
	if (err != nil) && health != nil && health.Status == domain.HealthCheckStatusPass {
		client.Close()
		return errors.New("open error: database not healthy")
	}

	il.client = client
	il.writeAPI = il.client.WriteAPI(org, bucketName) // Get non-blocking write client

	return nil
}

func (il *InfluxWriter) Write(p []byte) (int, error) {
	point := influxdb2.NewPointWithMeasurement("Messages").
		AddField("message", p).
		SetTime(time.Now())

	il.mx.Lock()
	il.writeAPI.WritePoint(point)
	// Flush writes
	il.writeAPI.Flush()
	il.mx.Unlock()

	return len(p), nil
}

func (il *InfluxWriter) Close() {
	il.client.Close()
}
