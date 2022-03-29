package models

import (
	"net"
	"sort"
)

//Estructura para la creacion del servidor
type Server struct {
	clients          map[net.Addr]*Client // Mapa de clientes en el servidor
	channels         map[string]*Channel  // Canales del servidor
	commands         chan Command         // Comando para ser analizado, es modificado por cada solicitud
	ClientOnlineReq  chan *Client
	clientOfflineReq chan *Client
}

/* Funcion
 * Nombre: NewServer
 * Descripcion: Funcion encargada de crear el servidor apartir de la estructura */
func NewServer() *Server {

	return &Server{
		clients:          make(map[net.Addr]*Client),
		channels:         make(map[string]*Channel),
		commands:         make(chan Command),
		ClientOnlineReq:  make(chan *Client),
		clientOfflineReq: make(chan *Client),
	}
}

/* Funcion
 * Nombre: Run
 * Descripcion: Ejecutara el comando que es asignado por la solicitud del cliente */
func (server *Server) Run() {
	for {

		select {

		case client := <-server.ClientOnlineReq: // Conectamos un cliente al servidor
			server.setClientOnline(client)

		case client := <-server.clientOfflineReq: // Desconectamos un cliente del servidor
			server.setClientOffline(client)

		case cmd := <-server.commands: // Comando solicitado por el cliente

			switch cmd.id { // Comando disponibles en el protocolo

			case JOIN: // Cliente ingresa a un canal
				server.joinChannel(cmd.sender.address, cmd.channel)

			case LEAVE: // Cliente sale de un canal
				server.leaveChannel(cmd.sender.address, cmd.channel)

			case MSG: // Envia un mensaje al servidor
				server.sendMessage(cmd.sender.address, cmd.channel, cmd.content, cmd.file)

			case CREATE: // Crea un canal nuevo
				server.createChannel(cmd.sender.address, cmd.channel)

			case LIST_CHN: // Lista los canales existentes
				server.listChannels(cmd.sender.address)

			case LIST_MSG: // Lista los mensajes de un canal
				server.listMessages(cmd.sender.address, cmd.channel)

			case LIST_USR: // Lista los cliente conectados de un canal
				server.listUsrChannel(cmd.sender.address, cmd.channel)
			}
		}
	}
}

/* Funcion: setClientOffline
 * Desconecta a un cliente del servidor
 * @param client cliente a desconectar */
func (server *Server) setClientOffline(c *Client) {

	if _, exists := server.clients[c.address]; exists { // Manejamos que el cliente que solicita desconectarse exista en el servidor

		delete(server.clients, c.address) // Lo eliminamos de los clientes

		for _, channel := range server.channels { //Lo eliminamos de los canales
			delete(channel.clients, c)
		}
	}
}

/* Funcion: setClientOnline
 * Conecta a un cliente a el servidor
 * @param client cliente a conectar */
func (server *Server) setClientOnline(client *Client) {

	if _, exists := server.clients[client.address]; exists {
		// Aqui se validara si el nombre de cliente esta disponible
	} else {
		server.clients[client.address] = client // Conectamos el cliente al servidor si el cliente no existe
	}
}

/* Funcion: joinChannel
 * Conecta a un cliente a un canal
 * @param sender direccion del emisor de la solicitud.
 * @param channelName nombre del canal a conectar */
func (server *Server) joinChannel(sender net.Addr, channelName string) {

	if client, ok := server.clients[sender]; ok { // Manejamos que el cliente que solicita exista en el servidor

		if channel, ok := server.channels[channelName]; ok { // Manejamos que el canal a conectar exista

			channel.clients[client] = true // Conectamos al cliente
		}
	}
}

/* Funcion: leaveChannel
 * Desconecta a un cliente de un canal
 * @param sender direccion del emisor de la solicitud.
 * @param channelName nombre del canal a desconectar */
func (server *Server) leaveChannel(sender net.Addr, channelName string) {

	if client, ok := server.clients[sender]; ok { // Manejamos que el cliente que solicita exista en el servidor

		if channel, ok := server.channels[channelName]; ok { // Manejamos que el canal a salir exista

			channel.clients[client] = true // Desconectamos al cliente
		}
	}
}

/* Funcion: sendMessage
 * Envia un mensaje al servidor
 * @param sender direccion del emisor de la solicitud.
 * @param channelName nombre del canal destinatario del mensaje
 * @param message mensaje a enviar */
