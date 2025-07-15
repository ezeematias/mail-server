# Mail Server API (Go)

API REST para enviar correos electrónicos usando Go y Gin.

## Configuración

Debes definir las siguientes variables de entorno:

- `SMTP_HOST`: Host del servidor SMTP (ej: smtp.gmail.com)
- `SMTP_PORT`: Puerto SMTP (ej: 587)
- `SMTP_USER`: Usuario SMTP
- `SMTP_PASS`: Contraseña SMTP
- `SMTP_FROM`: Correo remitente
- `API_PORT`: (opcional) Puerto de la API (por defecto 8080)

Puedes definirlas en tu terminal o en un archivo `.env` (usando herramientas como [direnv](https://direnv.net/) o [godotenv](https://github.com/joho/godotenv)).

## Instalación

```bash
go mod tidy
```

## Ejecución

```bash
go run main.go
```

La API estará disponible en `http://localhost:8080` (o el puerto que definas).

## Endpoint

### POST /send-email

**Body JSON:**
```json
{
  "to": "destinatario@dominio.com",
  "subject": "Asunto",
  "body": "Cuerpo del mensaje"
}
```

**Respuesta:**
- Éxito: `{ "success": true }`
- Error: `{ "success": false, "error": "mensaje" }` 