package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/smtp"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// EmailRequest representa el JSON de entrada
type EmailRequest struct {
	To      string `json:"to" binding:"required"`
	Subject string `json:"subject" binding:"required"`
	Body    string `json:"body" binding:"required"`
}

func main() {
	godotenv.Load() // Si no usas godotenv, puedes comentar esta línea

	r := gin.Default()

	// Middleware para habilitar CORS
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "https://ezequielunia.com.ar")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	r.POST("/send-email", func(c *gin.Context) {
		log.Println("[INFO] Nueva petición POST /send-email")
		var req EmailRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Println("[ERROR] JSON inválido:", err)
			c.JSON(400, gin.H{"error": "Datos inválidos"})
			return
		}
		log.Printf("[INFO] Datos recibidos: To=%s, Subject=%s\n", req.To, req.Subject)

		smtpHost := os.Getenv("SMTP_HOST")
		smtpPort := os.Getenv("SMTP_PORT")
		smtpUser := os.Getenv("SMTP_USER")
		smtpPass := os.Getenv("SMTP_PASS")

		log.Printf("[INFO] SMTP_HOST=%s, SMTP_PORT=%s, SMTP_USER=%s, SMTP_USER=%s\n", smtpHost, smtpPort, smtpUser, smtpUser)

		if smtpHost == "" || smtpPort == "" || smtpUser == "" || smtpPass == "" {
			log.Println("[ERROR] Faltan variables de entorno SMTP")
			c.JSON(500, gin.H{"error": "Faltan variables de entorno SMTP"})
			return
		}

		msg := []byte(fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=\"utf-8\"\r\n\r\n%s", smtpUser, req.To, req.Subject, req.Body))

		if smtpPort == "465" {
			// Conexión segura TLS (Hostinger, Gmail, etc.)
			log.Println("[INFO] Usando conexión TLS para SMTP (puerto 465)")
			serverAddr := fmt.Sprintf("%s:%s", smtpHost, smtpPort)
			tlsconfig := &tls.Config{
				InsecureSkipVerify: false,
				ServerName:         smtpHost,
			}
			conn, err := tls.Dial("tcp", serverAddr, tlsconfig)
			if err != nil {
				log.Println("[ERROR] No se pudo establecer conexión TLS:", err)
				c.JSON(500, gin.H{"success": false, "error": "No se pudo conectar al servidor SMTP (TLS)"})
				return
			}
			client, err := smtp.NewClient(conn, smtpHost)
			if err != nil {
				log.Println("[ERROR] No se pudo crear cliente SMTP:", err)
				c.JSON(500, gin.H{"success": false, "error": "No se pudo crear cliente SMTP"})
				return
			}
			defer client.Quit()
			// Autenticación
			auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)
			if err = client.Auth(auth); err != nil {
				log.Println("[ERROR] Autenticación SMTP fallida:", err)
				c.JSON(500, gin.H{"success": false, "error": "Autenticación SMTP fallida"})
				return
			}
			if err = client.Mail(smtpUser); err != nil {
				log.Println("[ERROR] MAIL FROM falló:", err)
				c.JSON(500, gin.H{"success": false, "error": "MAIL FROM falló"})
				return
			}
			if err = client.Rcpt(req.To); err != nil {
				log.Println("[ERROR] RCPT TO falló:", err)
				c.JSON(500, gin.H{"success": false, "error": "RCPT TO falló"})
				return
			}
			w, err := client.Data()
			if err != nil {
				log.Println("[ERROR] Error al obtener writer SMTP:", err)
				c.JSON(500, gin.H{"success": false, "error": "Error al obtener writer SMTP"})
				return
			}
			_, err = w.Write(msg)
			if err != nil {
				log.Println("[ERROR] Error escribiendo mensaje SMTP:", err)
				c.JSON(500, gin.H{"success": false, "error": "Error escribiendo mensaje SMTP"})
				return
			}
			w.Close()
			log.Println("[INFO] Correo enviado exitosamente (TLS)")
			c.JSON(200, gin.H{"success": true})
			return
		}

		// Para otros puertos (ej: 587), envío estándar
		log.Println("[INFO] Enviando correo (sin TLS explícito)...")
		auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)
		err := smtp.SendMail(smtpHost+":"+smtpPort, auth, smtpUser, []string{req.To}, msg)
		if err != nil {
			log.Println("[ERROR] Error enviando correo:", err)
			c.JSON(500, gin.H{"success": false, "error": "No se pudo enviar el correo"})
			return
		}
		log.Println("[INFO] Correo enviado exitosamente")
		c.JSON(200, gin.H{"success": true})
	})

	port := os.Getenv("API_PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("[INFO] Servidor escuchando en el puerto %s\n", port)
	r.Run(":" + port)
}
