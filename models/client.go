package models

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"log"
	"net"
)

// Estructura para la creacion de clientes
type Client struct {
	address    net.Addr
	connection net.Conn
	middlemane chan<- Command
	Online     chan<- *Client
	offline    chan<- *Client
}

/* Funcion
 * Nombre: NewClient
 * Descripcion: Funcion encargada de crear nuevos clientes apartir de la estructura */
func NewClient(connection net.Conn, server *Server) *Client {

	return &Client{
		address:    connection.RemoteAddr(),
		connection: connection,
		middlemane: server.commands,
	}
}

/* Funcion
 * Nombre: RequestReadHandle
 * Descripcion: Funcion encargada de manejar la lectura de solicitudes entrantes del cliente */
func (client *Client) RequestReadHandle() {

	connection := client.connection //Conexion perteniente al cliente con el servidor

	defer func() { //Funcion diferida que cerrara la conexion cada vez que handleRead termine
		if err := connection.Close(); err != nil {
			log.Println("Error closing connection: ", err)
		}
	}()

	for { // Ciclo para estar escuchando las solicitudes del cliente hasta que el rompa la conexion

		request, err := bufio.NewReader(connection).ReadBytes('\n') // buffer para leer las solicitudes entrantes del cliente. lector de solicitudes
		log.Println("Request: ", string(request))

		if err != nil { //Manejamos un posible error en la solicitud
			if err == io.EOF { //si el error es end-of-line (EOF) en el lector de solicitudes desconectamos el cliente
				client.offline <- client
			}
			client.writeError(err)
			break
		}

		client.requestHandler(request) //Si la solicitud es correcta, manejamos la solicitud
	}
}

/* Funciones
 * Nombre: requestHandler
 * Descripcion: Funcion encargada de manejar una solicitud del cliente
 * @request solicitud del cliente en bytes */
func (client *Client) requestHandler(request []byte) {

	// analizando la solicitud del cliente con la libreria bytes
	cmd := bytes.ToUpper(bytes.TrimSpace(bytes.Split(request, []byte(" "))[0])) //el comando (cmd) sera el primer corte de la solicitud segun los espacios en blanco
	args := bytes.TrimSpace(bytes.TrimPrefix(request, cmd))                     // Los argumentos (args) corresponderan al corte de la solicitud (request) menos el comando (cmd)

	if string(cmd) == "" { // Manejamos que el comando no sea vacio
		client.WriteResponse("FALSE")

	} else {

		switch string(cmd) { //Segun sea el comando

		case "JOIN": // Solicitud para entrar a un canal
			if err := client.joinChannel(args); err != nil {
				client.writeError(err)
			}

		case "LEAVE": // Solicitud para salir de un canal
			if err := client.leaveChannel(args); err != nil {
				client.writeError(err)
			}

		case "CREATE": //Solicitud para crear un canal
			if err := client.createChannel(args); err != nil {
				client.writeError(err)
			}

		case "LIST_CHN": // Solicitud para listar los canales existentes
			if err := client.listChannels(); err != nil {
				client.writeError(err)
			}

		case "MSG": //Solicitud para envio de un archivo a un canal existente
			if err := client.sendMsg(args); err != nil {
				client.writeError(err)
			}

		case "LIST_MSG": // Solicitud para listar los mensajes de un canal
			if err := client.listMsg(args); err != nil {
				client.writeError(err)
			}

		case "LIST_USR": // Solicitud para listar los clientes conectados en un canal
			if err := client.listUsrChannel(args); err != nil {
				client.writeError(err)
			}

		default: //Si el comando no es reconocido en las posibilidades
			client.WriteResponse("FALSE") //Informamos que es invalido
		}
	}
}

/** FUNCIONES PARA ASIGNAR AL INTERMEDIARIO DEL CLIENTE CON EL SERVIDOR UN NUEVO COMANDO **/
/** 																					 **/
/** @args: argumentos escritos por el cliente											 **/
/** return: @error: nil si fue correcta la creacion del comando, err si fallo.           **/

//Comando para que un cliente se conecte a un canal
func (client *Client) joinChannel(args []byte) error {

	channel, err := client.getArg(args, 0) // Manejamos que el primer argumento no sea vacio

	if err != nil { // Manejamos que el primer argumento no sea vacio
		return err
	}

	client.middlemane <- Command{ // Asignamos al intermediario el nuevo comando
		channel: string(channel),
		sender:  *client,
		id:      JOIN,
	}
	return nil
}

