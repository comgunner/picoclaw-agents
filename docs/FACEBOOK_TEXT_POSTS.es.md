# Facebook Text-Only Posts - PicoClaw

## ✅ La herramienta `facebook_post` SÍ soporta texto simple

La herramienta `facebook_post` en PicoClaw soporta **tanto publicaciones con imagen como solo texto**. El parámetro `image_path` es **opcional**.

## Cómo Publicar Solo Texto

### Opción 1: Sin parámetro `image_path`

```bash
./picoclaw agent -m "Usa facebook_post con message='¡Hola Mundo! Este es un post solo de texto'"
```

### Opción 2: Explícitamente vacío

```bash
./picoclaw agent -m "Usa facebook_post con message='Hola desde PicoClaw', image_path=''"
```

### Opción 3: Desde Telegram

```text
Publica en Facebook solo texto: "¡Gran anuncio! Sin imágenes."
```

## Configuración Requerida

### Mínimo (Page Token)

```json
{
  "tools": {
    "social_media": {
      "facebook": {
        "default_page_id": "123456789012345",
        "default_page_token": "EAABsbCS1iHgBO7ZCxZBZCqZL8ZCZCqZL8ZCZCqZL8ZCZCqZL8"
      }
    }
  }
}
```

### Recomendado (con auto-refresh)

```json
{
  "tools": {
    "social_media": {
      "facebook": {
        "default_page_id": "123456789012345",
        "default_page_token": "EAABsbCS1iHgBO7ZCxZBZCqZL8ZCZCqZL8ZCZCqZL8ZCZCqZL8",
        "app_id": "1234567890123456",
        "app_secret": "abcd1234efgh5678ijkl9012mnop3456",
        "user_token": "EAABsbCS1iHgBO7ZDyZDZDrZM9ZDZDrZM9ZDZDrZM9ZDZDrZM9"
      }
    }
  }
}
```

**Nota:** Con `app_id`, `app_secret`, y `user_token` configurados, la herramienta puede auto-refresh el page token cuando expira (Error 190).

## Permisos Requeridos

Para publicar en Facebook (texto o imagen), necesitas:

- ✅ `pages_manage_posts` - **Requerido** para publicar
- ✅ `pages_read_engagement` - Opcional, para leer insights
- ✅ `pages_manage_engagement` - Opcional, para gestionar comentarios

### ❌ NO necesitas:
- ❌ `publish_actions` - **Deprecado** desde 2018
- ❌ `publish_pages` - **Deprecado** desde 2020

## Errores Comunes

### Error: "publish_actions is deprecated"

**Causa:** Estás usando un token antiguo o mal configurado.

**Solución:**
1. Ve a Graph API Explorer: https://developers.facebook.com/tools/explorer
2. Selecciona tu app
3. Genera token con permisos: `pages_manage_posts`
4. Actualiza `default_page_token` en config.json

### Error: "(190) Invalid OAuth Access Token"

**Causa:** El token expiró.

**Solución:**
1. Configura `app_id`, `app_secret`, y `user_token` en config.json
2. La herramienta hará auto-refresh automáticamente
3. O genera un nuevo token manualmente

### Error: "Page ID no configurado"

**Causa:** Falta `default_page_id` en config.json.

**Solución:**
1. Ve a tu página de Facebook
2. Click en "Información" (About)
3. Copia el Page ID (número de 15-16 dígitos)
4. Agrega a config.json: `"default_page_id": "123456789012345"`

## Ejemplos de Uso

### Solo Texto

```bash
./picoclaw agent -m "Publica en Facebook: '¡Nuevo artículo disponible! Leanlo en nuestro blog'"
```

### Texto + Imagen

```bash
./picoclaw agent -m "Publica en Facebook la imagen /tmp/foto.jpg con mensaje '¡Gran anuncio!'"
```

### Texto + Comentario (sin imagen)

```bash
./picoclaw agent -m "Usa facebook_post con message='Post principal', comment='Comentario adicional'"
```

### Multi-página (texto)

```bash
./picoclaw agent -m "
  Publica 'Anuncio importante' en estas páginas:
  - page_id='123456789012345', page_token='EAAB...'
  - page_id='987654321098765', page_token='EAAB...'
"
```

## Flujo Interno

Cuando `image_path` está vacío o no se proporciona:

1. La herramienta llama a `FacebookPostTextOnly()` en lugar de `FacebookPost()`
2. Usa el endpoint `/feed` en lugar de `/photos`
3. Envía solo `message` y `access_token`
4. Maneja auto-refresh del token si está configurado
5. Agrega comentario si se proporciona

## Código de Referencia

```go
// pkg/tools/social_media.go
if imagePath == "" {
    postID, err = utils.FacebookPostTextOnly(
        callCtx, 
        pageID, pageToken, 
        t.appID, t.appSecret, t.userToken, 
        message, comment
    )
} else {
    req := utils.FBPostRequest{
        PageID:    pageID,
        PageToken: pageToken,
        AppID:     t.appID,
        AppSecret: t.appSecret,
        UserToken: t.userToken,
        Message:   message,
        ImagePath: imagePath,
        Comment:   comment,
    }
    postID, err = utils.FacebookPost(callCtx, req)
}
```

## Verificación

Para verificar que tu configuración es correcta:

```bash
# Test solo texto
./picoclaw agent -m "Usa facebook_post con message='Test de texto desde PicoClaw'"

# Debería responder:
# "Publicación en Facebook exitosa. Post ID: 123456789012345_987654321098765"
```

## Soporte

Si tenés problemas:

1. Verificá que `default_page_id` y `default_page_token` estén configurados
2. Verificá que el token tenga permiso `pages_manage_posts`
3. Probá primero con texto simple (sin imagen)
4. Revisá los logs para ver el error exacto de la API

---

**Nota:** Esta documentación es para PicoClaw v3.3+ con la herramienta `facebook_post` nativa.
