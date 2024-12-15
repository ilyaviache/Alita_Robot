#!/bin/bash

function setup() {
    echo "Setting up development environment..."
    go mod tidy
    go install github.com/air-verse/air@latest
}

function start_services() {
    docker-compose -f debug.docker-compose.yml up -d mongodb redis
}

function stop_services() {
    docker-compose -f debug.docker-compose.yml down
}

function run_local() {
    export $(cat .env.local | xargs)
    go run main.go
}

function run_docker() {
    docker-compose -f debug.docker-compose.yml up --build
}

function watch() {
    air -c .air.toml
}

case "$1" in
    "setup")
        setup
        ;;
    "services")
        start_services
        ;;
    "stop")
        stop_services
        ;;
    "local")
        setup
        start_services
        run_local
        ;;
    "docker")
        run_docker
        ;;
    "watch")
        setup
        start_services
        watch
        ;;
    *)
        echo "Usage: ./dev.sh [setup|services|stop|local|docker|watch]"
        echo "  setup: Install dependencies and prepare development environment"
        echo "  services: Start MongoDB and Redis in containers"
        echo "  stop: Stop all containers"
        echo "  local: Run app locally with containerized services"
        echo "  docker: Run everything in containers"
        echo "  watch: Run locally with hot reload"
        ;;
esac
