package main

import (
	"golang.org/x/crypto/ssh"
	"io"
	"log"
	"net"
)

func remoteForward(remoteConn net.Conn, localAddr string) {
	localConn, err := net.Dial("tcp", localAddr)
	if err != nil {
		log.Fatalf("net.Dial failed: %s", err)
	}

	// Copy localConn.Reader to sshConn.Writer
	go func() {
		_, err = io.Copy(remoteConn, localConn)
		if err != nil {
			log.Fatalf("io.Copy failed: %v", err)
		}
	}()

	// Copy sshConn.Reader to localConn.Writer
	go func() {
		_, err = io.Copy(localConn, remoteConn)
		if err != nil {
			log.Fatalf("io.Copy failed: %v", err)
		}
	}()
}

func main() {
	username := "root"
	password := "password"
	serverAddrString := "192.168.1.100:22"
	localAddrString := "localhost:9000"
	remoteAddrString := "localhost:9999"

	config := &ssh.ClientConfig{
		User:            username,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
	}

	sshClientConn, err := ssh.Dial("tcp", serverAddrString, config)
	if err != nil {
		log.Fatalf("ssh.Dial failed: %s", err)
	}

	sshListen, err := sshClientConn.Listen("tcp", remoteAddrString)
	if err != nil {
		log.Fatalf("ssh_client_conn listen failed: %s", err)
	}

	for {
		remoteConn, err := sshListen.Accept()
		if err != nil {
			log.Fatalf("listen.Accept failed: %v", err)
		}
		go remoteForward(remoteConn, localAddrString)
	}
}
