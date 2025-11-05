# PreyVPN - Arquitectura Multi-Plataforma

## Última actualización: 2025-11-04

---

## Visión General

PreyVPN es un cliente OpenVPN con interfaz gráfica que soporta múltiples plataformas mediante una arquitectura modular y abstracciones específicas por sistema operativo.

### Plataformas Soportadas

| Plataforma | Estado | Arquitecturas |
|------------|--------|---------------|
| **Linux** | ✅ Completo | amd64, arm64 |
| **Windows** | 🚧 En desarrollo | amd64, arm64 |
| **macOS** | 🚧 En desarrollo | amd64 (Intel), arm64 (Apple Silicon) |

---

## Estructura del Proyecto

```
binariovpnprey/
├── cmd/
│   └── preyvpn/
│       └── main.go                    # Entry point común para todas las plataformas
│
├── internal/
│   ├── core/
│   │   ├── manager.go                 # Management Interface (común)
│   │   └── openvpn.go                 # Wrapper que usa platform abstraction
│   │
│   ├── platform/                      # ⭐ Abstracciones por plataforma
│   │   ├── platform.go                # Interface común
│   │   ├── platform_linux.go          # Build tags para Linux
│   │   ├── platform_windows.go        # Build tags para Windows
│   │   ├── platform_darwin.go         # Build tags para macOS
│   │   │
│   │   ├── linux/
│   │   │   └── linux.go               # Implementación completa para Linux
│   │   │
│   │   ├── windows/
│   │   │   └── windows.go             # Stub con TODOs
│   │   │
│   │   └── darwin/
│   │       └── darwin.go              # Stub con TODOs
│   │
│   ├── ui/
│   │   ├── app.go                     # UI común (Fyne es cross-platform)
│   │   └── prompts.go
│   │
│   └── logs/
│       └── buffer.go
│
├── build/                             # Scripts de build por plataforma
│   ├── linux/
│   ├── windows/
│   └── darwin/
│
├── dist/                              # Binarios compilados
│   ├── linux-amd64/
│   ├── linux-arm64/
│   ├── windows-amd64/
│   ├── windows-arm64/
│   ├── darwin-amd64/
│   └── darwin-arm64/
│
├── configs/                           # Configuraciones por plataforma
│   ├── linux/
│   │   └── preyvpn.desktop           # Desktop entry para Linux
│   ├── windows/
│   │   └── README.md                 # Guía para iconos, manifests, etc.
│   └── darwin/
│       └── Info.plist                # App bundle info para macOS
│
├── Makefile                           # Build system multi-plataforma
├── go.mod
├── README.md
├── ARCHITECTURE.md                    # Este archivo
├── PreyVPN_Spec_MVP.md
└── TECHNICAL_CONTEXT.md
```

---

## Abstracción de Plataforma

### Interface `platform.Platform`

Define el contrato común que todas las plataformas deben implementar:

```go
type Platform interface {
    // Process management
    FindOpenVPN() (string, error)
    StartOpenVPN(config StartConfig) (*Process, error)
    StopOpenVPN(proc *Process) error

    // Privilege elevation
    RequiresElevation() bool
    ElevateCommand(path string, args []string) (string, []string, error)

    // Paths
    GetConfigDir() string
    GetDefaultConfigPath() string
    GetLogPath() string

    // Platform info
    Name() string
    Separator() string
}
```

### Selección Automática de Plataforma

El código usa **build tags** de Go para compilar solo la implementación correcta:

```go
// internal/platform/platform.go
func New() Platform {
    switch runtime.GOOS {
    case "linux":
        return NewLinux()
    case "windows":
        return NewWindows()
    case "darwin":
        return NewDarwin()
    }
}
```

Los archivos `platform_*.go` tienen build tags:
- `//go:build linux` → `platform_linux.go`
- `//go:build windows` → `platform_windows.go`
- `//go:build darwin` → `platform_darwin.go`

---

## Diferencias por Plataforma

### Linux (Completo)

| Aspecto | Implementación |
|---------|----------------|
| **OpenVPN Path** | `/usr/sbin/openvpn`, `/usr/bin/openvpn` |
| **Config Dir** | `~/.config/PreyVPN` (XDG spec) o `~/PreyVPN` (MVP) |
| **Log Path** | `~/.cache/PreyVPN/logs` |
| **Elevation** | `pkexec` (PolicyKit) |
| **Packaging** | .deb, .rpm, AppImage (futuro) |
| **Desktop Entry** | `configs/linux/preyvpn.desktop` |

**Dependencias:**
- `openvpn`
- `policykit-1` (pkexec)
- `libgl1-mesa-dev`, `xorg-dev` (para Fyne)

### Windows (En desarrollo)

