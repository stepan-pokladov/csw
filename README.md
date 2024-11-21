**To start**
- run `mkdir /tmp/logs && docker compose up` (will take a bit of time to start)

**To check logs**
- logs will be in `/tmp/logs`

**To look at graphs**
- login into grafana at `http://localhost:3000` with creds `admin:admin` and run queries like `http_requests_total`, `http_request_duration_seconds` or `http_errors_total`