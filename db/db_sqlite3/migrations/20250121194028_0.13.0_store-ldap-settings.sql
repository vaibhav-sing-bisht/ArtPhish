
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
-- Create a new table to store ldap records
CREATE TABLE IF NOT EXISTS `ldap`(
    id integer primary key autoincrement,
    user_id bigint,
    name varchar(255),
    protocol varchar(255),
    host varchar(255),
    username varchar(255),
    password varchar(255),
    base_dn varchar(255),
    query TEXT,
    attributes varchar(255),
    modified_date datetime  default CURRENT_TIMESTAMP,
    ignore_cert_errors BOOLEAN
);
-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE `ldap`;

