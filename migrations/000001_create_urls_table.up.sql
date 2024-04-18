CREATE TABLE IF NOT EXISTS urls (
  uid int generated always as identity primary key,
  uuid uuid NOT NULL,
  hash varchar(32) UNIQUE NOT NULL,
  original_url  text NOT NULL
);