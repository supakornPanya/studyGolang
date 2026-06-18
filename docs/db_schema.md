# Database Schema Design

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
        boolean can_read
        boolean can_write
        boolean can_update
        boolean can_delete
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
        unit64 tag
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
