# CorePOS

Core Point of Sale System — a multi-tenant SaaS POS backend built with Go.

## Tech Stack

- **Go** + **Gin** (HTTP framework)
- **GORM** (ORM) + **PostgreSQL**
- **Hexagonal Architecture** (Ports & Adapters)

## Project Structure

```
corepos/
├── cmd/api/                # Entry point (main.go)
├── internal/
│   ├── core/
│   │   ├── domain/         # Data models (Product, Order, etc.)
│   │   └── ports/          # Interfaces (Repository & Service)
│   ├── services/           # Business logic (Service implementations)
│   └── adapters/
│       ├── handlers/       # HTTP handlers (Gin)
│       ├── repositories/   # Database layer (GORM/PostgreSQL)
│       └── storage/        # Object storage (MinIO) — future
├── pkg/                    # Shared utilities
├── config/                 # Environment config loader
├── db/                     # SQL schema
├── .env                    # Environment variables
└── go.mod
```

## Getting Started

1. **Configure** — copy `.env` and set your PostgreSQL credentials
2. **Run** — `go run cmd/api/main.go`
3. **API Base URL** — `http://localhost:8080/api/v1/stores/:storeId`

## API Endpoints

| Method   | Path                                        | Description        |
|----------|---------------------------------------------|--------------------|
| `GET`    | `/health`                                   | Health check       |
| `GET`    | `/api/v1/stores/:storeId/products`          | List products      |
| `POST`   | `/api/v1/stores/:storeId/products`          | Create product     |
| `GET`    | `/api/v1/stores/:storeId/products/:id`      | Get product        |
| `PUT`    | `/api/v1/stores/:storeId/products/:id`      | Update product     |
| `DELETE` | `/api/v1/stores/:storeId/products/:id`      | Delete product     |
| `GET`    | `/api/v1/stores/:storeId/categories`        | List categories    |
| `POST`   | `/api/v1/stores/:storeId/categories`        | Create category    |
| `GET`    | `/api/v1/stores/:storeId/categories/:id`    | Get category       |
| `PUT`    | `/api/v1/stores/:storeId/categories/:id`    | Update category    |
| `DELETE` | `/api/v1/stores/:storeId/categories/:id`    | Delete category    |
| `GET`    | `/api/v1/stores/:storeId/orders`            | List orders        |
| `POST`   | `/api/v1/stores/:storeId/orders`            | Create order       |
| `GET`    | `/api/v1/stores/:storeId/orders/:id`        | Get order          |
| `PATCH`  | `/api/v1/stores/:storeId/orders/:id/void`   | Void order         |

## Database Schema

See [`db/schema.sql`](db/schema.sql) for the full PostgreSQL schema.