#!/bin/bash

# This is a temporary hacky script which should run all the components of the system using gnome-terminal sessions.
# This will be replaced with a proper docker setup later :)
kafkaDir=/home/alex/workspace/kafka-udemy-course/kakfa-2.4.0

# Launch ZooKeeper
gnome-terminal -t zookeeper -- $kafkaDir/bin/zookeeper-server-start.sh $kafkaDir/config/zookeeper.properties

# Wait 2 Seconds, then launch Kafka
sleep 2s
gnome-terminal -t kafka -- $kafkaDir/bin/kafka-server-start.sh $kafkaDir/config/server.properties

# Wait 5 Seconds, then launch the crawler and the link checker
sleep 5s

scriptDir=$(dirname -- "$(readlink -f -- "$BASH_SOURCE")")

gnome-terminal -t prawCrawler -- "$scriptDir/RunRickRollCrawler.sh"
gnome-terminal -t linkChecker -- "$scriptDir/RunLinkChecker.sh"