// Comando para que un cliente se desconecte de un canal
func (client *Client) leaveChannel(args []byte) error {

	channel, err := client.getArg(args, 0) // Obtenemos el primer argumento, correspondiente al nombre del canal a salir

	if err != nil { // Manejamos que el primer argumento no sea vacio
		return err
	}

	client.middlemane <- Command{ // Asignamos al intermediario el nuevo comando
		channel: string(channel),
		sender:  *client,
		id:      LEAVE,
	}
	return nil
}

// Comando para crear un canal
func (client *Client) createChannel(args []byte) error {

	channel, err := client.getArg(args, 0) // Obtenemos el primer argumento, correspondiente al nombre del canal a crear

	if err != nil { // Manejamos que el primer argumento no sea vacio
		return err
	}

	client.middlemane <- Command{ // Asignamos al intermediario el nuevo comando
		channel: string(channel),
		sender:  *client,
		id:      CREATE,
	}
	return nil
}

// Comando para listar los canales
func (client *Client) listChannels() error {

	client.middlemane <- Command{ // Asignamos al intermediario el nuevo comando
		sender: *client,
		id:     LIST_CHN,
	}
	return nil
}

// Comando para enviar un mensaje
func (client *Client) sendMsg(args []byte) error {

	channel, err := client.getArg(args, 0) // Obtenemos el primer argumento, correspondiente al nombre del canal destinatario

	if err != nil { // Manejamos que el primer argumento no sea vacio
		return err
	}

	message, err := client.getArg(args, 1) // Obtenemos el segundo argumento, correspondiente al mensaje

	if err != nil { // Manejamos que el segundo argumento no sea vacio
		return err
	}

	file, err := client.getArg(args, 2) // Obtenemos el tercer argumento, correspondiente a un archivo en base64

	if err != nil { // Manejamos que el tercer argumento no sea vacio
		return err
	}

	client.middlemane <- Command{ // Asignamos al intermediario el nuevo comando
		channel: string(channel),
		sender:  *client,
		content: message,
		file:    file,
		id:      MSG,
	}

	return nil
}

// Comando para listar los mensajes de un canal
func (client *Client) listMsg(args []byte) error {

	channel, err := client.getArg(args, 0) // Obtenemos el primer argumento, correspondiente al nombre del canal a listar los mensajes

	if err != nil { // Manejamos que el primer argumento no sea vacio
		return err
	}

	client.middlemane <- Command{ // Asignamos al intermediario el nuevo comando
		channel: string(channel),
		sender:  *client,
		id:      LIST_MSG,
	}

	return nil
}

// Comando para listar los usuarios de un canal.
func (client *Client) listUsrChannel(args []byte) error {

	channel, err := client.getArg(args, 0) // Obtenemos el primer argumento, correspondiente al nombre del canal a listar los usuarios

	if err != nil { // Manejamos que el primer argumento no sea vacio
		return err
	}

	client.middlemane <- Command{ // Asignamos al intermediario el nuevo comando
		channel: string(channel),
		sender:  *client,
		id:      LIST_USR,
	}

	return nil
}

/** FIN FUNCIONES PARA COMANDOS **/

/* Funcion
 * Nombre: getArg
 * Descripcion: Obtiene un argumento segun la posicion indicada en los argumentos totales
 * @args: argumentos escritos por el cliente.
 * @position: posicion deseada del argumento a tomar
 * return: @[]byte: argumento deseado en bytes
  		   @error:  nil si existe el argumento deseado, err si esta vacio. */
func (client *Client) getArg(args []byte, position int) ([]byte, error) {

	arg01 := bytes.Split(args, []byte(";;"))[position] //Tomamos un argumento segun la posicion solicitada de los argumentos

	if len(arg01) == 0 { //Manejamos que no este vacio el argumento solicitado

		return arg01, errors.New("empty arg01")
	}
	return arg01, nil
}

/* Funcion
 * Nombre: WriteResponse
 * Descripcion: Escribe a la conexion del cliente la respuesta del servidor
 * @res: respuesta dada */
func (client *Client) WriteResponse(res string) {

	if _, err := client.connection.Write([]byte(res + "\n")); err != nil {
		client.writeError(err)
		return
	}
	log.Println(res)
}

/* Funcion
 * Nombre: writeError
 * Descripcion: Escribe a la conexion del cliente un error surgido en el servidor
 * @error: error dado */
func (client *Client) writeError(e error) {

	if _, err := client.connection.Write([]byte("ERROR " + e.Error() + "\n")); err != nil {
		client.writeError(err)
		return
	}
	log.Println(e)
}
