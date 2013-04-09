package main

import (
	"bytes"
	"code.google.com/p/go.net/websocket"
	"compress/zlib"
	"encoding/json"
	zmq "github.com/alecthomas/gozmq"
	"io"
	"log"
	"net/http"
	//"os"
)

func main() {

	go relayListenerRoutine()

	http.Handle("/websocket/", websocket.Handler(handleWSConnection))
	err := http.ListenAndServe(":9000", nil)
	if err != nil {
		log.Fatal("%v", err)
	}
}

func relayListenerRoutine() {
	context, _ := zmq.NewContext()

	receiver, _ := context.NewSocket(zmq.SUB)
	receiver.SetSockOptString(zmq.SUBSCRIBE, "")
	//receiver.Connect("tcp://master.eve-emdr.com:8050")
	receiver.Connect("tcp://secondary.eve-emdr.com:8050")
	//receiver.Connect("tcp://relay-us-central-1.eve-emdr.com:8050")

	println("Listening on port 8050...")

	for {
		emdrMsg, emdrErr := receiver.Recv(0)

		if emdrErr != nil {
			println("EMDR error:", emdrErr.Error())
		}
		msgReader := bytes.NewReader(emdrMsg)

		r, zl_rr := zlib.NewReader(msgReader)
		if zl_rr != nil {
			println("ZL ERROR:", zl_rr.Error())
		}

		var out bytes.Buffer
		io.Copy(&out, r)
		r.Close()

		var f interface{}
		jsonErr := json.Unmarshal([]byte(out.String()), &f)
		if jsonErr != nil {
			println("JSON ERROR:", jsonErr.Error())
		}

		m := f.(map[string]interface{})

		resultType := m["resultType"]
		if resultType == "history" {
			continue
		}
		log.Printf("%T", m["rowsets"])

		t := m["rowsets"].([]interface{})
		for k, v := range t {
			log.Printf("%s:%T", k, v)
		}
	}

}

func handleWSConnection(ws *websocket.Conn) {
	go wsConnectionEchoRoutine(ws)
}

func wsConnectionEchoRoutine(ws *websocket.Conn) {

}
