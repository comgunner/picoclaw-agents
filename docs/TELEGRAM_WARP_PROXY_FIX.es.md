# Solución de Conexión a Telegram - Proxy Cloudflare WARP

**Última Actualización:** 31 de marzo de 2026  
**Estado:** ✅ Listo para Producción  
**Severidad:** CRÍTICA - Bot de Telegram completamente inoperable

---

## Resumen

Esta guía resuelve la **inoperabilidad completa del bot de Telegram** causada por bloqueo a nivel de ISP o red de los servidores de API de Telegram. La solución utiliza un túnel **Cloudflare WARP** (basado en WireGuard) con un **proxy SOCKS5 local** para bypassar restricciones de red.

**Problema:** El bot de Telegram muestra errores de timeout al conectar con los servidores de API de Telegram.

**Solución:** Enrutar el tráfico de Telegram a través de un túnel encriptado de Cloudflare WARP con proxy SOCKS5 local.

---

## Síntomas

### Logs de Error

```
ERROR Execution error getUpdates: 
  request call: fasthttp do request: error when dialing 149.154.166.110:443: 
  dialing to the given TCP address timed out

ERROR Error on get me: 
  telego: getMe: internal execution: request call: fasthttp do request: 
  error when dialing 149.154.166.110:443: 
  dialing to the given TCP address timed out
```

### Componentes Afectados

| Componente | Estado | Notas |
|-----------|--------|-------|
| **Bot de Telegram** | ❌ FALLIDO | Timeout de conexión a `149.154.166.110:443` |
| **Bot de Discord** | ✅ Funcionando | Mismo gateway, diferente librería |
| **Web UI** | ✅ Funcionando | Configuración de modelo funcional |
| **CLI** | ✅ Funcionando | Comandos funcionales |

### Causa Raíz

**Bloqueo de Red/ISP** - El sistema no puede conectar a los servidores de API de Telegram:

- **IP del Servidor:** `149.154.166.110` (API de Telegram)
- **Puerto:** `443` (HTTPS)
- **Error:** `dialing to the given TCP address timed out`

**Evidencia:**
1. Discord funciona → Gateway está funcional
2. Web UI funciona → Configuración de modelo es correcta
3. CLI funciona → Configuración es correcta
4. Telegram falla → **Bloqueo específico de red hacia API de Telegram**

---

## Arquitectura de la Solución

### Resumen

```
┌─────────────────────────────────────────────────────────────┐
│                    PicoClaw-Agents                          │
│                                                             │
│  Bot de Telegram ──→ Proxy SOCKS5 ──→ Túnel WARP ──→ Internet │
│                        (socat)         (Cloudflare)          │
│                   127.0.0.1:40000    Encriptado              │
└─────────────────────────────────────────────────────────────┘
```

### Beneficios

- ✅ **Bypassa bloqueo de red** - Tráfico enrutado vía Cloudflare
- ✅ **Encriptado** - Cifrado vía WireGuard
- ✅ **No disruptivo** - No afecta otros servicios
- ✅ **Persistente** - Sobrevive a reinicios (con systemd)
- ✅ **Proxy local** - Fácil de configurar en aplicaciones

---

## Guía de Instalación

### Instrucciones por Plataforma

Selecciona tu sistema operativo:

