/*
window.addEventListener("load", function (evt) {

  var output = document.getElementById("output");
  var input = document.getElementById("input");
  var ws;

  var print = function (message) {
    var d = document.createElement("div");
    d.innerHTML = message;
    output.appendChild(d);
  };

  document.getElementById("open").onclick = function (evt) {
    if (ws) {
      return false;
    }
    var loc = window.location, new_uri;
    if (loc.protocol === "https:") {
      new_uri = "wss:";
    } else {
      new_uri = "ws:";
    }
    new_uri += "//" + loc.host;
    new_uri += loc.pathname + "ws";
    ws = new WebSocket(new_uri);
    ws.onopen = function (evt) {
      print("OPEN");
    }
    ws.onclose = function (evt) {
      print("CLOSE");
      ws = null;
    }
    ws.onmessage = function (evt) {
      print("RESPONSE: " + evt.data);
    }
    ws.onerror = function (evt) {
      print("ERROR: " + evt.data);
    }
    return false;
  };

  document.getElementById("send").onclick = function (evt) {
    if (!ws) {
      return false;
    }
    print("SEND: " + input.value);
    ws.send(input.value);
    return false;
  };

  document.getElementById("close").onclick = function (evt) {
    if (!ws) {
      return false;
    }
    ws.close();
    return false;
  };
});
*/

new Vue({

    // Elemento del html donde trabajar
    el: '#app',

    // Al crearse ejecutamos la funcion para conectarnos con el WebSocket
    created() {
        this.connectToWebSocket();
        console.log("conectado")
    },

    // Datos para el elemento html correspondiente
    data: {

        ws: null,
        reqAndRes: [],
    },

    // Metodos
    methods: {
        connectToWebSocket() {
            console.log("conectando")
            if (this.ws) {
                return false;
            }
            var loc = window.location,
                new_uri;
            if (loc.protocol === "https:") {
                new_uri = "wss:";
            } else {
                new_uri = "ws:";
            }
            new_uri += "//" + loc.host;
            new_uri += loc.pathname + "ws";
            this.ws = new WebSocket(new_uri);
            this.ws.onopen = function(evt) {
                console.log("OPEN");
            }

            this.ws.onclose = function(evt) {
                console.log("close")
                ws = null;
            }

            this.ws.onmessage = (evt) => {

                let date = new Date();
                let arrayResponses = evt.data.split(";;");
                let div = document.getElementById('console');
                for (i in arrayResponses) {
                    this.reqAndRes.push({
                        text: arrayResponses[i],
                        date: date.toLocaleDateString() + " " + date.toLocaleTimeString()
                    })
                    div.scrollTop = '9999';
                }
            }
            this.ws.onerror = function(evt) {
                console.log("ERROR: " + evt.data);
            }
            return false;
        },

        on() {

            if (!this.ws) {
                return false;
            }
            this.ws.send("serverTcpOn");
        },

        off() {

            if (!this.ws) {
                return false;
            }
            this.ws.send("serverTcpOff");
        }
    }
})