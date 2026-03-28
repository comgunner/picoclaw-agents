# Scripts de PicoClaw-Agents

**Fecha:** 2026-03-27  
**Ubicación:** `scripts/`

---

## 📋 Descripción General

Este directorio contiene scripts de utilidad para el desarrollo, testing y despliegue de PicoClaw-Agents.

---

## 🔧 Scripts Disponibles

### 1. `build-macos-app.sh` ⭐

**Propósito:** Crea un bundle `.app` de macOS para el PicoClaw-Agents Launcher.

**Uso:**
```bash
./scripts/build-macos-app.sh <executable>
```

**Ejemplo:**
```bash
# Primero construye el launcher web
make build-launcher

# Luego crea el .app bundle
./scripts/build-macos-app.sh ./build/picoclaw-agents-launcher
```

**Output:**
- `./build/PicoClaw-Agents Launcher.app`

**Características:**
- Crea estructura de bundle macOS estándar
- Incluye ícono personalizado (`icon.icns`)
- Configura Info.plist con identificadores únicos
- Copia binarios: `picoclaw-agents-launcher` y `picoclaw-agents`
- Permite ejecución desde Finder o terminal (`open`)

**Requisitos:**
- macOS
- Binario `picoclaw-agents-launcher` construido en `web/build/`
- Binario `picoclaw-agents` construido en `build/`

---

### 2. `test-irc.sh`

**Propósito:** Inicia un servidor IRC Ergo local para testing del canal IRC.

**Uso:**
```bash
./scripts/test-irc.sh
```

**Output:**
- Contenedor Docker `picoclaw-test-ergo`
- Servidor IRC escuchando en `localhost:6667`

**Configuración Generada:**
```json
{
  "irc": {
    "enabled": true,
    "server": "localhost:6667",
    "tls": false,
    "nick": "picobot",
    "channels": ["#test"],
    "allow_from": [],
    "group_trigger": { "mention_only": true }
  }
}
```

**Requisitos:**
- Docker instalado
- Puerto 6667 disponible

**Para Detener:**
```bash
docker rm -f picoclaw-test-ergo
```

---

### 3. `test-docker-mcp.sh`

**Propósito:** Testea herramientas MCP en contenedor Docker (imagen full-featured).

**Uso:**
```bash
./scripts/test-docker-mcp.sh
```

**Tests Realizados:**
- ✅ `npx --version`
- ✅ `npm --version`
- ✅ `node --version`
- ✅ `git --version`
- ✅ `python3 --version`
- ✅ `uv --version`
- ✅ MCP server install (`@modelcontextprotocol/server-filesystem`)

**Requisitos:**
- Docker Compose
- Archivo `docker/docker-compose.full.yml`

---

### 4. `create-release.sh`

**Propósito:** Automatiza la creación de releases de PicoClaw-Agents.

**Uso:**
```bash
./scripts/create-release.sh
```

**Características:**
- Construye binarios para múltiples plataformas
- Genera assets para GitHub Releases
- Actualiza changelog

**Requisitos:**
- GoReleaser instalado
- Git tags configurados

---

### 5. `health_check.sh`

**Propósito:** Verifica el estado de salud del sistema y dependencias.

**Uso:**
```bash
./scripts/health_check.sh
```

**Verificaciones:**
- Go version
- Node.js version
- Docker status
- Configuración existente
- Workspace permissions

---

### 6. `update-changelog.sh`

**Propósito:** Actualiza el CHANGELOG.md automáticamente.

**Uso:**
```bash
./scripts/update-changelog.sh
```

**Características:**
- Extrae commits de Git
- Genera entradas de changelog
- Mantiene formato estándar

---

### 7. `update-readmes.sh`

**Propósito:** Actualiza los READMEs en múltiples idiomas.

**Uso:**
```bash
./scripts/update-readmes.sh
```

**Características:**
- Sincroniza cambios entre README.md y traducciones
- Mantiene consistencia de versiones
- Actualiza badges y status

---

## 📁 Archivos Adicionales

### `icon.icns`

**Propósito:** Ícono de la aplicación para macOS.

**Ubicación:** `scripts/icon.icns`

**Uso:** Utilizado por `build-macos-app.sh` para el bundle `.app`.

**Especificaciones:**
- Formato: Apple ICNS
- Incluye múltiples resoluciones (16x16 a 1024x1024)
- Compatible con macOS High Resolution

---

### `setup.iss`

**Propósito:** Script de instalación para Windows (Inno Setup).

**Uso:**
```bash
iscc scripts/setup.iss
```

**Output:**
- Instalador `.exe` para Windows
- Registro en PATH
- Accesos directos

---

## 🔐 Permisos

Todos los scripts `.sh` deben ser ejecutables:

```bash
chmod +x scripts/*.sh
```

**Verificación:**
```bash
ls -la scripts/
# -rwxr-xr-x para todos los .sh
```

---

## 🛠️ Desarrollo

### Agregar Nuevo Script

1. Crear archivo en `scripts/`
2. Agregar shebang: `#!/bin/bash` o `#!/bin/sh`
3. Hacer ejecutable: `chmod +x scripts/nuevo-script.sh`
4. Documentar en este README

### Convenciones

- **Shell:** Preferir `#!/bin/sh` para portabilidad
- **Errores:** Usar `set -e` para salir en error
- **Logging:** Usar `echo` con emojis para claridad (🧪, ✅, ❌, 📦)
- **Comentarios:** Documentar propósito y requisitos

---

## 📝 Historial de Cambios

### 2026-03-27

**Agregados:**
- ✅ `build-macos-app.sh` - Adaptado para `picoclaw-agents-launcher`
- ✅ `test-irc.sh` - Actualizado para usar `picoclaw-agents`
- ✅ `test-docker-mcp.sh` - Copiado del original
- ✅ `icon.icns` - Ícono de la aplicación
- ✅ `setup.iss` - Script de instalación Windows

**Modificaciones:**
- `build-macos-app.sh`: 
  - `APP_NAME` → `"PicoClaw-Agents Launcher"`
  - `APP_EXECUTABLE` → `"picoclaw-agents-launcher"`
  - Binario copiado: `./build/picoclaw-agents`
  - Info.plist actualizado con identificadores `com.picoclaw-agents`

- `test-irc.sh`:
  - Comando actualizado: `./build/picoclaw-agents gateway`
  - Directorio: `cd picoclaw-agents`

---

## 🚨 Solución de Problemas

### `build-macos-app.sh` falla

**Error:** `Error: ./web/build/picoclaw-agents-launcher not found`

**Solución:**
```bash
cd web/frontend && pnpm install && pnpm build:backend
cd ../backend && CGO_ENABLED=1 go build -o ../../web/build/picoclaw-agents-launcher .
```

### `test-irc.sh` no puede conectar

**Error:** `ERROR: Server did not start within 10s`

**Solución:**
```bash
# Verificar Docker
docker ps

# Verificar puerto
nc -zv localhost 6667

# Limpiar contenedor anterior
docker rm -f picoclaw-test-ergo
```

---

## 📞 Soporte

Para issues con los scripts:
1. Verificar permisos: `ls -la scripts/`
2. Revisar logs de error completos
3. Verificar dependencias (Docker, Go, Node.js)

---

*Documentación actualizada: 2026-03-27*  
*Scripts adaptados para: PicoClaw-Agents v3.4.5*
