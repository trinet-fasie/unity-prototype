#!/bin/sh

rabbitmqctl add_user creator Kasd9k1knksknjbbk
rabbitmqctl set_permissions -p "/" creator ".*" ".*" ".*"
