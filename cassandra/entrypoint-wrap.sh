#!/bin/bash
# shellcheck disable=SC2089
  CQL="CREATE KEYSPACE IF NOT EXISTS oauth WITH REPLICATION = {'class': 'SimpleStrategy', 'replication_factor': 1};USE oauth;CREATE TABLE IF NOT EXISTS access_tokens( access_token varchar PRIMARY KEY, user_id bigint, client_id bigint, expires bigint);create index if not exists access_tokens_user_id_index on access_tokens (user_id);"
  until echo $CQL | cqlsh cassandra-main; do
    echo "cqlsh: Cassandra is unavailable - retry later"
    sleep 2
  done &


exec /docker-entrypoint.sh "$@"