func (server *Server) sendMessage(senderAddress net.Addr, channelName string, message []byte, file []byte) {

	if _, ok := server.clients[senderAddress]; ok { // Manejamos que el cliente que solicita exista en el servidor

		if channel, ok := server.channels[channelName]; ok { // Manejamos que el canal destinatario exista

			channel.messages[string(message)] = NewMessage(senderAddress, message, file)
		}
	}
}

/* Funcion: createChannel
 * Crea un canal para el servidor
 * @param sender direccion del emisor de la solicitud.
 * @param channelName nombre del canal que se desea crear */
func (server *Server) createChannel(sender net.Addr, channelName string) {

	if client, ok := server.clients[sender]; ok { // Manejamos que el cliente que solicita exista en el servidor

		if _, ok := server.channels[channelName]; ok { // Manejamos que el canal a crear no exista

			client.WriteResponse("FALSE")

		} else {

			server.channels[channelName] = NewChannel(channelName) // Creamos el canal en el servidor
			server.listChannels(client.address)
		}
	}

}

/* Funcion: listChannels
 * Lista los canales del servidor
 * @param sender direccion del emisor de la solicitud */
func (server *Server) listChannels(sender net.Addr) {

	if client, ok := server.clients[sender]; ok { // Manejamos que el cliente que solicita exista en el servidor

		if len(server.channels) > 0 { // Verificamos si existen canales antes de proceder

			response := ""

			channels := make([]string, 0) // array de strings para ordenar los canales por su fecha de creacion

			for _, values := range server.channels {

				//llenamos el array de canales con: fecha_creacion,emisor,null
				channels = append(channels, values.date.Format("2006-01-02:15:04:05")+","+values.name+","+"null")
			}

			sort.Strings(channels) //Ordenamos el array de mensajes

			for _, key := range channels {

				// Juntamos los canales del array ordenado en la respuesta dividido por ;
				// fecha_creacion,emisor,null;fecha_creacion,emisor,null
				response = response + key + ";"
			}

			client.WriteResponse(response)

		}
	}
}

/* Funcion: listMessages
 * Lista los mensajes pertenecientes a un canal
 * @param sender direccion del emisor de la solicitud.
 * @param channelName nombre del canal que se listaran sus mensajes */
func (server *Server) listMessages(sender net.Addr, channelName string) {

	if client, ok := server.clients[sender]; ok { // Manejamos que el cliente que solicita exista en el servidor

		if channel, ok := server.channels[channelName]; ok { // Manejamos que el canal solicitado exista

			if len(channel.messages) > 0 { // Verificamos si el canal tiene mensajes antes de proceder

				response := ""

				messages := make([]string, 0) // array de strings para ordenar los mensajes por su fecha

				for _, values := range channel.messages {

					//llenamos el array de mensajes con: fecha_mensaje,emisor,mensaje,file
					messages = append(messages, values.date.Format("2006-01-02:15:04:05")+","+values.sender.String()+","+string(values.content)+","+string(values.file))
				}

				sort.Strings(messages) //Ordenamos el array de mensajes

				for _, key := range messages {

					// Juntamos los mensajes del array ordenado en la respuesta dividido por ;
					// fecha_mensaje,emisor,mensaje;fecha_mensaje,emisor,mensaje
					response = response + key + ";"
				}

				client.WriteResponse(response)

			}
		}
	}
}

/* Funcion: listUsrChannel
 * Lista los usuarios pertenecientes a un canal
 * @param sender direccion del emisor de la solicitud.
 * @param channelName nombre del canal que se listaran sus clientes */
func (server *Server) listUsrChannel(sender net.Addr, channelName string) {

	if client, ok := server.clients[sender]; ok { // Manejamos que el clientes que solicita exista en el servidor

		if channel, ok := server.channels[channelName]; ok { // Manejamos que el canal solicitado exista

			if len(channel.clients) > 0 { // Verificamos si el canal tiene clientes antes de proceder

				response := ""
				clients := make([]string, 0) // array de strings para ordenar los clientes por su direccion

				for key, values := range channel.clients {

					if values { // Si el clientes esta conectado lo agregamos a los clientes
						clients = append(clients, key.address.String())
					}
				}

				sort.Strings(clients) //Ordenamos el array de clientes conectados

				for _, key := range clients {

					// Juntamos los clientes del array ordenado en la respuesta dividido por ;
					// cliente01;cliente02
					response = response + key + ";"
				}

				client.WriteResponse(response)

			}
		}
	}
}
