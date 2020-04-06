package sink

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/influxdata/influxdb-client-go"
)

type InfluxDBSink struct {
	client influxdb2.InfluxDBClient
}

func NewInfluxDBSink(address, token string) Sink {
	sink := &InfluxDBSink{
		client: influxdb2.NewClient(address, token),
	}
	return sink
}

func (sink *InfluxDBSink) Put(data interface{}) error {
	writeApi := sink.client.WriteApiBlocking("my-org", "my-bucket") // TODO: get org and bucket from config

	switch d := data.(type) {
	case map[string]interface{}:
		point := influxdb2.NewPoint("measurement", nil, d, time.Now()) // TODO: get measurement from config
		return writeApi.WritePoint(context.Background(), point)
	case string:
		return writeApi.WriteRecord(context.Background(), d)
	default:
		return errors.New(fmt.Sprintf("InfluxDB sink does not support the type passed: %v", reflect.TypeOf(d)))
	}
}

func (sink *InfluxDBSink) PutMany(data *[]interface{}) error {
    size := uint(len(*data))
    if size == 0 { // TODO: check for coverage inc/dec
        return nil
    }

    sink.client.Options().SetBatchSize(size).SetUseGZip(true) // TODO: get gzip from config and check for performance changes
    writeApi := sink.client.WriteApiBlocking("my-org", "my-bucket") // TODO: get org and bucket from config

    switch d := (*data)[0].(type) {
    case map[string]interface{}:
        for _, entry := range *data {
            point := influxdb2.NewPoint("measurement", nil, entry.(map[string]interface{}), time.Now()) // TODO: get measurement from config
            if err := writeApi.WritePoint(context.Background(), point); err != nil {
                return err
            }
        }
    case string:
        for _, entry := range *data {
            if err := writeApi.WriteRecord(context.Background(), entry.(string)); err != nil {
                return err
            }
        }
    default:
        return errors.New(fmt.Sprintf("InfluxDB sink does not support the type passed: %v", reflect.TypeOf(d)))
    }

    return nil
}

func (sink *InfluxDBSink) Close() error {
	sink.client.Close()
	return nil
}
