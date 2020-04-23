package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
)

const (
	SERVER_HOST       = "192.168.8.101"
	SERVER_PORT       = 6062
	USER              = "chain"
	PASSWD            = "999000"
	USESSL            = false
	WALLET_PASSPHRASE = "WalletPassphrase"
)

func makeReqJson(arg ...interface{}) string {
	var stuReqJson rpcRequest

	stuReqJson.JsonRpc = "1.0"
	stuReqJson.Id = arg[0].(string)
	stuReqJson.Method = arg[1].(string)
	stuReqJson.Params = arg[2:]

	// fmt.Println("stuReqJson.Method=", stuReqJson.Method)
	// fmt.Println("stuReqJson.Params=", stuReqJson.Params)

	reqJson, err := json.Marshal(stuReqJson)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(stuReqJson.Method, "reqJson=", string(reqJson))
	return string(reqJson)
}

func decodeJson(resultInfo string) rpcResponse {
	var rpcresponse rpcResponse
	err := json.Unmarshal([]byte(resultInfo), &rpcresponse)
	if err != nil {
		log.Fatalln(err)
	}
	// fmt.Println("rpcresponse.Err=", rpcresponse.Err)
	// fmt.Println("rpcresponse.Id=", rpcresponse.Id)
	// fmt.Println("rpcresponse.Result=", string(rpcresponse.Result))
	return rpcresponse
}

func selectCommand(rpcresponse *rpcResponse) interface{} {

	return nil
}

func parseListunspent(rpcresponse *rpcResponse) (interface{}, interface{}, float64) {
	uint8Listunspent := rpcresponse.Result

	var listunspentInfo interface{}
	err := json.Unmarshal(uint8Listunspent, &listunspentInfo)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("listunspentInfo=", listunspentInfo)
	fmt.Println("len(listunspentInfo)=", len(listunspentInfo.([]interface{})))

	var slicInputTx []map[string]interface{}
	var slicListunspent []map[string]interface{}
	var sumAcount float64 = 0
	for i := 0; i < len(listunspentInfo.([]interface{})); i++ {
		mapInput := make(map[string]interface{})
		mapInput["txid"] = listunspentInfo.([]interface{})[i].(map[string]interface{})["txid"]
		mapInput["vout"] = listunspentInfo.([]interface{})[i].(map[string]interface{})["vout"]

		mapListunspent := make(map[string]interface{})
		mapListunspent["scriptPubKey"] = listunspentInfo.([]interface{})[i].(map[string]interface{})["scriptPubKey"]
		mapListunspent["amount"] = listunspentInfo.([]interface{})[i].(map[string]interface{})["amount"]
		mapListunspent["txid"] = listunspentInfo.([]interface{})[i].(map[string]interface{})["txid"]
		mapListunspent["vout"] = listunspentInfo.([]interface{})[i].(map[string]interface{})["vout"]
		mapListunspent["redeemScript"] = "522102a45ecbe752f01a863a89ac2267058b25d2273be19ee63bc59b4b7d6faf3091f72102bf39c73ee2e8e64c8f4af1ae3e85fd1e9a987b635c5024ddfbf7be41785eec3921036320132c8c66c2d8766e48d1c34487d2e41cd9d27a51a7ca53c9ece7a8eab0be53ae"

		slicInputTx = append(slicInputTx, mapInput)
		slicListunspent = append(slicListunspent, mapListunspent)
		sumAcount += listunspentInfo.([]interface{})[i].(map[string]interface{})["amount"].(float64)
	}
	fmt.Println("sumAcount=", sumAcount)
	return slicInputTx, slicListunspent, sumAcount
}
func listunspent(arg ...interface{}) string {
	minconf := arg[0]
	maxconf := arg[1]
	address := arg[2]
	isSafe := arg[3]
	queryOption := arg[4:][0] //此处可待优化

	reqJson := makeReqJson("listunspent", "listunspent", minconf, maxconf, address, isSafe, queryOption)
	return reqJson
}

func createrawtransaction(arg ...interface{}) string {
	// reqJson = `{"method":"createrawtransaction","params":[[{"txid":"0f0bcb033c9ff17405ef34eb966a2bf2119243532703d30c4aae6b69146e0e00","vout":2}], [{"STCwWTDAd3vYEG5MBL1R5mju8qu3vMTXD8":16},{"3PX9raYTM5MZRQahhxikPCsGcfvBvbFcWg":5.04251795}]],"id":1}`
	// {"method":"createrawtransaction","params":[[{"txid":"18bf31178fd79fa6a98d03f7a29c6acbb15760642dd812bfdef0c0f4257a0a00","vout":2},{"txid":"eb36490e5b034f9ffcdbf5674ea825379723aa643e4d76c4628f49f8673b1400","vout":2}],[{"3PX9raYTM5MZRQahhxikPCsGcfvBvbFcWg":3.0851359,"SWX9b3z4K47fLq8vK6x7FUM7XTj6r9Gr1V":39}]],"id":"createrawtransaction","jsonrpc":"1.0"}

	inputs := arg[0]
	outputs := arg[1]
	// fmt.Println("inputs=", inputs)
	// fmt.Println("outputs=", outputs)
	reqJson := makeReqJson("createrawtransaction", "createrawtransaction", inputs, outputs)
	return reqJson
}

