CREATE TABLE customers
(
    id       BIGSERIAL PRIMARY KEY,
    name     TEXT      NOT NULL,
    phone    TEXT      NOT NULL UNIQUE,
    password TEXT      NOT NULL,
    active   BOOLEAN   NOT NULL DEFAULT TRUE,
    created  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE managers
(
    id         BIGSERIAL PRIMARY KEY,
    name       TEXT      NOT NULL,
    phone      TEXT      NOT NULL UNIQUE,
    password   TEXT      NOT NULL,
    salary     INTEGER   NOT NULL check ( managers.salary > 0 ),
    plan       INTEGER   NOT NULL DEFAULT 0 CHECK ( managers.salary > 0 ),
    boss_id    BIGINT REFERENCES managers,
    department TEXT,
    roles      TEXT[]    NOT NULL DEFAULT '{}',
    active     BOOLEAN   NOT NULL DEFAULT TRUE,
    created    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE products
(
    id      BIGSERIAL PRIMARY KEY,
    name    TEXT      NOT NULL,
    price   INTEGER   NOT NULL CHECK ( products.price > 0 ),
    qty     INTEGER   NOT NULL DEFAULT 0 CHECK ( products.qty >= 0 ),
    active  BOOLEAN   NOT NULL DEFAULT TRUE,
    created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE customers_tokens
(
    token       TEXT      NOT NULL UNIQUE,
    customer_id BIGINT    NOT NULL REFERENCES customers,
    expire      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP + INTERVAL '1 hour',
    created     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE managers_tokens
(
    token      TEXT      NOT NULL UNIQUE,
    manager_id BIGINT    NOT NULL REFERENCES managers,
    expire     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP + INTERVAL '1 hour',
    created    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE sales
(
    id          BIGSERIAL PRIMARY KEY,
    manager_id  BIGINT    NOT NULL REFERENCES managers,
    customer_id BIGINT REFERENCES customers,
    created     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE sale_positions
(
    id         BIGSERIAL PRIMARY KEY,
    sale_id    BIGINT    NOT NULL REFERENCES sales,
    product_id BIGINT    NOT NULL REFERENCES products,
    name       TEXT      NOT NULL,
    price      INTEGER   NOT NULL CHECK ( sale_positions.price > 0 ),
    qty        INTEGER   NOT NULL DEFAULT 0 CHECK ( sale_positions.qty >= 0 ),
    created    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);