FROM alpine

WORKDIR /app

COPY --from=build:develop /app/cmd/main /app

CMD ["/app/main"]

