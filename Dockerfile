FROM docker.io/golang:1.21 as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN CGO_ENABLED=0 go build -v -o /app/kube-spammer ./cmd/kube-spammer

FROM scratch
WORKDIR /app
COPY --from=builder /app/kube-spammer /app/kube-spammer
CMD ["/app/kube-spammer"]
