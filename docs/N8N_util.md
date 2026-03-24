# n8n Utility & Integration Guide

Este documento describe cómo integrar PicoClaw con n8n para potenciar las automatizaciones y flujos de trabajo avanzados.

## Arquitectura de Integración

En el ecosistema PicoClaw (específicamente en la configuración `picoclaw-n8n-starter`), n8n y PicoClaw conviven en la misma red de Docker, permitiendo una comunicación rápida y segura.

- **n8n Internal URL**: `http://n8n:5678`
- **PicoClaw Gateway**: `http://picoclaw-gateway:18789`

## Uso del Skill `n8n_workflow`

El skill `n8n_workflow` permite al agente de PicoClaw generar archivos JSON listos para ser importados en n8n. El agente conoce la estructura de los nodos, conexiones y parámetros necesarios.

### Características Principales
- Generación de flujos de trabajo completos de principio a fin.
- Configuración de nodos de Base de Datos (PostgreSQL), Telegram, HTTP y Agentes AI.
- Guía de sanitización de credenciales.

## Guía de Webhooks (Mini-Tutorial)

Para conectar PicoClaw con n8n vía webhooks, sigue estos pasos fundamentados en la [documentación oficial de n8n](https://docs.n8n.io/integrations/builtin/core-nodes/n8n-nodes-base.webhook/):

### 1. Configuración en n8n
1. Añade un nodo **Webhook** a tu canvas.
2. Establece el **HTTP Method** a `POST`.
3. En **Authentication**, selecciona `Header Auth`.
4. Define el nombre del header como `X-Webhook-Secret` y asigna un valor seguro (ej. un UUID).
5. Copia la **Test URL** para pruebas o la **Production URL** para uso final.

### 2. Disparo desde PicoClaw
Puedes pedirle a PicoClaw que dispare el flujo usando la herramienta `exec`:

```bash
curl -X POST "http://n8n:5678/webhook-test/mi-ruta-segura" \
  -H "Content-Type: application/json" \
  -H "X-Webhook-Secret: MI_SUPER_SECRETO" \
  -d '{"status": "trigger", "message": "Hola desde PicoClaw"}'
```

> [!TIP]
> Si estás usando Docker, usa `http://n8n:5678`. Si PicoClaw accede desde fuera del contenedor, usa la URL pública.

## Ejemplo Avanzado: n8n-claw Agent

El flujo `n8n-claw-agent.json` (ubicado en `local_work/SKILLS_UPDATE/n8n-workflow/n8n-claw/workflows/`) es un ejemplo de cómo construir un agente autónomo dentro de n8n que interactúa con PicoClaw:

1. **Telegram Trigger**: Recibe mensajes.
2. **Postgres Load**: Carga "Soul" (personalidad) y configuración.
3. **AI Agent Node**: Procesa con LLM (Claude/GPT) y dispone de herramientas (Memory, HTTP, WorkflowBuilder).
4. **Telegram Reply**: Envía la respuesta de vuelta.

## Recursos Adicionales
- [Documentación de URLs de Webhook](https://docs.n8n.io/integrations/builtin/core-nodes/n8n-nodes-base.webhook/#webhook-urls)
- [Video Tutorial: n8n Webhooks](https://www.youtube.com/watch?v=7ekNNMmiNrM)
