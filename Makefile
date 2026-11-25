APP_SERVER_NAME = server
APP_CLIENT_NAME = client
APP_AGENT_NAME = agent


APP_SERVER_DIR = ./cmd/server
APP_CLIENT_DIR = ./cmd/client
APP_AGENT_DIR = ./cmd/agent

BUILD_DIR = ./build
TARGET_SERVER_BIN=$(BUILD_DIR)/$(APP_SERVER_NAME)
TARGET_CLIENT_BIN=$(BUILD_DIR)/$(APP_CLIENT_NAME)
TARGET_AGENT_BIN=$(BUILD_DIR)/$(APP_AGENT_NAME)

.PHONY: all build run proto clean showdb 


all: build

build: build_server build_client build_agent 

build_server:
	go build -o $(TARGET_SERVER_BIN) $(APP_SERVER_DIR)

build_client:
	go build -o $(TARGET_CLIENT_BIN) $(APP_CLIENT_DIR)

build_agent:
	go build -o $(TARGET_AGENT_BIN) $(APP_AGENT_DIR)

build_win_client:
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o $(TARGET_CLIENT_BIN).exe $(APP_CLIENT_DIR)

build_android_client:
	GOOS=android GOARCH=arm64 CGO_ENABLED=0 go build -o $(TARGET_CLIENT_BIN).apk $(APP_CLIENT_DIR)

runs: build_server
	$(TARGET_SERVER_BIN)

runc: build_client
	$(TARGET_CLIENT_BIN)

runa: build_agent
	$(TARGET_AGENT_BIN)

run: build
	@echo "==> Запуск сервера..."
	@$(TARGET_SERVER_BIN) & server_pid=$$!; \
	sleep 1; \
	echo "==> Запуск клиента..."; \
	$(TARGET_CLIENT_BIN); \
	echo "==> Остановка сервера..."; \
	kill $$server_pid 2>/dev/null || true



proto:
	protoc --go_out=. --go_opt=module=monitoring \
	--go-grpc_out=. --go-grpc_opt=module=monitoring \
	proto/metrics/v1/metrics.proto



LIMIT ?= 10
showdb:
	@sqlite3 -box database/monitoring.db "SELECT id, host, cpu_usage, ram_usage, datetime(timestamp, 'unixepoch', 'localtime') AS date_time FROM metrics ORDER BY id DESC LIMIT $(LIMIT);"


clean:
	rm -rf $(BUILD_DIR)

#save-projects-structure
sps:
	tree -I 'vendor|.git|bin' -a > project_structure.txt