| Aspecto | Implementación |
|---------|----------------|
| **OpenVPN Path** | `C:\Program Files\OpenVPN\bin\openvpn.exe` |
| **Config Dir** | `%APPDATA%\PreyVPN` |
| **Log Path** | `%LOCALAPPDATA%\PreyVPN\logs` |
| **Elevation** | UAC / `runas` / ShellExecute |
| **Packaging** | .msi, .exe installer (NSIS/WiX) |
| **Icon** | `configs/windows/preyvpn.ico` |

**TODOs:**
- [ ] Implementar elevación con UAC
- [ ] Manejar rutas de Windows correctamente
- [ ] Probar con OpenVPN GUI service
- [ ] Crear script de instalador NSIS

### macOS (En desarrollo)

| Aspecto | Implementación |
|---------|----------------|
| **OpenVPN Path** | `/usr/local/opt/openvpn/sbin/openvpn` (Homebrew) |
| **Config Dir** | `~/Library/Application Support/PreyVPN` |
| **Log Path** | `~/Library/Logs/PreyVPN` |
| **Elevation** | `osascript` (AppleScript) / SMJobBless |
| **Packaging** | .app bundle, .dmg |
| **Bundle Info** | `configs/darwin/Info.plist` |

**TODOs:**
- [ ] Implementar elevación con osascript
- [ ] Crear .app bundle correctamente
- [ ] Firmar código (para distribución)
- [ ] Probar en Apple Silicon (arm64)

---

## Build System

### Comandos Principales

```bash
# Desarrollo (plataforma actual)
make build          # Compilar para la plataforma actual
make run            # Compilar y ejecutar
make clean          # Limpiar archivos generados

# Multi-plataforma
make build-all      # Compilar para Linux, Windows, macOS (arch principal)
make build-all-arch # Compilar para todas las arquitecturas

# Específico por plataforma
make build-linux    # Linux amd64
make build-windows  # Windows amd64
make build-darwin   # macOS amd64 + arm64

# Utilidades
make info           # Mostrar información del sistema
make check-deps     # Verificar dependencias (Linux)
make help           # Ayuda completa
```

### Variables de Entorno

```bash
VERSION=v1.0.0 make build-release
```

---

## Cross-Compilation

Go soporta cross-compilation de forma nativa:

```bash
# Desde Linux, compilar para Windows
GOOS=windows GOARCH=amd64 go build -o preyvpn.exe cmd/preyvpn/main.go

# Desde Linux, compilar para macOS Apple Silicon
GOOS=darwin GOARCH=arm64 go build -o preyvpn cmd/preyvpn/main.go
```

### Limitaciones de Cross-Compilation

- **CGO**: Fyne requiere CGO, así que necesitas cross-compilers:
  - Linux → Windows: `mingw-w64`
  - Linux → macOS: `osxcross`
- **Pruebas**: Solo se puede probar en la plataforma nativa

---

## Flujo de Integración

### Añadir Soporte para Nueva Plataforma

1. **Crear implementación:** `internal/platform/<os>/<os>.go`
2. **Implementar interface:** Todos los métodos de `platform.Platform`
3. **Crear build tag:** `internal/platform/platform_<os>.go`
4. **Añadir target al Makefile:** `build-<os>`
5. **Configuración:** Añadir archivos en `configs/<os>/`
6. **Documentar:** Actualizar este archivo

### Probar en Múltiples Plataformas

```bash
# CI/CD debería probar en cada plataforma nativa
# Ejemplo con GitHub Actions:
strategy:
  matrix:
    os: [ubuntu-latest, windows-latest, macos-latest]
```

---

## Roadmap

### Corto Plazo (MVP - Linux)
- [x] Arquitectura multi-plataforma
- [x] Implementación completa para Linux
- [ ] Resolver problema de autenticación OTP
- [ ] Packaging básico (.deb)

### Mediano Plazo
- [ ] Implementación completa para Windows
- [ ] Implementación completa para macOS
- [ ] Auto-update system
- [ ] Instaladores nativos

### Largo Plazo
- [ ] Soporte para múltiples perfiles VPN
- [ ] Recordar usuario (keyring integration)
- [ ] Auto-reconexión
- [ ] Reglas polkit/UAC sin prompt

---

## Referencias

### Documentación Técnica
- [OpenVPN Management Interface](https://openvpn.net/community-resources/management-interface/)
- [Go Build Tags](https://pkg.go.dev/cmd/go#hdr-Build_constraints)
- [Fyne Cross-Platform](https://developer.fyne.io/started/)

### Herramientas de Packaging
- **Linux**: [fpm](https://github.com/jordansissel/fpm), AppImageKit
- **Windows**: [NSIS](https://nsis.sourceforge.io/), [WiX](https://wixtoolset.org/)
- **macOS**: [create-dmg](https://github.com/create-dmg/create-dmg)

---

**Última revisión:** 2025-11-04
**Mantenedor:** Equipo PreyVPN
