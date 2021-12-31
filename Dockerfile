FROM golang:1.17-bullseye as build-env
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /fakedns cmd/*.go
FROM scratch
EXPOSE 53
COPY --from=build-env /fakedns /fakedns
ENTRYPOINT ["/fakedns"]