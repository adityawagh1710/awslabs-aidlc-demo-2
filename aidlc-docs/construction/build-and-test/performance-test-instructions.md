# Performance Test Instructions

## Requirements
- p95 API response time < 200ms
- p95 search response time < 500ms
- Support hundreds of concurrent users
- Error rate < 1%

## Tool: k6

```bash
# Install k6
brew install k6   # macOS
# or: https://k6.io/docs/getting-started/installation/
```

## Load Test Script

```javascript
// k6-load-test.js
import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
  stages: [
    { duration: '30s', target: 50 },   // ramp up
    { duration: '2m',  target: 200 },  // sustained load
    { duration: '30s', target: 0 },    // ramp down
  ],
  thresholds: {
    http_req_duration: ['p(95)<200'],
    http_req_failed: ['rate<0.01'],
  },
};

const BASE = 'https://<domain>';
const TOKEN = '<test-jwt-token>';

export default function () {
  // List todos
  const res = http.get(`${BASE}/todos`, {
    headers: { Authorization: `Bearer ${TOKEN}` },
  });
  check(res, { 'todos 200': (r) => r.status === 200 });
  sleep(1);
}
```

## Run Load Test

```bash
k6 run k6-load-test.js
```

## Run Search Performance Test

```bash
# Modify script to hit /todos/search?q=test
# Threshold: p(95)<500
k6 run --env ENDPOINT=/todos/search?q=test k6-load-test.js
```

## Analyze Results

k6 outputs p50, p90, p95, p99 latencies and error rate.

**Pass criteria**:
- `http_req_duration p(95) < 200ms` ✅
- `http_req_failed rate < 1%` ✅

**If failing**: Check Prometheus/Grafana dashboards for bottlenecks (DB pool exhaustion, ES latency, Redis latency).
