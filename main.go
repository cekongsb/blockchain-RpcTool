package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

const (
	SERVER_HOST       = "192.168.0.143"
	SERVER_PORT       = 6062
	USER              = "chain"
	PASSWD            = "999000"
	USESSL            = false
	WALLET_PASSPHRASE = "WalletPassphrase"
)

func makeReqJson(arg ...interface{}) (string, error) {
	var stuReqJson rpcRequest

	stuReqJson.JsonRpc = "1.0"
	stuReqJson.Id = arg[0].(string)
	stuReqJson.Method = arg[1].(string)
	stuReqJson.Params = arg[2:]

	reqJson, err := json.Marshal(stuReqJson)
	if err != nil {
		fmt.Println("json.Marshal(stuReqJson) err:", err)
	}
	fmt.Println(stuReqJson.Method, "reqJson=", string(reqJson))
	return string(reqJson), err
}

func decodeJson(resultInfo string) (rpcResponse, error) {
	var rpcresponse rpcResponse
	err := json.Unmarshal([]byte(resultInfo), &rpcresponse)
	if err != nil {
		fmt.Println("json.Unmarshal([]byte(resultInfo), &rpcresponse) err:", err)
	}
	return rpcresponse, err
}

func parseListunspent(rpcresponse *rpcResponse) (interface{}, interface{}, float64) {
	uint8Listunspent := rpcresponse.Result

	var listunspentInfo interface{}
	err := json.Unmarshal(uint8Listunspent, &listunspentInfo)
	if err != nil {
		fmt.Println("json.Unmarshal(uint8Listunspent, &listunspentInfo):", err)
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

	reqJson, _ := makeReqJson("listunspent", "listunspent", minconf, maxconf, address, isSafe, queryOption)
	return reqJson
}

func createrawtransaction(arg ...interface{}) string {
	inputs := arg[0]
	outputs := arg[1]
	reqJson, _ := makeReqJson("createrawtransaction", "createrawtransaction", inputs, outputs)
	return reqJson
}

func parseCreaterawtransaction(rpcresponse *rpcResponse) (interface{}, error) {
	uint8Createrawtransaction := rpcresponse.Result
	var CreaterawtransactionInfo interface{}
	err := json.Unmarshal(uint8Createrawtransaction, &CreaterawtransactionInfo)
	if err != nil {
		fmt.Println("json.Unmarshal(uint8Createrawtransaction, &CreaterawtransactionInfo):", err)
	}
	return CreaterawtransactionInfo, err
}

func signrawtransactionwithkey(arg ...interface{}) string {
	txHex := arg[0]
	privKey := arg[1]
	listunspentInfo := arg[2]
	reqJson, _ := makeReqJson("signrawtransactionwithkey", "signrawtransactionwithkey", txHex, privKey, listunspentInfo)
	return reqJson
}

func parseSignrawtransactionwithkey(rpcresponse *rpcResponse) (interface{}, error) {
	uint8SignTxKey := rpcresponse.Result
	var signhex signHex
	err := json.Unmarshal(uint8SignTxKey, &signhex)
	if err != nil {
		fmt.Println("json.Unmarshal(uint8SignTxKey, &signhex) err:", err)
	}
	return signhex.Hex, err
}

func sendrawtransaction(arg ...interface{}) string {
	hex := arg[0]
	reqJson, _ := makeReqJson("sendrawtransaction", "sendrawtransaction", hex)
	return reqJson
}

func sendRpcRequest(client *rpcClient, reqJson string) (string, error) {
	returnJson, err := client.send(reqJson)
	if err != nil {
		fmt.Println("client.send(reqJson):", err)
	}
	return returnJson, err
}

func makeMultisigTx(distaddr string, fundaddr string, sendacount float64, signkey []string) error {
	rpcClient, err := newClient(SERVER_HOST, SERVER_PORT, USER, PASSWD, USESSL)
	if err != nil {
		fmt.Println("newClient:", err)
	}
	var respErr error
	//命令:listunspent
	// sliAddr := []string{"3PX9raYTM5MZRQahhxikPCsGcfvBvbFcWg"}
	sliAddr := []string{fundaddr}
	mapParam := map[string]interface{}{"minimumSumAmount": 600 /*, "minimumAmount": 0.05*/}

	reqListunspent := listunspent(100, 99999999, sliAddr, true, mapParam)
	returnJson, rpcErr := sendRpcRequest(rpcClient, reqListunspent)
	if rpcErr != nil {
		respErr = errors.New("rpc_command:listunspent send fail..")
	}
	log.Println("returnJson_list=", returnJson)

	rpcresponse, _ := decodeJson(returnJson)
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
	returnJson, rpcErr = sendRpcRequest(rpcClient, reqCreaterawtransaction)
	if rpcErr != nil {
		respErr = errors.New("rpc_command:createrawtransaction send fail..")
	}
	fmt.Println("returnJson_cretx=", returnJson)
	rpcresponse, _ = decodeJson(returnJson)
	txHex, _ := parseCreaterawtransaction(&rpcresponse)
	fmt.Println("txHex=", txHex)

	//命令:signrawtransactionwithkey
	// 签名1：
	// privKey1 := []string{"KyVzoDL4UbM3PmMaGgSZRWuADq3Au5R9wC8oBYajxtsY8eEtaMcv"}
	privKey1 := []string{signkey[0]}
	reqSignrawtransactionwithkey1 := signrawtransactionwithkey(txHex, privKey1, slicListunspent)
	returnJson, rpcErr = sendRpcRequest(rpcClient, reqSignrawtransactionwithkey1)
	if rpcErr != nil {
		respErr = errors.New("rpc_command:signrawtransactionwithkey send fail..")
	}
	fmt.Println("returnJson_sign1=", returnJson)
	rpcresponse, _ = decodeJson(returnJson)
	signHex1, _ := parseSignrawtransactionwithkey(&rpcresponse)
	//签名2：
	// privKey2 := []string{"Ky4CzdZ6VJsHFfH8DwFvk488XgzDb6g4ijqZmYxhL3M4iFyuZdAx"}
	privKey2 := []string{signkey[1]}
	reqSignrawtransactionwithkey2 := signrawtransactionwithkey(signHex1, privKey2, slicListunspent)
	returnJson, rpcErr = sendRpcRequest(rpcClient, reqSignrawtransactionwithkey2)
	if rpcErr != nil {
		respErr = errors.New("rpc_command:signrawtransactionwithkey send fail..")
	}
	fmt.Println("returnJson_sign2=", returnJson)
	rpcresponse, _ = decodeJson(returnJson)
	signHex2, _ := parseSignrawtransactionwithkey(&rpcresponse)
	fmt.Println("signHex2=", signHex2)

	//命令:sendrawtransaction
	reqSendrawtransaction := sendrawtransaction(signHex2)
	returnJson, rpcErr = sendRpcRequest(rpcClient, reqSendrawtransaction)
	if rpcErr != nil {
		respErr = errors.New("rpc_command:sendrawtransaction send fail..")
	}
	fmt.Println("returnJson_brocast=", returnJson)
	return respErr
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
		fmt.Println("ioutil.ReadAll(r.Body):", err1)
	}

	err2 := json.Unmarshal(ReadInfo, &txinfo)
	if err2 != nil {
		fmt.Println("json.Unmarshal err: ", err2)
	}

	distaddr := txinfo.DistAddr
	fundaddr := txinfo.FundAddr
	sendacount := txinfo.SendAcount
	signkey := txinfo.SignKey

	var msErr error
	if sendacount > 500 {
		slinum := int(sendacount / 500)
		fmt.Println("slinum:", slinum)
		for i := 0; i < slinum; i++ {
			msErr = makeMultisigTx(distaddr, fundaddr, 500, signkey)
			if msErr != nil {
				w.Write([]byte("makeMultisigTx fail"))
				return
			}
		}
		remainAcount := sendacount - float64(500*slinum)
		if remainAcount > 0 {
			msErr = makeMultisigTx(distaddr, fundaddr, remainAcount, signkey)
			if msErr != nil {
				w.Write([]byte("makeMultisigTx fail"))
				return
			}
		}
	} else {
		msErr = makeMultisigTx(distaddr, fundaddr, sendacount, signkey)
		if msErr != nil {
			w.Write([]byte("makeMultisigTx fail"))
			return
		}
	}
	w.Write([]byte("successfully complete transaction"))
}

func main() {
	http.HandleFunc("/", MultisigTx)
	err := http.ListenAndServe("192.168.0.188:9090", nil)
	if err != nil {
		log.Fatal("ListenAndServer: ", err)
	}
}
