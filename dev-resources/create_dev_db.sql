SELECT 'CREATE DATABASE pollfax_db'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'pollfax_db')\gexec
