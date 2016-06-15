package ssh

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"golang.org/x/crypto/ssh"
)

func ConnectSSH(host string, user string, keypath string) *ssh.Session {

	key, err := ioutil.ReadFile(keypath)
	if err != nil {
		log.Fatalf("unable to read private key: %v", err)
	}

	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatalf("unable to parse private key: %v", err)
	}

	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			// Use the PublicKeys method for remote authentication.
			ssh.PublicKeys(signer),
		},
	}

	i := 0
	var client *ssh.Client
	for i < 60 {
		// Connect to the remote server and perform the SSH handshake.
		client, err = ssh.Dial("tcp", host+":22", config)
		if err == nil {
			break
		} else {
			fmt.Printf("Waiting to establish SSH...\n")
			time.Sleep(10 * time.Second)
			i += 10
		}
	}

	session, err := client.NewSession()
	if err != nil {
		log.Fatalf("Failed to create session: " + err.Error())
	}
	return session
}

func CloseSSH(session *ssh.Session) {
	session.Close()
}

func IssueCommand(session *ssh.Session, command string) {
	var stdout bytes.Buffer
	session.Stdout = &stdout
	err := session.Run(command)
	if err != nil {
		log.Fatalf("Failed to run: " + err.Error())
	}
	fmt.Println(stdout.String())
}
