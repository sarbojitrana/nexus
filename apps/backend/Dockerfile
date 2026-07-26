FROM golang:1.25-alpine AS build
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /out/nexus ./cmd/nexus

FROM alpine:3.20
RUN apk add --no-cache ca-certificates
WORKDIR /app

COPY --from=build /out/nexus ./nexus
COPY --from=build /src/static ./static
COPY --from=build /src/templates ./templates

ENTRYPOINT ["./nexus"]
