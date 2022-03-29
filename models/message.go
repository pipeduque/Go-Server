package models

import (
	"net"
	"time"
)

// Estructura para la creacion de un mensaje
type Message struct {
	date    time.Time //Fecha del mensaje
	sender  net.Addr  //Direccion del emisor
	content []byte    //Contenido del mensaje
	file    []byte    //Base64 de un archivo
}

/* Funcion
 * Nombre: NewMessage
 * Descripcion: Funcion encargada de crear un nuevo mensaje apartir de la estructura */
func NewMessage(sender net.Addr, content []byte, file []byte) *Message {

	return &Message{
		date:    time.Now(),
		sender:  sender,
		content: content,
		file:    file,
	}
}
