package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/kelseyhightower/envconfig"
	"github.com/pipeduque/go-server/models"
	"github.com/pipeduque/go-server/tcpServer"
	"github.com/urfave/negroni"
)

var (
	config   configuration
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,

		// Chequeamos el origen de la solicitud, en este caso por trabajar en localhost se acepta simplemente
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

var server = models.NewServer()    //Establecemos el servidor para administrar el servicio
var tcp = tcpServer.NewTcpServer() // Establecemos el protocolo

type configuration struct {
	Debug         bool   `default:"true"`
	Scheme        string `default:"HTTP"`
	ListenAddress string `default:":8080"`
}

// Endpoint, necesitamos nuestro enrutador de respuesta y nuestro objeto de solicitud
func endpoint(writer http.ResponseWriter, reader *http.Request) {

	connection, err := upgrader.Upgrade(writer, reader, nil)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	defer connection.Close()

	// Bucle de escucha
	for {
		messageType, request, err := connection.ReadMessage()

		log.Println(string(request))

		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			break
		}

		if messageType != websocket.TextMessage {
			log.Fatal("Only text message are supported")
			break
		}

		log.Println(server.ServerOn)
		switch string(request) {
		case "serverTcpOn":
			if !server.ServerOn {
				server.ServerOn = true
				go startServerTcp(connection, writer, messageType)
				go tcp.Run(server)
				go server.Run() //corremos el servidor
			} else {
				server.ReqAndRes = "Server is on"
			}

		case "serverTcpOff":
			if server.ServerOn {
				server.ServerOn = false
				tcp.StopTcp()
				server.ReqAndRes = "Server off"
			}
		}
	}
}

func startServerTcp(connection *websocket.Conn, writer http.ResponseWriter, messageType int) {

	for server.ServerOn {
		response := server.ReqAndRes

		if response != "" {

			if err := connection.WriteMessage(messageType, []byte(response)); err != nil {
				log.Println(err)
				return
			}
			server.ReqAndRes = ""
		}
	}
}

/* Funcion
 * Nombre: main
 * Descripcion: Iniciamos nuestra */
func main() {

	// Variables de entorno para los ajustes de configuración
	err := envconfig.Process("SOCKETCAM", &config)
	if err != nil {
		log.Fatal(err.Error())
	}

	// Enrutador
	router := newRouter()
	n := negroni.Classic()

	n.UseHandler(router)

	// Informamos el inicio del servidor
	if config.Debug {
		log.Printf("==> PROTOCOL: %v", config.Scheme)
		log.Printf("==> ADDRESS: %v", config.ListenAddress)
	}

	// Dejamos al servidor escuchando en el puerto 8080 e informamos en caso de error
	log.Fatal(http.ListenAndServe(":8080", n))

}

/* Funcion
 * Nombre: newRouter
 * Descripcion: Constructor para todas las rutas */
func newRouter() *mux.Router {

	router := mux.NewRouter().StrictSlash(true)

	// Ruta de comunicacion cliente - servidor
	router.
		Methods("GET").
		Path("/ws").
		Name("Communication Channel").
		HandlerFunc(endpoint)

	// Ruta para enviar contenido al navegador
	router.
		Methods("GET").
		PathPrefix("/").
		Name("Static").
		Handler(http.FileServer(http.Dir("../static")))

	return router
}
