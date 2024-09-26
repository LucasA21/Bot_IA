# Telegram Stock Bot

Este es un bot de Telegram que gestiona un stock simple del hogar. El bot está conectado a la API de Cohere, lo que le permite interpretar comandos en lenguaje natural y almacenar artículos en una base de datos PostgreSQL. Además, puede generar y enviar un archivo `.txt` con el stock actual.

## Características

- **Añadir items al stock:** El bot interpreta comandos del usuario y agrega artículos al stock del hogar.
- **Listar el stock actual:** Puedes pedirle al bot que te envíe un archivo `.txt` con la lista de artículos en stock.
- **Inteligencia Artificial:** Gracias a la API de Cohere, el bot entiende comandos en lenguaje natural sin necesidad de seguir una sintaxis estricta.
- **Base de datos PostgreSQL:** Los datos del stock se almacenan de manera persistente en una base de datos.

## Ejemplo de uso

![Ejemplo del bot](src/assets/example.png)

## Requisitos

- [Go](https://golang.org/doc/install) 1.16 o superior.
- [PostgreSQL](https://www.postgresql.org/download/) para la base de datos.
- [Cohere API](https://docs.cohere.ai/docs) para el procesamiento de lenguaje natural.
- [Telegram Bot API](https://core.telegram.org/bots/api) para interactuar con el bot.

