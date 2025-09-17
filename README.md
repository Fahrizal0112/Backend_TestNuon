# Transaction API

Transaction API adalah REST API yang dibangun dengan Go (Golang) menggunakan framework Gin untuk mengelola data transaksi. API ini mendukung operasi CRUD, upload CSV, dan pencarian data transaksi dengan berbagai filter.

## 🚀 Fitur

- ✅ CRUD Operations untuk transaksi
- 📁 Upload data transaksi melalui file CSV
- 🔍 Pencarian dan filtering data transaksi
- 📄 Pagination untuk data yang besar
- 🗄️ Database PostgreSQL dengan GORM ORM
- 🌐 CORS support untuk frontend integration
- 📊 Auto migration database

## 🛠️ Tech Stack

- **Language**: Go 1.24.1
- **Framework**: Gin Web Framework
- **Database**: PostgreSQL
- **ORM**: GORM
- **Environment**: godotenv

## 📋 Prerequisites

- Go 1.24.1 atau lebih tinggi
- PostgreSQL
- Git

## ⚙️ Installation

1. **Clone repository**
   ```bash
   git clone <repository-url>
   cd transaction-api
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Setup environment variables**
   
   Buat file `.env` di root directory:
   ```env
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=postgres
   DB_PASSWORD=your_password
   DB_NAME=transaction_db
   DB_SSLMODE=disable
   PORT=8080
   ```

4. **Setup PostgreSQL Database**
   ```sql
   CREATE DATABASE transaction_db;
   ```

5. **Run the application**
   ```bash
   go run main.go
   ```

   Server akan berjalan di `http://localhost:8080`

## 📊 Database Schema

### Transaction Table

| Field | Type | Description |
|-------|------|-------------|
| id | uint | Primary key (auto increment) |
| msisdn | string | Nomor telepon (indexed) |
| trx_id | string | Transaction ID (unique) |
| trx_date | timestamp | Tanggal transaksi (indexed) |
| item | string | Nama item/produk |
| voucher_code | string | Kode voucher |
| status | int | Status transaksi (0/1, indexed) |
| created_at | timestamp | Waktu dibuat |
| updated_at | timestamp | Waktu diupdate |
| deleted_at | timestamp | Soft delete timestamp |

## 🔌 API Endpoints

### Base URL
```
http://localhost:8080/api/v1
```

### Endpoints

#### 1. Get All Transactions
```http
GET /transactions
```

**Query Parameters:**
- `page` (int): Halaman (default: 1)
- `limit` (int): Jumlah data per halaman (default: 1000)
- `search` (string): Pencarian global di semua field
- `msisdn` (string): Filter berdasarkan nomor telepon
- `status` (int): Filter berdasarkan status (0 atau 1)
- `item` (string): Filter berdasarkan item
- `start_date` (string): Filter tanggal mulai (format: YYYY-MM-DD)
- `end_date` (string): Filter tanggal akhir (format: YYYY-MM-DD)

**Example:**
```bash
GET /api/v1/transactions?page=1&limit=10&search=ISAT&status=1
```

#### 2. Get Transaction by ID
```http
GET /transactions/:id
```

#### 3. Create Transaction
```http
POST /transactions
```

**Request Body:**
```json
{
  "msisdn": "62856034xxxxx",
  "trx_id": "unique_transaction_id",
  "trx_date": "2025-09-05T00:51:20Z",
  "item": "GPAY_TI_127",
  "voucher_code": "ZZQEHY4EDCFQ59M0",
  "status": 1
}
```

#### 4. Upload CSV
```http
POST /transactions/upload
```

**Form Data:**
- `csv_file`: File CSV dengan format yang sesuai

**CSV Format:**
```csv
msisdn,trx_id,trx_date,item,voucher_code,status
62856034xxxxx,d3d95b2a094e34f6911266ccfdf6d126,2025-09-05 00:51:20,GPAY_TI_127,ZZQEHY4EDCFQ59M0,0
```

#### 5. Clear All Transactions
```http
DELETE /transactions/clear
```

## 📝 Response Format

### Success Response
```json
{
  "data": [...],
  "pagination": {
    "page": 1,
    "limit": 10,
    "total": 100
  }
}
```

### Error Response
```json
{
  "error": "Error message"
}
```

## 🧪 Testing

### Using cURL

1. **Get all transactions:**
   ```bash
   curl -X GET "http://localhost:8080/api/v1/transactions?page=1&limit=10"
   ```

2. **Search transactions:**
   ```bash
   curl -X GET "http://localhost:8080/api/v1/transactions?search=ISAT&status=1"
   ```

3. **Upload CSV:**
   ```bash
   curl -X POST \
     http://localhost:8080/api/v1/transactions/upload \
     -H "Content-Type: multipart/form-data" \
     -F "csv_file=@data.csv"
   ```

### Using Postman

1. Import collection dengan base URL: `http://localhost:8080/api/v1`
2. Test semua endpoints sesuai dokumentasi di atas

## 📁 Project Structure

```
.
├── .env                    # Environment variables
├── main.go                 # Application entry point
├── go.mod                  # Go modules
├── go.sum                  # Go modules checksum
├── data.csv               # Sample CSV data
├── config/
│   └── database.go        # Database configuration
├── controllers/
│   └── transaction_controller.go  # API controllers
├── models/
│   └── transaction.go     # Database models
├── routes/
│   └── routes.go         # API routes
└── utils/
    └── csv_loader.go     # CSV utility functions
```

## 🔧 Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|----------|
| DB_HOST | Database host | localhost |
| DB_PORT | Database port | 5432 |
| DB_USER | Database username | postgres |
| DB_PASSWORD | Database password | - |
| DB_NAME | Database name | transaction_db |
| DB_SSLMODE | SSL mode | disable |
| PORT | Server port | 8080 |

### CORS Configuration

API ini dikonfigurasi untuk menerima request dari:
- `http://localhost:3000`
- `http://127.0.0.1:3000`

## 🚨 Error Handling

API mengembalikan HTTP status code yang sesuai:

- `200` - Success
- `201` - Created
- `206` - Partial Content (CSV upload dengan error)
- `400` - Bad Request
- `404` - Not Found
- `500` - Internal Server Error

## 📈 Performance

- **Database Indexing**: Field `msisdn`, `trx_date`, dan `status` diindex untuk performa query yang optimal
- **Pagination**: Mendukung pagination untuk menangani dataset besar
- **Connection Pooling**: Menggunakan GORM connection pooling
```