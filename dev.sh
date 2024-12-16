#!/bin/bash

function setup() {
    echo "Setting up development environment..."
    go mod tidy
    go install github.com/air-verse/air@latest
}

function start_services() {
    echo "Starting MongoDB and Redis services..."
    docker-compose -f debug.docker-compose.yml up -d mongodb redis
}

function stop_services() {
    echo "Stopping all services..."
    docker-compose -f debug.docker-compose.yml down
}

function run_local() {
    echo "Running app locally..."
    export $(cat .env.local | xargs)
    export DEBUG=true
    go run main.go
}

function run_docker() {
    echo "Running in Docker..."
    docker-compose -f debug.docker-compose.yml up --build
}

function watch() {
    echo "Running with hot reload..."
    export DEBUG=true
    air -c .air.toml
}

function dev() {
    echo "Starting development environment..."
    setup
    start_services
    watch
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
    "dev")
        dev
        ;;
    *)
        echo "Usage: ./dev.sh [setup|services|stop|local|docker|watch|dev]"
        echo "  setup: Install dependencies and prepare development environment"
        echo "  services: Start MongoDB and Redis in containers"
        echo "  stop: Stop all containers"
        echo "  local: Run app locally with containerized services"
        echo "  docker: Run everything in containers"
        echo "  watch: Run locally with hot reload"
        echo "  dev: Setup, start services and run with hot reload"
        ;;
esac
