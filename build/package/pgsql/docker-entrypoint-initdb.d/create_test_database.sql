SELECT 'CREATE DATABASE "gophkeeper_test"'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'gophkeeper_test')\gexec
