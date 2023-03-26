FROM golang:alpine
COPY cmd/ /app/cmd/
COPY internal/ /app/internal/
COPY go.mod /app/
COPY go.sum /app/
WORKDIR /app
RUN go build -o pa-ctrl cmd/main.go

FROM alpine
COPY --from=0 /app/pa-ctrl /pa-ctrl
CMD /pa-ctrl
