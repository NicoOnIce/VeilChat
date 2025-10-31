# VeilChat
VeilChat is a peer to peer messaging platform, with end to end encryption, for private and decentralized conversations.


## There is an official server running
You don't need to make your own rendezvous server, I have already made one and I am currently running it so you can use it for free!
If you are interested in coding your own rendezvous server or a modified client, contact me for the source code of the official server.


## Contact me
Discord: https://discord.gg/5reyKquvfc

Discord account: hi.im.nico


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
