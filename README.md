# apiscaleready

go mod tidy
go run .

GLAB_CONCURRENCY=100000 GLAB_NO_KEEP_ALIVE=true GLAB_DURATION=30s GLAB_MODE=wg go run .


GLAB_TARGET_URL=https://httpbin.org/post \
GLAB_METHOD=POST \
GLAB_JSON_BODY='{"test":42}' \
GLAB_CONCURRENCY=800 \
GLAB_RATE=10000 \
GLAB_DURATION=20s \
GLAB_MODE=semaphore \
GLAB_PROXY_FILE=proxies.txt \
go run .

