#!/bin/bash
# Script de inicio rápido para desarrollo con Docker

set -e

# Colores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Funciones de utilidad
info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

success() {
    echo -e "${GREEN}✅ $1${NC}"
}

warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

error() {
    echo -e "${RED}❌ $1${NC}"
    exit 1
}

# Verificar dependencias
check_dependencies() {
    info "Verificando dependencias..."

    if ! command -v docker &> /dev/null; then
        error "Docker no está instalado. Instálalo desde: https://docs.docker.com/get-docker/"
    fi

    if ! command -v docker-compose &> /dev/null && ! docker compose version &> /dev/null; then
        error "Docker Compose no está instalado"
    fi

    if ! command -v task &> /dev/null; then
        warning "Task (Taskfile) no está instalado"
        info "Recomendamos instalarlo para mejor experiencia: https://taskfile.dev/installation/"
        info "Por ahora usaremos docker-compose directamente"
        USE_TASK=false
    else
        USE_TASK=true
    fi

    success "Todas las dependencias están instaladas"
}

# Configurar X11
setup_x11() {
    info "Configurando permisos X11 para Docker..."
    xhost +local:docker &> /dev/null || warning "No se pudo configurar X11 (puede que no tengas entorno gráfico)"
    success "X11 configurado"
}

# Construir imagen
build_image() {
    info "Construyendo imagen Docker..."
    export UID=$(id -u)
    export GID=$(id -g)
    docker-compose build
    success "Imagen construida"
}

# Mostrar ayuda
show_help() {
    echo "PreyVPN - Script de Desarrollo Docker"
    echo ""
    echo "Uso: ./dev.sh [comando]"
    echo ""
    echo "Comandos disponibles:"
    echo "  setup         - Configurar entorno (primera vez)"
    echo "  up            - Iniciar servicios en background"
    echo "  dev           - Iniciar en modo desarrollo (foreground)"
    echo "  down          - Detener servicios"
    echo "  logs          - Ver logs"
    echo "  shell         - Abrir shell en el container"
    echo "  build         - Reconstruir imagen Docker"
    echo "  build-binary  - Compilar binario (NO requiere Go instalado)"
    echo "  clean         - Limpiar archivos generados"
    echo "  help          - Mostrar esta ayuda"
    echo ""
    if [ "$USE_TASK" = true ]; then
        echo "💡 Task está instalado. Puedes usar 'task --list' para ver más comandos"
    else
        echo "💡 Instala Task para más funcionalidades: https://taskfile.dev/installation/"
    fi
}

# Comando: setup
cmd_setup() {
    info "🚀 Configurando entorno de desarrollo..."
    check_dependencies
    setup_x11
    build_image
    success "Entorno configurado. Usa './dev.sh up' o './dev.sh dev' para iniciar"
}

# Comando: up
cmd_up() {
    setup_x11
    info "Iniciando servicios en background..."
    docker-compose up -d
    success "Servicios iniciados. Usa './dev.sh logs' para ver logs"
}

# Comando: dev
cmd_dev() {
    setup_x11
    info "Iniciando modo desarrollo (Ctrl+C para detener)..."
    docker-compose up
}

# Comando: down
cmd_down() {
    info "Deteniendo servicios..."
    docker-compose down
    success "Servicios detenidos"
}

# Comando: logs
cmd_logs() {
    docker-compose logs -f
}

# Comando: shell
cmd_shell() {
    info "Abriendo shell en el container..."
    docker-compose exec preyvpn /bin/bash
}

# Comando: build (imagen Docker)
cmd_build() {
    build_image
}

# Comando: build-binary (compilar binario sin instalar dependencias)
cmd_build_binary() {
    info "🔨 Compilando binario usando Docker (no requiere Go instalado)..."
    mkdir -p dist

    info "Construyendo imagen de compilación..."
    docker build -f Dockerfile.build -t preyvpn-builder --target builder . || error "Falló la construcción de la imagen"

    info "Compilando binario..."
    docker run --rm -v "$(pwd)/dist:/output" preyvpn-builder sh -c "cp /build/preyvpn /output/ && chmod +x /output/preyvpn" || error "Falló la compilación"

    success "Binario compilado exitosamente en ./dist/preyvpn"
    ls -lh dist/preyvpn
    file dist/preyvpn
}

# Comando: clean
cmd_clean() {
    info "Limpiando archivos generados..."
    rm -rf tmp/ build-errors.log
    docker-compose down -v
    success "Limpieza completada"
}

# Main
main() {
    cd "$(dirname "$0")"

    case "${1:-help}" in
        setup)
            cmd_setup
            ;;
        up)
            cmd_up
            ;;
        dev|start)
            cmd_dev
            ;;
        down|stop)
            cmd_down
            ;;
        logs)
            cmd_logs
            ;;
        shell|sh|bash)
            cmd_shell
            ;;
        build|rebuild)
            cmd_build
            ;;
        build-binary|compile)
            cmd_build_binary
            ;;
        clean)
            cmd_clean
            ;;
        help|--help|-h)
            show_help
            ;;
        *)
            error "Comando desconocido: $1\n\nUsa './dev.sh help' para ver comandos disponibles"
            ;;
    esac
}

main "$@"
