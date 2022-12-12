SELECT 'CREATE DATABASE "gophkeeper"'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'gophkeeper')\gexec
