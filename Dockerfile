FROM golang:1.26.4 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . . 
# ENV GOOS=linux // we are already using a linux builder
ENV GOARCH=amd64

RUN CGO_ENABLED=0 go build -o wacdo-backend .

# Stage 2
FROM alpine:latest

COPY --from=builder /app/wacdo-backend /wacdo-backend

CMD ["/wacdo-backend"]

EXPOSE 8080
