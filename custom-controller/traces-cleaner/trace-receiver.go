package tracescleaner

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/receiver"
	"go.opentelemetry.io/collector/receiver/otlpreceiver"
	"go.uber.org/zap"
)

type traceReceiver struct {
	host   component.Host
	cancel context.CancelFunc

	logger       *zap.Logger
	nextConsumer consumer.Traces
	config       *Config
	params       receiver.CreateSettings

	otelRec receiver.Traces
}

func (tr *traceReceiver) Start(ctx context.Context, host component.Host) error {
	tr.host = host
	ctx = context.Background()
	ctx, tr.cancel = context.WithCancel(ctx)

	tr.logger.Info("Starting trace receiver", zap.String("name", "tracescleaner"))
	otel := otlpreceiver.NewFactory()
	rec, err := otel.CreateTracesReceiver(ctx, tr.params, &tr.config.OTEL, tr.nextConsumer)
	if err != nil {
		return err
	}

	serv := http.Server{
		Addr: ":4318",
	}

	client := http.Client{}

	serv.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		body := string(b)

		for _, e := range tr.config.Exclude {
			re, _ := regexp.Compile(e)
			if re.MatchString(body) {
				tr.logger.Info("Excluding trace", zap.String("body", body))
				return
			}
		}

		r.URL = &url.URL{
			Path:   tr.config.OTEL.Protocols.HTTP.TracesURLPath,
			Host:   tr.config.OTEL.Protocols.HTTP.Endpoint,
			Scheme: "http",
		}

		r.RequestURI = ""
		r.Body = io.NopCloser(strings.NewReader(body))

		resp, err := client.Do(r.WithContext(ctx))

		if err != nil {
			tr.logger.Error("Failed to forward trace", zap.Error(err))
			b, _ = io.ReadAll(resp.Body)
			tr.logger.Error("Failed to forward trace", zap.String("body", string(b)))
			return
		}

		resp.Body.Close()

		tr.logger.Info("Trace forwarded", zap.String("path", r.URL.Path))
	})

	go func() {
		serv.ListenAndServe()
	}()

	tr.otelRec = rec

	// if err = rec.Start(ctx, host); err != nil {
	// 	return err
	// }

	return nil
}

func (tr *traceReceiver) Shutdown(ctx context.Context) error {
	if tr.cancel != nil {
		tr.cancel()
	}

	if tr.otelRec != nil {
		if err := tr.otelRec.Shutdown(ctx); err != nil {
			return err
		}
	}

	return nil
}
