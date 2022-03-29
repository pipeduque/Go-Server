package models

import "time"

// Estructura para la creacion de un canal
type Channel struct {
	name     string              // Nombre del canal
	date     time.Time           // Fecha de creacion
	clients  map[*Client]bool    // Clientes
	messages map[string]*Message // Mensajes del canal
}

/* Funcion
 * Nombre: NewChannel
 * Descripcion: Funcion encargada de crear un canal apartir de la estructura */
func NewChannel(nameChannel string) *Channel {

	return &Channel{
		name:     nameChannel,
		date:     time.Now(),
		clients:  make(map[*Client]bool),
		messages: make(map[string]*Message),
	}
}
