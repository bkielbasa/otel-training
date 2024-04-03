# OTEL

Useful links:

 * https://docs.newrelic.com/docs/more-integrations/open-source-telemetry-integrations/opentelemetry/get-started/opentelemetry-set-up-your-app/#review-settings
 * https://opentelemetry.io/docs/specs/otel/configuration/sdk-environment-variables/
 * https://opentelemetry.io/docs/languages/js/exporters/#otlp-dependencies

## Task 1

Add autoinstrumentation to the NodeJS app (api).

### Step 1

Install dependencies:

```sh
npm install @opentelemetry/sdk-node \
  @opentelemetry/api \
  @opentelemetry/auto-instrumentations-node \
  @opentelemetry/sdk-metrics \
  @opentelemetry/sdk-trace-node

```

### Step 2

Create `instrumentation.js` file.

```js
const { NodeSDK } = require('@opentelemetry/sdk-node');
const { ConsoleSpanExporter } = require('@opentelemetry/sdk-trace-node');
const {
  getNodeAutoInstrumentations,
} = require('@opentelemetry/auto-instrumentations-node');
const {
  PeriodicExportingMetricReader,
  ConsoleMetricExporter,
} = require('@opentelemetry/sdk-metrics');

const sdk = new NodeSDK({
  traceExporter: new ConsoleSpanExporter(),
  metricReader: new PeriodicExportingMetricReader({
    exporter: new ConsoleMetricExporter(),
  }),
  instrumentations: [getNodeAutoInstrumentations()],
});

sdk.start();

```

... and run it!


```sh
node --require ./instrumentation.js app.js
```

Right now, you should see all metrics and traces in the stdout.

### Step 3

Configure env variables

```
OTEL_EXPORTER_OTLP_ENDPOINT=https://otlp.eu01.nr-data.net
OTEL_EXPORTER_OTLP_HEADERS=api-key=API_KEY
OTEL_SERVICE_NAME=api
```

### Step 4

Change the configuration to use OTEL exporter

```js
require('dotenv').config()

const { NodeSDK } = require('@opentelemetry/sdk-node');
const { OTLPTraceExporter } = require('@opentelemetry/exporter-trace-otlp-http');
const {
  getNodeAutoInstrumentations,
} = require('@opentelemetry/auto-instrumentations-node');
const {
  PeriodicExportingMetricReader,
} = require('@opentelemetry/sdk-metrics');
const { OTLPMetricExporter } = require('@opentelemetry/exporter-metrics-otlp-http');

const otlpTraceExporter = new OTLPTraceExporter({});
const otlpMetricExporter = new OTLPMetricExporter({});

const { diag, DiagConsoleLogger, DiagLogLevel } = require('@opentelemetry/api');

// For troubleshooting, set the log level to DiagLogLevel.DEBUG
diag.setLogger(new DiagConsoleLogger(), DiagLogLevel.INFO);

const sdk = new NodeSDK({
  traceExporter: otlpTraceExporter,
  metricReader: new PeriodicExportingMetricReader({
    exporter: otlpMetricExporter,
    exportIntervalMillis: 60000, // Export interval in milliseconds. Adjust as needed.
  }),
  instrumentations: [getNodeAutoInstrumentations()],
});

sdk.start();

```

### Step 5

Go to newrelic dashboard and see your traces there!
