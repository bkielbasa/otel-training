package tracescleaner

import (
	"context"
	"regexp"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/ptrace"
	"go.opentelemetry.io/collector/processor"
	"go.uber.org/zap"
)

type ProcessorConfig struct {
	Exclude []string `mapstructure:"exclude"`
}

type tracesProcessor struct {
	config       *ProcessorConfig
	logger       *zap.Logger
	nextConsumer consumer.Traces
	params       processor.CreateSettings
}

func (tp *tracesProcessor) Start(ctx context.Context, host component.Host) error {
	return nil
}

func (tp *tracesProcessor) Shutdown(ctx context.Context) error {
	return nil
}

func (tp *tracesProcessor) ConsumeTraces(ctx context.Context, td ptrace.Traces) error {
	es := td.ResourceSpans()

	for i := 0; i < es.Len(); i++ {
		e := es.At(i)
		tp.logger.Info("ResourceSpans", zap.Any("attributes", e.Resource().Attributes().AsRaw()))
		for j := 0; j < e.ScopeSpans().Len(); j++ {
			ss := e.ScopeSpans().At(j)

			tp.logger.Info("ScopeSpans", zap.Any("attributes", ss.Scope().Attributes().AsRaw()))

			spanIDsToRemove := make([]pcommon.SpanID, 0)

			for k := 0; k < ss.Spans().Len(); k++ {
				s := ss.Spans().At(k)

				tp.logger.Info("Span", zap.Any("attributes", s.Attributes().AsRaw()))
				tp.logger.Info("Span name", zap.Any("attributes", s.Name()))

				for _, ex := range tp.config.Exclude {
					r, _ := regexp.Compile(ex)
					if r.MatchString(s.Name()) {
						spanIDsToRemove = append(spanIDsToRemove, s.SpanID())
						tp.logger.Info("Exclude span", zap.Any("span", s.Name()))
						continue
					}
				}
			}

			ss.Spans().RemoveIf(func(s ptrace.Span) bool {
				for _, id := range spanIDsToRemove {
					if s.SpanID() == id {
						return true
					}
				}
				return false
			})
		}
	}

	tp.logger.Info("ConsumeTraces", zap.Any("traces", td.SpanCount()))
	tp.nextConsumer.ConsumeTraces(ctx, td)
	return nil
}

func (tp *tracesProcessor) Capabilities() consumer.Capabilities {
	return consumer.Capabilities{MutatesData: true}
}

func createTracesProcessor(ctx context.Context, params processor.CreateSettings, baseCfg component.Config, consumer consumer.Traces) (processor.Traces, error) {
	cfg := baseCfg.(*ProcessorConfig)
	logger := params.Logger

	return &tracesProcessor{
		config:       cfg,
		logger:       logger,
		nextConsumer: consumer,
		params:       params,
	}, nil
}

func createProcessorDefaultConfig() component.Config {
	return &ProcessorConfig{}
}

func NewProcessorFactory() processor.Factory {
	typeStr, _ := component.NewType("tracescleaner")

	return processor.NewFactory(
		typeStr,
		createProcessorDefaultConfig,
		processor.WithTraces(createTracesProcessor, component.StabilityLevelAlpha),
	)
}
