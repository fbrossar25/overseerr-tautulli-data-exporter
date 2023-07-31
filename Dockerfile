# build app in docker
FROM golang:1.19 AS base

WORKDIR /app

# download dependencies
COPY go.mod go.sum ./
RUN go mod download

# copy source code
COPY . .

# compile
RUN CGO_ENABLED=0 GOOS=linux go build -o /overseerr-tautulli-data-exported

# run app in docker
FROM golang:1.19 AS production

WORKDIR /app

COPY --from=base /app .

EXPOSE 8080

# run
CMD ["/overseerr-tautulli-data-exported"]