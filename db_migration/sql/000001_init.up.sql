CREATE TABLE users (
    id int PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    email varchar(255) NOT NULL,
    name varchar(255) NOT NULL,
    username varchar(255) UNIQUE NOT NULL,
    password_hash varchar(255) NOT NULL
);

CREATE TABLE microgreens_family (
    id int PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    name varchar(255) UNIQUE NOT NULL,
    description varchar(255)
);

CREATE TABLE microgreens_list (
    id int PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    name varchar(255) NOT NULL,
    description varchar(255)
);

CREATE TABLE microgreens_item (
    id int PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    name varchar(255) NOT NULL,
    description varchar(255),
    price numeric NOT NULL,
    microgreens_family_id int REFERENCES microgreens_family(id) ON DELETE CASCADE NOT NULL
);

CREATE TABLE microgreens_list_items (
    id int PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    microgreens_list_id int REFERENCES microgreens_list(id) ON DELETE CASCADE NOT NULL,
    microgreens_item_id int REFERENCES microgreens_item(id) ON DELETE CASCADE NOT NULL
);

CREATE TABLE users_microgreens_lists (
    id int PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    user_id int REFERENCES users(id) ON DELETE CASCADE NOT NULL,
    microgreens_list_id int REFERENCES microgreens_list(id) ON DELETE CASCADE NOT NULL
);

CREATE TABLE microgreens_family_items (
    id int PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    microgreens_family_id int REFERENCES microgreens_family(id) ON DELETE CASCADE NOT NULL,
    microgreens_item_id int REFERENCES microgreens_item(id) ON DELETE CASCADE NOT NULL
);

