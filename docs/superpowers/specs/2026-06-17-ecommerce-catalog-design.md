# E-Commerce Database & API Design Specification

**Date**: 2026-06-17  
**Status**: Approved

---

## 1. Goal Description
To design and implement a relational database schema for an e-commerce catalog featuring hierarchical categories (subcategories), product tags, and a direct shopping cart model.

---

## 2. Database Schema Design
We are implementing **Approach 2 (Hierarchical Catalog with Many-to-Many Tags)** with a **Direct Shopping Cart** association (linking cart items directly to users without an intermediate carts table for simplicity and performance).

### Schema ERD Diagram
```mermaid
erDiagram
    USERS ||--o{ CART_ITEMS : "has items in"
    PRODUCTS ||--o{ CART_ITEMS : "added to"
    CATEGORIES ||--o{ PRODUCTS : "contains"
    CATEGORIES ||--o{ CATEGORIES : "parent of"
    PRODUCTS }o--o{ TAGS : "has"

    USERS {
        uint64 id PK
        string username UK
        string password
        boolean is_admin
        boolean is_seller
        boolean is_customer
    }

    CART_ITEMS {
        uint64 id PK
        uint64 user_id FK
        uint64 product_id FK
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

---

## 3. Key Architecture Decisions

### A. Hierarchical Categories
The `categories` table utilizes a self-referencing foreign key (`parent_id`) pointing to `categories.id`. This enables nested categories (e.g., `Apparel -> Men's -> Shoes`).

### B. Direct Cart Items
Since each user has exactly 1 active cart, `cart_items` connects directly to the `user_id`, reducing table joins and storage overhead.

### C. Permissions (User Roles)
Authentication permissions are migrated to a role-based access model (`IsAdmin`, `IsSeller`, `IsCustomer`) validated by the HTTP JWT middleware.

---

## 4. Verification Plan
- Compiles successfully with Uber Fx injection.
- Endpoint validation:
  - `POST /products/` (Requires Admin/Seller role)
  - `GET /products/` (Accessible by all verified roles)
  - `PUT /products/:id` (Requires Admin/Seller role)
  - `DELETE /products/:id` (Requires Admin/Seller role)
