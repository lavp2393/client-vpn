# PreyVPN - Guía Rápida para Desarrollo

Cliente OpenVPN con GUI que maneja autenticación multi-factor (usuario + contraseña + OTP).

## 🚀 Inicio Rápido

### Compilar el binario (Linux)

Solo necesitas **Docker** instalado. NO requiere Go ni dependencias.

```bash
# Con Taskfile
task build-docker

# O con script
./dev.sh build-binary
```

**El binario compilado estará en:** `dist/preyvpn`

### Ejecutar el binario

```bash
# Dar permisos si es necesario
chmod +x dist/preyvpn

# Ejecutar
./dist/preyvpn
```

## 📋 Requisitos del Sistema

### Para compilar
- Docker (solo eso)

### Para ejecutar el binario
- OpenVPN instalado: `sudo apt install openvpn`
- Archivo de configuración en: `~/PreyVPN/prey-prod.ovpn`
- Sudo configurado para OpenVPN (opcional, facilita uso):
  ```bash
  echo "$USER ALL=(ALL) NOPASSWD: /usr/sbin/openvpn" | sudo tee /etc/sudoers.d/preyvpn-openvpn
  sudo chmod 0440 /etc/sudoers.d/preyvpn-openvpn
  ```

## 🛠️ Comandos Principales

### Compilación

```bash
# Compilar binario para Linux
task build-docker                    # Con Taskfile
./dev.sh build-binary                # Con script

# Ver comandos disponibles
task --list                          # Si tienes Task
./dev.sh help                        # Con script
```

### Desarrollo con Hot-Reload

```bash
# Iniciar entorno de desarrollo
task dev                             # Inicia container con hot-reload
task logs                            # Ver logs en tiempo real
task down                            # Detener todo

# Abrir shell en el container
task exec-sh
```

### Gestión de Containers

```bash
task up              # Iniciar en background
task down            # Detener y limpiar
task restart         # Reiniciar servicios
task ps              # Ver containers corriendo
```

## 📁 Estructura del Proyecto

```
binariovpnprey/
├── cmd/preyvpn/main.go          # Punto de entrada
├── internal/
│   ├── core/
│   │   ├── manager.go           # Gestión OpenVPN + PTY para prompts interactivos
│   │   └── openvpn.go           # Wrapper de proceso
│   ├── ui/
│   │   ├── app.go               # GUI principal (Fyne)
│   │   └── prompts.go           # Modales de autenticación
│   └── platform/
│       └── linux/               # Implementación específica de Linux
├── dist/                        # ⭐ Binarios compilados (aquí está preyvpn)
├── Dockerfile                   # Desarrollo con hot-reload
├── Dockerfile.build            # Compilación limpia
├── docker-compose.yml          # Servicios de desarrollo
├── Taskfile.yml                # Comandos automatizados
└── dev.sh                      # Script alternativo
```

## 🔧 Problema Actual: OTP

**Estado:** La aplicación NO está capturando correctamente el prompt del OTP.

**Contexto:**
- OpenVPN con static-challenge requiere: username → password → OTP
- El manager.go usa PTY (pseudo-terminal) para capturar prompts
- El problema está en que OpenVPN no está mostrando los prompts interactivos o el parser no los detecta

**Archivos relevantes:**
- `internal/core/manager.go`: Gestiona la comunicación con OpenVPN vía PTY
- `internal/ui/app.go`: Maneja los modales de la GUI

**Documentación:**
- `TECHNICAL_CONTEXT.md`: Análisis completo del problema OTP
- `PreyVPN_Spec_MVP.md`: Especificación original

## 🧪 Testing

```bash
# Compilar y probar localmente
task build-docker
./dist/preyvpn

# O con Docker en desarrollo (hot-reload)
task dev
# Edita archivos .go → se recompila automáticamente
```

## 📚 Documentación Completa

- **`BUILD.md`**: Guía detallada de compilación
- **`DOCKER-README.md`**: Documentación del entorno Docker
- **`TECHNICAL_CONTEXT.md`**: Análisis del problema OTP
- **`ARCHITECTURE.md`**: Arquitectura multi-plataforma

## ⚡ Workflow Típico

```bash
# 1. Clonar el repo
git clone <repo-url>
cd binariovpnprey

# 2. Compilar
task build-docker

# 3. Probar
./dist/preyvpn

# 4. Desarrollar (con hot-reload)
task dev
# Edita código → ve cambios en tiempo real

# 5. Limpiar
task down
task clean
```

## 🐛 Debug

### Ver logs del container
```bash
task logs
```

### Shell interactivo
```bash
task exec-sh
# Dentro del container puedes:
go build ./cmd/preyvpn
sudo openvpn --version
```

### Logs de OpenVPN en la app
Los logs aparecen en la ventana de la aplicación con formato:
```
[stdout] Enter Auth Username:
[stdout] Enter Auth Password:
[stdout] CHALLENGE: Your OTP
```

### Problemas comunes

**"Cannot connect to Docker daemon"**
```bash
sudo systemctl start docker
```

**"Permission denied" en el binario**
```bash
chmod +x dist/preyvpn
```

**GUI no aparece**
```bash
xhost +local:docker
task restart
```

## 🤝 Contribuir

1. Crear rama para tu feature
2. Editar código
3. Probar con `task dev` (hot-reload)
4. Compilar versión final: `task build-docker`
5. Commit y push

## 📞 Contacto

Para problemas o preguntas, revisar primero:
- `TECHNICAL_CONTEXT.md` para el problema del OTP
- `BUILD.md` para compilación
- `DOCKER-README.md` para desarrollo

---

**Última actualización:** 2025-11-05
**Estado:** MVP en desarrollo - problema de OTP pendiente
