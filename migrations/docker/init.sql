-- ---------- 1. Create roles ----------
DO $$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_roles WHERE rolname = 'db_migrator') THEN
        CREATE ROLE db_migrator LOGIN PASSWORD '12345678';
    END IF;

    IF NOT EXISTS (SELECT FROM pg_roles WHERE rolname = 'db_rw') THEN
        CREATE ROLE db_rw LOGIN PASSWORD '12345678';
    END IF;

    IF NOT EXISTS (SELECT FROM pg_roles WHERE rolname = 'db_test') THEN
        CREATE ROLE db_test LOGIN PASSWORD 'test';
    END IF;
END
$$;

-- ---------- 2. Create databases ----------
CREATE DATABASE greenlight_test;

-- ============================================================
-- Apply permissions for each database
-- ============================================================

------------------------------------------------------------
-- Setup for greenlight
------------------------------------------------------------
\connect greenlight

CREATE EXTENSION IF NOT EXISTS citext;

-- Remove unsafe default privileges
REVOKE ALL ON DATABASE greenlight FROM PUBLIC;
REVOKE CREATE ON SCHEMA public FROM PUBLIC;

-- Prisma requires the db_migrator to own the schema
ALTER SCHEMA public OWNER TO db_migrator;

-- db_Migrator full control
GRANT USAGE, CREATE ON SCHEMA public TO db_migrator;

-- Read/Write application user
GRANT USAGE ON SCHEMA public TO db_rw;
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO db_rw;
GRANT USAGE, SELECT, UPDATE ON ALL SEQUENCES IN SCHEMA public TO db_rw;

-- Default privileges for FUTURE objects created by db_migrator
ALTER DEFAULT PRIVILEGES FOR ROLE db_migrator IN SCHEMA public
    GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO db_rw;

ALTER DEFAULT PRIVILEGES FOR ROLE db_migrator IN SCHEMA public
    GRANT USAGE, SELECT, UPDATE ON SEQUENCES TO db_rw;
    --
-- Allow connections
GRANT CONNECT ON DATABASE greenlight TO db_migrator, db_rw;

-------------------------------
-- Setup for greenlight_test database
-------------------------------
\connect greenlight_test

CREATE EXTENSION IF NOT EXISTS citext;

-- Grant full control to db_test
GRANT ALL PRIVILEGES ON DATABASE greenlight_test TO db_test;

REVOKE ALL ON DATABASE greenlight_test FROM PUBLIC;
REVOKE CREATE ON SCHEMA public FROM PUBLIC;

ALTER SCHEMA public OWNER TO db_test;
GRANT ALL PRIVILEGES ON SCHEMA public TO db_test;

-- Default privileges for objects created in the schema
ALTER DEFAULT PRIVILEGES FOR ROLE db_test IN SCHEMA public
    GRANT ALL PRIVILEGES ON TABLES TO db_test;

ALTER DEFAULT PRIVILEGES FOR ROLE db_test IN SCHEMA public
    GRANT ALL PRIVILEGES ON SEQUENCES TO db_test;

GRANT CONNECT ON DATABASE greenlight_test TO db_migrator, db_test;
