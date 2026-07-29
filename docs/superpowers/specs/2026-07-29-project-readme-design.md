# Design Spec: Ecommerce Backend README.md

This specification details the structure, contents, and visual representations that will be implemented in the root `README.md` for the Go Ecommerce-Backend project.

## 1. Goal Description
Create a comprehensive, premium-grade `README.md` at the root of the `Ecommerce-Backend` workspace. This documentation will serve as the onboarding guide and system architecture map for developers. It will detail the technology stack, application architecture (Clean Architecture + Uber Fx), database schema, environment configuration, setup instructions, and the API routes/permissions matrix.

---

## 2. Proposed Outline & Content Specs

The `README.md` will consist of the following main sections:

### Section 1: Header & Technical Stack
* **Project Title:** Ecommerce Backend API
* **Description:** An extensible, high-performance ecommerce REST API built in Go utilizing clean architectural boundaries, automated database migrations, and caching.
* **Tech Stack Badges/Details:**
  * **Language:** Go 1.26+
  * **HTTP Web Framework:** Gin Gonic v1.12
  * **Dependency Injection:** Uber Fx v1.24
  * **Database ORM:** GORM v1.31
  * **Relational Database:** PostgreSQL 15
  * **Cache Store:** Redis 7 (for Cart persistence)
  * **Schema Migrations:** Atlas CLI & Goose Format migrations
  * **API Docs:** Swaggo/Gin-Swagger (OpenAPI 2.0 / Swagger 2.0)
  * **Structured Logging:** Go `log/slog`

### Section 2: Architecture & Directory Structure
* **Application Design Pattern:** Clean Architecture.
  * **Entity:** Standard domain data models (User, Product, Category, Tag, CartItems) independent of external frameworks.
  * **Repository:** Database operations. PostgreSQL/GORM handles persistent relational data; Redis handles cache/ephemeral cart items.
  * **Delivery (HTTP Handlers & Middlewares):** Manages requests, JSON binding, token validation, and authorization check.
* **DI Lifecycle (Uber Fx):** Explain how Fx wires handlers, repositories, and databases automatically.
* **Directory Layout Map:** Detailed folder tree explaining:
  * `/cmd/api` - Main entry point.
  * `/cmd/atlas-loader` - Loader utility for Atlas external schema representation.
  * `/internal/db` - Database initialization and migrations.
  * `/internal/delivery/http` - Handlers and route registration.
  * `/internal/delivery/http/middleware` - Auth, CORS, and logging middlewares.
  * `/internal/domain/entity` - Domain structs.
  * `/internal/domain/repository` - Repository interfaces.
  * `/internal/repository` - GORM/Redis concrete implementations.
  * `/pkg/auth` - Shared JWT utility library.
* **Mermaid Architecture Diagram:**
  ```mermaid
  graph TD
      Client[HTTP Client] -->|Request| Router[Gin Engine Router]
      Router -->|1. Logging & CORS| Middleware[Middlewares]
      Middleware -->|2. Authentication| AuthMiddleware[AuthMiddleware]
      AuthMiddleware -->|3. Role & Ownership Check| PermissionMiddleware[Permission/Ownership Middleware]
      PermissionMiddleware -->|4. Forward| Handler[Http Handlers]
      Handler -->|Query / Save| Repositories[Repositories Interface]
      Repositories -->|GORM postgres| PG[(PostgreSQL)]
      Repositories -->|go-redis client| Redis[(Redis Caching)]
  ```

### Section 3: Database Schema Design
* **ERD Representation (Mermaid):**
  ```mermaid
  erDiagram
      USERS ||--o{ CART_ITEMS : "has items in"
      PRODUCTS ||--o{ CART_ITEMS : "added to"
      CATEGORIES ||--o{ PRODUCTS : "contains"
      CATEGORIES ||--o{ CATEGORIES : "parent of"
      PRODUCTS }o--o{ TAGS : "has (many-to-many)"

      USERS {
          uint64 id PK
          string username UK
          string password_hash
          string role
      }

      CART_ITEMS {
          uint64 id PK
          uint64 user_id FK
          uint64 list_product_id FK
          int quantity
      }

      PRODUCTS {
          uint64 id PK
          string sku UK
          string name
          string description
          numeric price
          int stock
          uint64 category_id FK
          uint64 created_by FK
      }

      CATEGORIES {
          uint64 id PK
          string name UK
          string description
          uint64 parent_id FK
      }

      TAGS {
          uint64 id PK
          string name UK
      }
  ```

