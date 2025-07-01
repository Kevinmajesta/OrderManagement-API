# Tahap pertama: Pembangunan
FROM golang:latest AS builder

# Set lingkungan kerja
WORKDIR /app

# Mengcopy go.mod dan go.sum agar dependensi dapat di-cache
COPY go.mod go.sum ./

# Mengunduh dependensi
RUN go mod download

# Mengcopy kode sumber ke dalam kontainer
COPY . .

# Membuat binary aplikasi
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/myapp cmd/app/main.go

# Tahap kedua: Produksi
FROM alpine:latest

WORKDIR /app

# Mengcopy binary dari tahap pembangunan
COPY --from=builder /app/myapp .
COPY --from=builder /app/.env .

# Menjalankan aplikasi
EXPOSE 8080

CMD ["./myapp"]
