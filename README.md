# be_evindo

Contoh backend Go dengan Clean Architecture.

## Struktur

```text
cmd/api/main.go              # composition root: merakit dependency
internal/entity/user/         # entity dan contract usecase/repository
internal/repository/         # contract repository tingkat aplikasi
internal/model/              # bentuk data persistence
internal/postgres/           # implementasi repository dengan PostgreSQL
internal/usecase/            # aturan bisnis
internal/handler/            # HTTP request/response
internal/router/             # route HTTP
migrations/                  # schema database
```

Arah dependency utama:

```text
handler -> usecase -> entity.Repository <- mysql
router  -> handler
main    -> semua implementasi konkret
```

Entity dan usecase tidak bergantung pada PostgreSQL atau HTTP. `main` bertugas melakukan dependency injection.

## Menjalankan

1. Jalankan migration pada database PostgreSQL.
2. Set environment variable `POSTGRES_PASSWORD` sesuai password database Anda, atau set `POSTGRES_DSN` langsung:

```text
POSTGRES_DSN=postgres://postgres:password@localhost:5432/be_evindo?sslmode=disable
```

3. Jalankan server:

```bash
go run ./cmd/api
```

Endpoint awal:

- `GET /health`
- `GET /users/{id}`
- `POST /users` dengan body `{"name":"Budi","email":"budi@example.com","password":"secret"}`
- `POST /products/` dengan body `{"sku":"SKU-001","name":"Produk 1","unit":"pcs","is_active":1}`
- `GET /products/{id}`

Contoh pengujian di Postman:

```text
POST http://localhost:8080/products/
Content-Type: application/json
```

```json
{
	"sku": "SKU-001",
	"name": "Produk 1",
	"unit": "pcs",
	"is_active": 1
}
```