func parseCreaterawtransaction(rpcresponse *rpcResponse) interface{} {
	uint8Createrawtransaction := rpcresponse.Result
	var CreaterawtransactionInfo interface{}
	err := json.Unmarshal(uint8Createrawtransaction, &CreaterawtransactionInfo)
	if err != nil {
		log.Fatalln(err)
	}
	// fmt.Println("CreaterawtransactionInfo=", CreaterawtransactionInfo)
	return CreaterawtransactionInfo
}

func signrawtransactionwithkey(arg ...interface{}) string {

	txHex := arg[0]
	privKey := arg[1]
	listunspentInfo := arg[2]
	reqJson := makeReqJson("signrawtransactionwithkey", "signrawtransactionwithkey", txHex, privKey, listunspentInfo)
	// fmt.Println("reqJson=", reqJson)
	return reqJson
}

func parseSignrawtransactionwithkey(rpcresponse *rpcResponse) interface{} {
	uint8SignTxKey := rpcresponse.Result
	var signhex signHex
	err := json.Unmarshal(uint8SignTxKey, &signhex)
	if err != nil {
		log.Fatalln(err)
	}
	return signhex.Hex
}

func sendrawtransaction(arg ...interface{}) string {
	hex := arg[0]
	reqJson := makeReqJson("sendrawtransaction", "sendrawtransaction", hex)
	// fmt.Println("reqJson=", reqJson)
	return reqJson
}

func sendRpcRequest(client *rpcClient, reqJson string) string {
	returnJson, err := client.send(reqJson)
	if err != nil {
		log.Fatalln(err)
	}
	return returnJson
}

func main() {
	rpcClient, err := newClient(SERVER_HOST, SERVER_PORT, USER, PASSWD, USESSL)
	if err != nil {
		log.Fatalln(err)
	}

	//命令:listunspent
	sliAddr := []string{"3PX9raYTM5MZRQahhxikPCsGcfvBvbFcWg"}
	mapParam := map[string]interface{}{"minimumSumAmount": 520, "minimumAmount": 1}

	reqListunspent := listunspent(1, 99999, sliAddr, true, mapParam)
	returnJson := sendRpcRequest(rpcClient, reqListunspent)
	log.Println("returnJson=", returnJson)

	rpcresponse := decodeJson(returnJson)
	inputs, slicListunspent, sumAcount := parseListunspent(&rpcresponse)
	fmt.Println("inputs=", inputs)

	//命令:createrawtransaction
	distAddr := "SYYD1P2bXYxvWRM9Eh1mANcp3JS6HS5eLF"
	var distAcount float64 = 0
	changeAddr := "3PX9raYTM5MZRQahhxikPCsGcfvBvbFcWg"
	distAcount = 500
	tempChangeAcount := sumAcount - distAcount - 0.001
	changeAcount, _ := strconv.ParseFloat(fmt.Sprintf("%.8f", tempChangeAcount), 64)

	mapDistout := make(map[string]interface{})
	mapChangeout := make(map[string]interface{})
	mapDistout[distAddr] = distAcount
	mapChangeout[changeAddr] = changeAcount
	outputs := []map[string]interface{}{}
	outputs = append(outputs, mapDistout, mapChangeout)

	reqCreaterawtransaction := createrawtransaction(inputs, outputs)
	returnJson = sendRpcRequest(rpcClient, reqCreaterawtransaction)
	fmt.Println("returnJson_cretx=", returnJson)
	rpcresponse = decodeJson(returnJson)
	txHex := parseCreaterawtransaction(&rpcresponse)
	fmt.Println("txHex=", txHex)

	//命令:signrawtransactionwithkey
	// 签名1：
	privKey1 := []string{"KyVzoDL4UbM3PmMaGgSZRWuADq3Au5R9wC8oBYajxtsY8eEtaMcv"}
	reqSignrawtransactionwithkey1 := signrawtransactionwithkey(txHex, privKey1, slicListunspent)
	returnJson = sendRpcRequest(rpcClient, reqSignrawtransactionwithkey1)
	fmt.Println("returnJson_sign1=", returnJson)
	rpcresponse = decodeJson(returnJson)
	signHex1 := parseSignrawtransactionwithkey(&rpcresponse)
	//签名2：
	privKey2 := []string{"Ky4CzdZ6VJsHFfH8DwFvk488XgzDb6g4ijqZmYxhL3M4iFyuZdAx"}
	reqSignrawtransactionwithkey2 := signrawtransactionwithkey(signHex1, privKey2, slicListunspent)
	returnJson = sendRpcRequest(rpcClient, reqSignrawtransactionwithkey2)
	fmt.Println("returnJson_sign2=", returnJson)
	rpcresponse = decodeJson(returnJson)
	signHex2 := parseSignrawtransactionwithkey(&rpcresponse)
	fmt.Println("signHex2=", signHex2)

	//命令:sendrawtransaction
	reqSendrawtransaction := sendrawtransaction(signHex2)
	returnJson = sendRpcRequest(rpcClient, reqSendrawtransaction)
	fmt.Println("returnJson_brocast=", returnJson)
}

