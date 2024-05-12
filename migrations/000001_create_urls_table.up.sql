CREATE TABLE IF NOT EXISTS urls (
  uid int generated always as identity primary key,
  creator_id uuid,
  uuid uuid NOT NULL,
  short_url text NOT NULL,
  original_url  text UNIQUE NOT NULL
);