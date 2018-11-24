FROM golang:1.11 as build-dev

ENV GO111MODULE=on
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

COPY ./ ./
RUN GOOS=linux GOARCH=386 go build -o goapp

FROM alpine
ENV PORT=$PORT
COPY --from=build-dev /src/goapp /app/
CMD /app/goapp