# Configuración para Windows

## Archivos necesarios

### preyvpn.ico
Icono de la aplicación (formato .ico)

- Resoluciones recomendadas: 16x16, 32x32, 48x48, 256x256
- Usar herramientas como ImageMagick para convertir desde PNG:
  ```bash
  magick convert icon.png -define icon:auto-resize=256,128,64,48,32,16 preyvpn.ico
  ```

### preyvpn.manifest
Manifest para configurar UAC (User Account Control)

Ejemplo básico:
```xml
<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<assembly xmlns="urn:schemas-microsoft-com:asm.v1" manifestVersion="1.0">
  <assemblyIdentity
    version="1.0.0.0"
    processorArchitecture="*"
    name="PreyVPN"
    type="win32"
  />
  <description>PreyVPN Client</description>
  <trustInfo xmlns="urn:schemas-microsoft-com:asm.v3">
    <security>
      <requestedPrivileges>
        <requestedExecutionLevel level="asInvoker" uiAccess="false"/>
      </requestedPrivileges>
    </security>
  </trustInfo>
</assembly>
```

## Instalador (futuro)

### NSIS (Nullsoft Scriptable Install System)
- Script de instalador: `installer.nsi`
- Comandos para compilar:
  ```cmd
  makensis installer.nsi
  ```

### WiX Toolset
- Para crear instaladores .msi más robustos
- Soporte para upgrades, patches, etc.

## Registry

Configuración en el registro de Windows para auto-start (opcional):
```
HKEY_CURRENT_USER\Software\Microsoft\Windows\CurrentVersion\Run
  PreyVPN = "C:\Program Files\PreyVPN\preyvpn.exe"
```
