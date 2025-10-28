######################################################################
#
#                  This software was created by
#                  The Online Anonymity Project
#
#               This software is being maintained by
#                   The Online Anonymity Project
#
#                           Licensing
#           Please read our licneses in the LICENSE file.
#                  This software is protected by a
#           "MIT extended - Non-Comercial Use Only" license
#
#                             Notice
#        Please do not use this software comercially, or for profit.
#    This software was created purely for public use, by everyone equally.
#
######################################################################

import socket
import json
import threading
import time
import os

from prompt_toolkit import PromptSession
from prompt_toolkit.patch_stdout import patch_stdout
from threading import Thread
import time

from cryptography.hazmat.primitives.asymmetric import rsa, padding
from cryptography.hazmat.primitives import serialization, hashes
import base64

PORT = 5000

connected = False

privateKey = rsa.generate_private_key(public_exponent=65537, key_size=2048)
publicKey = privateKey.public_key()

encryptionKey = publicKey

def keyToBase64(publicKey) -> str:
    der = publicKey.public_bytes(
        encoding=serialization.Encoding.DER,
        format=serialization.PublicFormat.SubjectPublicKeyInfo
    )
    return base64.urlsafe_b64encode(der).decode('utf-8')

def getKey(b64Str):
    der = base64.urlsafe_b64decode(b64Str)
    return serialization.load_der_public_key(der)

def encrypt(message):
    if isinstance(message, str):
        message = message.encode('utf-8')

    ciphertext = encryptionKey.encrypt(
        message,
        padding.OAEP(
            mgf=padding.MGF1(algorithm=hashes.SHA256()),
            algorithm=hashes.SHA256(),
            label=None
        )
    )
    encoded_ciphertext = base64.b64encode(ciphertext).decode('utf-8')

    return encoded_ciphertext

def decrypt(message):
    message = base64.b64decode(message)
    plaintext = privateKey.decrypt(
        message,
        padding.OAEP(
            mgf=padding.MGF1(algorithm=hashes.SHA256()),
            algorithm=hashes.SHA256(),
            label=None
        )
    )

    return plaintext.decode('utf-8')

publicKeyb64 = keyToBase64(publicKey)
print(f"SYS -> Using public key: {publicKeyb64}\n")

servers = {
    "Official server": "87.106.13.106:9999"
}

username = ""
password = ""

state = {
    "peer_endpoint": None,
    "event": threading.Event(),
    "lock": threading.Lock()
}

requests = {

}

heartBeats = []

def clear():
    if os.name == "nt":
        os.system("cls")
    else:
        os.system("clear")

def send_json(sock, addr, obj):
    sock.sendto(json.dumps(obj).encode("utf-8"), addr)

def receiver_loop(sock):
    global encryptionKey, connected, state, heartBeats, requests
    while True:

        try:
            data, addr = sock.recvfrom(4096)
        except Exception as e:
            print("SYS -> Receiver socket error:", e)
            break

        parsed = None
        try:
            parsed = json.loads(data.decode("utf-8"))
        except Exception:
            pass

        if isinstance(parsed, dict) and parsed.get("cmd") == "peer":
            p = parsed.get("addr", {})
            ip = p.get("ip")
            port = p.get("port")

            try:
                peerAddr = (ip, int(port))

                def attemptConnection():
                    clear()
                    print("SYS -> Attempting connection to", peerAddr)
                    connect(sock, peerAddr)

                    heartBeats.append(peerAddr)
                    threading.Thread(target=heartBeatLoop, args=(sock, peerAddr, {"cmd": "heartBeat"}), daemon=True).start()
                    print("SYS -> Heartbeat loop started for", peerAddr)

                    state["peer_endpoint"] = peerAddr
                    print("SYS -> This is the start of your conversation with", peerAddr)
                    print("\n")
                
                threading.Thread(target=attemptConnection).start()
                
            except Exception as e:
                print("SYS -> Received malformed peer inforomation:", parsed, "error:", e)
                continue
        
        elif isinstance(parsed, dict) and parsed.get("cmd") == "key":
            encryptionKey = getKey(parsed.get("encKey"))
            connected = True

            print(f"SYS -> Recieved peer's public key")
        
        elif isinstance(parsed, dict) and parsed.get("cmd") == "msg":
            message = decrypt(parsed.get("msg"))

            print(f"{addr[0]} -> {message}")

        elif isinstance(parsed, dict) and parsed.get("cmd") == "heartBeat":
            pass
        
        elif isinstance(parsed, dict) and parsed.get("cmd") == "registered":
            print(f"SYS -> Registered self as", parsed.get("name"), "on server", addr)
        
        elif isinstance(parsed, dict) and parsed.get("cmd") == "smsg":
            print(f"SYS ->", parsed.get("msg"))

        elif isinstance(parsed, dict) and parsed.get("cmd") == "peerreq":
            print(f"SYS -> {parsed.get('username')} ({parsed.get('addr').get('ip')}, {parsed.get('addr').get('port')}) requested to connect with you.")
            print(f"SYS -> To estabilish a connection, run the command \"accept {parsed.get('username')}\"")

            requests[parsed.get("username")] = (parsed.get("addr").get("ip"), parsed.get("addr").get("port"))
        
        elif isinstance(parsed, dict) and parsed.get("cmd") == "connected":
            print(f"SYS -> You have been authorized on server", addr)

        elif isinstance(parsed, dict) and parsed.get("cmd") == "ok":
            pass

        else:
            try:
                if parsed is not None:
                    if parsed != "0" and parsed != b"0": 
                        print("SYS -> [RECV JSON]", parsed, "from", addr)
                else:
                    if data == b"punch":
                        connected = True
                    else:
                        print("SYS -> [RECV RAW]", data, "from", addr)
            except Exception as e:
                print("Receiver print error:", e)

