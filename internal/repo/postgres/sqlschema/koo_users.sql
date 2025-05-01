CREATE TABLE koo_users (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4 (),
  is_subscribed BOOLEAN NOT NULL DEFAULT FALSE,
  first_name VARCHAR(255) NOT NULL
)
;