### Section 4: Local Development & Setup
* **Prerequisites:** Links/commands to install:
  * Docker & Docker Compose
  * Go (v1.26+)
  * Atlas CLI
* **Step-by-step Setup:**
  1. **Clone & Directory:** Explain how to get into the workspace.
  2. **Infrastructure Launch:** Running the docker-compose stack:
     ```bash
     docker compose up -d
     ```
  3. **Configuration:** Copying/defining environmental variables:
     * Display a markdown table mapping all `.env` parameters (e.g. `JWT_SECRET`, `RUN_MIGRATIONS`, `DB_HOST`, `DB_PORT`, `REDIS_ADDR`, etc.) and their purpose.
  4. **Database Migrations:** Explain Atlas & GORM integration.
     * Document standard commands provided by the `Makefile`:
       * `make migrate-apply`: Run pending migrations.
       * `make migrate-status`: Get migration status.
       * `make migrate-diff name=<migration_name>`: Generate a schema diff migration.
       * `make migrate-hash`: Re-hash migration directory history.
  5. **Run Application:** Running the service:
     ```bash
     go run cmd/api/main.go
     ```

### Section 5: API Documentation & Security Flow
* **Swagger Documentation:**
  * Endpoint: `GET http://localhost:8080/swagger/index.html`
  * Explain how to regenerate docs using:
    ```bash
    swag init -g cmd/api/main.go
    ```
* **Authentication & Authorization Policy:**
  * JWT Bearer Token validation.
  * Middleware explanation:
    * `RequirePermission(roles)`: Enforces role restrictions (`admin`, `seller`, `customer`).
    * `RequireOwnership(roles, ownerResolver)`: Enforces that if the user belongs to a specific role, they must also own the resource (e.g., product created by the seller).
* **Endpoint Permissions Matrix Table:**
  | Method | Endpoint | Description | Auth Required | Allowed Roles | Middleware Checks |
  |--------|----------|-------------|---------------|---------------|-------------------|
  | POST   | `/register` | Register a new user | No | Anyone | None |
  | POST   | `/login` | Log in & obtain token | No | Anyone | None |
  | GET    | `/products` | List all products | Yes | admin, seller, customer | RequirePermission |
  | GET    | `/products/:id` | Get details of a product | Yes | admin, seller, customer | RequirePermission |
  | POST   | `/products` | Create a new product | Yes | admin, seller | RequirePermission |
  | PUT    | `/products/:id` | Update an existing product | Yes | admin, seller | RequirePermission & ownership check |
  | DELETE | `/products/:id` | Delete a product | Yes | admin, seller | RequirePermission & ownership check |
  | POST   | `/cart/add` | Add items to cart | Yes | admin, seller, customer | RequirePermission |
  | GET    | `/cart` | Fetch user's cart | Yes | admin, seller, customer | RequirePermission |
  | POST   | `/product-categories` | Create a category | Yes | admin | RequirePermission |
  | GET    | `/product-categories` | List all categories | Yes | admin, seller, customer | RequirePermission |
  | GET    | `/product-categories/:id` | Get a category by ID | Yes | admin, seller, customer | RequirePermission |
  | PUT    | `/product-categories/:id` | Update a category | Yes | admin | RequirePermission |
  | DELETE | `/product-categories/:id` | Delete a category | Yes | admin | RequirePermission |
  | POST   | `/product-tags` | Create a tag | Yes | admin | RequirePermission |
  | GET    | `/product-tags` | List all tags | Yes | admin, seller, customer | RequirePermission |
  | GET    | `/product-tags/:id` | Get a tag by ID | Yes | admin, seller, customer | RequirePermission |
  | PUT    | `/product-tags/:id` | Update a tag | Yes | admin | RequirePermission |
  | DELETE | `/product-tags/:id` | Delete a tag | Yes | admin | RequirePermission |

---

## 3. Verification Plan
Once the `README.md` is written, we will verify:
* **Syntax & Links:** Ensure all absolute or relative file links correctly resolve to existing project files (e.g. referencing internal files using local project schemas).
* **Mermaid Rendering:** Verify that the Mermaid architecture and database ERD diagrams parse and render correctly without syntax errors.
* **Completeness:** Verify that all setup steps are accurate by comparing with database files (`internal/db/db.go`, `internal/db/redis.go`), Docker configuration (`docker-compose.yml`), and task triggers (`Makefile`).
