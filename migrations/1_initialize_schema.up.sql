-- uuid v4 generator
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE user_record(
  id         uuid primary key default uuid_generate_v4(),
  name       TEXT      NOT NULL,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  deleted_at TIMESTAMP
);
