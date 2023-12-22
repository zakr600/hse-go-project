FROM alpine

WORKDIR /app

COPY --from=build:develop /app/cmd/main /app
COPY --from=build:develop /app/configs /app

CMD ["/app/main"]

