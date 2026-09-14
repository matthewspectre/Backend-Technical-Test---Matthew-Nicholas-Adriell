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
```




