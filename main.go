package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"

	"golang.org/x/net/websocket"
)

func cp(w io.Writer, r io.Reader, errc chan<- error) {
	_, err := io.Copy(w, r)
	errc <- err
}

func chat(a, b io.ReadWriteCloser) {
	fmt.Fprintln(a, "Found one! Say hi.")
	fmt.Fprintln(b, "Found one! Say hi.")
	errc := make(chan error, 1)
	go cp(a, b, errc)
	go cp(b, a, errc)
	if err := <-errc; err != nil {
		log.Println(err)
	}
	a.Close()
	b.Close()
}

var partner = make(chan io.ReadWriteCloser)

func match(c io.ReadWriteCloser) {
	fmt.Fprint(c, "Waiting for a partner...")
	select {
	case partner <- c:
		return
	case p := <-partner:
		chat(p, c)
	}
}

const listenAddr = "jkaho.github.io"

func main() {
	http.HandleFunc("/", rootHandler)
	http.Handle("/socket", websocket.Handler(socketHandler))
	log.Println("http server running at: jkaho.github.io")
	err := http.ListenAndServe(listenAddr, nil)
	if err != nil {
		log.Fatal(err)
	}
}

type socket struct {
	io.ReadWriter
	done chan bool
}

func (s socket) Close() error {
	s.done <- true
	return nil
}

func socketHandler(ws *websocket.Conn) {
	s := socket{ws, make(chan bool)}
	go match(s)
	<-s.done
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	rootTemplate.Execute(w, listenAddr)
}

var rootTemplate = template.Must(template.ParseFiles("index.html"))

// func rootHandler(w http.ResponseWriter, r *http.Request) {
// 	rootTemplate.Execute(w, listenAddr)
// }

// var rootTemplate = template.Must(template.New("root").Parse(`
// <!DOCTYPE html>
// <html>
// <head>
// 	<meta charset="utf-8" />
// 	<title>Go chat</title>
// 	<link rel="stylesheet" href="https://use.fontawesome.com/releases/v5.8.1/css/all.css"
//     integrity="sha384-50oBUHEmvpQ+1lW4y57PTFmhCaXp0ML5d60M1M7uH2+nqUivzIebhndOJK28anvf" crossorigin="anonymous" />
// 	<style>
// 		body {
// 			box-sizing: border-box;
// 			font-family: Arial, Helvetica, sans-serif;
// 			padding: 20px;
// 		}
// 		#msgFormDiv {
// 			position: fixed;
// 			bottom: 20px;
// 			padding: 10px;
// 			width: 90%;
// 			border: 1px solid rgb(182, 90, 105);
// 			border-radius: 40px;
// 		}
// 		#msgForm {
// 			display: flex;
// 		}
// 		#msgBar {
// 			border: none;
// 			flex: 1;
// 		}
// 		#msgBar:focus {
// 			outline: none;
// 		}
// 		#msgSend {
// 			background: pink;
// 			color: white;
// 			padding: 10px;
// 			border-radius: 100%;
// 			border: none;
// 			margin-left: 5px;
// 		}
// 		#msgSend:hover {
// 			cursor: pointer;
// 		}
// 		#chatBoxOuter {
// 			display: none;
// 		}
// 		#chatBox {
// 			height: 400px;
// 			overflow: auto;
// 		}
// 		.msgLine {
// 			width: 100%;
// 			padding: 20px 0;
// 		}
// 		.msgBubble {
// 			padding: 5px 10px;
// 			border-radius: 10px;
// 		}
// 		.myMsgInner {
// 			float: right;
// 			background: pink;
// 		}
// 		.theirMsgInner {
// 			float: left;
// 			background: rgb(250, 250, 250);
// 		}
// 	</style>
// </head>
// <body>
// <script>
// 	const partnerFound = false;
// 	websocket = new WebSocket("ws://{{.}}/socket");
// 	websocket.addEventListener("message", function(e) {
// 		const msg = e.data;
// 		if (msg === "Found one! Say hi.\n") {
// 			document.getElementById("startBox").style.display = "none";
// 			document.getElementById("chatBoxOuter").style.display = "block";
// 			return;
// 		}
// 		const theirMsg = document.createElement("div");
// 		theirMsg.classList.add("msgLine");
// 		const theirMsgInner = document.createElement("div");
// 		theirMsgInner.classList.add("theirMsgInner");
// 		theirMsgInner.classList.add("msgBubble");
// 		theirMsgInner.textContent = msg;
// 		theirMsg.appendChild(theirMsgInner);
// 		document.getElementById("chatBox").appendChild(theirMsg);
// 	});
// 	websocket.addEventListener("close", function(e) {
// 		console.log("Connection has closed");
// 		const closeMsg = document.createElement("div");
// 		closeMsg.textContent = "Your partner has left the chat.";
// 		document.getElementById("chatBox").appendChild(closeMsg);
// 		document.getElementById("serverMsg").textContent = "";
// 	});
// </script>
// <div id="startBox"><p id="serverMsg">Waiting for a chat partner...</p></div>
// <div id="chatBoxOuter">
// 	<div id="chatBox"></div>
// 	<div id="msgFormDiv">
// 		<form id="msgForm">
// 				<input id="msgBar" type="text" placeholder="Start typing here..."/>
// 				<button id="msgSend"><i class="fas fa-paper-plane"></i></button>
// 		</form>
// 	</div>
// </div>
// <script defer>
// 	document.getElementById("msgForm").addEventListener("submit", function(e) {
// 		e.preventDefault();
// 		const msg = document.getElementById("msgBar").value;
// 		websocket.send(msg);
// 		const myMsg = document.createElement("div");
// 		myMsg.classList.add("msgLine");
// 		const myMsgInner = document.createElement("div");
// 		myMsgInner.classList.add("myMsgInner");
// 		myMsgInner.classList.add("msgBubble");
// 		myMsgInner.textContent = msg;
// 		myMsg.appendChild(myMsgInner);
// 		document.getElementById("chatBox").appendChild(myMsg);
// 		document.getElementById("msgBar").value = "";
// 	});
// </script>
// </body>
// </html>
// `))
