# Creating traces cleaner

Make sure you have jaeger up and running

```sh
docker run -d --name jaeger \
  -e COLLECTOR_OTLP_ENABLED=true \
  -p 16686:16686 \
  -p 14317:4317 \
  -p 14318:4318 \
  jaegertracing/all-in-one:1.41
```

Instal `ocb` tool to generate some code

```sh
curl --proto '=https' --tlsv1.2 -fL -o ocb \
https://github.com/open-telemetry/opentelemetry-collector/releases/download/cmd%2Fbuilder%2Fv0.97.0/ocb_0.97.0_linux_amd64
chmod +x ocb
```

When you want to generate a trace, use this

```sh
go install github.com/open-telemetry/opentelemetry-collector-contrib/cmd/telemetrygen@latest

# and generate a single trace
telemetrygen traces --otlp-insecure --traces 1
```

We want to achieve a new receiver like this

```yaml
receivers:
  tracescleaner: # this line represents the ID of your receiver
    exclude:
      - dupa
      - ([\d]{3})
```

## Config

To do so, we have to describe the config in traces-cleaner/config.go

```go
type Config struct {
	Exclude []string `mapstructure:"exclude"`
}
```

We can write a validation function just to make sure that all patterns are valid regexps:

```go
package tracescleaner

func (cfg *Config) Validate() error {
	if len(cfg.Exclude) == 0 {
		return errors.New("exclude list is empty")
	}

	for _, e := range cfg.Exclude {
		_, err := regexp.Compile(e)
		if err != nil {
			return err
		}
	}

	return nil
}
```

Your task is to use the `_, err := regexp.Compile(e)` function to iterate over all parameters and check if they are valid regexps.

## Factory

Create a new `factory.go` file inside the new receiver with the following content


```go
package tracescleaner

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/receiver"
)

func createDefaultConfig() component.Config {
	return &Config{}
}

func createTracesReceiver(_ context.Context, params receiver.CreateSettings, baseCfg component.Config, consumer consumer.Traces) (receiver.Traces, error) {
	return nil, nil
}

func NewFactory() receiver.Factory {
	typeStr, _ := component.NewType("tracescleaner")

	return receiver.NewFactory(
		typeStr,
		createDefaultConfig,
		receiver.WithTraces(createTracesReceiver, component.StabilityLevelAlpha))
}
```

# Trace receiver

Create a new file: `trace-receiver.go`.

```go
package tracescleaner

import (
	"context"

	"go.opentelemetry.io/collector/component"
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

	return nil
}

func (tr *traceReceiver) Shutdown(ctx context.Context) error {
	tr.cancel()
	return nil
}
```

We are ready to update the function

```go
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
```

Update `components.go` file to add our receiver

```go
	factories.Receivers, err = receiver.MakeFactoryMap(
		otlpreceiver.NewFactory(),
		tracescleaner.NewFactory(),
	)
```

and update `config.yaml` to add our receiver.

```yaml
  tracescleaner: # this line represents the ID of your receiver
    exclude:
      - dupa
      - ([\d]{3})
```

and regenerate it!

```sh
./ocb --config builder-config.yaml
go run ./otelcol-dev --config config.yaml
```

OK, we can generate new traces! But we don't want to do it...

Your task is to use the one in https://github.com/open-telemetry/opentelemetry-collector/blob/main/receiver/otlpreceiver/factory.go as the base in our receiver and skip those traces that doesn't match our patterns.