CREATE TABLE "authors" (
  "id" uuid PRIMARY KEY DEFAULT generate_uuid_v7(),
  "name" varchar(100) NOT NULL
);
CREATE INDEX idx_authors_on_name ON authors("name");
