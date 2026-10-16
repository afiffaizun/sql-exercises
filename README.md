# Aplikasi CRUD Data Produk

Aplikasi CRUD sederhana untuk mengelola data produk/barang dengan tampilan modern.

**Stack:** Go + Gin (backend) · Vue 3 + Vite + Tailwind (frontend) · MariaDB (database) · dbx (DB management UI) · Docker Compose.

## Arsitektur

```
Browser ──> Vue 3 (Nginx :3000) ──/api──> Go + Gin (:8080) ──> MariaDB (:3306)
                                                                    ▲
                                              dbx UI (:4224) ───────┘
```

## Menjalankan

```bash
docker compose up -d --build
```

- Aplikasi  : http://localhost:3000
- dbx (DB UI): http://localhost:4224

## Menghubungkan dbx ke database

dbx berjalan di dalam container, jadi **isi Host dengan nama service `mysql`**,
bukan `localhost` (localhost di dalam container dbx menunjuk ke dirinya sendiri
dan menyebabkan `Connection refused (os error 111)`).

1. Buka http://localhost:4224
2. Tambah koneksi baru:
   - Host: `mysql` (nama service Docker; jangan pakai `localhost`)
   - Port: `3306`
   - User: `root`
   - Password: `root`
   - Database: `crud_produk`
3. Jelajahi / edit tabel `produks` langsung dari dbx.

Alternatif dari sisi host: Host `host.docker.internal` Port `3306`
(service `dbx` sudah diberi `extra_hosts: host-gateway`).

## Endpoint API

| Method | Endpoint            | Keterangan                          |
| ------ | ------------------- | ----------------------------------- |
| GET    | `/api/produk`       | Daftar produk (`?q=`, `?kategori=`) |
| GET    | `/api/produk/:id`   | Detail produk                       |
| POST   | `/api/produk`       | Tambah produk                       |
| PUT    | `/api/produk/:id`   | Ubah produk                         |
| DELETE | `/api/produk/:id`   | Hapus produk                        |
| GET    | `/api/stats`        | Statistik dashboard                 |
| GET    | `/api/health`       | Health check                        |

## Struktur

```
.
├── backend/          # Go + Gin + GORM
├── frontend/         # Vue 3 + Vite + Tailwind
├── docker-compose.yml
└── docker-compose.release.yml   # contoh deploy dbx standalone
```

## Development tanpa Docker

Backend (butuh MariaDB lokal):
```bash
cd backend
DB_HOST=localhost DB_USER=root DB_PASSWORD=root DB_NAME=crud_produk go run .
```

Frontend:
```bash
cd frontend
npm install
npm run dev   # proxy /api ke http://localhost:8080
```

## Menghentikan

```bash
docker compose down          # stop container
docker compose down -v       # stop + hapus data
```
