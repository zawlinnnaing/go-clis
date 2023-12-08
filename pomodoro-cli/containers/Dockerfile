FROM golang:1.20 as builder
RUN mkdir /pomodro-cli
WORKDIR /pomodoro-cli
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -tags=containers
RUN ls


FROM scratch
WORKDIR /
COPY --from=builder /pomodoro-cli/pomodoro-cli .



ENTRYPOINT ["/pomodoro-cli"]