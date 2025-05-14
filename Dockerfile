# menggunakan image resmi Go untuk build
FROM golang:1.24-alpine

# membuat folder kerja di dalam container
WORKDIR /app

COPY go.mod ./
# download dependency Go
RUN go mod download

# copy semua file ke dalam container
COPY . .

# build aplikasi
RUN go build -o awb-scraper-app

EXPOSE 8080

# set default command saat container dijalankan
CMD ["./awb-scraper-app"]