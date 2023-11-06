FROM golang:1.21.3 as build

WORKDIR /app

COPY . .

RUN go mod download
RUN go mod tidy
# RUN cd /app/cmd && go build -o go-auth

# FROM golang:1.21.3

# WORKDIR /app

# COPY --from=build /app/cmd/go-auth /app/

CMD go run cmd/main.go

# CMD ./go_auth
# ENTRYPOINT ["./app/go-auth"]