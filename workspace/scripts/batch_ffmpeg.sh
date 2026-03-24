#!/bin/bash
# Procesamiento por lotes de videos con FFmpeg
# Uso: ./batch_ffmpeg.sh BATCH_ID [--input DIR] [--output DIR]
# 
# Demuestra el contrato con QueueManager:
# - Recibe BATCH_ID como primer argumento
# - Escribe estado a /tmp/picoclaw_queue_{BATCH_ID}.json
# - El QueueManager monitorea sin intervención

set -e

BATCH_ID="${1:-UNKNOWN}"
CLEAN_ID="${BATCH_ID#\#}"
STATE_FILE="/tmp/picoclaw_queue_${CLEAN_ID}.json"

INPUT_DIR="./videos/input"
OUTPUT_DIR="./videos/output"

# Parsear argumentos
while [[ $# -gt 0 ]]; do
    case $1 in
        --input)
            INPUT_DIR="$2"
            shift 2
            ;;
        --output)
            OUTPUT_DIR="$2"
            shift 2
            ;;
        *)
            shift
            ;;
    esac
done

# Función para reportar estado
report_state() {
    local status="$1"
    local progress="$2"
    local message="$3"
    
    cat > "$STATE_FILE" << EOF
{
  "status": "$status",
  "progress": $progress,
  "message": "$message",
  "timestamp": $(date +%s)
}
EOF
    
    echo "[$BATCH_ID] $message"
}

# Crear directorio de salida
mkdir -p "$OUTPUT_DIR"

# Contar videos
VIDEO_COUNT=$(find "$INPUT_DIR" -type f \( -name "*.mp4" -o -name "*.mov" -o -name "*.avi" \) 2>/dev/null | wc -l)

if [ "$VIDEO_COUNT" -eq 0 ]; then
    report_state "failed" 0 "No se encontraron videos en $INPUT_DIR"
    exit 1
fi

report_state "processing" 0 "Encontrados $VIDEO_COUNT videos. Iniciando procesamiento..."

# Procesar cada video
PROCESSED=0
for video in "$INPUT_DIR"/*.{mp4,mov,avi}; do
    [ -f "$video" ] || continue
    
    filename=$(basename "$video")
    output_file="$OUTPUT_DIR/processed_${filename}"
    
    echo "Procesando: $filename"
    
    # Convertir a H.264 (FFmpeg)
    if command -v ffmpeg &> /dev/null; then
        ffmpeg -i "$video" \
               -c:v libx264 -preset fast -crf 23 \
               -c:a aac -b:a 128k \
               -y "$output_file" 2>/dev/null || true
    else
        # Si FFmpeg no está instalado, simular
        echo "Simulando procesamiento de $filename (FFmpeg no disponible)"
        cp "$video" "$output_file"
    fi
    
    PROCESSED=$((PROCESSED + 1))
    PROGRESS=$((PROCESSED * 100 / VIDEO_COUNT))
    
    report_state "processing" "$PROGRESS" "Procesado $PROCESSED/$VIDEO_COUNT: $filename"
done

report_state "completed" 100 "Procesamiento completado: $PROCESSED videos"
echo "[$BATCH_ID] ✅ Completado"
