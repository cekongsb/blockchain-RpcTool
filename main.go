package main

import (
	"log"
)

const (
	SERVER_HOST       = "192.168.0.125"
	SERVER_PORT       = 9380
	USER              = "chain"
	PASSWD            = "999000"
	USESSL            = false
	WALLET_PASSPHRASE = "WalletPassphrase"
)

func main() {
	rpcClient, err := newClient(SERVER_HOST, SERVER_PORT, USER, PASSWD, USESSL)
	if err != nil {
		log.Fatalln(err)
	}
	//生成一个新地址
	reqJson := "{\"method\":\"getnewaddress\",\"params\":[\"labelName002\"],\"id\":1}"
	returnJson, err2 := rpcClient.send(reqJson)
	if err2 != nil {
		log.Fatalln(err2)
	}
	log.Println("returnJson:", returnJson)
}