erDiagram
    USERS ||--o{ BOOKINGS : "membuat (1:N)"
    ROOMS ||--o{ BOOKINGS : "dipesan pada (1:N)"

    USERS {
        bigint id PK "AUTO_INCREMENT"
        varchar name "NOT NULL"
        varchar email UK "UNIQUE, NOT NULL"
        varchar password_hash "bcrypt, NOT NULL"
        enum role "USER | ADMIN"
        datetime created_at
    }

    ROOMS {
        bigint id PK "AUTO_INCREMENT"
        varchar name "NOT NULL"
        varchar code UK "UNIQUE, NOT NULL"
        text description
        int capacity "NOT NULL"
        varchar location "NOT NULL"
        enum status "ACTIVE | INACTIVE"
        datetime created_at
        datetime updated_at
    }

    BOOKINGS {
        bigint id PK "AUTO_INCREMENT"
        bigint user_id FK "ON DELETE RESTRICT"
        bigint room_id FK "ON DELETE RESTRICT"
        varchar title "NOT NULL"
        text purpose "NOT NULL"
        date booking_date "NOT NULL"
        time start_time "NOT NULL"
        time end_time "NOT NULL"
        enum status "CONFIRMED | CANCELLED"
        datetime created_at
        datetime updated_at
    }