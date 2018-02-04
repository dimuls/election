# election
This is example election service based on ethereum smart contract with web interface and command line control app. Written in [go](https://golang.org/) and [solidity](https://github.com/ethereum/solidity) using [geth](https://github.com/ethereum/go-ethereum/)

# Download

You can download it from [releases page](https://github.com/someanon/election/releases). Now is only windows x64 binary published.

# How to build

## Prerequirements 

* [go](https://golang.org/dl/)
* [glide](https://glide.sh/)
* [geth](https://github.com/ethereum/go-ethereum/wiki/Building-Ethereum)
* For Debians `apt install build-essential`
* For Red Hats `yum groupinstall 'Development Tools'` 
* For Windows building environment [like this](https://github.com/orlp/dev-on-windows/wiki/Installing-GCC--&-MSYS2)

## Building

Install dependencies first:
```
$ glide install
```

Next build it:
```
$ make
```

Now you will get `election` (or `election.exe` for windows) binary.

# How to use it

Almost all need information you can get using `-h` flag:
 
```
NAME:
   election.exe - Election smart contract controller

USAGE:
   election.exe [global options] command [command options] [arguments...]

VERSION:
   0.1.0

AUTHOR:
   Vadim Chernov <dimuls@yandex.ru>

COMMANDS:
     deploy      deploy election smart contract to the ethereum block chain
     add-voters  add voters to the election smart contract
     web-server  run web server with ui to vote in the election smart contract
     help, h     Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```
 
Besides you can learn sources yourself or [contact to me](#contacts).

# Contacts

Feel free to contact me:

* Email: dimuls@yandex.ru
* Telegram: @dimuls
* Skype: dimuls

