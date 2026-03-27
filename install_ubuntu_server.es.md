# Instalación de PicoClaw-Agents en Ubuntu Server

Esta guía detalla los pasos necesarios para instalar PicoClaw-Agents en un entorno Ubuntu Server.

## Métodos de Instalación

**IMPORTANTE:** Existen dos métodos de instalación:

1. **Desde Releases de GitHub** (Recomendado si hay releases disponibles)
2. **Desde Código Fuente** (Si no hay releases o para desarrollo)

### Verificar si hay Releases Disponibles

Antes de comenzar, verifica si existen releases en:
```bash
curl -s https://api.github.com/repos/comgunner/picoclaw-agents/releases/latest | grep '"tag_name"'
# Si retorna algo como "v3.7.0", hay releases disponibles
# Si retorna "Not Found", usa el Método 2 (desde fuente)
```

---

## Método 1: Desde Releases de GitHub (Si están disponibles)

**Requisitos:**
- `wget` o `curl` para descargar

```bash
sudo apt update
sudo apt install -y wget curl
```

### Pasos de Instalación

1. **Crear directorio de instalación:**

   ```bash
   sudo mkdir -p /opt/picoclaw-agents
   sudo chown -R $USER:$USER /opt/picoclaw-agents
   cd /opt/picoclaw-agents
   ```

2. **Descargar el último release:**

   ```bash
   # Descargar tar.gz más reciente
   wget https://github.com/comgunner/picoclaw-agents/releases/latest/download/picoclaw-agents_Linux_arm64.tar.gz
   
   # O para versión específica (ej: v3.7.0)
   # wget https://github.com/comgunner/picoclaw-agents/releases/download/v3.7.0/picoclaw-agents_Linux_arm64.tar.gz
   ```

3. **Extraer y configurar:**

   ```bash
   # Extraer tar.gz
   tar -xzf picoclaw-agents_Linux_arm64.tar.gz
   
   # El binario se extrae como 'picoclaw-agents'
   chmod +x picoclaw-agents
   
   # Crear symlink para comando global
   sudo ln -sf /opt/picoclaw-agents/picoclaw-agents /usr/local/bin/picoclaw-agents
   ```

4. **Verificar instalación:**

   ```bash
   picoclaw-agents --version
   # o
   ./picoclaw-agents --help
   ```

---

## Método 2: Desde Código Fuente (Desarrollo o si no hay releases)

**Requisitos:**
- Go 1.25.7 o superior
- Git
- Make

### 1. Instalar Go

```bash
# Descargar e instalar Go 1.25.7+
wget https://go.dev/dl/go1.25.7.linux-arm64.tar.gz
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.25.7.linux-arm64.tar.gz
rm go1.25.7.linux-arm64.tar.gz

# Añadir al PATH permanentemente
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.profile
source ~/.profile

# Verificar instalación
go version
# Debe mostrar: go version go1.25.7 linux/arm64
```

### 2. Clonar y Compilar

```bash
# Clonar repositorio
cd /opt
git clone https://github.com/comgunner/picoclaw-agents.git
sudo chown -R $USER:$USER /opt/picoclaw-agents
cd picoclaw-agents

# Descargar dependencias
make deps

# Compilar binario
make build
# El binario se genera en: build/picoclaw-agents
```

### 3. Instalar Binario

```bash
# Copiar binario a directorio de instalación
sudo mkdir -p /opt/picoclaw-agents
cp build/picoclaw-agents /opt/picoclaw-agents/
chmod +x /opt/picoclaw-agents/picoclaw-agents

# Crear symlink para comando global
sudo ln -sf /opt/picoclaw-agents/picoclaw-agents /usr/local/bin/picoclaw-agents
```

### 4. Verificar Instalación

```bash
picoclaw-agents --version
# o
./picoclaw-agents --help
```

---

## Configuración del Servidor

El archivo principal de configuración (`config.json`) debe colocarse en el directorio home del usuario que ejecutará PicoClaw-Agents.

### 1. Ubicación en Ubuntu Server

```bash
~/.picoclaw/config.json
```
*(Esto equivale a `/home/TU_USUARIO/.picoclaw/config.json` o `/root/.picoclaw/config.json` si lo corres como root).*

### 2. Copiando tu configuración desde Mac/Local

Si ya configuraste tus bots (como Telegram) y tus API Keys localmente, puedes transferir tu archivo actual al servidor usando `scp`:

```bash
# Ejecuta esto desde tu Mac/Local:
scp ~/.picoclaw/config.json usuario@ip_de_tu_servidor:~/.picoclaw/config.json
```

*Nota: Es posible que necesites crear la carpeta oculta en el servidor antes de copiarlo (`mkdir -p ~/.picoclaw`).*

### 3. O usar el comando `onboard` para autogenerarlo

Si prefieres no copiar el archivo, puedes dejar que PicoClaw-Agents genere uno nuevo con los valores por defecto:

```bash
picoclaw-agents onboard
```

Esto creará el archivo `~/.picoclaw/config.json` con toda la estructura lista. Luego edita el archivo con `nano` o `vim` para agregar:
- Tu `api_key` de Deepseek/OpenAI/Anthropic en `model_list`
- Tu `token` de Telegram en `channels > telegram`
- Tu ID de usuario en `allow_from`

---

## Ejecutar como Servicio (Systemd)

Para que PicoClaw-Agents se ejecute en segundo plano y se reinicie automáticamente:

### 1. Crear el archivo del servicio

