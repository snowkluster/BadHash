#!/bin/sh

Green="\e[32m"
Red="\033[0;31m"
NC='\033[0m'

if [ "$1" = "stop" ]; then
    echo "${Red}Stopping All Containers${NC}"
    sudo docker compose down --rmi local
    exit 0
fi

echo "${Green}Starting All Containers${NC}"
sudo docker compose up -d

go run main.go