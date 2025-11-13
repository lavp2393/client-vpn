.PHONY: all build run clean install deps help
.PHONY: build-all build-linux build-windows build-darwin
.PHONY: build-all-arch clean-dist

# Variables
BINARY_NAME=navtunnel
BUILD_DIR=bin
DIST_DIR=dist
MAIN_PATH=cmd/navtunnel/main.go
VERSION?=dev
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS=-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME)
LDFLAGS_RELEASE=$(LDFLAGS) -s -w

# Detectar el sistema operativo y arquitectura actual
GOOS := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)

all: deps build

# Instalar dependencias
deps:
	@echo "📦 Instalando dependencias..."
	go mod download
	go mod tidy

# Compilar el binario para la plataforma actual
build: deps
	@echo "🔨 Compilando $(BINARY_NAME) para $(GOOS)/$(GOARCH)..."
	@mkdir -p $(BUILD_DIR)
	go build -ldflags="$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "✅ Binario creado en $(BUILD_DIR)/$(BINARY_NAME)"

# Compilar para distribución (sin símbolos de debug)
build-release: deps
	@echo "🔨 Compilando $(BINARY_NAME) para distribución ($(GOOS)/$(GOARCH))..."
	@mkdir -p $(BUILD_DIR)
	go build -ldflags="$(LDFLAGS_RELEASE)" -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "✅ Binario de distribución creado en $(BUILD_DIR)/$(BINARY_NAME)"

# ========================================
# Multi-platform builds
# ========================================

# Compilar para todas las plataformas
build-all: build-linux build-windows build-darwin
	@echo "✅ Compilación completada para todas las plataformas"

# Compilar todas las arquitecturas para todas las plataformas
build-all-arch: deps
	@echo "🌍 Compilando para todas las plataformas y arquitecturas..."
	@$(MAKE) build-linux-amd64
	@$(MAKE) build-linux-arm64
	@$(MAKE) build-windows-amd64
	@$(MAKE) build-windows-arm64
	@$(MAKE) build-darwin-amd64
	@$(MAKE) build-darwin-arm64
	@echo "✅ Compilación completada para todas las plataformas"

# Linux builds
build-linux: build-linux-amd64

build-linux-amd64: deps
	@echo "🐧 Compilando para Linux (amd64)..."
	@mkdir -p $(DIST_DIR)/linux-amd64
	GOOS=linux GOARCH=amd64 go build -ldflags="$(LDFLAGS_RELEASE)" -o $(DIST_DIR)/linux-amd64/$(BINARY_NAME) $(MAIN_PATH)
	@echo "✅ $(DIST_DIR)/linux-amd64/$(BINARY_NAME)"

build-linux-arm64: deps
	@echo "🐧 Compilando para Linux (arm64)..."
	@mkdir -p $(DIST_DIR)/linux-arm64
	GOOS=linux GOARCH=arm64 go build -ldflags="$(LDFLAGS_RELEASE)" -o $(DIST_DIR)/linux-arm64/$(BINARY_NAME) $(MAIN_PATH)
	@echo "✅ $(DIST_DIR)/linux-arm64/$(BINARY_NAME)"

# Windows builds
build-windows: build-windows-amd64

build-windows-amd64: deps
	@echo "🪟 Compilando para Windows (amd64)..."
	@mkdir -p $(DIST_DIR)/windows-amd64
	GOOS=windows GOARCH=amd64 go build -ldflags="$(LDFLAGS_RELEASE)" -o $(DIST_DIR)/windows-amd64/$(BINARY_NAME).exe $(MAIN_PATH)
	@echo "✅ $(DIST_DIR)/windows-amd64/$(BINARY_NAME).exe"

build-windows-arm64: deps
	@echo "🪟 Compilando para Windows (arm64)..."
	@mkdir -p $(DIST_DIR)/windows-arm64
	GOOS=windows GOARCH=arm64 go build -ldflags="$(LDFLAGS_RELEASE)" -o $(DIST_DIR)/windows-arm64/$(BINARY_NAME).exe $(MAIN_PATH)
	@echo "✅ $(DIST_DIR)/windows-arm64/$(BINARY_NAME).exe"

# macOS builds
build-darwin: build-darwin-amd64 build-darwin-arm64

build-darwin-amd64: deps
	@echo "🍎 Compilando para macOS (amd64 - Intel)..."
	@mkdir -p $(DIST_DIR)/darwin-amd64
	GOOS=darwin GOARCH=amd64 go build -ldflags="$(LDFLAGS_RELEASE)" -o $(DIST_DIR)/darwin-amd64/$(BINARY_NAME) $(MAIN_PATH)
	@echo "✅ $(DIST_DIR)/darwin-amd64/$(BINARY_NAME)"

