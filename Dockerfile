FROM golang:1.11 as build-dev

ENV GO111MODULE=on
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

COPY ./ ./
RUN GOOS=linux GOARCH=386 go build -o goapp

FROM alpine
ENV MONGO_HOST=$MONGO_HOST
ENV MONGO_USER=$MONGO_USER
ENV MONGO_PASS=$MONGO_PASS
ENV PORT=$PORT
COPY --from=build-dev /src/goapp /app/
RUN /app/goapp