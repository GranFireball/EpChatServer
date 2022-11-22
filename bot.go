package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"sync"
)

func mustCopy(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatal(err)
	}
}

func initializeBot(wg *sync.WaitGroup, conn net.Conn) {
	conn.Write([]byte("/trocarNick Bot\n"))
	defer wg.Done() // Diminuo o contador em 1
}

func main() {
	wg := new(sync.WaitGroup)
	wg.Add(1) // Adiciono 1 ao contador
	conn, err := net.Dial("tcp", "localhost:3000")
	fmt.Println("Connected!")
	if err != nil {
		log.Fatal(err)
	}

	done := make(chan struct{})

	go initializeBot(wg, conn)
	wg.Wait() // Bloqueio a execução até o contador chegar a 0

	go func() {
		io.Copy(os.Stdout, conn)
		log.Println("done")
		done <- struct{}{} // sinaliza para a gorrotina principal
	}()
	mustCopy(conn, os.Stdin)
	conn.Close()
	<-done // espera a gorrotina terminar
}
