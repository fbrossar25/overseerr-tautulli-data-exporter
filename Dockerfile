# build app in docker
FROM golang:1.19 AS base

ARG DOCKER_TAG
ENV DOCKER_TAG=$DOCKER_TAG

WORKDIR /app/src

# download dependencies
COPY src/go.mod src/go.sum ./
RUN go mod download

# copy source code
COPY src/* ./

# compile
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/exporter

# run app in docker
FROM golang:1.19 AS production

ARG DOCKER_TAG

ENV CONF_DIR=/app/conf \
    LOG_DIR=/app/logs \
    DOCKER_TAG=$DOCKER_TAG

# TODO when finished
# ENV GIN_MODE=release

WORKDIR /app

COPY --from=base /app .

VOLUME /app/conf /app/logs /app/data
EXPOSE 8090

# run
CMD ["/app/bin/exporter"]