# Aplikasi Scraping AWB - Golang

Aplikasi Scraping AWB ini dibuat dengan menggunakan bahasa pemrograman Go (Golang) untuk melakukan scraping informasi pelacakan (tracking) dari nomor AWB dan mengembalikan data dalam format JSON. Aplikasi dapat dijalankan dengan Docker untuk kemudahan deployment.

---

## Fitur

- Scraping halaman HTML berdasarkan nomor AWB
- Menampilkan riwayat pelacakan pengiriman dalam format JSON
- Penanganan error untuk koneksi jaringan dan nomor tidak valid
- Aplikasi ringan dan siap dijalankan menggunakan Docker

---

## Requirement

- Go version 1.24.3
- Docker

---

## Cara Menjalankan Aplikasi

### 1. Clone Repository

```bash
git clone https://github.com/ahmadalkhaidar/golang-web-scraping.git
cd golang-web-scraping
```
### 2. Build And Run Docker Image

```bash
docker build -t golang-web-scraping .
docker run -p 8080:8080 golang-web-scraping
```

## Endpoint
`GET /awb/{awb_number}`

## Aplikasi berjalan di http://localhost:8080
contoh: http://localhost:8080/awb/325244024003
