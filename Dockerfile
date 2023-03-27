FROM pa-ctrl:base
COPY cmd/ /app/cmd/
COPY internal/ /app/internal/
WORKDIR /app
RUN go build -o /pa-ctrl cmd/main.go

FROM alpine
COPY --from=0 /pa-ctrl /pa-ctrl

CMD /pa-ctrl
