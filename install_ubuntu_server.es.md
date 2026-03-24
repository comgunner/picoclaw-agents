# Instalación de PicoClaw en Ubuntu Server

Esta guía detalla los pasos necesarios para compilar e instalar PicoClaw en un entorno Ubuntu Server.

## Requisitos Previos

Asegúrate de tener instalados los siguientes componentes antes de comenzar:

- **Git:** Para clonar el repositorio.
- **Go:** PicoClaw está escrito en Go, por lo que requieres el compilador de Go instalado y configurado en tu sistema.
- **Make:** Herramienta utilizada para gestionar el proceso de compilación.

Puedes instalar las dependencias básicas del sistema operativo con:

```bash
sudo apt update
sudo apt install -y build-essential git wget tar
```

### Instalación de Go (Golang)

PicoClaw requiere Go 1.25.7 o superior. La forma recomendada de instalarlo directamente en Ubuntu Server es:

```bash
# Descargar e instalar Go 1.25.7
wget https://go.dev/dl/go1.25.7.linux-amd64.tar.gz
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.25.7.linux-amd64.tar.gz
rm go1.25.7.linux-amd64.tar.gz

# Añadir al PATH permanentemente
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.profile
source ~/.profile

# Verificar instalación
go version
```

## Pasos de Instalación

1. **Clonar el repositorio y entrar al directorio:**
   (Si instalas globalmente en `/opt`, asegúrate de dar permisos a tu usuario).

   ```bash
   cd /opt
   sudo git clone https://github.com/comgunner/picoclaw.git
   sudo chown -R $USER:$USER /opt/picoclaw
   cd picoclaw
   ```

2. **Descargar e instalar las dependencias de Go:**
   Esto descargará todas las librerías necesarias especificadas en el archivo `go.mod`.

   ```bash
   make deps
   ```

3. **Compilar el proyecto:**
   Generaremos el binario directamente en la carpeta raíz de `/opt/picoclaw`.

   ```bash
   go build -v -o picoclaw ./cmd/picoclaw
   chmod +x picoclaw
   ls -lh picoclaw
   ```
Una vez finalizado el proceso de instalación, puedes verificar que PicoClaw funciona correctamente ejecutando:

```bash
./picoclaw --help
```

## Configuración del Servidor

El archivo principal de configuración (`config.json`) debe colocarse en el directorio home del usuario que ejecutará PicoClaw.

1. **Ubicación en Ubuntu Server:**
   ```bash
   ~/.picoclaw/config.json
   ```
   *(Esto equivale a `/home/TU_USUARIO/.picoclaw/config.json` o `/root/.picoclaw/config.json` si lo corres como root).*

2. **Copiando tu configuración de Mac a Ubuntu:**
   Dado que ya configuraste tus bots (como Telegram) y tus API Keys (Deepseek) localmente en tu Mac, puedes simplemente transferir tu archivo actual al servidor usando `scp`:

   ```bash
   # Ejecuta esto desde tu Mac:
   scp ~/.picoclaw/config.json usuario@ip_de_tu_servidor:~/.picoclaw/config.json
   ```

   *Nota: Es posible que necesites crear la carpeta oculta en el servidor antes de copiarlo (`mkdir -p ~/.picoclaw`).*

3. **O usar el comando `onboard` para autogenerarlo:**
   Si prefieres no copiar el archivo, puedes dejar que PicoClaw genere uno nuevo con los valores por defecto (incluyendo la estructura multiagente `deepseek-chat` que dejamos configurada en el código fuente). Solo ejecuta:

   ```bash
   ./picoclaw onboard
   ```
   Esto creará el archivo `~/.picoclaw/config.json` con toda la estructura lista. A partir de ahí, solo tienes que editar el archivo usando `nano` o `vim` para pegar:
   - Tu `api_key` de Deepseek en la sección `model_list`.
   - Tu `token` de Telegram en la sección `channels > telegram`.
   - Tu ID de usuario en la lista `allow_from`.

## Ejecutar como Servicio (Systemd)

Para que PicoClaw se ejecute en segundo plano y se reinicie automáticamente si el servidor se apaga o ocurre un error, puedes configurarlo como un servicio de Systemd.

1. **Crear el archivo del servicio:**
   Este comando creará el archivo automáticamente usando tu usuario actual (`$USER`):

   ```bash
   cat <<EOF | sudo tee /etc/systemd/system/picoclaw.service
   [Unit]
   Description=PicoClaw Gateway Service
   After=network.target

   [Service]
   # Directorio donde vive el código
   WorkingDirectory=/opt/picoclaw

   # Ruta completa al binario compilado
   ExecStart=/opt/picoclaw/picoclaw gateway

   # El usuario que lo ejecuta
   User=$USER
   Group=$(id -gn)

   # Reiniciar automáticamente si falla
   Restart=always
   RestartSec=5

   [Install]
   WantedBy=multi-user.target
   EOF
   ```

2. **Habilitar e iniciar el servicio:**
   ```bash
   sudo systemctl daemon-reload
   sudo systemctl enable picoclaw.service
   sudo systemctl start picoclaw.service
   ```

4. **Ver estado y logs:**
   ```bash
   sudo systemctl status picoclaw.service
   journalctl -u picoclaw.service -f
   ```
