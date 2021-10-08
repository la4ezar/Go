#!/bin/bash

IFS=$'\n' read -ra arr -d '' </conf/migrator.conf

DB_USER="${arr[0]}"
DB_HOST="${arr[1]}"
DB_NAME="${arr[2]}"
DB_PORT="${arr[3]}"
DB_PASSWORD="${arr[4]}"
DB_SSL="${arr[5]}"
DIRECTION="${arr[6]}"

for var in DB_USER DB_HOST DB_NAME DB_PORT DB_PASSWORD DB_SSL DIRECTION; do
    if [ -z "${!var}" ] ; then
        echo "ERROR: $var is not set"
        discoverUnsetVar=true
    fi
done

if [ "${discoverUnsetVar}" = true ] ; then
    exit 1
fi

if [[ "${DIRECTION}" == "up" ]]; then
    echo "Migration UP"
elif [[ "${DIRECTION}" == "down" ]]; then
    echo "Migration DOWN"
else
    echo "ERROR: DIRECTION variable accepts only two values: up or down"
    exit 1
fi

echo '# WAITING FOR CONNECTION WITH DATABASE #'
for i in {1..30}
do
    pg_isready -U "$DB_USER" -h "$DB_HOST" -p "$DB_PORT" -d "$DB_NAME"
    if [ $? -eq 0 ]
    then
        dbReady=true
        break
    fi
    sleep 1
done

if [ "${dbReady}" != true ] ; then
    echo '# COULD NOT ESTABLISH CONNECTION TO DATABASE #'
    exit 1
fi

DB_NAME_SSL="$DB_NAME?sslmode=${DB_SSL}"

# stops the execution of a script if a command or pipeline has an error
set -e

CONNECTION_STRING="postgres://$DB_USER:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME_SSL"

CMD="migrate -path migrations/ -database "$CONNECTION_STRING" ${DIRECTION}"
echo '# STARTING MIGRATION #'
if [[ "${NON_INTERACTIVE}" == "true" ]]; then
    yes | $CMD
else
    $CMD
fi

set +e