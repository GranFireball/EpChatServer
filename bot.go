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

func main() {
	conn, err := net.Dial("tcp", "localhost:3000")
	fmt.Println("Connected!")
	if err != nil {
		log.Fatal(err)
	}

	done := make(chan struct{})

	cont := new(sync.WaitGroup)
	cont.Add(1)

	go bot(cont, conn)
	cont.Wait()

	go func() {
		io.Copy(os.Stdout, conn)
		log.Println("done")
		done <- struct{}{} // sinaliza para a gorrotina principal
	}()
	mustCopy(conn, os.Stdin)
	conn.Close()
	<-done // espera a gorrotina terminar
}

func bot(cont *sync.WaitGroup, conn net.Conn) {
	conn.Write([]byte("Bot: "))
	defer cont.Done()
}
