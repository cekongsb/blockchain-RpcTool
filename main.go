package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

const (
	SERVER_HOST       = "192.168.0.125"
	SERVER_PORT       = 9380
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
		mapListunspent["redeemScript"] = "5221027ae577fa58d2309dee977ed5acc7d630af967409a7bb2630dba864031eecf05c21030c465532f86d29fcba1e87edc22225ab5220636ff4f594fecb046f164ef27718210371cb6debd35cd607f37ea5debefe4941b2d99c8ac2983d2d7f66206ff687069153ae"

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
	inputs := arg[0]
	outputs := arg[1]
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
	return CreaterawtransactionInfo
}

func signrawtransactionwithkey(arg ...interface{}) string {
	txHex := arg[0]
	privKey := arg[1]
	listunspentInfo := arg[2]
	reqJson := makeReqJson("signrawtransactionwithkey", "signrawtransactionwithkey", txHex, privKey, listunspentInfo)
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
	return reqJson
}

func sendRpcRequest(client *rpcClient, reqJson string) string {
	returnJson, err := client.send(reqJson)
	if err != nil {
		log.Fatalln(err)
	}
	return returnJson
}

func makeMultisigTx(distaddr string, fundaddr string, sendacount float64, signkey []string) {
	rpcClient, err := newClient(SERVER_HOST, SERVER_PORT, USER, PASSWD, USESSL)
	if err != nil {
		log.Fatalln(err)
	}
	//命令:listunspent
	// sliAddr := []string{"3PX9raYTM5MZRQahhxikPCsGcfvBvbFcWg"}
	sliAddr := []string{fundaddr}
	mapParam := map[string]interface{}{"minimumSumAmount": 1200 /*, "minimumAmount": 0.05*/}

	reqListunspent := listunspent(100, 99999999, sliAddr, true, mapParam)
	returnJson := sendRpcRequest(rpcClient, reqListunspent)
	log.Println("returnJson_list=", returnJson)

	rpcresponse := decodeJson(returnJson)
	inputs, slicListunspent, sumAcount := parseListunspent(&rpcresponse)
	fmt.Println("inputs=", inputs)

	//命令:createrawtransaction
	// distAddr := "SWX9b3z4K47fLq8vK6x7FUM7XTj6r9Gr1V"
	distAddr := distaddr
	var distAcount float64 = 0
	var fee float64 = 0
	// changeAddr := "3PX9raYTM5MZRQahhxikPCsGcfvBvbFcWg"
	changeAddr := fundaddr
	// distAcount = 500
	distAcount = sendacount
	fee = 0.002
	tempChangeAcount := sumAcount - distAcount - fee
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
	// privKey1 := []string{"KyVzoDL4UbM3PmMaGgSZRWuADq3Au5R9wC8oBYajxtsY8eEtaMcv"}
	privKey1 := []string{signkey[0]}
	reqSignrawtransactionwithkey1 := signrawtransactionwithkey(txHex, privKey1, slicListunspent)
	returnJson = sendRpcRequest(rpcClient, reqSignrawtransactionwithkey1)
	fmt.Println("returnJson_sign1=", returnJson)
	rpcresponse = decodeJson(returnJson)
	signHex1 := parseSignrawtransactionwithkey(&rpcresponse)
	//签名2：
	// privKey2 := []string{"Ky4CzdZ6VJsHFfH8DwFvk488XgzDb6g4ijqZmYxhL3M4iFyuZdAx"}
	privKey2 := []string{signkey[1]}
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

func MultisigTx(w http.ResponseWriter, r *http.Request) {
	var txinfo txInfo
	defer func() {
		if x := recover(); x != nil {
			fmt.Errorf("unexpected error:%v", x)
		}
	}()
	defer r.Body.Close()

	ReadInfo, err1 := ioutil.ReadAll(r.Body)
	if err1 != nil {
		log.Fatal(err1)
	}

	err2 := json.Unmarshal(ReadInfo, &txinfo)
	if err2 != nil {
		log.Fatal("err2: ", err2)
	}

	distaddr := txinfo.DistAddr
	fundaddr := txinfo.FundAddr
	sendacount := txinfo.SendAcount
	signkey := txinfo.SignKey

	if sendacount > 1000 {
		slinum := int(sendacount / 1000)
		fmt.Println("slinum:", slinum)
		for i := 0; i < slinum; i++ {
			makeMultisigTx(distaddr, fundaddr, 1000, signkey)
		}
		remainAcount := sendacount - float64(1000*slinum)
		if remainAcount > 0 {
			makeMultisigTx(distaddr, fundaddr, remainAcount, signkey)
		}
	} else {
		makeMultisigTx(distaddr, fundaddr, sendacount, signkey)
	}

	w.Write([]byte("successfully complete transaction"))
}

func main() {
	http.HandleFunc("/", MultisigTx)
	err := http.ListenAndServe("192.168.0.181:9090", nil)
	if err != nil {
		log.Fatal("ListenAndServer: ", err)
	}
}
