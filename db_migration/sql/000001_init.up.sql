CREATE TABLE users (
    id bigint PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    email varchar(255) NOT NULL,
    name varchar(255) NOT NULL,
    username varchar(255) UNIQUE NOT NULL,
    password_hash varchar(255) NOT NULL,
    registered_at timestamp with time zone NOT NULL default now(),
    last_visit_at timestamp with time zone NOT NULL default now()
);

CREATE TABLE microgreens_family (
    id bigint PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    name varchar(255) UNIQUE NOT NULL,
    description text
);

CREATE TABLE microgreens_list (
    id bigint PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    name varchar(255) NOT NULL,
    description varchar(255)
);

CREATE TABLE microgreens_item (
    id bigint PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    name varchar(255) NOT NULL,
    description text,
    price numeric CHECK (price > 0) NOT NULL,
    microgreens_family_id bigint REFERENCES microgreens_family(id) ON DELETE CASCADE NOT NULL
);

CREATE TABLE microgreens_list_items (
    id bigint PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    microgreens_list_id bigint REFERENCES microgreens_list(id) ON DELETE CASCADE NOT NULL,
    microgreens_item_id bigint REFERENCES microgreens_item(id) ON DELETE CASCADE NOT NULL
);

CREATE TABLE users_microgreens_lists (
    id bigint PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    user_id bigint REFERENCES users(id) ON DELETE CASCADE NOT NULL,
    microgreens_list_id bigint REFERENCES microgreens_list(id) ON DELETE CASCADE NOT NULL
);

CREATE TABLE microgreens_family_items (
    id bigint PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    microgreens_family_id bigint REFERENCES microgreens_family(id) ON DELETE CASCADE NOT NULL,
    microgreens_item_id bigint REFERENCES microgreens_item(id) ON DELETE CASCADE NOT NULL
);

CREATE TABLE refresh_sessions (
    id bigint PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    user_id bigint REFERENCES users(id) ON DELETE CASCADE NOT NULL,
    refresh_token varchar(64) UNIQUE NOT NULL,
    expires_at timestamp with time zone NOT NULL,
    created_at timestamp with time zone NOT NULL default now()
);

