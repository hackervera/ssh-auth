package main

import (
    "golang.org/x/crypto/ssh"
    "net"
    "log"
    "strings"
    "os"
    "fmt"
	"io/ioutil"
	"path/filepath"
    "github.com/Unknwon/com"
)

func main(){
	config := &ssh.ServerConfig{
		PublicKeyCallback: func(conn ssh.ConnMetadata, key ssh.PublicKey) (*ssh.Permissions, error) {

            //log.Printf(conn)
            keyStr := strings.TrimSpace(string(ssh.MarshalAuthorizedKey(key)))
            return &ssh.Permissions{Extensions: map[string]string{"key-id": keyStr}}, nil
		},
	}

	
	keyPath := filepath.Join("./ssh-auth.rsa")
	if !com.IsExist(keyPath) {
		os.MkdirAll(filepath.Dir(keyPath), os.ModePerm)
		_, stderr, err := com.ExecCmd("ssh-keygen", "-f", keyPath, "-t", "rsa", "-N", "")
		if err != nil {
			panic(fmt.Sprintf("Fail to generate private key: %v - %s", err, stderr))
		}
		log.Printf("SSH: New private key is generateed: %s", keyPath)
	}

	privateBytes, err := ioutil.ReadFile(keyPath)
	if err != nil {
		panic("SSH: Fail to load private key")
	}
	private, err := ssh.ParsePrivateKey(privateBytes)
	if err != nil {
		panic("SSH: Fail to parse private key")
	}
	config.AddHostKey(private)
    listener, err := net.Listen("tcp", "0.0.0.0:8080")
    if err != nil {
        log.Fatal("Can't listen on port")
    }
	go listen(config, listener)
	

    select{}
}


func listen(config *ssh.ServerConfig, listener net.Listener){
    for {
        conn, err := listener.Accept()
        if err != nil {
            log.Printf("Can't accept conn")
            continue
        }
        log.Printf("user connected")
        sshConn, chans, reqs, err := ssh.NewServerConn(conn, config)
        if err != nil {
            log.Printf("Can't create new server: ", err)
            continue
        }
        if chans != nil {}
        if reqs != nil {}
        log.Printf("New SSH connection from %s (%s) %s %s", sshConn.RemoteAddr(), sshConn.ClientVersion(), sshConn.User(), sshConn.Permissions.Extensions["key-id"])
        sshConn.Close()
    }
}
