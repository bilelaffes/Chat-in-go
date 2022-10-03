package Server

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"strings"
)

type IServer interface {
	registration(string, string, net.Conn)
	connect(string, string)
	disconnect(string)
}

type Client struct {
	passWord string
	connect  net.Conn
}

type Server struct {
	clients map[string]*Client // tableau de clients
}

func checkPassWord(passWord string) error {
	if len(passWord) > 20 {
		return errors.New("Error taille password: must be at least 20 characters")
	}
	return nil
}

func (s *Server) setPassWord(name string, passWord string) error {
	if err := checkPassWord(passWord); err != nil {
		return err
	}
	s.clients[name].passWord = passWord
	return nil
}

func createClient(passWord string, conn net.Conn) (*Client, error) {
	var client Client

	if err := checkPassWord(passWord); err != nil {
		return nil, err
	}

	client.passWord = passWord
	client.connect = conn
	return &client, nil

}

func (s *Server) addClient(name string, client *Client) {
	s.clients[name] = client
}

func (s *Server) registration(name string, passWord string, conn net.Conn) string {
	client, err := createClient(passWord, conn)
	if err == nil {
		s.addClient(name, client)
		return name + " is registred\n"
	}
	return fmt.Sprintf("%s", err)
}

func (s *Server) connect(name string, passWord string) string {
	if s.clients[name].passWord == passWord {
		return name + " is connected\n"
	}
	return name + " error in password\n"
}

func (s *Server) disconnect(name string) string {
	if _, exist := s.clients[name]; exist {
		delete(s.clients, name)
		return name + " is disconnected\n"
	}
	return "Error in disconnect: " + name + " not found\n"
}

func (s *Server) sendMessage(message string, conn net.Conn) {
	var sender string

	for k, c := range s.clients {
		if c.connect == conn {
			sender = k
			fmt.Printf("%s Sending message", sender)
			break
		}
	}

	for _, c := range s.clients {

		if c.connect != conn {
			_, err := c.connect.Write([]byte(sender + " : " + message + "\n")) // on envoie un message à chaque client
			if err != nil {
				fmt.Println("Erreur to send message to client")
			}
		}
	}
}

func gestionErreur(err error) {
	if err != nil {
		panic(err)
	}
}

const (
	IPServer   = "127.0.0.01"
	PORTServer = "3569"
)

func read(s *Server, conn net.Conn) string {

	var ret string = ""
	message, err := bufio.NewReader(conn).ReadString('\n')
	gestionErreur(err)

	funcArg := strings.Split(message, ":")
	funcArg[len(funcArg)-1] = strings.Split(funcArg[len(funcArg)-1], "\n")[0]

	switch funcArg[0] {
	case "registration":
		ret = s.registration(funcArg[1], funcArg[2], conn)
		break
	case "connect":
		ret = s.connect(funcArg[1], funcArg[2])
		break
	case "disconnect":
		ret = s.disconnect(funcArg[1])
		break
	case "sendMessage":
		s.sendMessage(funcArg[1], conn)
		break
	}

	return ret
}

func InitServer() {

	fmt.Println("Lancement du serveur ...")

	s := &Server{make(map[string]*Client)}

	ln, err := net.Listen("tcp", fmt.Sprintf("%s:%s", IPServer, PORTServer))
	gestionErreur(err)

	for {
		conn, err := ln.Accept()
		gestionErreur(err)
		fmt.Println("Un client est connecté depuis", conn.RemoteAddr())

		go func() { // création de notre goroutine quand un client est connecté

			for {
				ret := read(s, conn)
				fmt.Println(ret)

				for _, c := range s.clients {

					if c.connect != conn {
						_, err := c.connect.Write([]byte(ret)) // on envoie un message à chaque client
						if err != nil {
							fmt.Println("Erreur to send message to client")
						}
					}
				}
			}
		}()
	}
}