// func main() {
// 	rpcClient, err := newClient(SERVER_HOST, SERVER_PORT, USER, PASSWD, USESSL)
// 	if err != nil {
// 		log.Fatalln(err)
// 	}
// 	//生成一个新地址
// 	// reqJson := "{\"method\":\"getnewaddress\",\"params\":[\"labelName002\"],\"id\":1}"
// 	// reqJson := "{\"method\":\"listunspent\",\"params\":[],\"id\":1}"

// 	// reqJson := "{\"method\":\"listunspent\",\"params\":[1, 9999999, [], true, {\"minimumSumAmount\":100}],\"id\":1}"
// 	// reqJson := `{"method":"listunspent","params":[1, 9999999, [], true, {"minimumSumAmount":40}],"id":1}`
// 	reqJson := `{"method":"listunspent","params":[1, 9999999, ["3PX9raYTM5MZRQahhxikPCsGcfvBvbFcWg"], true, {"minimumSumAmount":40}],"id":1}`

// 	returnJson, err2 := rpcClient.send(reqJson)
// 	if err2 != nil {
// 		log.Fatalln(err2)
// 	}
// 	log.Println("returnJson:", returnJson)

// 	var rpcresponse rpcResponse
// 	err3 := json.Unmarshal([]byte(returnJson), &rpcresponse)
// 	if err3 != nil {
// 		log.Fatalln(err3)
// 	}
// 	log.Println(string(rpcresponse.Result))

// 	var resultInfo interface{}
// 	err4 := json.Unmarshal(rpcresponse.Result, &resultInfo)
// 	if err4 != nil {
// 		log.Fatalln(err4)
// 	}
// 	log.Println(resultInfo)
// 	log.Println("len(resultInfo) =", len(resultInfo.([]interface{})))
// 	for i := 0; i < len(resultInfo.([]interface{})); i++ {
// 		log.Println("address:", resultInfo.([]interface{})[i].(map[string]interface{})["address"])
// 		log.Println("amount:", resultInfo.([]interface{})[i].(map[string]interface{})["amount"])
// 		log.Println("scriptPubKey:", resultInfo.([]interface{})[i].(map[string]interface{})["scriptPubKey"])
// 		log.Println("txid:", resultInfo.([]interface{})[i].(map[string]interface{})["txid"])
// 		log.Println("vout:", resultInfo.([]interface{})[i].(map[string]interface{})["vout"])
// 	}

// 	reqJson = `{"method":"createrawtransaction","params":[[{"txid":"0f0bcb033c9ff17405ef34eb966a2bf2119243532703d30c4aae6b69146e0e00","vout":2}], [{"STCwWTDAd3vYEG5MBL1R5mju8qu3vMTXD8":16},{"3PX9raYTM5MZRQahhxikPCsGcfvBvbFcWg":5.04251795}]],"id":1}`

// 	returnJson2, err5 := rpcClient.send(reqJson)
// 	if err5 != nil {
// 		log.Fatalln(err5)
// 	}
// 	log.Println("returnJson2:", returnJson2)

