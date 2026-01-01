-- table roles
CREATE TABLE roles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL, --super admin, admin, staff
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- table users
CREATE TABLE users (
    id UUID PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    role_id INT NOT NULL REFERENCES roles(id),
    is_active BOOLEAN DEFAULT TRUE,
    deleted_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);


-- table sessions
CREATE TABLE sessions (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id),
    expires_at TIMESTAMPTZ NOT NULL,
    revoked_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_sessions_user_id ON sessions(user_id); --get all user sessions (ex: to revoke all)
CREATE INDEX idx_sessions_expires_at ON sessions(expires_at); --use to clean all revoked sessions


-- table warehouses
CREATE TABLE warehouses (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    location TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- table racks
CREATE TABLE racks (
    id SERIAL PRIMARY KEY,
    warehouse_id INT NOT NULL REFERENCES warehouses(id),
    code VARCHAR(50) NOT NULL,
    description TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE (warehouse_id, code) --must diferent per warehouse_id
);

-- table categories
CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);


-- table items
CREATE TABLE items (
    id UUID PRIMARY KEY,
    name VARCHAR(150) NOT NULL,
    category_id INT NOT NULL REFERENCES categories(id),
    rack_id INT NOT NULL REFERENCES racks(id),
    stock INT NOT NULL DEFAULT 1 CHECK (stock >= 0),
    min_stock INT NOT NULL DEFAULT 1 CHECK (min_stock >= 0),
    price NUMERIC(12,2) NOT NULL CHECK (price > 0),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_items_stock ON items(stock);
CREATE INDEX idx_items_category_id ON items(category_id);
CREATE INDEX idx_items_rack_id ON items(rack_id);
CREATE INDEX idx_items_deleted_at ON items(deleted_at);


-- table sales
CREATE TABLE sales (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id),
    total_amount NUMERIC(14,2) NOT NULL CHECK (total_amount > 0),
    status VARCHAR(20) NOT NULL DEFAULT 'COMPLETED',
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_sales_created_at ON sales(created_at);
CREATE INDEX idx_sales_user_id ON sales(user_id);


-- table stock ajustments
CREATE TABLE stock_adjustments (
    id UUID PRIMARY KEY,
    item_id UUID NOT NULL REFERENCES items(id),
    user_id UUID NOT NULL REFERENCES users(id),
    change INT NOT NULL CHECK (change <> 0),
    reason TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_stock_adjustments_item_id ON stock_adjustments(item_id);
CREATE INDEX idx_stock_adjustments_user_id ON stock_adjustments(user_id);
CREATE INDEX idx_stock_adjustments_created_at ON stock_adjustments(created_at);
