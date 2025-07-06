CREATE TABLE "articles" (
  "id" uuid PRIMARY KEY DEFAULT generate_uuid_v7(),
  "author_id" uuid NOT NULL,
  "title" varchar(255) NOT NULL,
  "body" text NOT NULL,
  "created_at" TIMESTAMPTZ(0) NOT NULL DEFAULT NOW(),
  FOREIGN KEY ("author_id") REFERENCES authors("id")
);
CREATE INDEX idx_articles_on_author_id ON articles("author_id");
