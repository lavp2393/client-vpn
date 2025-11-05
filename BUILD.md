# Guía de Compilación - PreyVPN

Esta guía explica cómo compilar PreyVPN **sin necesidad de instalar Go ni dependencias** en tu máquina.

## 🎯 Compilación con Docker (Recomendado)

**Ventajas:**
- ✅ NO requiere instalar Go
- ✅ NO requiere instalar dependencias de Fyne (libgl, xorg-dev, etc.)
- ✅ Entorno reproducible
- ✅ Funciona en cualquier máquina con Docker

### Requisito Único

Solo necesitas **Docker** instalado:

```bash
# Verificar que Docker está instalado
docker --version
```

Si no tienes Docker: https://docs.docker.com/get-docker/

---

## 📦 Opción 1: Compilar con Taskfile

Si tienes [Task](https://taskfile.dev/installation/) instalado:

```bash
# Compilar binario para desarrollo
task build-docker

# O compilar versión optimizada para distribución
task build-docker-release
```

El binario estará en `./dist/preyvpn`

---

## 📦 Opción 2: Compilar con script dev.sh

```bash
# Compilar binario
./dev.sh build-binary
```

El binario estará en `./dist/preyvpn`

---

## 📦 Opción 3: Compilar con Docker directamente

```bash
# 1. Crear directorio de salida
mkdir -p dist

# 2. Construir imagen de compilación
docker build -f Dockerfile.build -t preyvpn-builder --target builder .

# 3. Compilar y extraer binario
docker run --rm -v $(pwd)/dist:/output preyvpn-builder \
    sh -c "cp /build/preyvpn /output/ && chmod +x /output/preyvpn"

# 4. Verificar el binario
ls -lh dist/preyvpn
file dist/preyvpn
```

---

## ⏱️ Tiempos de Compilación

| Acción | Primera vez | Siguientes veces |
|--------|-------------|------------------|
| Construir imagen | ~5-7 min | ~10 seg (cache) |
| Compilar binario | ~3-5 min | ~10 seg (cache) |
| **Total** | **~8-12 min** | **~20 seg** |

**Nota:** La primera vez toma más tiempo porque Docker descarga las imágenes base y compila todas las dependencias. Las siguientes compilaciones son **mucho más rápidas** gracias al cache de Docker.

---

## 🚀 Ejecutar el Binario Compilado

```bash
# Verificar que existe
ls -lh dist/preyvpn

# Ejecutar
./dist/preyvpn
```

**Requisitos para ejecutar:**
- OpenVPN instalado: `sudo apt install openvpn`
- Archivo de configuración en: `~/PreyVPN/prey-prod.ovpn`

---

## 🔧 Compilación para Múltiples Plataformas

### Linux (nativo)

```bash
# AMD64 (Intel/AMD de 64 bits)
task build-docker

# ARM64 (Raspberry Pi 4, servidores ARM)
docker build -f Dockerfile.build -t preyvpn-builder \
    --build-arg GOARCH=arm64 --target builder .
docker run --rm -v $(pwd)/dist:/output preyvpn-builder \
    sh -c "cp /build/preyvpn /output/preyvpn-arm64 && chmod +x /output/preyvpn-arm64"
```

### Windows (cross-compilation desde Linux)

```bash
# Requiere mingw-w64 en la imagen
docker build -f Dockerfile.build -t preyvpn-builder-windows \
    --build-arg GOOS=windows --build-arg GOARCH=amd64 --target builder .
docker run --rm -v $(pwd)/dist:/output preyvpn-builder-windows \
    sh -c "cp /build/preyvpn.exe /output/"
```

### macOS (cross-compilation desde Linux)

```bash
# Requiere osxcross en la imagen
docker build -f Dockerfile.build -t preyvpn-builder-darwin \
    --build-arg GOOS=darwin --build-arg GOARCH=amd64 --target builder .
docker run --rm -v $(pwd)/dist:/output preyvpn-builder-darwin \
    sh -c "cp /build/preyvpn /output/preyvpn-darwin"
```

---

## 🐛 Solución de Problemas

### Error: "Cannot connect to the Docker daemon"

```bash
# Verificar que Docker está corriendo
sudo systemctl start docker

# O en macOS/Windows
# Abrir Docker Desktop
```

### Error: "permission denied" al ejecutar el binario

```bash
chmod +x dist/preyvpn
```

### El binario no se creó

```bash
# Ver logs de compilación
docker build -f Dockerfile.build -t preyvpn-builder --target builder . 2>&1 | tee build.log
```

### Limpiar cache de Docker

Si necesitas recompilar desde cero:

```bash
# Limpiar cache de build
docker builder prune -a

# O eliminar la imagen y reconstruir
docker rmi preyvpn-builder
task build-docker
```

---

## 📊 Comparación: Docker vs Local

| Aspecto | Compilación Docker | Compilación Local |
|---------|-------------------|-------------------|
| **Instalación Go** | ❌ No requerido | ✅ Requerido |
| **Dependencias** | ❌ No requerido | ✅ Requerido |
| **Primera compilación** | ~8-12 min | ~5-7 min |
| **Siguientes compilaciones** | ~20 seg | ~10 seg |
| **Reproducibilidad** | ✅ 100% | ⚠️ Depende del entorno |
| **Tamaño del binario** | ~27 MB | ~27 MB |

---

## 💡 Tips

1. **Cache de Docker**: La primera compilación toma tiempo, pero las siguientes son rápidas gracias al cache de layers.

2. **Compilar en background**:
   ```bash
   task build-docker > build.log 2>&1 &
   tail -f build.log
   ```

3. **Verificar el binario**:
   ```bash
   # Ver información del archivo
   file dist/preyvpn

   # Ver tamaño
   ls -lh dist/preyvpn

   # Ver dependencias dinámicas
   ldd dist/preyvpn
   ```

4. **Optimizar tamaño**:
   ```bash
   # Usar build-docker-release que incluye strip
   task build-docker-release

   # Reduce el binario de ~35MB a ~27MB
   ```

---

## 🔗 Recursos

- **Dockerfile.build**: Configuración del entorno de compilación
- **Taskfile.yml**: Comandos automatizados
- **dev.sh**: Script alternativo para compilación
- **DOCKER-README.md**: Documentación del entorno de desarrollo

---

## ❓ Preguntas Frecuentes

### ¿Puedo compilar sin Docker?

Sí, pero necesitarás instalar:
- Go 1.22+
- Dependencias de Fyne: `sudo apt install libgl1-mesa-dev xorg-dev`
- OpenVPN: `sudo apt install openvpn`

Ver [README.md](README.md) para instrucciones de compilación local.

### ¿El binario funciona en cualquier distro de Linux?

El binario está compilado para Linux genérico y debería funcionar en:
- Ubuntu 20.04+
- Debian 11+
- Fedora 35+
- Arch Linux
- Otras distros con glibc 2.31+

### ¿Puedo distribuir el binario compilado?

Sí, el binario en `dist/preyvpn` es autocontenido y puede distribuirse a otros usuarios de Linux. Solo necesitan tener OpenVPN instalado.

### ¿Cómo actualizar las dependencias?

```bash
# Actualizar go.mod
go get -u ./...
go mod tidy

# Reconstruir imagen sin cache
docker build --no-cache -f Dockerfile.build -t preyvpn-builder --target builder .
```

---

**¿Más preguntas?** Consulta [DOCKER-README.md](DOCKER-README.md) o abre un issue.
