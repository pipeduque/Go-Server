package tcpServer

/* flag: bandera que sirve de utilidad para la linea de comandos
 * ftm: Imprementa Entradas / Salidas similares a C */
import (
	"flag"
	"log"
	"net"

	"github.com/pipeduque/go-server/models"
)

type TcpServer struct {
	listener net.Listener
	adrress  string
	network  string
}

/* Funcion
 * Nombre: NewTcpServer
 * Descripcion: Iniciamos un oyente TCP, para recibir una conexicion TCP con el cliente */
func NewTcpServer() *TcpServer {

	var address string //Variable para la direccion de escucha del el servidor
	var network string //Variable para el protocolo de red

	flag.StringVar(&address, "e", ":3000", "Service Endpoint [ip address]") //Bandera que analiza la direccion, vinculado a la variable address
	flag.StringVar(&network, "n", "tcp", "network protocol [tpc]")          //Bandera que analiza el protocolo de red, vinculado a la variable network
	flag.Parse()                                                            //Analizamos las banderas

	switch network { //Validamos que el protocolo de red sea soportado
	case "tcp", "tcp4", "tcp6":
	default:
		log.Fatalln("Unsupported network protocol: ", network)
	}

	//Conexion sockets
	listen, err := net.Listen(network, address) //Creamos el oyente para el protocolo de red proporcionado y la dirección de host

	if err != nil { // Manejamos un posible error al crear el oyente
		log.Fatal("Failed to create listener ", err)
	}

	return &TcpServer{
		listener: listen,
		adrress:  address,
		network:  network,
	}
}

/* Funcion
 * Nombre: Run
 * Descripcion: Corremos nuestro protocolo tcp, para recibir conexiones con clientes */
func (tcpServer *TcpServer) Run(server *models.Server) {

	defer tcpServer.listener.Close()

	server.ReqAndRes = "Server started (" + tcpServer.network + ") " + tcpServer.adrress //Informamos la iniciacion del servidor

	//Ciclo de conexion - Maneja las solicitudes entrantes
	for server.ServerOn {
		connection, err := tcpServer.listener.Accept() //Usamos el oyente con el punto aceptar para crear la conexion, que se bloqueará hasta que llegue una conexion con el cliente

		if err != nil { //Manejamos un posible error al crear la conexion
			log.Fatal(err)
			if err := connection.Close(); err != nil {
				log.Fatal("Failed to close listener:", err)
			}
			continue
		}

		server.ReqAndRes = "Connected to " + connection.RemoteAddr().String()

		client := models.NewClient(connection, server) // Referenciamos al nuevo cliente que creo la conexion
		server.ClientOnlineReq <- client               // Lo conectamos al servidor

		go client.RequestReadHandle() //LLamamos al manejador de lectura de solicitudes del cliente
	}

}

func (tcpServer *TcpServer) StopTcp() {

	tcpServer.listener.Close()
}
