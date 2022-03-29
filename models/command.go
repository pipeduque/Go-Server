package models

type ID int

// Comandos disponibles en el protocolo personalizado
const (
	REG      ID = iota
	JOIN        // Cliente ingresa a un canal
	LEAVE       // Cliente sale de un canal
	MSG         // Envia un mensaje
	CREATE      // Crea un canal
	LIST_CHN    // Lista los canales
	LIST_MSG    // Lista los mensajes de un canal
	LIST_USR    // Lista los usuarios de un canal
)

// Estructura para la creacion de un comando
type Command struct {
	id      ID     // Identificador del comando
	channel string // Nombre del canal a crear si es el caso
	sender  Client // Emisor del comando
	content []byte // Contenido de un mensaje
	file    []byte // Base64 de un archivo
}
