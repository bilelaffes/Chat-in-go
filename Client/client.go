package Client

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"sync"
)

func gestionErreur(err error) {
	if err != nil {
		panic(err)
	}
}

const (
	IPClient   = "127.0.0.01" // IP local
	PORTClient = "3569"       // Port utilisé
)

func read(conn net.Conn) {

	message, _ := bufio.NewReader(conn).ReadString('\n')
	fmt.Print(message)
}

func InitClient() {

	var wg sync.WaitGroup

	// Connexion au serveur
	conn, _ := net.Dial("tcp", fmt.Sprintf("%s:%s", IPClient, PORTClient))

	wg.Add(2)

	go func() { // goroutine dédiée à l'entrée utilisateur
		defer wg.Done()
		for {
			reader := bufio.NewReader(os.Stdin)
			text, _ := reader.ReadString('\n')

			conn.Write([]byte(text))
		}
	}()

	go func() { // goroutine dédiée à la reception des messages du serveur
		defer wg.Done()
		for {
			read(conn)
		}
	}()

	wg.Wait()

}
