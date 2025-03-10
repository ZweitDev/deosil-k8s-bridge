//go:build ignore
// +build ignore

package main

import (
	"deosil-k8s-bridge/lib/k8s"
	"deosil-k8s-bridge/lib/k8s/minio"
	"encoding/json"
	"flag"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"k8s.io/client-go/kubernetes"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

var upgrader = websocket.Upgrader{}

func handler(clientset *kubernetes.Clientset) http.HandlerFunc {
return func (writer http.ResponseWriter, reader *http.Request) {
	c, err := upgrader.Upgrade(writer, reader, nil)

	if err != nil {
		log.Print("upgrade:", err)
		return
	}

	defer c.Close()

	log.Printf("Starting WebSocket Handler for new session")
	
	// TODO: Do not do this upon every client connecting...
	minio.CreatePVC(clientset)

	for {
		mt, message, err := c.ReadMessage()

		podsResult, err := k8s.GetPods(clientset, "kube-system")

		if err != nil {
			log.Println("read:", err)
			break
		}

		log.Printf("recv: %s", message)
		err = c.WriteMessage(mt, message)

		if err != nil {
			log.Println("write:", err)
			break
		}

		// Convert nodesResult ([]string) to JSON ([]byte)
		podsJSON, err := json.Marshal(podsResult)
		if err != nil {
			log.Println("error encoding nodes to JSON:", err)
			break
		}
		err = c.WriteMessage(mt, podsJSON)
		if err != nil {
			log.Println("write:", err)
			break
		}

		log.Printf("Wrote pods to client")
	}
}
}

func home(w http.ResponseWriter, r *http.Request) {
	homeTemplate.Execute(w, "ws://"+r.Host+"/echo")
}

func main() {
	clientset, err  := k8s.GetKubernetesClient()
	if err != nil {
		log.Fatalf("Failed to create Kubernetes client: %v", err)
	}
	flag.Parse()
	log.SetFlags(0)
	
	http.HandleFunc("/echo", handler(clientset))
	http.HandleFunc("/", home)
	log.Printf("Starting Deosil Server")
	log.Fatal(http.ListenAndServe(*addr, nil))
}

var homeTemplate = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<script>  
window.addEventListener("load", function(evt) {

    var output = document.getElementById("output");
    var input = document.getElementById("input");
    var ws;

    var print = function(message) {
        var d = document.createElement("div");
        d.textContent = message;
        output.appendChild(d);
        output.scroll(0, output.scrollHeight);
    };

    document.getElementById("open").onclick = function(evt) {
        if (ws) {
            return false;
        }
        ws = new WebSocket("{{.}}");
        ws.onopen = function(evt) {
            print("OPEN");
        }
        ws.onclose = function(evt) {
            print("CLOSE");
            ws = null;
        }
        ws.onmessage = function(evt) {
            print("RESPONSE: " + evt.data);
        }
        ws.onerror = function(evt) {
            print("ERROR: " + evt.data);
        }
        return false;
    };

    document.getElementById("send").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        print("SEND: " + input.value);
        ws.send(input.value);
        return false;
    };

    document.getElementById("close").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        ws.close();
        return false;
    };

});
</script>
</head>
<body>
<table>
<tr><td valign="top" width="50%">
<p>Click "Open" to create a connection to the server, 
"Send" to send a message to the server and "Close" to close the connection. 
You can change the message and send multiple times.
<p>
<form>
<button id="open">Open</button>
<button id="close">Close</button>
<p><input id="input" type="text" value="Hello world!">
<button id="send">Send</button>
</form>
</td><td valign="top" width="50%">
<div id="output" style="max-height: 70vh;overflow-y: scroll;"></div>
</td></tr></table>
</body>
</html>
`))
