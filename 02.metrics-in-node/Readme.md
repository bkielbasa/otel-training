# OTEL

Create a new meter

```js
const opentelemetry = require("@opentelemetry/api");
const meter = opentelemetry.metrics.getMeter('my-meter');

const requestCounter = meter.createCounter('request_sent_counter', {
  description: 'The number of requests we sent to external services'
});

requestCounter.add(1, { 'action.type': 'create' });
```

