#!/usr/bin/env bash


containers="theangrydev_ruby \
	theangrydev_python \
	theangrydev_rust \
	theangrydev_node"
docker stop $containers && docker rm $containers
