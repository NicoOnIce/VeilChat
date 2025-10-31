////////////////////////////////////////////////////////////////////////////////
//
//                  This software was created by
//                  The Online Anonymity Project
//
//               This software is being maintained by
//                  The Online Anonymity Project
//
//                           Licensing
//           Please read our licneses in the LICENSE file.
//                  This software is protected by a
//           "MIT extended - Non-Comercial Use Only" license
//
//                             Notice
//        Please do not use this software comercially, or for profit.
//    This software was created purely for public use, by everyone equally.
//
//////////////////////////////////////////////////////////////////////////////

package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/chzyer/readline"
	"golang.org/x/term"
)

const PORT = 5000

var (
	connected     = false
	privateKey    *rsa.PrivateKey
	publicKey     *rsa.PublicKey
	encryptionKey *rsa.PublicKey

	publicKeyB64 string

	servers = map[string]string{
		"Official server": "87.106.13.106:9999",
	}

	username string
	password string

	state = struct {
		peerEndpoint *net.UDPAddr
		lock         sync.Mutex
	}{}

	requests   = make(map[string]*net.UDPAddr)
	heartBeats = make([]*net.UDPAddr, 0)
	allow      = net.UDPAddr{}

	headers  []string = []string{"Version: VeilChat go-1.0.0", "Architecture: Unknown", "Github: https://github.com/NicoOnIce/VeilChat-Secure-P2P-Chat", "--------------------------------------------------------------------------"}
	messages []string

	_, height, _ = term.GetSize(int(os.Stdout.Fd()))
)

var rl, _ = readline.NewEx(&readline.Config{
	Prompt:          "> ",
	HistoryFile:     "/tmp/chat_history.tmp",
	InterruptPrompt: "^C",
	EOFPrompt:       "exit",
})

func keyToBase64(pub *rsa.PublicKey) (string, error) {
	der, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(der), nil
}

func getKey(b64 string) (*rsa.PublicKey, error) {
	der, err := base64.URLEncoding.DecodeString(b64)
	if err != nil {
		return nil, err
	}
	pk, err := x509.ParsePKIXPublicKey(der)
	if err != nil {
		return nil, err
	}
	rsaPub, ok := pk.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("parsed key is not RSA public key")
	}
	return rsaPub, nil
}

func encrypt(message string) (string, error) {
	if encryptionKey == nil {
		return "", errors.New("peer public key not set")
	}
	ct, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, encryptionKey, []byte(message), nil)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(ct), nil
}

func decrypt(message string) (string, error) {
	ct, err := base64.StdEncoding.DecodeString(message)
	if err != nil {
		return "", err
	}
	pt, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, ct, nil)
	if err != nil {
		return "", err
	}
	return string(pt), nil
}

func sendJSON(conn *net.UDPConn, addr *net.UDPAddr, obj map[string]interface{}) error {
	b, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	_, err = conn.WriteToUDP(b, addr)
	return err
}

