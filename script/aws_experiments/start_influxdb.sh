#!/bin/bash

DATA_DIR=${DATA_DIR:-/disk}
INFLUXDB_VERSION=${INFLUXDB_VERSION:-1.0.0-alpine}
INFLUXDB_CONF=${INFLUXDB_CONF:-/tmp/influxdb.conf}
INFLUXDB_META_DIR=${INFLUXDB_META_DIR:-${DATA_DIR}/1/influxdb/meta}
INFLUXDB_DATA_DIR=${INFLUXDB_DATA_DIR:-${DATA_DIR}/1/influxdb/data}
INFLUXDB_WAL_DIR=${INFLUXDB_WAL_DIR:-${DATA_DIR}/2/influxdb/wal}

docker stop influxdb &> /dev/null

# Generate default config
echo "Generating default config ${INFLUXDB_CONF}"
docker run --rm influxdb:${INFLUXDB_VERSION} influxd config > $INFLUXDB_CONF

echo "Data dirs: ${DATA_DIR} ${INFLUXDB_DATA_DIR}"
docker run -d \
       --name influxdb \
       -p 8083:8083 \
       -p 8086:8086 \
       -v ${INFLUXDB_DATA_DIR}:/var/lib/influxdb/data \
       -v ${INFLUXDB_META_DIR}:/var/lib/influxdb/meta \
       -v ${INFLUXDB_WAL_DIR}:/var/lib/influxdb/wal \
       -v ${INFLUXDB_CONF}:/etc/influxdb/influxdb.conf:ro \
       influxdb:${INFLUXDB_VERSION} -config /etc/influxdb/influxdb.conf