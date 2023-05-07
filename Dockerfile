FROM pa-ctrl:base
COPY cmd/pa-ctrl /app/cmd/pa-ctrl
COPY internal/ /app/internal/
WORKDIR /app
RUN go build -o /pa-ctrl cmd/pa-ctrl/main.go

FROM alpine
COPY --from=0 /pa-ctrl /pa-ctrl

CMD /pa-ctrl