func clearConsole() {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	} else {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func printHeaders() {
	for _, h := range headers {
		fmt.Println(h)
	}
	fmt.Println()
}

func safeClear() {
	messages = []string{}
}

func safeUpdate(line string) {
	messages = append(messages, line)

	clearConsole()
	printHeaders()

	maxMessages := height - len(headers) - 2
	start := 0
	if len(messages) > maxMessages {
		start = len(messages) - maxMessages
		// fmt.Println(start)
		// fmt.Println(len(messages))
	}

	for i := start; i < len(messages); i++ {
		fmt.Println(messages[i])
	}
}

// func refresh() {
// 	moveCursorTopLeft()
// 	fmt.Print("\033[J")

// 	for _, h := range headers {
// 		fmt.Println(h)
// 	}
// 	fmt.Println("")

// 	maxMessages := height - len(headers) - 2
// 	start := 0
// 	if len(messages) > maxMessages {
// 		start = len(messages) - maxMessages
// 	}
// 	for i := start; i < len(messages); i++ {
// 		fmt.Println(messages[i])
// 	}

// 	rl.Refresh()
// }

func moveCursorTopLeft() {
	fmt.Print("\033[H")
}

func safePrint(message string) {
	messages = append(messages, message)

	moveCursorTopLeft()
	fmt.Print("\033[J")

	for _, h := range headers {
		fmt.Println(h)
	}
	fmt.Println("")

	maxMessages := height - len(headers) - 2
	start := 0
	if len(messages) > maxMessages {
		start = len(messages) - maxMessages
	}
	for i := start; i < len(messages); i++ {
		fmt.Println(messages[i])
	}

	rl.Refresh()
}

func receiverLoop(conn *net.UDPConn) {
	buf := make([]byte, 8192)
	for {
		n, addr, err := conn.ReadFromUDP(buf)
		if err != nil {
			safePrint(fmt.Sprintf("SYS -> Receiver error: %v", err))
			return
		}

		data := buf[:n]
		var parsed map[string]interface{}
		if err := json.Unmarshal(data, &parsed); err == nil {
			cmd := fmt.Sprintf("%v", parsed["cmd"])

			switch cmd {
			case "peer":
				p, _ := parsed["addr"].(map[string]interface{})
				ip, _ := p["ip"].(string)
				portF := p["port"]
				port := int(portF.(float64))
				peerAddr := &net.UDPAddr{IP: net.ParseIP(ip), Port: port}

				headers[1] = "Peer: " + ip + ":" + strconv.Itoa(port)
				headers[2] = "Status: Connecting"

				safePrint(fmt.Sprintf("SYS -> Connecting to %s", peerAddr))
				go connect(conn, peerAddr)
				allow = *peerAddr
				state.lock.Lock()
				state.peerEndpoint = peerAddr
				state.lock.Unlock()

			case "key":
				if addr.IP.Equal(allow.IP) && addr.Port == allow.Port {
					encKeyStr, _ := parsed["encKey"].(string)
					if pk, err := getKey(encKeyStr); err == nil {
						encryptionKey = pk
						connected = true
						rl.SetPrompt("Me -> ")

						safeClear()
						safePrint(fmt.Sprintf("SYS -> Connected to %s", addr))
					}
				}

			case "msg":
				if addr.IP.Equal(allow.IP) && addr.Port == allow.Port {
					msgEnc, _ := parsed["msg"].(string)
					if msg, err := decrypt(msgEnc); err == nil {
						safePrint(fmt.Sprintf("%s -> %s", addr.IP.String(), msg))
					}
				}

			case "registered":
				name, _ := parsed["name"].(string)
				safePrint(fmt.Sprintf("SYS -> Registered as %s on server %s", name, addr))

			case "smsg":
				m, _ := parsed["msg"].(string)
				safePrint(fmt.Sprintf("SYS -> %s", m))

			case "peerreq":
				uname, _ := parsed["username"].(string)
				p, _ := parsed["addr"].(map[string]interface{})
				ip, _ := p["ip"].(string)
				portF := p["port"]
				port := int(portF.(float64))
				peerAddr := &net.UDPAddr{IP: net.ParseIP(ip), Port: port}
				safePrint(fmt.Sprintf("SYS -> %s (%s:%d) requested connection. Use 'accept %s'", uname, ip, port, uname))
				requests[uname] = peerAddr

			case "connected":
				headers[1] = "Server: " + addr.String()
				headers[2] = "Status: Authorized"
				safePrint(fmt.Sprintf("SYS -> Authorized on server %s", addr))

			case "heartBeat":
				timeSinceLastHeartbeat = 0
			}
		} else {
			if string(data) == "punch" && addr.IP.Equal(allow.IP) && addr.Port == allow.Port {
				connected = true
			} // else {
			// 	safePrint(fmt.Sprintf("SYS -> [RAW] %s from %s", string(data), addr))
			// }
		}
	}
}

func punchLoop(conn *net.UDPConn, peer *net.UDPAddr) {
	for !connected {
		conn.WriteToUDP([]byte("punch"), peer)
		time.Sleep(500 * time.Millisecond)
	}
	for i := 0; i < 5; i++ {
		conn.WriteToUDP([]byte("punch"), peer)
		time.Sleep(500 * time.Millisecond)
	}
}

func containsAddr(slice []*net.UDPAddr, addr *net.UDPAddr) bool {
	for _, a := range slice {
		if a.IP.Equal(addr.IP) && a.Port == addr.Port {
			return true
		}
	}
	return false
}

func heartBeatLoop(conn *net.UDPConn, peer *net.UDPAddr, hb map[string]interface{}) {
	ticker := time.NewTicker(20 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		if containsAddr(heartBeats, peer) {
			if err := sendJSON(conn, peer, hb); err != nil {
				safePrint("SYS -> Failed heartbeat: " + err.Error())
			}
		} else {
			return
		}
	}
}

var timeSinceLastHeartbeat = 0

func timeKeeper() {
	for {
		ticker := time.NewTicker(1 * time.Second)
		for range ticker.C {
			if state.peerEndpoint != nil {
				timeSinceLastHeartbeat++
				headers[2] = "Last Heartbeat: " + strconv.Itoa(timeSinceLastHeartbeat) + " seconds ago"
				fmt.Print("\033[s")
				fmt.Print("\033[?25l")
				fmt.Print("\033[H")
				fmt.Print("\033[2B")
				fmt.Print("\033[16C")
				fmt.Print(strconv.Itoa(timeSinceLastHeartbeat) + " ")
				fmt.Print("\033[u")
				fmt.Print("\033[?25h")
				rl.Refresh()
			}
		}
	}
}

// func addToMessages(message string) {
// 	messages = append(messages, message)
// }

func connect(conn *net.UDPConn, peer *net.UDPAddr) {
	punchLoop(conn, peer)
	sendJSON(conn, peer, map[string]interface{}{"cmd": "key", "encKey": publicKeyB64})

	heartBeats = append(heartBeats, peer)
	go heartBeatLoop(conn, peer, map[string]interface{}{"cmd": "heartBeat"})

	headers[0] = "Peer: " + peer.String()
	headers[1] = "Status: Connected"
	headers[2] = "Last Heartbeat: 0 seconds ago"
}

func main() {
	var err error
	privateKey, err = rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		safePrint("SYS -> RSA key gen failed: " + err.Error())
		return
	}
	publicKey = &privateKey.PublicKey
	publicKeyB64, _ = keyToBase64(publicKey)

	addr := net.UDPAddr{IP: net.ParseIP("0.0.0.0"), Port: PORT}
	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		safePrint("SYS -> UDP bind failed: " + err.Error())
		return
	}
	defer conn.Close()
	safePrint("SYS -> Local UDP socket bound to " + conn.LocalAddr().String())

	go timeKeeper()
	go receiverLoop(conn)

	serverSelected := false
	var serverAddr *net.UDPAddr

	defer rl.Close()

	internalConnected := false

	for {
		state.lock.Lock()
		peer := state.peerEndpoint
		state.lock.Unlock()

		if internalConnected {
			line, _ := rl.Readline()
			line = strings.TrimSpace(line)
			safeUpdate("Me -> " + line)

			if line == "/disconnect" {
				safePrint("SYS -> Disconnected from peer")
				internalConnected = false
				state.peerEndpoint = nil
				allow = net.UDPAddr{}
				heartBeats = heartBeats[:len(heartBeats)-1]
				rl.SetPrompt("> ")
				headers[0] = "Version: VeilChat go-1.0.0"
				headers[1] = "Server: " + serverAddr.String()
				headers[2] = "Status: Authorized"
				continue
			}
			if peer != nil {
				if enc, err := encrypt(line); err == nil {
					sendJSON(conn, peer, map[string]interface{}{"cmd": "msg", "msg": enc})
				}
			}
			continue
		}

		if !serverSelected {
			keys := make([]string, 0, len(servers))
			for k := range servers {
				keys = append(keys, k)
			}
			safePrint("=== Select a server ===")
			for i, k := range keys {
				safePrint(fmt.Sprintf("[%d] %s -> %s", i, k, servers[k]))
			}

			line, _ := rl.Readline()
			line = strings.TrimSpace(line)
			safeUpdate("> " + line)

			idx, _ := strconv.Atoi(line)
			server := servers[keys[idx]]

			sa := strings.Split(server, ":")
			port, _ := strconv.Atoi(sa[1])
			serverAddr = &net.UDPAddr{IP: net.ParseIP(sa[0]), Port: port}

			safeClear()
			safePrint("=== Login ===\n[0] connect\n[1] signup")
			modeLine, _ := rl.Readline()
			modeLine = strings.TrimSpace(modeLine)
			safeUpdate("> " + modeLine)
			safePrint("Username: ")
			username, _ = rl.Readline()
			safeUpdate("> " + username)
			safePrint("Password: ")
			password, _ = rl.Readline()
			safeUpdate("> " + password)
			safeClear()

			headers[1] = "Server: " + peer.String()
			headers[2] = "Status: Authorizing"

			if modeLine == "1" {
				safePrint(fmt.Sprintf("SYS -> Creating account on %s", serverAddr))
				sendJSON(conn, serverAddr, map[string]interface{}{"cmd": "create", "username": username, "password": password})
			}

			safePrint(fmt.Sprintf("SYS -> Connecting to server %s", serverAddr))
			sendJSON(conn, serverAddr, map[string]interface{}{"cmd": "connect", "username": username, "password": password})
			serverSelected = true
			heartBeats = append(heartBeats, serverAddr)
			go heartBeatLoop(conn, serverAddr, map[string]interface{}{"cmd": "heartBeat", "username": username, "password": password})
		} else {
			line, _ := rl.Readline()
			line = strings.TrimSpace(line)
			safeUpdate("> " + line)

			parts := strings.Split(line, " ")
			switch parts[0] {
			case "connect":
				rl.SetPrompt("Me -> ")
				internalConnected = true
				sendJSON(conn, serverAddr, map[string]interface{}{"cmd": "get", "username": username, "password": password, "peer": parts[1]})
			case "accept":
				uname := parts[1]
				peerAddr := requests[uname]
				safePrint(fmt.Sprintf("SYS -> Connecting to %s", peerAddr))

				internalConnected = true

				allow = *peerAddr
				connect(conn, peerAddr)

				heartBeats = append(heartBeats, peerAddr)
				go heartBeatLoop(conn, peerAddr, map[string]interface{}{"cmd": "heartBeat"})

				state.peerEndpoint = peerAddr
				rl.SetPrompt("Me -> ")
			case "disconnect":
				safePrint(fmt.Sprintf("SYS -> Disconnected from server %s", serverAddr))
				serverSelected = false
				heartBeats = make([]*net.UDPAddr, 0)
				rl.SetPrompt("> ")
			case "clear":
				safeClear()
				safePrint("> clear")
			}
		}
	}
}
