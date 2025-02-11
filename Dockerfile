FROM golang:1.23-alpine AS builder
COPY ./ /src/github.com/adroll/ecs-ship/
RUN cd  /src/github.com/adroll/ecs-ship/ \
    && go test ./... \
    && CGO_ENABLED=0 go build -a -ldflags '-extldflags "-static"' .

FROM alpine:3
COPY --from=builder /src/github.com/adroll/ecs-ship/ecs-ship /usr/bin/
CMD [ "ecs-ship" ]
