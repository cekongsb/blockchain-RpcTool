package main

import (
	"encoding/json"
	"log"
)

const (
	SERVER_HOST       = "192.168.0.143"
	SERVER_PORT       = 6062
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
	// reqJson := `{"method":"listunspent","params":[1, 9999999, [], true, {"minimumSumAmount":40}],"id":1}`
	reqJson := `{"method":"listunspent","params":[1, 9999999, ["3PX9raYTM5MZRQahhxikPCsGcfvBvbFcWg"], true, {"minimumSumAmount":40}],"id":1}`

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
	log.Println("len(resultInfo) =", len(resultInfo.([]interface{})))
	for i := 0; i < len(resultInfo.([]interface{})); i++ {
		log.Println("address:", resultInfo.([]interface{})[i].(map[string]interface{})["address"])
		log.Println("amount:", resultInfo.([]interface{})[i].(map[string]interface{})["amount"])
		log.Println("scriptPubKey:", resultInfo.([]interface{})[i].(map[string]interface{})["scriptPubKey"])
		log.Println("txid:", resultInfo.([]interface{})[i].(map[string]interface{})["txid"])
		log.Println("vout:", resultInfo.([]interface{})[i].(map[string]interface{})["vout"])
	}

	reqJson = `{"method":"createrawtransaction","params":[[{"txid":"0f0bcb033c9ff17405ef34eb966a2bf2119243532703d30c4aae6b69146e0e00","vout":2}], [{"STCwWTDAd3vYEG5MBL1R5mju8qu3vMTXD8":16},{"3PX9raYTM5MZRQahhxikPCsGcfvBvbFcWg":5.04251795}]],"id":1}`

	returnJson2, err5 := rpcClient.send(reqJson)
	if err5 != nil {
		log.Fatalln(err5)
	}
	log.Println("returnJson2:", returnJson2)
	// CreateTx(rpcClient)

	reqJson = `{"method":"signrawtransactionwithkey","params":["0200000001000e6e14696bae4a0cd3032753439211f22b6a96eb34ef0574f19f3c03cb0b0f02000000b500473044022041c97005d3719ca346bb710239a526d1204eda700b43cfd60131f6158005003702204533181e29e2863fe2382af1f40f3c18120e70a52cb990e53df714586c30f4fd01004c69522102a45ecbe752f01a863a89ac2267058b25d2273be19ee63bc59b4b7d6faf3091f72102bf39c73ee2e8e64c8f4af1ae3e85fd1e9a987b635c5024ddfbf7be41785eec3921036320132c8c66c2d8766e48d1c34487d2e41cd9d27a51a7ca53c9ece7a8eab0be53aeffffffff0200105e5f000000001976a91440dbbbb627083697c314f63c737f7e951555cf4a88ac93450e1e0000000017a914ef753a876709cc4a834d52d2090425e7525007908700000000",["KyVzoDL4UbM3PmMaGgSZRWuADq3Au5R9wC8oBYajxtsY8eEtaMcv"],[{"txid":"0f0bcb033c9ff17405ef34eb966a2bf2119243532703d30c4aae6b69146e0e00","vout":2,"scriptPubKey":"a914ef753a876709cc4a834d52d2090425e75250079087","redeemScript":"522102a45ecbe752f01a863a89ac2267058b25d2273be19ee63bc59b4b7d6faf3091f72102bf39c73ee2e8e64c8f4af1ae3e85fd1e9a987b635c5024ddfbf7be41785eec3921036320132c8c66c2d8766e48d1c34487d2e41cd9d27a51a7ca53c9ece7a8eab0be53ae","amount":21.04261795}]],"id":1}`
	returnJson3, err6 := rpcClient.send(reqJson)
	if err6 != nil {
		log.Fatalln(err6)
	}
	log.Println("returnJson3:", returnJson3)

	reqJson = `{"method":"signrawtransactionwithkey","params":["0200000001000e6e14696bae4a0cd3032753439211f22b6a96eb34ef0574f19f3c03cb0b0f02000000b500473044022041c97005d3719ca346bb710239a526d1204eda700b43cfd60131f6158005003702204533181e29e2863fe2382af1f40f3c18120e70a52cb990e53df714586c30f4fd01004c69522102a45ecbe752f01a863a89ac2267058b25d2273be19ee63bc59b4b7d6faf3091f72102bf39c73ee2e8e64c8f4af1ae3e85fd1e9a987b635c5024ddfbf7be41785eec3921036320132c8c66c2d8766e48d1c34487d2e41cd9d27a51a7ca53c9ece7a8eab0be53aeffffffff0200105e5f000000001976a91440dbbbb627083697c314f63c737f7e951555cf4a88ac93450e1e0000000017a914ef753a876709cc4a834d52d2090425e7525007908700000000",["Ky4CzdZ6VJsHFfH8DwFvk488XgzDb6g4ijqZmYxhL3M4iFyuZdAx"],[{"txid":"0f0bcb033c9ff17405ef34eb966a2bf2119243532703d30c4aae6b69146e0e00","vout":2,"scriptPubKey":"a914ef753a876709cc4a834d52d2090425e75250079087","redeemScript":"522102a45ecbe752f01a863a89ac2267058b25d2273be19ee63bc59b4b7d6faf3091f72102bf39c73ee2e8e64c8f4af1ae3e85fd1e9a987b635c5024ddfbf7be41785eec3921036320132c8c66c2d8766e48d1c34487d2e41cd9d27a51a7ca53c9ece7a8eab0be53ae","amount":21.04261795}]],"id":1}`
	returnJson4, err7 := rpcClient.send(reqJson)
	if err7 != nil {
		log.Fatalln(err7)
	}
	log.Println("returnJson4:", returnJson4)

	reqJson = `{"method":"sendrawtransaction","params":["0200000001000e6e14696bae4a0cd3032753439211f22b6a96eb34ef0574f19f3c03cb0b0f02000000fc00473044022041c97005d3719ca346bb710239a526d1204eda700b43cfd60131f6158005003702204533181e29e2863fe2382af1f40f3c18120e70a52cb990e53df714586c30f4fd0147304402203d12bd93756aad8d43550e1d6bc1b6e6929d208ef2d7d32dc18988e6a5ff521a02200bf5e1e68f73a676eeea601d80f713b69707d388c486f29d913db56772482881014c69522102a45ecbe752f01a863a89ac2267058b25d2273be19ee63bc59b4b7d6faf3091f72102bf39c73ee2e8e64c8f4af1ae3e85fd1e9a987b635c5024ddfbf7be41785eec3921036320132c8c66c2d8766e48d1c34487d2e41cd9d27a51a7ca53c9ece7a8eab0be53aeffffffff0200105e5f000000001976a91440dbbbb627083697c314f63c737f7e951555cf4a88ac93450e1e0000000017a914ef753a876709cc4a834d52d2090425e7525007908700000000"],"id":1}`
	returnJson5, err8 := rpcClient.send(reqJson)
	if err8 != nil {
		log.Fatalln(err8)
	}
	log.Println("returnJson8:", returnJson5)
}

// func CreateTx(rpcClient *rpcClient) {
// 	reqJson := `{"method":"createrawtransaction","params":["[{"txid":"0f0bcb033c9ff17405ef34eb966a2bf2119243532703d30c4aae6b69146e0e00","vout":2}]","[{"STCwWTDAd3vYEG5MBL1R5mju8qu3vMTXD8":16},{"3PX9raYTM5MZRQahhxikPCsGcfvBvbFcWg":5.04251795‬}]"],"id":1}`
// 	returnJson, err2 := rpcClient.send(reqJson)
// 	if err2 != nil {
// 		log.Fatalln(err2)
// 	}
// 	log.Println("returnJson:", returnJson)
// }
