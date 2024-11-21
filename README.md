**To start**
- run `mkdir /tmp/logs && docker compose up` (will take a bit of time to start), logs will be in `/tmp/logs`

**To check logs**
- get container ID with `docker ps` and open a shell inside with `docker exec -ti CONTAINER_ID sh`, logs are in `logs` subfolder

**To look at graphs**
- login into grafana at `http://localhost:3000` with creds `admin:admin` and run queries like `http_requests_total`, `http_request_duration_seconds` or `http_errors_total`