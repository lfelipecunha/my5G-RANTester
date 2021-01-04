#!/bin/bash

if [ $# -lt 1 ]; then
    echo "Usage: $0 LOGS_PATH"
    exit -1
fi


BASE_DIR=$1

mkdir $BASE_DIR
rm $BASE_DIR/*
i=0

sudo docker-compose logs > ./$BASE_DIR/full.log

total=$(grep -iE '\UE\[|\[UE' ./$BASE_DIR/full.log | sed 's/\[//g' | sed 's/.*UEimsi-20893\([0-9]\+\)[^0-9].*/\1/' | sort -nr | head -1)
echo "Total: $total"
while [ $i -le $total ]; do
    echo -n '.'
    grep -E "imsi-20893(0+)$i[^0-9]" ./$BASE_DIR/full.log > ./$BASE_DIR/UE_$i.log
    i=$((i+1))
done

echo "Total finished: $(grep "REGISTRATION FINISHED" ./$BASE_DIR/UE* | wc -l) OF $total"
