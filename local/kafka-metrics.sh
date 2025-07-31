#!/bin/bash

while true; do
  # Get topic information
  docker exec kafka kafka-topics --bootstrap-server localhost:9092 --list | while read topic; do
    if [ -n "$topic" ]; then
      # Get partition info
      partitions=$(docker exec kafka kafka-topics --bootstrap-server localhost:9092 --describe --topic "$topic" --quiet | wc -l)
      
      # Get consumer group lag
      docker exec kafka kafka-consumer-groups --bootstrap-server localhost:9092 --list | while read group; do
        if [ -n "$group" ]; then
          lag=$(docker exec kafka kafka-consumer-groups --bootstrap-server localhost:9092 --describe --group "$group" --topic "$topic" 2>/dev/null | tail -n +2 | awk '{sum += $6} END {print sum+0}')
          if [ -n "$lag" ] && [ "$lag" != "0" ]; then
            echo "kafka.consumer.lag,topic=$topic,group=$group value=$lag"
          fi
        fi
      done
      
      echo "kafka.topic.partitions,topic=$topic value=$partitions"
    fi
  done
  
  sleep 30
done
