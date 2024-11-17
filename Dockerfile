FROM golang:1.23 as build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY ./ ./
RUN go build -o rxm-server ./cmd/backend

FROM gcr.io/distroless/base-debian12 AS run
WORKDIR /
COPY --from=build /app/rxm-server /rxm-server
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/rxm-server"]