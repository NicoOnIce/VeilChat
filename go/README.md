# VeilChat
VeilChat is a peer to peer messaging platform, with end to end encryption, for private and decentralized conversations.

## There is an official server running
You don't need to make your own rendezvous server, I have already made one and I am currently running it so you can use it for free!
If you are interested in coding your own rendezvous server or a modified client, contact me for the source code of the official server.

## Contact me
Discord: https://discord.gg/5reyKquvfc

Discord account: hi.im.nico

## How does VeilChat work?
VeilChat using UDP hole punching with STUN in order to create a connection between two machines, which are both behind NATs. This means that even without port forwarding two machines can still communicate directly with each other.

Upon a connection between client A and client B, client A will send over its encryption key (encryption key A) to client B and so will client B, to client A. These keys are called the public keys. After both client A and client B have received each other's encryption keys, communication can begin. When client A types out a message and attempts to send it to client B, the message will first be encrypted with encryption key B. When client B receives the encrypted data, it will be decrypted and output to the console.


## Why use VeilChat?
VeilChat is decentralized, meaning that instead of putting your trust in a company or server, you put your trust into the other client (who ever you are talking to.) This means that unless either you, or who ever you are speaking to reveal your private messages, it is almost impossible for anyone to read what ever you send to them, even if they got access to our routing server.

On top of this, all messages are end to end encrypted. Even if some one managed to find the data being transfered between you and who ever you are talking to, it would be almost impossible for them to see what you are sending to each other.


## Want to contribute?
To contribute to VeilChat, create a pull request with the edits you have made. Your request will be examined and tested before being pushed to the main thread.


## Licensing notices
This software is under an "extended MIT - Non-Commercial Use Only" license.

This software uses third-party libraries under their respective licenses:
- cryptography (Apache 2.0)
- prompt_toolkit (BSD 3-Clause)


## How to run VeilChat
Currently, VeilChat must be run through python. It is compatible with any Python 3.X version (it has been tested on python 3.9, 3.10 and 3.14).

Steps to run VeilChat from release
- Download the latest release of VeilChat from [here](https://github.com/NicoOnIce/VeilChat/releases)
- Run it
- Once connected and authorized by a server, run the command `connect <CLIENT_NAME>` in order to connect and talk to another client.
- You can also use the `/disconnect` command, when in a chat in order to disconnect from the current client

Steps to run VeilChat from source
- Download GOlang 1.25.3 by clicking [here](https://go.dev/dl/)
- Then, download the VeilChat source code as a zip, and unzip it.
- Go into the ./go directory of the source code
- If on Windows, run the install.bat file, in order to install the required dependencies. Otherwise, run the command `go mod tidy`
- If on Windows, run the build.bat file to build the application. Otherwise, run the commands `cd <DIRECTORY_TO_VEILCHAT>` & `go build -o client main.go` in the same terminal

- Signup with a username and password. Once connected and authorized by the official server, you can connect and talk to others!
- You can use the `connect <username>` command to connect to another client. You can use the `/disconnect` command, when in a chat with a peer, to disconnect from the peer and go back to the official server. You can use the `clear` command to clear the console.

**Known bugs: When clients connect, first message doesn't go through**
**Packet loss: sometimes, packets don't go through! This isn't a bug, but I will work on making it more reliable.**
**If you encounter a bug while using VeilChat, please report it!**
