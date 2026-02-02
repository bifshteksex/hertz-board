@echo off
REM Run HertzBoard API Gateway with local configuration
set CONFIG_PATH=configs/config.yaml
go run cmd/api-gateway/main.go
