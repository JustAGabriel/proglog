from golang:1.21-alpine as build
workdir /go/src/proglog
copy . .
run CGO_ENABLED=0 go build -o /go/bin/proglog ./internal/cmd/proglog

from scratch
copy --from=build /go/bin/proglog /bin/proglog
entrypoint ["/bin/proglog"]