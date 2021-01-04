#!/bin/bash

if [ $# -lt 1 ]; then
    echo "Usage: $0 LOGS_PATH"
    exit -1
fi


BASE_DIR=$1

total=$(grep -iE '\UE\[|\[UE' ./$BASE_DIR/full.log | sed 's/\[//g' | sed 's/.*UEimsi-20893\([0-9]\+\)[^0-9].*/\1/' | sort -nr | head -1)
i=1
echo "UEs not complete:"
while [ $i -le $total ]; do
    has=$(grep "REGISTRATION FINISHED" ./$BASE_DIR/UE_$i.log | wc -l)
    if [ $has -lt 1 ] ; then
        cause="SCTP Connection"
        if [ $(grep "REGISTRATION REQUEST" ./$BASE_DIR/UE_$i.log| wc -l) -ge 1 ]; then
            cause="others"
        fi
        echo "    $i - $cause"
    fi
    i=$((i+1))
done

