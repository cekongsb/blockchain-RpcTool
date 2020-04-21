package main

import (
	"encoding/json"
	"log"
)

const (
	SERVER_HOST       = "192.168.8.107"
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
	// reqJson := "{\"method\":\"getnewaddress\",\"params\":[\"labelName002\"],\"id\":1}"
	// reqJson := "{\"method\":\"listunspent\",\"params\":[],\"id\":1}"

	// reqJson := "{\"method\":\"listunspent\",\"params\":[1, 9999999, [], true, {\"minimumSumAmount\":100}],\"id\":1}"
	reqJson := `{"method":"listunspent","params":[1, 9999999, [], true, {"minimumSumAmount":2}],"id":1}`

	returnJson, err2 := rpcClient.send(reqJson)
	if err2 != nil {
		log.Fatalln(err2)
	}
	log.Println("returnJson:", returnJson)

	var rpcresponse rpcResponse
	err3 := json.Unmarshal([]byte(returnJson), &rpcresponse)
	if err3 != nil {
		log.Fatalln(err3)
	}
	log.Println(string(rpcresponse.Result))

	var resultInfo interface{}
	err4 := json.Unmarshal(rpcresponse.Result, &resultInfo)
	if err4 != nil {
		log.Fatalln(err4)
	}
	log.Println(resultInfo)
	// log.Println(resultInfo)
}
