package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

type client chan<- string // canal de mensagem

var (
	entering = make(chan client)
	leaving  = make(chan client)
	messages = make(chan string)
	canal    = make(map[string]client)
	msgPV    = make(chan string)
)

func broadcaster() {
	clients := make(map[client]bool) // todos os clientes conectados
	for {
		select {
		case msg := <-messages:
			// broadcast de mensagens. Envio para todos
			for cli := range clients {
				cli <- msg
			}
		case cli := <-entering:
			clients[cli] = true
		case cli := <-leaving:
			delete(clients, cli)
			close(cli)
		case msg := <-msgPV:
			texto := strings.Split(msg, " ")
			msgEnv := false
			msgBot := ""

			for cli, _ := range clients {
				if msgEnv == false {
					if cli == canal[texto[1]] && texto[2] == "Bot" {
						msgBot = msgInv(texto[3])
						canal[texto[1]] <- texto[2] + ": " + msgBot
						msgEnv = true
					} else if cli == canal[texto[2]] && texto[2] != "Bot" {
						canal[texto[2]] <- texto[1] + ": " + texto[3]
						msgEnv = true
					}
				}
			}
		}
	}
}

func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg)
	}
}

func handleConn(conn net.Conn) {
	ch := make(chan string)
	go clientWriter(conn, ch)

	apelido := conn.RemoteAddr().String()
	ch <- "vc é " + apelido
	messages <- apelido + " chegou!"
	entering <- ch
	canal[apelido] = ch
	input := bufio.NewScanner(conn)
	for input.Scan() {
		texto := strings.Split(input.Text(), " ")
		if apelido == "Bot" {
			canal[apelido] <- "Você é um Bot! Não Faz Nada!"
			//Comando trocar Nome, deve ser digitado: /trocarNick Nome
		} else if texto[0] == "/trocarNick" {
			messages <- apelido + " tornou-se: " + texto[1]
			apelido = texto[1]
			canal[apelido] = ch
			//Comando Sair, deve ser digitado: /sair
		} else if texto[0] == "/sair" {
			leaving <- ch
			messages <- apelido + " se foi "
			return
			//Comando mensagem Privada, deve ser digitado: /msgPV SeuNome NomeDestinatario Mensagem
		} else if texto[0] == "/msgPV" {
			msgPV <- apelido + ":" + input.Text()
			//Para enviar mensagem para todos, basta digitar a mensagem
		} else {
			messages <- apelido + ":" + input.Text()
		}
	}

	conn.Close()
}

func main() {
	fmt.Println("Iniciando servidor...")
	listener, err := net.Listen("tcp", "localhost:3000")
	if err != nil {
		log.Fatal(err)
	}

	go broadcaster()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}

func msgInv(texto1 string) (texto2 string) {
	for _, i := range texto1 {
		texto2 = string(i) + texto2
	}
	return
}