build-darwin-arm64: deps
	@echo "🍎 Compilando para macOS (arm64 - Apple Silicon)..."
	@mkdir -p $(DIST_DIR)/darwin-arm64
	GOOS=darwin GOARCH=arm64 go build -ldflags="$(LDFLAGS_RELEASE)" -o $(DIST_DIR)/darwin-arm64/$(BINARY_NAME) $(MAIN_PATH)
	@echo "✅ $(DIST_DIR)/darwin-arm64/$(BINARY_NAME)"

# ========================================
# Packaging (futuro)
# ========================================

package-linux: build-linux
	@echo "📦 Empaquetando para Linux..."
	@echo "⚠️  Packaging no implementado aún"
	@echo "TODO: Crear .deb, .rpm, AppImage"

package-windows: build-windows
	@echo "📦 Empaquetando para Windows..."
	@echo "⚠️  Packaging no implementado aún"
	@echo "TODO: Crear instalador .msi o .exe"

package-darwin: build-darwin
	@echo "📦 Empaquetando para macOS..."
	@echo "⚠️  Packaging no implementado aún"
	@echo "TODO: Crear .app bundle y .dmg"

# ========================================
# Utilidades
# ========================================

# Ejecutar la aplicación
run: build
	@echo "🚀 Ejecutando $(BINARY_NAME)..."
	./$(BUILD_DIR)/$(BINARY_NAME)

# Limpiar archivos generados
clean:
	@echo "🧹 Limpiando..."
	rm -rf $(BUILD_DIR)
	go clean

# Limpiar distribuciones
clean-dist:
	@echo "🧹 Limpiando distribuciones..."
	rm -rf $(DIST_DIR)

# Limpiar todo
clean-all: clean clean-dist
	@echo "✅ Limpieza completa"

# Instalar el binario en el sistema (solo Linux/macOS)
install: build
	@echo "📥 Instalando $(BINARY_NAME) en /usr/local/bin..."
	sudo cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/
	@echo "✅ Instalación completada. Ejecuta '$(BINARY_NAME)' desde cualquier lugar."

# Desinstalar el binario del sistema
uninstall:
	@echo "🗑️  Desinstalando $(BINARY_NAME)..."
	sudo rm -f /usr/local/bin/$(BINARY_NAME)
	@echo "✅ Desinstalación completada."

# Verificar dependencias del sistema (Linux)
check-deps:
	@echo "🔍 Verificando dependencias del sistema..."
	@which openvpn > /dev/null || (echo "❌ OpenVPN no está instalado. Instala con: sudo apt install openvpn" && exit 1)
	@which pkexec > /dev/null || (echo "❌ pkexec no está instalado. Instala con: sudo apt install policykit-1" && exit 1)
	@which go > /dev/null || (echo "❌ Go no está instalado. Instala desde https://golang.org/dl/" && exit 1)
	@echo "✅ Todas las dependencias están instaladas"

# Preparar el directorio de configuración
setup-config:
	@echo "📁 Creando directorio de configuración..."
	@mkdir -p ~/.config/NavTunnel
	@echo "✅ Directorio ~/.config/NavTunnel creado"
	@echo "💡 La aplicación te pedirá seleccionar tu archivo .ovpn al iniciar"

# Mostrar información del sistema
info:
	@echo "ℹ️  Información del sistema:"
	@echo "  GOOS:    $(GOOS)"
	@echo "  GOARCH:  $(GOARCH)"
	@echo "  Go:      $(shell go version)"
	@echo "  Version: $(VERSION)"

# Mostrar ayuda
help:
	@echo "NavTunnel - Makefile Multi-Platform"
	@echo ""
	@echo "📦 Desarrollo:"
	@echo "  make deps           - Instalar dependencias de Go"
	@echo "  make build          - Compilar el binario para la plataforma actual"
	@echo "  make build-release  - Compilar para distribución (optimizado)"
	@echo "  make run            - Compilar y ejecutar"
	@echo "  make clean          - Limpiar archivos generados"
	@echo ""
	@echo "🌍 Multi-platform:"
	@echo "  make build-all      - Compilar para todas las plataformas (main arch)"
	@echo "  make build-all-arch - Compilar para todas las plataformas y arquitecturas"
	@echo "  make build-linux    - Compilar para Linux (amd64)"
	@echo "  make build-windows  - Compilar para Windows (amd64)"
	@echo "  make build-darwin   - Compilar para macOS (amd64 + arm64)"
	@echo ""
	@echo "📦 Packaging (futuro):"
	@echo "  make package-linux  - Crear paquetes para Linux"
	@echo "  make package-windows- Crear instalador para Windows"
	@echo "  make package-darwin - Crear bundle para macOS"
	@echo ""
	@echo "🛠️  Utilidades:"
	@echo "  make install        - Instalar en el sistema"
	@echo "  make uninstall      - Desinstalar del sistema"
	@echo "  make check-deps     - Verificar dependencias del sistema"
	@echo "  make setup-config   - Crear directorio de configuración"
	@echo "  make info           - Mostrar información del sistema"
	@echo "  make clean-all      - Limpiar todo (bin + dist)"
	@echo "  make help           - Mostrar esta ayuda"
