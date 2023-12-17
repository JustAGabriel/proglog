from golang:1.21-alpine as build
workdir /go/src/proglog
copy . .
run CGO_ENABLED=0 go build -o /go/bin/proglog ./internal/cmd/proglog
run GRPC_HEALTH_PROBE_VERSION=v0.4.13 && \
    wget -qO/go/bin/grpc_health_probe \
    https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 && \
    chmod +x /go/bin/grpc_health_probe

# todo: config dir with all the certs etc. was missing during runtime.
# embed into final binary?
from scratch

arg PROGLOG_CONFIG=/etc/proglog

copy --from=build /go/src/proglog/test ${PROGLOG_CONFIG} 
copy --from=build /go/bin/proglog /bin/proglog
copy --from=build /go/bin/grpc_health_probe /bin/grpc_health_probe

# define location of TLS certs etc.
env CONFIG_DIR=${PROGLOG_CONFIG}

entrypoint ["/bin/proglog"]