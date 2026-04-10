#!/usr/bin/env python3
"""
Generación por lotes de imágenes con CUDA
Uso: python batch_cuda_gen.py BATCH_ID [--count N] [--model MODEL]

Este script demuestra el contrato con QueueManager:
- Recibe BATCH_ID como primer argumento
- Escribe estado a /tmp/picoclaw_queue_{BATCH_ID}.json
- El QueueManager monitorea este archivo sin intervención
"""

import sys
import json
import time
from pathlib import Path

def get_state_file(batch_id: str) -> Path:
    """Obtiene la ruta del archivo de estado."""
    clean_id = batch_id.lstrip('#')
    return Path(f"/tmp/picoclaw_queue_{clean_id}.json")

def report_state(batch_id: str, status: str, progress: int, message: str, result=None):
    """Reporta el estado actual al QueueManager."""
    state = {
        "status": status,
        "progress": progress,
        "message": message,
        "timestamp": time.time()
    }
    if result:
        state["result"] = result

    state_file = get_state_file(batch_id)
    state_file.write_text(json.dumps(state, indent=2))
    print(f"[{batch_id}] {message}")

def generate_image(prompt: str, index: int) -> str:
    """
    Simula generación de imagen.
    En producción, reemplazar con llamada a CUDA/Diffusion.
    """
    print(f"  Generando imagen {index + 1}...")
    time.sleep(2)  # Simular trabajo pesado
    return f"image_{index + 1:03d}.png"

def main():
    if len(sys.argv) < 2:
        print("Uso: python batch_cuda_gen.py BATCH_ID [--count N] [--model MODEL]")
        sys.exit(1)

    batch_id = sys.argv[1]

    # Parsear argumentos opcionales
    count = 5
    model = "sdxl"

    i = 2
    while i < len(sys.argv):
        if sys.argv[i] == "--count":
            count = int(sys.argv[i + 1])
            i += 2
        elif sys.argv[i] == "--model":
            model = sys.argv[i + 1]
            i += 2
        else:
            i += 1

    try:
        # Reportar inicio
        report_state(batch_id, "processing", 0,
                    f"Iniciando generación de {count} imágenes (modelo: {model})...")

        # Generar imágenes
        generated_images = []
        for i in range(count):
            image_path = generate_image("landscape", i)
            generated_images.append(image_path)

            # Reportar progreso
            progress = (i + 1) * 100 // count
            report_state(
                batch_id,
                "processing",
                progress,
                f"Imagen {i + 1}/{count} completada",
                {"images_so_far": generated_images}
            )

        # Reportar completado
        report_state(
            batch_id,
            "completed",
            100,
            f"Generación completada: {count} imágenes",
            {
                "images": generated_images,
                "model": model,
                "output_dir": "workspace/images/"
            }
        )

        print(f"[{batch_id}] ✅ Completado")
        return 0

    except Exception as e:
        # Reportar error
        report_state(
            batch_id,
            "failed",
            0,
            f"Error: {str(e)}",
            {"error": str(e)}
        )
        print(f"[{batch_id}] ❌ Error: {str(e)}")
        return 1

if __name__ == "__main__":
    exit_code = main()
    sys.exit(exit_code)
