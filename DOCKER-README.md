# NavTunnel - Entorno de Desarrollo Docker

Este proyecto está configurado con Docker, docker-compose y Taskfile para desarrollo rápido y reproducible.

## 📋 Requisitos Previos

- **Docker** (20.10+)
- **Docker Compose** (v2.0+)
- **Task** (Taskfile CLI): https://taskfile.dev/installation/
  ```bash
  # Ubuntu/Debian
  sudo sh -c "$(curl --location https://taskfile.dev/install.sh)" -- -d -b /usr/local/bin
  ```
- **X11** corriendo en tu host (para la GUI)

## 🚀 Inicio Rápido

### 1. Setup inicial (primera vez)

```bash
# Configurar el entorno
task setup
```

Esto:
- Configura permisos X11 para Docker
- Construye la imagen Docker
- Prepara el entorno

### 2. Iniciar desarrollo

```bash
# Modo desarrollo con hot-reload (logs en foreground)
task dev

# O en background
task up
```

La aplicación se reconstruirá automáticamente cuando cambies archivos `.go`.

### 3. Ver logs

```bash
# Logs en tiempo real
task logs

# Ver últimas 100 líneas
task exec CMD="tail -100 build-errors.log"
```

## 🏗️ Compilación sin Dependencias Locales

**Importante:** Puedes compilar NavTunnel usando Docker **SIN necesidad de instalar Go ni dependencias** en tu PC.

```bash
# Compilar binario (NO requiere Go instalado)
task build-docker

# O con dev.sh
./dev.sh build-binary

# El binario estará en ./dist/navtunnel
./dist/navtunnel
```

Ver [BUILD.md](BUILD.md) para documentación completa.

---

## 📚 Comandos Disponibles

### Compilación

```bash
task build-docker         # Compilar sin instalar dependencias
task build-docker-release # Compilar versión optimizada
```

### Gestión de Containers

```bash
task up              # Iniciar servicios en background
task down            # Detener y eliminar containers
task restart         # Reiniciar servicio
task stop            # Detener servicio
task start           # Iniciar servicio
task ps              # Listar containers
```

### Desarrollo

```bash
task dev             # Modo desarrollo con hot-reload
task compile         # Compilar binario manualmente
task run             # Ejecutar aplicación
task test            # Ejecutar tests
task fmt             # Formatear código
task vet             # Analizar código
task tidy            # Limpiar go.mod
```

### Logs y Debug

```bash
task logs            # Ver logs del servicio
task logs-all        # Ver logs de todos los servicios
task exec-sh         # Abrir shell en el container
task top             # Ver procesos corriendo
task stats           # Ver uso de recursos
```

### Construcción

```bash
task build           # Construir imagen Docker
task rebuild         # Reconstruir sin cache
task pull            # Actualizar imágenes base
```

### Limpieza

```bash
task clean           # Limpiar binarios y cache
task clean-all       # Limpieza profunda (containers + volúmenes + imágenes)
```

### Utilidades

```bash
task x11-fix         # Arreglar permisos X11
task vpn-config      # Verificar configuración VPN
task help            # Lista completa de comandos
```

## 🔧 Ejecutar Comandos Personalizados

```bash
# Ejecutar cualquier comando en el container
task exec CMD="go version"
task exec CMD="ls -la /app"
task exec CMD="sudo openvpn --version"

# Abrir shell interactivo
task exec-sh

# Dentro del container, puedes:
developer@container:/app$ go build ./cmd/navtunnel
developer@container:/app$ sudo openvpn --config ~/NavTunnel/tu-archivo.ovpn
```

## 🖥️ Desarrollo con GUI

El container está configurado para usar el X11 del host:

1. **Permisos X11**: Se configuran automáticamente con `task setup` o `task x11-fix`
2. **Display**: La variable `$DISPLAY` se pasa automáticamente al container
3. **Hot-reload**: Air detecta cambios y recompila automáticamente

### Problemas comunes de X11

Si la GUI no aparece:

```bash
# 1. Verificar que X11 permite conexiones de Docker
xhost +local:docker

# 2. Verificar DISPLAY
echo $DISPLAY

# 3. Reiniciar el container
task restart
```

## 📁 Estructura de Volúmenes

```
Host                          → Container
──────────────────────────────────────────────────
./                            → /app (código fuente)
~/NavTunnel/                    → /home/developer/NavTunnel (configs VPN)
/tmp/.X11-unix                → /tmp/.X11-unix (X11 socket)
```

## 🔐 Permisos y Privilegios

- El container corre con **tu UID/GID** para evitar problemas de permisos
- Tiene `CAP_NET_ADMIN` para OpenVPN
- Tiene acceso a `/dev/net/tun` para crear túneles VPN
- El usuario `developer` tiene sudo NOPASSWD solo para `/usr/sbin/openvpn`

## 🛠️ Workflow de Desarrollo Típico

```bash
# 1. Iniciar entorno de desarrollo
task dev

# 2. En otra terminal, ver logs en tiempo real
task logs

# 3. Editar código en tu IDE favorito (en el host)
# Air detectará cambios automáticamente y recompilará

# 4. Si necesitas ejecutar comandos manualmente
task exec-sh

# 5. Al terminar
task down
```

## 🐛 Debugging

### Ver errores de compilación

```bash
# En el container
cat /app/build-errors.log

# O desde el host
tail -f build-errors.log
```

### Ejecutar tests específicos

```bash
task exec CMD="go test -v ./internal/core/..."
```

### Depuración con Delve

```bash
task exec-sh

# Dentro del container
go install github.com/go-delve/delve/cmd/dlv@latest
dlv debug ./cmd/navtunnel
```

## 📦 Compilar Binario de Producción

```bash
# Compilar dentro del container
task compile

# El binario estará en ./tmp/navtunnel
# Copiarlo al host si es necesario
cp tmp/navtunnel ./navtunnel-linux-amd64
```

## 🔄 Actualizar Dependencias

```bash
# Agregar nueva dependencia
task exec CMD="go get github.com/some/package"

# Limpiar dependencias no usadas
task tidy

# Reconstruir imagen con nuevas dependencias
task rebuild
```

## 🚨 Solución de Problemas

### Container no inicia

```bash
# Ver logs de error
docker-compose logs

# Verificar que no hay containers corriendo
docker-compose ps

# Limpiar y reiniciar
task clean-all
task setup
```

### Problemas de permisos

```bash
# Reconstruir con tu UID/GID
docker-compose build --build-arg USER_ID=$(id -u) --build-arg GROUP_ID=$(id -g)
```

### Hot-reload no funciona

```bash
# Verificar que Air está corriendo
task logs

# Reiniciar Air
task restart

# Si sigue sin funcionar, ejecutar manualmente
task exec CMD="air -c .air.toml"
```

## 📖 Recursos

- [Taskfile](https://taskfile.dev/)
- [Docker Compose](https://docs.docker.com/compose/)
- [Air (Hot Reload)](https://github.com/air-verse/air)
- [Fyne GUI](https://developer.fyne.io/)

## 💡 Tips

1. **Usa `task` sin argumentos** para ver todos los comandos disponibles
2. **Mantén el container corriendo** con `task up` y usa `task exec` para comandos
3. **Los cambios en el código se reflejan inmediatamente** gracias a Air
4. **Los binarios compilados se guardan en `./tmp`** (ignorados por git)
5. **Usa `task exec-sh`** para explorar el container interactivamente

---

**¿Problemas?** Abre un issue o consulta los logs con `task logs`
