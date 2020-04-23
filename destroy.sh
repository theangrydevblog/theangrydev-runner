#!/usr/bin/env bash


containers="theangrydev_ruby \
	theangrydev_python \
	theangrydev_node"
docker stop $containers && docker rm $containers
