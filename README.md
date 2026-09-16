#Project masih on progress

# Project Overview
Aplikasi ini dirancang untuk mengelola proses inventory dan procurement, yang mencakup pengelolaan produk, supplier, warehouse, Purchase Request (PR), Purchase Order (PO), 
serta proses penerimaan barang (Goods Receipt).


## 🛠️ Tech Stack

| Technology                 | Description                                                                          |
| -------------------------- | ------------------------------------------------------------------------------------ |
| 🐹 **Golang**              | Bahasa pemrograman utama untuk backend                                               |
| 🌐 **Gin**                 | Framework untuk membangun RESTful API                                                |
| 🐘 **PostgreSQL**          | Relational database untuk menyimpan data aplikasi                                    |
| 🔗 **GORM**                | ORM untuk mengelola interaksi dengan database                                        |
| 🔐 **JWT**                 | Authentication dan authorization berbasis token                                      |
| 🏗️ **Clean Architecture** | Arsitektur untuk memisahkan business logic, use case, repository, dan delivery layer |
| 📡 **REST API**            | Interface komunikasi antara client dan backend                                       |
| 🧪 **Go Testing**          | Automated testing untuk memastikan business logic berjalan sesuai kebutuhan          |
| 📬 **Postman**             | Pengujian dan validasi REST API                                                      |
| 🐙 **Git & GitHub**        | Version control dan repository management                                            |

### 🏛️ Architecture

Project ini menerapkan **Clean Architecture** dengan struktur layer sebagai berikut:

```text
Clean Architecture Golang:
1. Entity
    2. Repository
        3. Model
            4. mysql
                5. usecase
                    6. handler
                        7. router
                            8. main
Struktur layer:

1. **Entity** — Mendefinisikan struktur dan aturan dasar data/domain.
2. **Repository** — Mendefinisikan interface untuk akses dan pengelolaan data.
3. **Model** — Merepresentasikan struktur data yang digunakan untuk database.
4. **PostgreSQL** — Mengimplementasikan akses database melalui repository.
5. **Usecase** — Menangani business logic dan alur proses aplikasi.
6. **Handler** — Menangani request dan response HTTP/API.
7. **Router** — Mendefinisikan endpoint dan menghubungkan request ke handler.
8. **Main** — Melakukan inisialisasi aplikasi, dependency, database, router, dan menjalankan server.

Alur aplikasi:
Entity → Repository → Model → Database → Usecase → Handler → Router → Main
```

## Database Design

```text
1. Master Data
-------------------------------------------------------------------------
**users**
Menyimpan data pengguna sistem.
Role yang tersedia:

USER
APPROVER
User digunakan untuk mencatat siapa yang membuat purchase request dan siapa yang menerima barang.

**product**
Menyimpan data barang yang tersedia.

Kolom penting:

- id: primary key
- sku: kode unik barang
- name: nama barang
- unit: satuan barang
- is_active: status aktif barang
- created_at dan updated_at: informasi waktu data dibuat dan diubah
SKU dibuat UNIQUE agar satu barang tidak memiliki kode yang sama.

**supplier**
Menyimpan data pemasok barang.

Relasinya digunakan oleh tabel purchase_orders melalui supplier_id.

**warehouse**
Menyimpan data gudang yang ada.

Relasinya digunakan oleh:

- inventory
- purchase_requests
- purchase_orders
- goods_receipts
- inventory_movements
-------------------------------------------------------------------------

2. Purchase Request
**purchase_requests** menyimpan permintaan barang sebelum purchase order.
Status yang diperbolehkan:
- DRAFT
- SUBMITTED
- APPROVED
- REJECTED

Setiap purchase request memiliki banyak barang melalui tabel purchase_request_items.

**purchase_request_items**
Tabel ini menjadi detail barang yang diminta.
Relasi :
- satu purchase request memiliki banyak item
- satu product dapat muncul di banyak purchase request

**quantity** harus lebih besar dari nol.
-------------------------------------------------------------------------

3. Purchase Order
**purchase_orders** menyimpan pesanan pembelian yang dibuat berdasarkan purchase request.
relasi :
- purchase_request_id: request sumber
  purchase_request_id diberi constraint UNIQUE, sehingga satu purchase request hanya dapat memiliki satu purchase order.
- supplier_id: supplier yang dipilih
- warehouse_id: gudang tujuan

Status purchase order:
- DRAFT
- ORDERED
- PARTIALLY_RECEIVED
- RECEIVED
- CANCELLED

**purchase_order_items**
Menyimpan detail barang dalam purchase order.

memiliki informasi kuantitas:
- ordered_quantity: jumlah yang dipesan
- received_quantity: jumlah yang sudah diterima

Constraint:
- ordered_quantity > 0
- received_quantity >= 0
- kombinasi purchase_order_id dan product_id harus unik, sehinga satu produk tidak boleh muncul dua kali dalam purchase order yang sama.
-------------------------------------------------------------------------

4. Goods Receipt
**goods_receipts** mencatat penerimaan barang dari supplier.

Status goods receipt:
- POSTED
- CANCELLED

**goods_receipt_items**
Menyimpan detail barang yang diterima.

Relasi:
- goods receipt
- purchase order item
- product

Constraint:
- UNIQUE (goods_receipt_id, product_id) mencegah produk yang sama dicatat dua kali dalam satu penerimaan.
-------------------------------------------------------------------------

5. Inventory
**inventory** menyimpan stok barang berdasarkan kombinasi barang dan gudang.

Relasi:
- product_id
- warehouse_id

Constraint:
- UNIQUE (product_id, warehouse_id), sehingga satu produk hanya memiliki satu record stok untuk setiap gudang.
-------------------------------------------------------------------------

6. Inventory Movement
**inventory_movements** berfungsi sebagai histori perubahan stok.

- Informasi yang dicatat:
- gudang
- produk
- jenis pergerakan
- jumlah
- referensi transaksi
- waktu perubahan

 jenis movement dilakukan melalui PURCHASE_RECEIPT


```
<img width="775" height="795" alt="image" src="https://github.com/user-attachments/assets/03e1d564-8100-4c51-bdac-79984230a906" />

