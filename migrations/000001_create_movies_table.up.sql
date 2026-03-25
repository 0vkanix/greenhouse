CREATE TABLE IF NOT EXISTS movies (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
  title text NOT NULL,
  year integer NOT NULL,
  runtime integer NOT NULL,
  genres text[] NOT NULL,
  version integer NOT NULL DEFAULT 1
);
