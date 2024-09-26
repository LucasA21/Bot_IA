package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/GiulianoPoeta99/telegram_go.git/src/IA"
	"github.com/GiulianoPoeta99/telegram_go.git/src/db"
	stock "github.com/GiulianoPoeta99/telegram_go.git/src/models/producto"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error al cargar archivo .env")
	}

	telegramBotToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	cohereApiKey := os.Getenv("COHERE_API_KEY")
	databaseURL := os.Getenv("DATABASE_URL")

	if telegramBotToken == "" || cohereApiKey == "" || databaseURL == "" {
		log.Fatal("Asegúrate de establecer TELEGRAM_BOT_TOKEN, COHERE_API_KEY y DATABASE_URL en tu entorno.")
	}

	conn := db.ConnectToDB()
	defer conn.Close(context.Background())

	bot, err := tgbotapi.NewBotAPI(telegramBotToken)
	if err != nil {
		log.Panic(err.Error())
	}

	bot.Debug = true
	log.Printf("Bot autorizado en cuenta %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		userMessage := update.Message.Text
		userID := update.Message.From.ID

		if userMessage == "flipo" {
			// Ruta a la imagen específica
			imagePath := "src/assets/flipo.png"

			// Abrir la imagen
			file, err := os.Open(imagePath)
			if err != nil {
				log.Printf("Error al abrir la imagen: %v", err)
				continue
			}

			// Enviar la imagen al usuario
			msg := tgbotapi.NewPhoto(update.Message.Chat.ID, tgbotapi.FileReader{
				Name:   "flipo.png",
				Reader: file,
			})

			if _, err := bot.Send(msg); err != nil {
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Hubo un error al enviar la imagen."))
				log.Printf("Error al enviar la imagen: %v", err)
			}

			// Cerrar el archivo explícitamente después de usarlo
			file.Close()

			continue
		}

		coherePrompt := fmt.Sprintf("El usuario dice: '%s'. Interpreta este mensaje y devuelve una estructura JSON clara con las claves 'accion', 'producto' y 'cantidad'. Las acciones válidas son 'agregar', 'quitar', 'actualizar' o 'enviar archivo stock'. Si no puedes interpretar el mensaje, responde 'accion': 'desconocida'.", userMessage)
		cohereResponse, err := IA.GetCohereResponse(coherePrompt, cohereApiKey)
		if err != nil {
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Hubo un error al procesar tu solicitud."))
			continue
		}

		switch cohereResponse.Accion {
		case "agregar":
			err := stock.AgregarAlStock(conn, userID, cohereResponse.Producto, cohereResponse.Cantidad)
			if err != nil {
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Hubo un error al agregar al stock."))
			} else {
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Se ha agregado %d %s al stock.", cohereResponse.Cantidad, cohereResponse.Producto)))
			}
		case "quitar":
			// Implementar lógica para quitar productos
		case "actualizar":
			// Implementar lógica para actualizar productos
		case "enviar archivo stock":
			fileName, err := stock.GenerarArchivoStock(conn, userID)
			if err != nil {
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Hubo un error al generar el archivo del stock."))
				continue
			}

			// Abrir el archivo para enviarlo
			file, err := os.Open(fileName)
			if err != nil {
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Hubo un error al abrir el archivo."))
				continue
			}

			// Enviar el archivo al usuario
			msg := tgbotapi.NewDocument(update.Message.Chat.ID, tgbotapi.FileReader{
				Name:   fileName,
				Reader: file,
			})
			if _, err := bot.Send(msg); err != nil {
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Hubo un error al enviar el archivo."))
			}

			// Eliminar el archivo después de enviarlo para evitar acumulación
			err = os.Remove(fileName)
			if err != nil {
				log.Printf("Error al eliminar el archivo: %v", err)
			}

		default:
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "¡Hola! Soy SHObot, tu asistente para gestionar el stock. Puedes pedirme que agregue, actualice o quite artículos de tu inventario. Además, puedo enviarte un archivo con los productos actuales."))
		}
	}
}
