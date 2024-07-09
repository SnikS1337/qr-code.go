package backend_modules

import (
	"QR_CODE_GO/site_files"
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"log"
	"net/http"
)

// Addr set local address in main.go
var Addr = flag.String("addr", "localhost:8080", "http service address")

// set websocket
var upgrader = websocket.Upgrader{}

func Echo(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer func(c *websocket.Conn) {
		err := c.Close()
		if err != nil {
			log.Printf("close err: %v", err)
		}
	}(c)

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)

		// Echo the message back to the client
		err = c.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

func Home(w http.ResponseWriter, r *http.Request, globalSessions *sessions.Manager) {
	if globalSessions == nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	session := globalSessions.SessionStart(w, r)
	visited, ok := session.Get("visited").(bool)
	if !ok {
		session.Set("visited", true)
	}

	i := &site_files.IndexHTML

	tmpl, err := template.ParseFS(i, "index.html")
	if err != nil {
		log.Printf("Error parsing template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"WebSocketUrl": "ws://" + r.Host + "/echo",
		"SessionID":    session.SessionID(),
		"Visited":      visited,
	}
	tmpl.Execute(w, data)
}

func QrPostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", 405)
		return
	}

	globalSessions, ok := r.Context().Value("globalSessions").(*Manager)
	if !ok || globalSessions == nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	session := globalSessions.SessionStart(w, r)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", 500)
		return
	}
	defer r.Body.Close()

	_, err = fmt.Fprintf(w, "Received POST request. Body: %s\n", body)
	if err != nil {
		log.Printf("Error receiving POST: %v", err)
		http.Error(w, "Error writing response", 500)
		return
	}

	log.Printf("Received POST request. Body: %s\n", body)
	session.Set("lastPostBody", string(body))
}
