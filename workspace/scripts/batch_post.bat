@echo off
REM Subida por lotes a redes sociales
REM Uso: batch_post.bat BATCH_ID --platforms facebook,twitter --dir images
REM
REM Demuestra el contrato con QueueManager:
REM - Recibe BATCH_ID como primer argumento
REM - Escribe estado a C:\temp\picoclaw_queue_{BATCH_ID}.json

setlocal enabledelayedexpansion

set BATCH_ID=%1
set CLEAN_ID=%BATCH_ID:#=%
set STATE_FILE=C:\temp\picoclaw_queue_%CLEAN_ID%.json

REM Parsear argumentos
set PLATFORMS=facebook
set IMAGE_DIR=images

:parse_args
shift
if "%~1"=="" goto end_parse
if /i "%~1"=="--platforms" set PLATFORMS=%~2& shift & shift & goto parse_args
if /i "%~1"=="--dir" set IMAGE_DIR=%~2& shift & shift & goto parse_args
shift
goto parse_args

:end_parse

REM Crear directorio temporal si no existe
if not exist "C:\temp" mkdir "C:\temp"

REM Contar imágenes
setlocal enabledelayedexpansion
set COUNT=0
for %%f in ("%IMAGE_DIR%\*.jpg" "%IMAGE_DIR%\*.png") do set /a COUNT+=1

if %COUNT%==0 (
    (
        echo {
        echo   "status": "failed",
        echo   "progress": 0,
        echo   "message": "No se encontraron imágenes en %IMAGE_DIR%"
        echo }
    ) > %STATE_FILE%
    echo [%BATCH_ID%] No se encontraron imágenes
    exit /b 1
)

REM Reportar inicio
(
    echo {
    echo   "status": "processing",
    echo   "progress": 0,
    echo   "message": "Encontradas %COUNT% imágenes. Iniciando subida a %PLATFORMS%..."
    echo }
) > %STATE_FILE%

echo [%BATCH_ID%] Encontradas %COUNT% imágenes. Iniciando subida...

REM Subir cada imagen (ejemplo simplificado)
set PROCESSED=0
for %%f in ("%IMAGE_DIR%\*.jpg" "%IMAGE_DIR%\*.png") do (
    set /a PROCESSED+=1
    set /a PROGRESS=PROCESSED * 100 / COUNT

    echo Subiendo: %%~nxf

    REM Aquí iría la lógica real de subida
    REM Por ejemplo: python scripts/upload_to_facebook.py "%%f"
    REM timeout /t 1 /nobreak

    (
        echo {
        echo   "status": "processing",
        echo   "progress": !PROGRESS!,
        echo   "message": "Subida !PROCESSED!/%COUNT%: %%~nxf"
        echo }
    ) > %STATE_FILE%
)

REM Reportar completado
(
    echo {
    echo   "status": "completed",
    echo   "progress": 100,
    echo   "message": "Subida completada: %PROCESSED% imágenes a %PLATFORMS%"
    echo }
) > %STATE_FILE%

echo [%BATCH_ID%] ✅ Completado
