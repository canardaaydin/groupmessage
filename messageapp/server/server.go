package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
	"strings"
)


func main() {

	//server begin



	delivery := new(Delivery)
	rpc.Register(delivery)

	ip := os.Args[1] + ":1234"


	l, e := net.Listen("tcp", ":1234")

	if e != nil {
		log.Fatal("listen error:", e)
	}

	// similar to server code
	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				continue
			}
			rpc.ServeConn(conn)
		}
	}()


	//server end

	//client begin

	file, err := os.Open("group.txt")

	if err != nil {
		log.Fatalf("failed to open %s", err.Error())

	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var text []string

	for scanner.Scan() {

		if scanner.Text() != ip {
			text = append(text, scanner.Text())
		}
	}

	file.Close()



	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("\n-> ")
		txt, _ := reader.ReadString('\n')
		// convert CRLF to LF
		txt = strings.Replace(txt, "\n", "", -1)


		//new message broadcasts to every peer


		var clients []*rpc.Client

		for _, ip := range text {

			fmt.Print(ip)
			client, err := rpc.Dial("tcp", ip)
			if err != nil {
				log.Fatal("dialing:", err)
			}
			clients = append(clients, client)
		}

		for _, cli := range clients{

			args := &Args{txt, ip}
			resp := new(Response)
			divCall := cli.Go("Delivery.MessagePost", args, resp, nil)
			_ = <-divCall.Done // will be equal to divCall
			cli.Close()

		}




	}
	//client end

}

type Args struct {
	Content string
	Sender string
}

type Response struct {
	Content string
}

type Delivery int


func (t *Delivery) MessagePost(message *Args, response *Response) error {
	fmt.Printf("Incoming message from %s : %s", message.Sender, message.Content)
	return nil
}