def punch_loop(sock, peer_addr):
    while not connected:
        try:
            sock.sendto(b"punch", peer_addr)
        except Exception as e:
            print("SYS -> Punch send error:", e)
            break
        time.sleep(0.5)

    for i in range(0, 10):
        try:
            sock.sendto(b"punch", peer_addr)
        except Exception as e:
            print("SYS -> Punch send error:", e)
            break
        time.sleep(0.5)

def heartBeatLoop(sock, peer_addr, heartBeatObject):
    while peer_addr in heartBeats:
        try:
            send_json(sock,peer_addr, heartBeatObject)
        except:
            print("SYS -> Failed to send heartbeat packet to peer")
        
        time.sleep(20)

def connect(sock, peer_endpoint):
    punch_loop(sock, peer_endpoint)

    try:
        send_json(sock, peer_endpoint, {"cmd": "key", "encKey": publicKeyb64})
    except Exception as e:
        print("SYS -> Send error:", e)

def main():
    global connected, username, password, state

    server_addr = []

    sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
    sock.bind(("0.0.0.0", PORT))
    local = sock.getsockname()
    print(f"SYS -> Local UDP socket bound to {local}")

    threading.Thread(target=receiver_loop, args=(sock, ), daemon=True).start()

    serverSelected = False
    session = PromptSession()

    with patch_stdout():
        while True:
            if state["peer_endpoint"]:
                message = session.prompt("Me -> ")

                if message == "/disconnect":
                    clear()

                    message = encrypt("Has left the chat")
                    send_json(sock, state["peer_endpoint"], {"cmd": "msg", "msg": message})

                    state["peer_endpoint"] = None
                    del heartBeats[-1]

                else:
                    message = encrypt(message)
                    send_json(sock, state["peer_endpoint"], {"cmd": "msg", "msg": message})
            elif not serverSelected:
                print("\n=== Choose a routing server ===")
                for key in servers.keys():
                    print(f"[{list(servers.keys()).index(key)}] {key} - {servers[key]}")

                option = session.prompt("> ")
                parts = option.split(" ")
                server = servers[list(servers.keys())[int(parts[0])]]

                server_addr = server.split(":")
                server_addr = (server_addr[0], int(server_addr[1]))

                print("\n=== Login ===")
                print("[0] connect\n[1] signup")
                option = session.prompt("> ")

                if option == "0":
                    username = session.prompt("Username: ")
                    password = session.prompt("Password: ")

                    print("SYS -> Connecting to server", server_addr)
                    send_json(sock, server_addr, {"cmd":"connect", "username": username, "password": password})
                else:
                    username = session.prompt("Username: ")
                    password = session.prompt("Password: ")

                    print("SYS -> Creating account on server", server_addr)
                    send_json(sock, server_addr, {"cmd": "create", "username": username, "password": password})
                    print("SYS -> Connecting to server", server_addr)
                    send_json(sock, server_addr, {"cmd":"connect", "username": username, "password": password})

                heartBeats.append(server_addr)
                threading.Thread(target=heartBeatLoop, args=(sock, server_addr, {"cmd": "heartBeat", "username": username, "password": password})).start()
                serverSelected = True
            else:
                command = session.prompt("> ")
                parts = command.split(" ")

                if parts[0] == "connect":
                    print("SYS -> Sending connection request...")
                    send_json(sock, server_addr, {"cmd": "get", "username": username, "password": password, "peer": parts[1]})
                elif parts[0] == "disconnect":
                    print("SYS -> Disconnecting from server", server_addr)
                    del heartBeats[heartBeats.index(server_addr)]
                    print("SYS -> Disconnected from server", server_addr)

                    clear()

                    server_addr = ("0", 0)

                    serverSelected = False
                elif parts[0] == "accept":
                    username = parts[1]
                    peerAddr = requests[username]
                    clear()

                    print("SYS -> Attempting connection to", peerAddr)
                    connect(sock, peerAddr)

                    heartBeats.append(peerAddr)
                    threading.Thread(target=heartBeatLoop, args=(sock, peerAddr, {"cmd": "heartBeat"}), daemon=True).start()
                    print("SYS -> Heartbeat loop started for", peerAddr)

                    state["peer_endpoint"] = peerAddr
                    print("SYS -> This is the start of your conversation with", peerAddr)
                    print("\n")
                

if __name__ == "__main__":
    main()
