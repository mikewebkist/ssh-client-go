package main

import (
	"bufio"
	"code.google.com/p/go.crypto/ssh"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	fmt.Print("Remote host? (Default=localhost): ")
	server := scanConfig()
	if server == "" {
		server = "localhost"
	}
	fmt.Print("Port? (Default=22): ")
	port := scanConfig()
	if port == "" {
		port = "22"
	}
	server = server + ":" + port
	fmt.Print("Username? (Default=root): ")
	user := scanConfig()
	if user == "" {
		user = "root"
	}
	fmt.Print("Password?: ")
	p := scanConfig()

	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{ssh.Password(p)},
	}
	conn, err := ssh.Dial("tcp", server, config)
	if err != nil {
		panic("Failed to dial: " + err.Error())
	}
	defer conn.Close()

	// Each ClientConn can support multiple interactive sessions,
	// represented by a Session.
	session, err := conn.NewSession()
	if err != nil {
		panic("Failed to create session: " + err.Error())
	}
	defer session.Close()

	// Set IO
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	in, _ := session.StdinPipe()
	out, _ := session.StdoutPipe()

	// Start remote shell
	if err := session.Shell(); err != nil {
		log.Fatalf("failed to start shell: %s", err)
	}

	fmt.Fprint(in, "unset HISTFILE\n")

	// Accepting commands
	for {
		reader := bufio.NewReader(os.Stdin)
		str, _ := reader.ReadString('\n')
		fmt.Printf("[%s@%s] $ ", user, server)
		fmt.Fprint(in, str)
	}

}

func scanConfig() string {
	config, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	config = strings.Trim(config, "\n")
	return config
}
