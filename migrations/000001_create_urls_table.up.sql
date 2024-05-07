CREATE TABLE IF NOT EXISTS urls (
  uid int generated always as identity primary key,
  uuid uuid NOT NULL,
  hash varchar(32) NOT NULL,
  original_url  text UNIQUE NOT NULL
);