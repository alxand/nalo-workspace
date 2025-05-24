# DailyLog API

A Go Fiber + GORM backend API for managing daily logs, secured with JWT authentication.

## Features

- CRUD operations on daily logs and related entities
- JWT-based authentication and authorization
- PostgreSQL database integration via GORM
- Swagger/OpenAPI documentation at `/swagger/index.html`
- Validation on incoming JSON payloads
- Dockerized for easy deployment
- Automated testing included

## Requirements

- Docker & Docker Compose (recommended)
- Go 1.20+ (for local dev)
- PostgreSQL (if running without Docker)

## Setup and Run

### Using Docker (Recommended)

```bash
docker-compose up --build


# Nalo workspace API

[![Go Report Card](https://goreportcard.com/badge/github.com/alxand/nalo-workspace)](https://goreportcard.com/report/github.com/alxand/nalo-workspace)
[![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/alxand/nalo-workspace/ci.yml?branch=main&label=CI)](https://github.com/alxand/nalo-workspace/actions/workflows/ci.yml)
[![Docker Pulls](https://img.shields.io/docker/pulls/alexanderadade/nalo-workspace)](https://hub.docker.com/r/alexanderadade/nalo-workspace)
