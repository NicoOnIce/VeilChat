# VeilChat
VeilChat is a peer to peer messaging platform, with end to end encryption, for private and decentralized conversations.

## There is an offical server running
You don't need to make your own rendezvous server, I have already made one! I am currently running it so you can use VeilChat for free!
If you are interested in coding your own rendezvous server or a modded client, contact me for the source code of the offical server.

## Contact me
Discord: https://discord.gg/5reyKquvfc

Discord account: hi.im.nico

## NOTICE
The I will no longer be maintaining, or updating the python version of VeilChat. There may be the occasional update, as it is currently in a very buggy state, but after most bugs have been patched, there won't be any UI or functionality updates performed on it. Feel free to contribute to it if you'd like.

I will be maintaining and updating the GOlang version of VeilChat, both in UI aand in functionality. Feel free to contribute to it if you'd like.

## How does VeilChat work?
VeilChat using UDP hole punching with STUN in order to create a connection between two machines, which are both behind NATs. This means that even without port forwarding two machines can still communicate directly with each other.

Upon a connection between client A and client B, client A will send over its encryption key (encryption key A) to client B and so will client B, to client A. These keys are called the public keys. After both client A and cliebt B have recieved each other's encryption keys, communication can begin. When client A types out a message and attempts to send it to client B, the message will first be encrypted with encryption key B. When client B recieves the encrypted data, it will be decrypted and output to the console.


## Why use VeilChat?
VeilChat is decentralized, meaning that instead of putting your trust in a company or server, you put your trust into the other client (who ever you are talking to.) This means that unless either you, or who ever you are speaking to reveal your private messages, it is almost impossible for anyone to read what ever you send to them, even if they got access to our routing server.

Ontop of this, all messages are end to end encrypted. Even if some one managed to find the data being transfered between you and who ever you are talking to, it would be almost impossible for them to see what you are sending to each other.