```bash
cat <<EOF | sudo tee /etc/systemd/system/picoclaw-agents.service
[Unit]
Description=PicoClaw-Agents Gateway Service
After=network.target

[Service]
# Directorio donde vive el binario
WorkingDirectory=/opt/picoclaw-agents

# Ruta completa al binario
ExecStart=/opt/picoclaw-agents/picoclaw-agents gateway

# El usuario que lo ejecuta
User=$USER
Group=$(id -gn)

# Reiniciar automáticamente si falla
Restart=always
RestartSec=5

# Logs en journal
StandardOutput=journal
StandardError=journal
SyslogIdentifier=picoclaw-agents

[Install]
WantedBy=multi-user.target
EOF
```

### 2. Habilitar e iniciar

```bash
sudo systemctl daemon-reload
sudo systemctl enable picoclaw-agents.service
sudo systemctl start picoclaw-agents.service
```

### 3. Ver estado y logs

```bash
# Ver estado
sudo systemctl status picoclaw-agents.service

# Ver logs en tiempo real
journalctl -u picoclaw-agents.service -f

# Ver logs de las últimas 2 horas
journalctl -u picoclaw-agents.service --since "2 hours ago"
```

### 4. Comandos útiles

```bash
# Detener servicio
sudo systemctl stop picoclaw-agents.service

# Reiniciar servicio
sudo systemctl restart picoclaw-agents.service

# Ver logs sin follow
journalctl -u picoclaw-agents.service -n 50
```

---

## Actualización

### Si usas Releases (Método 1)

```bash
# 1. Detener servicio
sudo systemctl stop picoclaw-agents.service

# 2. Descargar nueva versión
cd /opt/picoclaw-agents
wget https://github.com/comgunner/picoclaw-agents/releases/latest/download/picoclaw-agents_Linux_arm64.tar.gz

# 3. Extraer (reemplaza el binario anterior)
tar -xzf picoclaw-agents_Linux_arm64.tar.gz
chmod +x picoclaw-agents

# 4. Reiniciar servicio
sudo systemctl start picoclaw-agents.service

# 5. Verificar versión
picoclaw-agents --version
```

### Si usas Código Fuente (Método 2)

```bash
# 1. Detener servicio
sudo systemctl stop picoclaw-agents.service

# 2. Actualizar código
cd /opt/picoclaw-agents
git pull

# 3. Recompilar
make build

# 4. Reemplazar binario
cp build/picoclaw-agents /opt/picoclaw-agents/picoclaw-agents
chmod +x /opt/picoclaw-agents/picoclaw-agents

# 5. Reiniciar servicio
sudo systemctl start picoclaw-agents.service

# 6. Verificar versión
picoclaw-agents --version
```

---

## Comandos Disponibles

```bash
# One-shot query (consulta única)
picoclaw-agents agent -m "¿Qué tiempo hace?"

# Modo interactivo
picoclaw-agents agent

# Gateway (bots de Telegram, Discord, etc.)
picoclaw-agents gateway

# Listar skills disponibles
picoclaw-agents skills list

# Instalar skill desde ClawHub
picoclaw-agents skills install --slug skill-name

# Ver versión
picoclaw-agents --version

# Setup inicial
picoclaw-agents onboard
```

---

## Solución de Problemas

### El servicio no inicia

```bash
# Ver logs de error
journalctl -u picoclaw-agents.service -n 100 --no-pager

# Verificar que el binario existe
ls -lh /opt/picoclaw-agents/picoclaw-agents

# Verificar permisos
sudo chmod +x /opt/picoclaw-agents/picoclaw-agents
```

### Error de configuración

```bash
# Validar config.json
picoclaw-agents onboard --dry-run

# Regenerar config
rm ~/.picoclaw/config.json
picoclaw-agents onboard
```

### Verificar que el gateway está corriendo

```bash
# Ver proceso
ps aux | grep picoclaw-agents

# Ver puerto escuchando (si usa webhook)
sudo netstat -tulpn | grep picoclaw
```

### Error: "command not found: picoclaw-agents"

```bash
# Verificar symlink
ls -lh /usr/local/bin/picoclaw-agents

# Si no existe, crear manualmente
sudo ln -sf /opt/picoclaw-agents/picoclaw-agents /usr/local/bin/picoclaw-agents

# O usar ruta completa
/opt/picoclaw-agents/picoclaw-agents --version
```

---

## Enlaces Útiles

- **Releases:** https://github.com/comgunner/picoclaw-agents/releases
- **Documentación:** https://github.com/comgunner/picoclaw-agents/tree/main/docs
- **Issues:** https://github.com/comgunner/picoclaw-agents/issues
- **Discusión:** https://github.com/comgunner/picoclaw-agents/discussions

---

## Notas Importantes

### Nombre del Binario

- **Después de extraer el tar.gz:** El binario se llama `picoclaw-agents` (sin sufijo de plataforma)
- **En los comandos:** Siempre usa `picoclaw-agents` (el symlink apunta al binario correcto)

### Arquitectura del Servidor

- **ARM64 (Raspberry Pi, AWS Graviton):** Usa `picoclaw-agents_Linux_arm64.tar.gz`
- **AMD64 (Intel/AMD tradicional):** Usa `picoclaw-agents_Linux_amd64.tar.gz`

Para verificar tu arquitectura:
```bash
uname -m
# aarch64 o arm64 = ARM
# x86_64 = AMD64
```

### Go No Requerido (Método 1)

Si usas el binario pre-compilado de Releases, **NO necesitas instalar Go**. Go solo es necesario si compilas desde el código fuente (Método 2).
