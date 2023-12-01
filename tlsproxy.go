// Generate private key (.key)
// Key considerations for algorithm "RSA" ≥ 2048-bit
//
//	openssl genrsa -out server.key 2048
//
// Key considerations for algorithm "ECDSA" (X25519 || ≥ secp384r1)
// https://safecurves.cr.yp.to/
// List ECDSA the supported curves (openssl ecparam -list_curves)
//
//	openssl ecparam -genkey -name secp384r1 -out server.key
//
// Generation of self-signed(x509) public key (PEM-encodings .pem|.crt) based on the private (.key)
//
//	openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650
package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"syscall"
)

var targetAddress string
var printDebug bool

func copyConn(reader io.Reader, writer io.Writer, done chan bool) {
	buffer := make([]byte, 512)
	r := bufio.NewReader(reader)

	for {
		n, err := r.Read(buffer)
		if err != nil {
			if err != io.EOF && printDebug {
				log.Printf("Read error: %v", err)
			}
			break
		}

		if n > 0 {
			_, err = writer.Write(buffer[0:n])
			if err != nil && printDebug {
				log.Printf("Write error: %v", err)
				break
			}
		}
	}
	done <- true
}

func handleConnection(clientConn net.Conn, targetPort int) {
	if printDebug {
		log.Printf("Handle new connection %v", clientConn.RemoteAddr())
	}
	defer clientConn.Close()

	target := fmt.Sprintf("%s:%d", targetAddress, targetPort)
	targetConn, err := net.Dial("tcp", target)
	if err != nil {
		log.Printf("Failed to connect to %s: %v", target, err)
		return
	}
	defer targetConn.Close()
	done := make(chan bool, 1)

	go copyConn(clientConn, targetConn, done)
	go copyConn(targetConn, clientConn, done)

	<-done
	if printDebug {
		log.Printf("Closed connection %v", clientConn.RemoteAddr())
	}
}

func startSocket(clientPort int, targetPort int, certFile string, keyFile string) {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		log.Fatal("LoadingCertKeyPair", err)
	}

	config := &tls.Config{Certificates: []tls.Certificate{cert}}

	ln, err := tls.Listen("tcp", fmt.Sprintf(":%d", clientPort), config)
	if err != nil {
		log.Fatal("ListenAndServeSocket", err)
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handleConnection(conn, targetPort)
	}
}

func main() {
	log.Printf("TLS proxy version %s", Version)

	debugEnv := strings.ToLower(os.Getenv("TLSPROXY_DEBUG"))

	printDebug = debugEnv == "1" || debugEnv == "true" || debugEnv == "enable" || debugEnv == "on"

	targetEnv := os.Getenv("TLSPROXY_TARGET")
	if targetEnv == "" {
		log.Fatal("Error: missing or empty TLSPROXY_TARGET env var\n")
	}
	targetAddress = targetEnv

	portsEnv := os.Getenv("TLSPROXY_PORTS")
	if portsEnv == "" {
		log.Fatal("Error: missing or empty TLSPROXY_PORTS env var\n")
	}

	certFileEnv := os.Getenv("TLSPROXY_CERT_FILE")
	if certFileEnv == "" {
		certFileEnv = filepath.Join("certificates", "server.crt")
		log.Printf("Missing or empty TLSPROXY_CERT_FILE, using '%s' as default value", certFileEnv)
	}

	keyFileEnv := os.Getenv("TLSPROXY_KEY_FILE")
	if keyFileEnv == "" {
		keyFileEnv = filepath.Join("certificates", "server.key")
		log.Printf("Missing or empty TLSPROXY_KEY_FILE, using '%s' as default value", keyFileEnv)
	}

	r := regexp.MustCompile(`([0-9]+)(:([0-9]+))?`)
	match := r.FindAllStringSubmatch(portsEnv, -1)

	// Loop over all matches and start listeners for each one
	for _, m := range match {
		if m[3] == "" {
			port, _ := strconv.Atoi(m[1])
			go startSocket(port, port, certFileEnv, keyFileEnv)
		} else {
			clientPort, _ := strconv.Atoi(m[1])
			targetPort, _ := strconv.Atoi(m[3])
			go startSocket(clientPort, targetPort, certFileEnv, keyFileEnv)
		}
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	done := make(chan bool, 1)

	go func() {
		sig := <-sigs
		log.Printf("Quit (%v)", sig)
		done <- true
	}()

	<-done
}
