SET NAMES utf8;

--DROP TABLE IF EXISTS users;
--DROP TABLE IF EXISTS ads;
CREATE TABLE users
(
    id                SERIAL      PRIMARY KEY,
    first_name        varchar(80) NOT NULL,
    last_name         varchar(80) NOT NULL,
    email             varchar(80) UNIQUE NOT NULL,
    password_hash     text        NOT NULL,
    telephone         varchar(80),
    about             text,
    reg_time timestamp   DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TABLE ads
(
    id             SERIAL       PRIMARY KEY,
    title          varchar(80)  NOT NULL,
    price          integer      CONSTRAINT positive_price CHECK (price > 0),
    country        varchar(80),
    city           varchar(80),
    subway_station varchar(80),
    images_folder  varchar(256),
    -- when deleting user we should delete his ads
    owner_ad       integer      REFERENCES users (id) ON DELETE CASCADE NOT NULL,
    description_ad text,
    -- do we need this tocken if there is id already?
    -- tocken         varchar(16) UNIQUE NOT NULL,
    creation_time  timestamp    DEFAULT CURRENT_TIMESTAMP NOT NULL
);

