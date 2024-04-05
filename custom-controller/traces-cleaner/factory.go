package tracescleaner

import (
	"context"
	"errors"
	"fmt"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/configgrpc"
	"go.opentelemetry.io/collector/config/confighttp"
	"go.opentelemetry.io/collector/config/confignet"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/receiver"
	"go.opentelemetry.io/collector/receiver/otlpreceiver"
)

const (
	grpcPort = 4317
	httpPort = 4318

	defaultTracesURLPath  = "/v1/traces"
	defaultMetricsURLPath = "/v1/metrics"
	defaultLogsURLPath    = "/v1/logs"
)

func createDefaultConfig() component.Config {
	return &Config{
		OTEL: otlpreceiver.Config{
			otlpreceiver.Protocols{
				GRPC: &configgrpc.ServerConfig{
					NetAddr: confignet.AddrConfig{
						Endpoint:  EndpointForPort(grpcPort),
						Transport: confignet.TransportTypeTCP,
					},
					// We almost write 0 bytes, so no need to tune WriteBufferSize.
					ReadBufferSize: 512 * 1024,
				},
				HTTP: &otlpreceiver.HTTPConfig{
					ServerConfig: &confighttp.ServerConfig{
						Endpoint: EndpointForPort(httpPort),
					},
					TracesURLPath:  defaultTracesURLPath,
					MetricsURLPath: defaultMetricsURLPath,
					LogsURLPath:    defaultLogsURLPath,
				},
			},
		},
	}
}

func EndpointForPort(port int) string {
	host := "localhost"

	return fmt.Sprintf("%s:%d", host, port)
}

func createTracesReceiver(_ context.Context, params receiver.CreateSettings, baseCfg component.Config, consumer consumer.Traces) (receiver.Traces, error) {
	if consumer == nil {
		return nil, errors.New("nil nextConsumer")
	}

	logger := params.Logger
	traceCleanerCfg := baseCfg.(*Config)

	traceRcvr := &traceReceiver{
		logger:       logger,
		nextConsumer: consumer,
		config:       traceCleanerCfg,
		params:       params,
	}

	return traceRcvr, nil
}

func NewFactory() receiver.Factory {
	typeStr, _ := component.NewType("tracescleaner")

	return receiver.NewFactory(
		typeStr,
		createDefaultConfig,
		receiver.WithTraces(createTracesReceiver, component.StabilityLevelAlpha),
	)
}