// 	reqJson = `{"method":"signrawtransactionwithkey","params":["0200000001000e6e14696bae4a0cd3032753439211f22b6a96eb34ef0574f19f3c03cb0b0f02000000b500473044022041c97005d3719ca346bb710239a526d1204eda700b43cfd60131f6158005003702204533181e29e2863fe2382af1f40f3c18120e70a52cb990e53df714586c30f4fd01004c69522102a45ecbe752f01a863a89ac2267058b25d2273be19ee63bc59b4b7d6faf3091f72102bf39c73ee2e8e64c8f4af1ae3e85fd1e9a987b635c5024ddfbf7be41785eec3921036320132c8c66c2d8766e48d1c34487d2e41cd9d27a51a7ca53c9ece7a8eab0be53aeffffffff0200105e5f000000001976a91440dbbbb627083697c314f63c737f7e951555cf4a88ac93450e1e0000000017a914ef753a876709cc4a834d52d2090425e7525007908700000000",["KyVzoDL4UbM3PmMaGgSZRWuADq3Au5R9wC8oBYajxtsY8eEtaMcv"],[{"txid":"0f0bcb033c9ff17405ef34eb966a2bf2119243532703d30c4aae6b69146e0e00","vout":2,"scriptPubKey":"a914ef753a876709cc4a834d52d2090425e75250079087","redeemScript":"522102a45ecbe752f01a863a89ac2267058b25d2273be19ee63bc59b4b7d6faf3091f72102bf39c73ee2e8e64c8f4af1ae3e85fd1e9a987b635c5024ddfbf7be41785eec3921036320132c8c66c2d8766e48d1c34487d2e41cd9d27a51a7ca53c9ece7a8eab0be53ae","amount":21.04261795}]],"id":1}`
// 	returnJson3, err6 := rpcClient.send(reqJson)
// 	if err6 != nil {
// 		log.Fatalln(err6)
// 	}
// 	log.Println("returnJson3:", returnJson3)

// 	reqJson = `{"method":"signrawtransactionwithkey","params":["0200000001000e6e14696bae4a0cd3032753439211f22b6a96eb34ef0574f19f3c03cb0b0f02000000b500473044022041c97005d3719ca346bb710239a526d1204eda700b43cfd60131f6158005003702204533181e29e2863fe2382af1f40f3c18120e70a52cb990e53df714586c30f4fd01004c69522102a45ecbe752f01a863a89ac2267058b25d2273be19ee63bc59b4b7d6faf3091f72102bf39c73ee2e8e64c8f4af1ae3e85fd1e9a987b635c5024ddfbf7be41785eec3921036320132c8c66c2d8766e48d1c34487d2e41cd9d27a51a7ca53c9ece7a8eab0be53aeffffffff0200105e5f000000001976a91440dbbbb627083697c314f63c737f7e951555cf4a88ac93450e1e0000000017a914ef753a876709cc4a834d52d2090425e7525007908700000000",["Ky4CzdZ6VJsHFfH8DwFvk488XgzDb6g4ijqZmYxhL3M4iFyuZdAx"],[{"txid":"0f0bcb033c9ff17405ef34eb966a2bf2119243532703d30c4aae6b69146e0e00","vout":2,"scriptPubKey":"a914ef753a876709cc4a834d52d2090425e75250079087","redeemScript":"522102a45ecbe752f01a863a89ac2267058b25d2273be19ee63bc59b4b7d6faf3091f72102bf39c73ee2e8e64c8f4af1ae3e85fd1e9a987b635c5024ddfbf7be41785eec3921036320132c8c66c2d8766e48d1c34487d2e41cd9d27a51a7ca53c9ece7a8eab0be53ae","amount":21.04261795}]],"id":1}`
// 	returnJson4, err7 := rpcClient.send(reqJson)
// 	if err7 != nil {
// 		log.Fatalln(err7)
// 	}
// 	log.Println("returnJson4:", returnJson4)

// 	reqJson = `{"method":"sendrawtransaction","params":["0200000001000e6e14696bae4a0cd3032753439211f22b6a96eb34ef0574f19f3c03cb0b0f02000000fc00473044022041c97005d3719ca346bb710239a526d1204eda700b43cfd60131f6158005003702204533181e29e2863fe2382af1f40f3c18120e70a52cb990e53df714586c30f4fd0147304402203d12bd93756aad8d43550e1d6bc1b6e6929d208ef2d7d32dc18988e6a5ff521a02200bf5e1e68f73a676eeea601d80f713b69707d388c486f29d913db56772482881014c69522102a45ecbe752f01a863a89ac2267058b25d2273be19ee63bc59b4b7d6faf3091f72102bf39c73ee2e8e64c8f4af1ae3e85fd1e9a987b635c5024ddfbf7be41785eec3921036320132c8c66c2d8766e48d1c34487d2e41cd9d27a51a7ca53c9ece7a8eab0be53aeffffffff0200105e5f000000001976a91440dbbbb627083697c314f63c737f7e951555cf4a88ac93450e1e0000000017a914ef753a876709cc4a834d52d2090425e7525007908700000000"],"id":1}`
// 	returnJson5, err8 := rpcClient.send(reqJson)
// 	if err8 != nil {
// 		log.Fatalln(err8)
// 	}
// 	log.Println("returnJson8:", returnJson5)
// }
