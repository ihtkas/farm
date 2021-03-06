#!/bin/sh
create_cassandra_keyspace() {
    cqlsh -e "CREATE KEYSPACE IF NOT EXISTS farm_seller WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 }"
    cqlsh -e "CREATE KEYSPACE IF NOT EXISTS farm_buyer WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 }"
}

create_postgres_database() {
    psql -p 5434 -U postgres -c "create database farm_seller"
}