- [**Linux (Ubuntu/Debian)**](#linux-ubuntudebian)
- [**macOS**](#macos)
- [**Windows**](#windows)

---

### Linux (Ubuntu/Debian)

#### Paso 1: Instalar el Cliente Cloudflare WARP

```bash
# Actualizar lista de paquetes
sudo apt update

# Instalar prerrequisitos
sudo apt install gnupg lsb-release -y

# Agregar clave GPG de Cloudflare
curl -fsSL https://pkg.cloudflareclient.com/pubkey.gpg | \
  sudo gpg --yes --dearmor --output /usr/share/keyrings/cloudflare-warp-archive-keyring.gpg

# Agregar repositorio de Cloudflare
echo "deb [signed-by=/usr/share/keyrings/cloudflare-warp-archive-keyring.gpg] \
  https://pkg.cloudflareclient.com/ $(lsb_release -cs) main" | \
  sudo tee /etc/apt/sources.list.d/cloudflare-client.list

# Instalar cliente WARP
sudo apt update && sudo apt install cloudflare-warp -y
```

#### Paso 2: Registrar Dispositivo con Cloudflare

```bash
# Registrar tu dispositivo (output: "Success")
warp-cli registration new
```

#### Paso 3: Configurar Modo Proxy

```bash
# Configurar en modo proxy (crítico para seguridad del servidor)
warp-cli mode proxy

# Esto crea un proxy SOCKS5 local en 127.0.0.1:40000
```

#### Paso 4: Conectar a WARP

```bash
# Habilitar el túnel
warp-cli connect

# Verificar estado de conexión
warp-cli status

# Output esperado: "Status: Connected"
```

#### Paso 5: Instalar y Configurar socat

```bash
# Instalar socat
sudo apt install socat -y

# Crear puente de proxy SOCKS5
# Esto reenvía externo:40001 → interno:127.0.0.1:40000
socat TCP-LISTEN:40001,fork,bind=0.0.0.0 TCP:127.0.0.1:40000 &

# Hacerlo persistente (agregar a /etc/rc.local o crear servicio systemd)
```

---

### macOS

#### Paso 1: Instalar Cloudflare WARP

**Opción A: Usando Homebrew (Recomendado)**

```bash
# Instalar Homebrew si no está instalado
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"

# Instalar Cloudflare WARP
brew install --cask cloudflare-warp
```

**Opción B: Descarga Directa**

1. Descargar desde: https://1.1.1.1/
2. Abrir el archivo `.dmg`
3. Arrastrar "Cloudflare WARP" a Aplicaciones
4. Abrir desde la carpeta Aplicaciones

#### Paso 2: Registrar Dispositivo con Cloudflare

```bash
# Abrir Terminal y registrar
warp-cli registration new

# Output esperado: "Success"
```

#### Paso 3: Configurar Modo Proxy

```bash
# Configurar en modo proxy
warp-cli mode proxy

# Esto crea un proxy SOCKS5 local en 127.0.0.1:40000
```

#### Paso 4: Conectar a WARP

```bash
# Habilitar el túnel
warp-cli connect

# Verificar estado de conexión
warp-cli status

# Output esperado: "Status: Connected"
```

#### Paso 5: Instalar y Configurar socat

```bash
# Instalar socat usando Homebrew
brew install socat

# Crear puente de proxy SOCKS5
socat TCP-LISTEN:40001,fork,bind=0.0.0.0 TCP:127.0.0.1:40000 &

# Hacerlo persistente (usando LaunchAgents)
# Ver sección de Persistencia abajo
```

#### Persistencia en macOS (LaunchAgents)

```bash
# Crear LaunchAgent para socat
cat > ~/Library/LaunchAgents/com.user.socat-proxy.plist << 'EOF'
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>com.user.socat-proxy</string>
    <key>ProgramArguments</key>
    <array>
        <string>/usr/local/bin/socat</string>
        <string>TCP-LISTEN:40001,fork,bind=0.0.0.0</string>
        <string>TCP:127.0.0.1:40000</string>
    </array>
    <key>RunAtLoad</key>
    <true/>
    <key>KeepAlive</key>
    <true/>
</dict>
</plist>
EOF

# Cargar el LaunchAgent
launchctl load ~/Library/LaunchAgents/com.user.socat-proxy.plist

# Verificar que está corriendo
launchctl list | grep socat
```

---

### Windows

#### Paso 1: Instalar Cloudflare WARP

**Opción A: Usando Winget (Recomendado)**

```powershell
# Instalar Cloudflare WARP usando winget
winget install Cloudflare.CloudflareWARP

# O usando Chocolatey
choco install cloudflare-warp
```

**Opción B: Descarga Directa**

1. Descargar desde: https://1.1.1.1/
2. Ejecutar el instalador (`Cloudflare_WARP_setup.exe`)
3. Seguir el asistente de instalación
4. Lanzar Cloudflare WARP desde el menú Inicio

#### Paso 2: Registrar Dispositivo con Cloudflare

```powershell
# Abrir PowerShell como Administrador
cd "C:\Program Files\Cloudflare\Cloudflare WARP"

# Registrar tu dispositivo
.\warp-cli.exe registration new

# Output esperado: "Success"
```

#### Paso 3: Configurar Modo Proxy

```powershell
# Configurar en modo proxy
.\warp-cli.exe mode proxy

# Esto crea un proxy SOCKS5 local en 127.0.0.1:40000
```

#### Paso 4: Conectar a WARP

```powershell
# Habilitar el túnel
.\warp-cli.exe connect

# Verificar estado de conexión
.\warp-cli.exe status

# Output esperado: "Status: Connected"
```

#### Paso 5: Instalar y Configurar socat para Windows

**Opción A: Usando WSL (Recomendado)**

```bash
# Instalar WSL si no está instalado
wsl --install

# En terminal de WSL, instalar socat
sudo apt update && sudo apt install socat -y

# Crear puente de proxy SOCKS5
socat TCP-LISTEN:40001,fork,bind=0.0.0.0 TCP:127.0.0.1:40000 &
```

**Opción B: Usando Cygwin**

```bash
# Instalar Cygwin con paquete socat
# https://www.cygwin.com/

# En terminal de Cygwin
socat TCP-LISTEN:40001,fork,bind=0.0.0.0 TCP:127.0.0.1:40000 &
```

**Opción C: Usando Alternativa (netsh)**

```powershell
# Crear proxy de puerto usando netsh (incluido en Windows)
netsh interface portproxy add v4tov4 listenport=40001 listenaddress=0.0.0.0 connectport=40000 connectaddress=127.0.0.0

# Verificar proxy
netsh interface portproxy show all

# Eliminar proxy (si es necesario)
netsh interface portproxy delete v4tov4 listenport=40001 listenaddress=0.0.0.0
```

#### Persistencia en Windows

**Para WARP:**
- WARP para Windows se auto-inicia por defecto
- Verificar en Administrador de Tareas → Pestaña Inicio

**Para socat (WSL):**

```bash
# Agregar al bashrc de WSL
echo "socat TCP-LISTEN:40001,fork,bind=0.0.0.0 TCP:127.0.0.1:40000 &" >> ~/.bashrc

# O crear servicio systemd en WSL
```

**Para proxy netsh:**

```powershell
# Crear tarea programada para proxy de puerto
$action = New-ScheduledTaskAction -Execute "netsh" -Argument "interface portproxy add v4tov4 listenport=40001 listenaddress=0.0.0.0 connectport=40000 connectaddress=127.0.0.0"
$trigger = New-ScheduledTaskTrigger -AtStartup
Register-ScheduledTask -TaskName "WARP Proxy Bridge" -Action $action -Trigger $trigger -RunLevel Highest
```

---

## Configuración

### Configuración de PicoClaw-Agents

Editar `~/.picoclaw/config.json`:

**Opción A: Proxy Local (localhost)**

```json
{
  "channels": {
    "telegram": {
      "enabled": true,
      "token": "YOUR_TELEGRAM_BOT_TOKEN",
      "proxy": "socks5://127.0.0.1:40000",
      "allow_from": [
        "YOUR_USER_ID"
      ]
    }
  }
}
```

**Opción B: Proxy Remoto (otra IP en la red)**

```json
{
  "channels": {
    "telegram": {
      "enabled": true,
      "token": "YOUR_TELEGRAM_BOT_TOKEN",
      "proxy": "socks5://192.168.0.100:40001",
      "allow_from": [
        "YOUR_USER_ID"
      ]
    }
  }
}
```

**Notas de Configuración:**

| Setting | Valor | Descripción |
|---------|-------|-------------|
| **Host** | `127.0.0.1` | Proxy local (WARP + socat en el mismo servidor) |
| **Port** | `40000` | Puerto por defecto del proxy WARP |
| **Protocol** | `SOCKS5` | Protocolo del proxy |
| **Full URL** | `socks5://127.0.0.1:40000` | URL completa del proxy |

**Importante:**
- **Proxy local:** Usa `127.0.0.1:40000` cuando WARP + socat corren en el mismo servidor
- **Proxy remoto:** Usa `192.168.x.x:40001` cuando el proxy corre en otro equipo de la red
- **Token:** Reemplaza `YOUR_TELEGRAM_BOT_TOKEN` con tu token real de @BotFather
- **User ID:** Reemplaza `YOUR_USER_ID` con tu ID numérico de Telegram

### Alternativa: Variables de Entorno

```bash
# Agregar a ~/.bashrc o ~/.zshrc
export ALL_PROXY=socks5://127.0.0.1:40000
export HTTPS_PROXY=socks5://127.0.0.1:40000

# O temporal para una sesión
export ALL_PROXY=socks5://127.0.0.1:40000
./build/picoclaw-agents-launcher
```

---

## Verificación

### Paso 1: Test del Proxy SOCKS5

```bash
# Test de conectividad del proxy
curl -v -x socks5://127.0.0.1:40000 https://www.google.com

# Esperado: Conexión exitosa vía proxy
```

### Paso 2: Test de API de Telegram

```bash
# Reemplazar YOUR_BOT_TOKEN con tu token real
curl -v -x socks5://127.0.0.1:40000 \
  "https://api.telegram.org/botYOUR_BOT_TOKEN/getMe"

# Output esperado:
# {"ok":true,"result":{"id":123456789,"is_bot":true,"first_name":"...","username":"..."}}
```

### Paso 3: Verificar Logs de la Aplicación

```bash
# Reiniciar aplicación
pkill -f picoclaw
./build/picoclaw-agents-launcher &

# Esperar 15 segundos
sleep 15

# Verificar logs de Telegram
tail -100 ~/.picoclaw/logs/launcher.log | grep -i telegram

# Output esperado:
# [INFO] telegram: Telegram bot connected {username=your_bot_username}
```

### Paso 4: Test del Bot en Telegram

1. Abrir Telegram
2. Buscar tu bot
3. Enviar `/start`
4. El bot debería responder inmediatamente

---

## Comparativa de Rendimiento

| Métrica | Antes (Bloqueado) | Después (WARP) |
|--------|-----------------|--------------|
| **API de Telegram** | ❌ Timeout | ✅ < 500ms |
| **Respuesta del Bot** | ❌ Sin respuesta | ✅ Instantáneo |
| **Conexión** | ❌ Bloqueado | ✅ Enrutado vía Cloudflare |
| **Otros Servicios** | ✅ Funcionando | ✅ Still working |

---

## Troubleshooting

### Problema: WARP No Conecta

```bash
# Verificar estado de registro
warp-cli registration status

# Si no está registrado, registrar de nuevo
warp-cli registration new

# Verificar firewall
sudo ufw status
sudo ufw allow out to any port 2408 proto udp
```

### Problema: Proxy No Funciona

```bash
# Verificar si WARP está corriendo
warp-cli status

# Verificar si el puerto 40000 está escuchando
netstat -tlnp | grep 40000

# Reiniciar WARP
warp-cli disconnect
warp-cli connect
```

### Problema: socat No Reenvía

```bash
# Verificar si socat está corriendo
ps aux | grep socat

# Reiniciar socat
pkill socat
socat TCP-LISTEN:40001,fork,bind=0.0.0.0 TCP:127.0.0.1:40000 &

# Verificar si el puerto está escuchando
netstat -tlnp | grep 40001
```

### Problema: Aún Hay Timeouts

```bash
# Test de conexión directa (debería fallar)
curl -v https://api.telegram.org/botYOUR_TOKEN/getMe

# Test vía proxy (debería funcionar)
curl -v -x socks5://127.0.0.1:40000 https://api.telegram.org/botYOUR_TOKEN/getMe

# Si el test del proxy falla, verificar estado de WARP
warp-cli status
warp-cli disconnect
warp-cli connect
```

---

## Consideraciones de Seguridad

### Lo que Esto Hace

- ✅ **Encripta tráfico** vía WireGuard
- ✅ **Enruta vía Cloudflare** network
- ✅ **Bypassa bloqueo del ISP**
- ✅ **Proxy local** (solo accesible desde localhost)

### Lo que Esto NO Hace

- ❌ **No anonimiza** tu identidad
- ❌ **No oculta** el token del bot
- ❌ **No bypassa** límites de la API de Telegram
- ❌ **No previene** rate limits

### Mejores Prácticas

1. **Mantener proxy local** - No exponer a la red pública
2. **Usar firewall** - Bloquear acceso externo al puerto 40000
3. **Monitorear logs** - Vigilar actividad inusual
4. **Actualizar regularmente** - Mantener el cliente WARP actualizado

```bash
# Regla de firewall para bloquear acceso externo
sudo ufw deny from any to any port 40000
sudo ufw allow from 127.0.0.1 to any port 40000
```

---

## Persistencia (Auto-start en Boot)

### Opción 1: Servicio systemd para WARP

```bash
# Crear servicio systemd
sudo nano /etc/systemd/system/warp-client.service

# Agregar contenido:
[Unit]
Description=Cloudflare WARP Client
After=network.target

[Service]
Type=simple
ExecStart=/opt/cloudflare-warp/warp-client
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target

# Habilitar e iniciar
sudo systemctl enable warp-client
sudo systemctl start warp-client
```

### Opción 2: Servicio systemd para socat

```bash
# Crear servicio systemd
sudo nano /etc/systemd/system/socks5-proxy.service

# Agregar contenido:
[Unit]
Description=SOCKS5 Proxy Bridge for WARP
After=warp-client.service
Requires=warp-client.service

[Service]
Type=simple
ExecStart=/usr/bin/socat TCP-LISTEN:40001,fork,bind=0.0.0.0 TCP:127.0.0.1:40000
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target

# Habilitar e iniciar
sudo systemctl enable socks5-proxy
sudo systemctl start socks5-proxy
```

---

## Referencias

### Documentación Oficial

- **Cloudflare WARP:** https://developers.cloudflare.com/cloudflare-one/connections/connect-devices/warp/
- **Telegram Bot API:** https://core.telegram.org/bots/api
- **Telegram Proxy:** https://core.telegram.org/api/proxy

### Problemas Relacionados

- **Bloqueo de API de Telegram:** Común en algunas regiones/ISPs
- **Cloudflare WARP:** Tier gratuito disponible, pago para equipos
- **Proxy SOCKS5:** Protocolo estándar para conexiones proxy

---

## Resumen

### Problema
- Servidores de API de Telegram bloqueados/con timeout
- El bot no podía conectarse a Telegram
- Otros servicios (Discord, Web UI, CLI) funcionaban correctamente

### Solución
- Instalar Cloudflare WARP (túnel WireGuard)
- Configurar en modo proxy (127.0.0.1:40000)
- Usar socat para puentear puertos si es necesario
- Configurar canal de Telegram para usar proxy SOCKS5

### Resultado
- ✅ Bot de Telegram conectado exitosamente
- ✅ Todos los comandos funcionales
- ✅ Sin impacto en otros servicios
- ✅ Conexión encriptada y segura

---

**Autor:** @comgunner  
**Repositorio:** https://github.com/comgunner/picoclaw-agents  
**Última Actualización:** 31 de marzo de 2026
