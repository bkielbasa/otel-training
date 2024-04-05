package tracescleaner

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.uber.org/zap"
)

type traceReceiver struct {
	host   component.Host
	cancel context.CancelFunc

	logger       *zap.Logger
	nextConsumer consumer.Traces
	config       *Config
}

func (tr *traceReceiver) Start(ctx context.Context, host component.Host) error {
	tr.host = host
	ctx = context.Background()
	ctx, tr.cancel = context.WithCancel(ctx)

	<-ctx.Done()

	return nil
}

func (tr *traceReceiver) Shutdown(ctx context.Context) error {
	if tr.cancel != nil {
		tr.cancel()
	}
	return nil
}
