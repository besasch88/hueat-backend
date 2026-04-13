#!/bin/sh

# This script run all the migrations. It is used in development environment and it contains hardcoded parameters to connect to DB.
migrate -path "./scripts/migrations" -database "postgres://hueat:iUkcBQj2o_PpfX*uaXx7@127.0.0.1:54322/hueat?sslmode=disable" up
