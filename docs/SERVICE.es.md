# Gestión de Servicio del Sistema Operativo

Instala PicoClaw-Agents como un servicio del sistema que se inicia automáticamente al arrancar.

## Plataformas Compatibles

| Plataforma | Tipo de Servicio | Inicio automático | Requiere sudo |
|----------|-------------|------------|--------------|
| **Linux** | systemd (usuario) | Sí (al iniciar sesión) | No |
| **macOS** | launchd (LaunchAgent) | Sí (al iniciar sesión) | No |
| **Windows** | Tarea Programada (ONLOGON) | Sí (al iniciar sesión) | No |

> **Nota:** Los servicios de nivel de usuario se inician cuando inicias sesión. Para servicios de nivel de sistema (que inician antes de iniciar sesión), usa `sudo` con configuración de nivel de sistema (actualmente no soportado).

## Inicio Rápido

### Instalar

```bash
# Instalar con configuración predeterminada
picoclaw-agents service install

# Instalar en un puerto personalizado, escuchando en todas las interfaces
picoclaw-agents service install --port 9999 --public

# Previsualizar lo que se instalaría (dry run)
picoclaw-agents service install --dry-run
```

### Gestionar el Servicio

```bash
# Verificar estado
picoclaw-agents service status

# Ver logs
picoclaw-agents service logs              # Últimas 50 líneas
picoclaw-agents service logs -n 100       # Últimas 100 líneas
picoclaw-agents service logs -f           # Seguir en tiempo real

# Detener e iniciar
picoclaw-agents service stop
picoclaw-agents service start
picoclaw-agents service restart

# Desinstalar
picoclaw-agents service uninstall
```

## Opciones de Instalación

| Parámetro | Predeterminado | Descripción |
|------|---------|-------------|
| `--port` | `18800` | Número de puerto del Gateway |
| `--public` | `false` | Escuchar en todas las interfaces (0.0.0.0) en lugar de localhost |
| `--config` | `~/.picoclaw/config.json` | Ruta al archivo de configuración |
| `--dry-run` | `false` | Mostrar lo que se haría sin realizar cambios |

## Detalles por Plataforma

### Linux (systemd)

**Ubicación del archivo unit:** `~/.config/systemd/user/picoclaw-agents.service`

```bash
# Comandos manuales de systemd (si es necesario)
systemctl --user status picoclaw-agents.service
systemctl --user restart picoclaw-agents.service
journalctl --user -u picoclaw-agents.service -f
```

**Requisitos:** La sesión de usuario de systemd debe estar disponible. Verifica con:
```bash
systemctl --user is-system-running
```

### macOS (launchd)

**Ubicación del plist:** `~/Library/LaunchAgents/com.picoclaw.agents.gateway.plist`

**Archivos de log:**
- stdout: `~/Library/Logs/picoclaw-agents/gateway.out.log`
- stderr: `~/Library/Logs/picoclaw-agents/gateway.err.log`

```bash
# Comandos manuales de launchctl (si es necesario)
launchctl list | grep com.picoclaw.agents.gateway
launchctl bootout gui/$(id -u)/com.picoclaw.agents.gateway  # detener
launchctl bootstrap gui/$(id -u)/ ~/Library/LaunchAgents/com.picoclaw.agents.gateway.plist  # iniciar
```

**Nota:** Los servicios launchd solo funcionan en una sesión GUI. Si te conectas por SSH, el servicio no se iniciará automáticamente.

### Windows (Tarea Programada)

**Nombre de la tarea:** `PicoClaw-Agents Gateway`

**Archivos creados:**
- Script wrapper: `%APPDATA%\PicoClaw\gateway.cmd`
- Archivo PID: `%APPDATA%\PicoClaw\gateway.pid`
- Logs: `%APPDATA%\PicoClaw\logs\gateway.log`

```powershell
# Comandos manuales de schtasks (si es necesario)
schtasks /Query /TN "PicoClaw-Agents Gateway" /FO LIST
schtasks /End /TN "PicoClaw-Agents Gateway"    # detener
schtasks /Run /TN "PicoClaw-Agents Gateway"     # iniciar
schtasks /Delete /TN "PicoClaw-Agents Gateway" /F  # desinstalar
```

**Nota:** La tarea programada se ejecuta al iniciar sesión (ONLOGON). Para iniciar antes de iniciar sesión, usa `schtasks /SC ONSTART` manualmente.

## ⚠️ Advertencias Importantes

### Prevenir Conexiones Duplicadas

**El servicio y un gateway iniciado manualmente NO DEBEN ejecutarse simultáneamente.** Si ambos están activos:

- **Telegram:** Dos conexiones al mismo token del bot → mensajes procesados dos veces, límites de tasa de la API, posibles baneos
- **Discord:** Dos instancias del bot → conflictos de intents, mensajes duplicados
- **Conflicto de puerto:** Ambos procesos intentan vincularse al mismo puerto

El instalador del servicio verifica si el puerto del gateway ya está en uso y se negará a instalar si detecta un gateway en ejecución.

**Si ves este error:**
```
Error: gateway is already running on port 18800.
```

**Solución:** Detén primero el gateway manual:
```bash
# Si se ejecuta por CLI, presiona Ctrl+C
# O encuentra y mata el proceso:
picoclaw-agents gateway stop   # si está disponible
# O:
kill <PID>
```

### Servicio vs Gateway Manual

| | Servicio | Manual (`picoclaw-agents gateway`) |
|--|---------|-------------------------------------|
| Inicio automático | ✅ Al iniciar sesión | ❌ No |
| Sobrevive al cierre de sesión | ✅ Sí (macOS/Linux) | ❌ No |
| Logs a archivo | ✅ Sí | ❌ Solo stdout |
| Reinicio tras fallo | ✅ Sí | ❌ No |
| Puede ejecutarse simultáneamente | ❌ NO | ❌ NO |

**Regla:** Usa UNO u otro, nunca ambos al mismo tiempo.

## Solución de Problemas

### El servicio no se inicia

1. Verifica si el puerto está en uso:
   ```bash
   lsof -i :18800    # macOS/Linux
   netstat -ano | findstr 18800  # Windows
   ```

2. Verifica que la ruta del binario no haya cambiado:
   ```bash
   picoclaw-agents service status
   ```

3. Verifica los logs de error:
   ```bash
   # Linux
   journalctl --user -u picoclaw-agents.service -n 50

   # macOS
   tail -50 ~/Library/Logs/picoclaw-agents/gateway.err.log

   # Windows
   type %APPDATA%\PicoClaw\logs\gateway.log
   ```

### Servicio instalado pero no en ejecución

```bash
# Intenta iniciarlo
picoclaw-agents service start

# Verifica el estado
picoclaw-agents service status
```

### PicoClaw movido o reinstalado

Si moviste el binario a una ruta diferente después de instalar el servicio, el servicio apuntará a la ruta antigua. Desinstala y vuelve a instalar:

```bash
picoclaw-agents service uninstall
picoclaw-agents service install
```
