#!/bin/bash
screen -S -d -m  mongod --nojournal --smallfiles --fork --syslog
