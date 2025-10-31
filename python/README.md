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
- Download python 3.14 by clicking [here](https://www.python.org/downloads/release/python-3140/)
- If you are on Windows, When installing python 3.14, make sure to click add to PATH in the installer
- Then, download the VeilChat source code as a zip, and unzip it.
- If on Windows, run the install.bat file, in order to install the required dependencies. Otherwise, run the command `pip install -r requirements.txt`
- If on Windows, run the run.bat file to run the application. Otherwise, run the commands `cd <DIRECTORY_TO_VEILCHAT>` & `python client.py` in the same terminal
- 
- Once connected and authorized by a server, run the command `connect <CLIENT_NAME>` in order to connect and talk to another client.
- You can also use the `/disconnect` command, when in a chat in order to disconnect from the current client.

**If you encounter a bug while using VeilChat, please report it!**

