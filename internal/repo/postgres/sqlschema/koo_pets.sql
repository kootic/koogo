CREATE TABLE koo_pets (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4 (),
  owner_id uuid NOT NULL REFERENCES koo_users (id)
)
;
