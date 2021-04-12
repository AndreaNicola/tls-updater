FROM golang:1.15 as build
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    go build \
    -ldflags="-w -s" \
    -o app .

FROM scratch as bin
COPY --from=build /app/app .
CMD ["./app"]
EXPOSE 8080
USER 1001