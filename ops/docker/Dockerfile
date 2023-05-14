FROM golang:alpine as base

ARG SSH_PRIVATE_KEY

# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git openssh-client

RUN mkdir /root/.ssh/
RUN touch /root/.ssh/id_rsa
RUN chmod 600 /root/.ssh/*
RUN echo "${SSH_PRIVATE_KEY}" > /root/.ssh/id_rsa

RUN echo "StrictHostKeyChecking no " > /root/.ssh/config
RUN touch /root/.ssh/known_hosts
RUN ssh-keyscan  git.casinomodule.org >> /root/.ssh/known_hosts
RUN git config --global url."ssh://git@git.casinomodule.org:10022".insteadOf https://git.casinomodule.org

ENV GOPRIVATE="git.casinomodule.org"
ENV CGO_ENABLED=0

WORKDIR /app
# Set the current working directory inside the container

FROM base as application-dev

RUN apk --no-cache add entr
RUN apk --no-cache add npm

RUN npm install -g nodemon

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and the go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the working Directory inside the container
COPY . .

RUN echo " \
    (while true; do nodemon --exec go run cmd/todo --signal SIGTERM || sleep 10; done) \
    " > /entrypoint.sh \
    && chmod +x /entrypoint.sh

ENTRYPOINT /entrypoint.sh

FROM base as build

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and the go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the working Directory inside the container
COPY . .
# Build the Go app
RUN go build -o /bin/todo ./cmd

# Start a new stage from scratch
FROM alpine:latest as application
RUN apk --no-cache add ca-certificates

COPY --from=build /bin/todo /bin/todo
COPY --from=build /app/ssl /ssl

ENTRYPOINT ["/bin/todo"]
