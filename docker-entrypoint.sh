#!/bin/bash
./migrator --uri=postgres:postgres@localhost:5430/postgres --migrations-path=./migrations
./billing --config=./config/dev.yaml