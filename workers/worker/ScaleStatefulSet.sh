#!/bin/bash

# Redis server connection details
REDIS_HOST="localhost"
REDIS_PORT="6379"

# Key to retrieve from Redis
KEY="replicas_16d53fd2-9f0c-48a4-921a-251009f9768a"

prefix="replicas_"
NAMESPACE="${KEY#$prefix}"

# Retrieve value from Redis
value=$(redis-cli -h $REDIS_HOST -p $REDIS_PORT GET $KEY)

# Check if the value is empty
if [ -z "$value" ]; then
    echo "Value for key '$KEY' not found in Redis."
else
    echo "Value for key '$KEY': $value"
fi

# Set replicas for StatefulSet
echo "Set scale for StatefulSet $NAMESPACE to $value"

kubectl -n $NAMESPACE scale sts moodle --replicas=$value
