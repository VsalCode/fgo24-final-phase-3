# Flowchart

```mermaid
 flowchart TD
    A((Start)) --> R[/User Register/]
    R --> B[/User Login/]
    B --> C{Valid?}
    C -->|No| R
    C -->|Yes| D[Pilih Operasi]
    D --> E[/Lihat Stok Produk/]
    D --> F[/Transaksi Baru/]
    F --> G{Type?}
    G -->|IN| H[/Input Stok Masuk/]
    G -->|OUT| I[/Input Stok Keluar/]
    H --> J[Update Quantity]
    I --> K{Stok Cukup?}
    K -->|No| L[/Tolak Transaksi/]
    K -->|Yes| J
    J --> M[Catat Transaksi]
    M --> N[Update Stok Produk]
    N --> O[/Tampilkan Laporan/]
    O --> P((end))
    E --> P
```

# ERD

```mermaid
erDiagram
    user {
        id int PK
        name string
        email string UK
        password_hash string
        phone string
        created_at timestamp
        updated_at timestamp
    }

    product_category {
        id int PK
        name string
        description string
        created_at timestamp
        updated_at timestamp
    }

    product {
        id int PK
        code_product string UK
        name string
        image_url string
        purchase_price int
        selling_price int
        quantity int
        user_id int FK
        category_id int FK
        created_at timestamp
        updated_at timestamp
    }

    transaction {
        id int PK
        product_id int FK
        user_id int FK
        type enum "IN, OUT"
        quantity_change int
        created_at timestamp
        updated_at timestamp
    }

    user ||--o{ product : "creates"
    product_category ||--o{ product : "contains"
    user ||--o{ transaction : "performs"
    product ||--o{ transaction : "has"
```