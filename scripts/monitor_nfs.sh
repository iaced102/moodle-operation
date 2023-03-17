#!/bin/bash
a=$(ssh root@103.148.57.178 "du -s /srv/nfs4/moodle/*")


# map a to dict
declare -A dict
while read -r line; do
	key=$(echo $line | awk '{print $2}')
	value=$(echo $line | awk '{print $1}')
	dict[$key]=$value
done <<< "$a"

# print dict
# for i in "${!dict[@]}"; do
# 	echo "$i ${dict[$i]}"
# done

# delete all from collection
mongo "mongodb://localhost:27017/moodle"  --quiet --eval "db.nfs_tracking.deleteMany({})"
# insert all to collection as array with id = 1
mongo "mongodb://localhost:27017/moodle"  --quiet --eval "db.nfs_tracking.insertOne({\"id\":1, \"data\":[$(for i in "${!dict[@]}"; do echo "{ \"path\": \"$i\", \"size\": ${dict[$i]} },"; done)]})"
