#!/bin/bash


pg_restore -d ${POSTGRES_DB} -U ${POSTGRES_USER} /docker-entrypoint-initdb.d/dvdrental.tar