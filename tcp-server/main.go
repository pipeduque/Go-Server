package main

/* flag: bandera que sirve de utilidad para la linea de comandos
 * ftm: Imprementa Entradas / Salidas similares a C */
import (
	"flag"
	"log"
	"net"

	"github.com/pipeduque/go-server/models"
)

/* Funcion
 * Nombre: main
 * Descripcion: En la funcion main iniciamos un oyente TCP, para recibir una conexicion TCP con el cliente */
func main() {

	var address string //Variable para la direccion de escucha del el servidor
	var network string //Variable para el protocolo de red

	flag.StringVar(&address, "e", ":3000", "Service Endpoint [ip address]") //Bandera que analiza la direccion, vinculado a la variable address
	flag.StringVar(&network, "n", "tcp", "network protocol [tpc]")          //Bandera que analiza el protocolo de red, vinculado a la variable network
	flag.Parse()                                                            //Analizamos las banderas

	server := models.NewServer() //Establecemos el servidor para administrar el servicio y lo corremos
	go server.Run()

	switch network { //Validamos que el protocolo de red sea soportado
	case "tcp", "tcp4", "tcp6":
	default:
		log.Fatalln("Unsupported network protocol: ", network)
	}

	//Conexion sockets
	listener, err := net.Listen(network, address) //Creamos el oyente para el protocolo de red proporcionado y la dirección de host

	if err != nil { // Manejamos un posible error al crear el oyente
		log.Fatal("Failed to create listener ", err)
	}
	defer listener.Close()

	log.Printf("Server started (%s) %s", network, address) //Informamos la iniciacion del servidor

	//Ciclo de conexion - Maneja las solicitudes entrantes
	for {

		connection, err := listener.Accept() //Usamos el oyente con el punto aceptar para crear la conexion, que se bloqueará hasta que llegue una conexion con el cliente

		if err != nil { //Manejamos un posible error al crear la conexion
			log.Println(err)
			if err := connection.Close(); err != nil {
				log.Println("Failed to close listener:", err)
			}
			continue
		}

		log.Println("Connected to", connection.RemoteAddr())

		client := models.NewClient(connection, server) // Referenciamos al nuevo cliente que creo la conexion
		server.ClientOnlineReq <- client               // Lo conectamos al servidor

		go client.RequestReadHandle() //LLamamos al manejador de lectura de solicitudes del cliente
	}
}
