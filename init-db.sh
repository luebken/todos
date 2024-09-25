#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 -h localhost --username "postgres" <<-EOSQL
    CREATE DATABASE todos;
    \c todos
    CREATE TABLE todos (
        item TEXT PRIMARY KEY
    );
    INSERT INTO todos (item) VALUES
    ('Buy groceries'),
    ('Finish homework'),
    ('Clean the house');
EOSQL