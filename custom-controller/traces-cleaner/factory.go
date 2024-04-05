package tracescleaner

import (
	"context"
	"errors"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/receiver"
)

func createDefaultConfig() component.Config {
	return &Config{}
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
	}

	return traceRcvr, nil
}

func NewFactory() receiver.Factory {
	typeStr, _ := component.NewType("tracescleaner")

	return receiver.NewFactory(
		typeStr,
		createDefaultConfig,
		receiver.WithTraces(createTracesReceiver, component.StabilityLevelAlpha))
}
