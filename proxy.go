package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"

	"github.com/serialx/hashring"
)

type KvPair struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type JsonFile struct {
	ExistedKeys    []KvPair `json:"existedKey"`
	NotExistedKeys []KvPair `json:"notExistedKey,omitempty"`
}

var ring *hashring.HashRing

/*
 * This handler handles the set request from client
 */
func setHandler(w http.ResponseWriter, req *http.Request) {
	//decode json file into slice
	kvPairs := decodeJSONFromReq(req)
	multiJsons := make(map[string][]KvPair)

	for _, pair := range kvPairs {

		//use key to do the hash and determine which server it should go
		server, _ := ring.GetNode(pair.Key)

		//put all the pairs which belong to the same server into same
		//pairs slice
		_, ok := multiJsons[server]
		if ok {
			multiJsons[server] = append(multiJsons[server], KvPair{
				pair.Key,
				pair.Value,
			})
		} else {
			multiJsons[server] = make([]KvPair, 0)
			multiJsons[server] = append(multiJsons[server], KvPair{
				pair.Key,
				pair.Value,
			})
		}
	}

	//Encode every key-value data into new Json file and send
	//to each server
	for server, data := range multiJsons {

		//make a new request from proxy to server
		var middleData []byte
		middleData = encodeJSON(data)
		resp := makeSetRequest(server, middleData)

		//get the result from the server
		if resp.StatusCode == 200 {
			w.Write([]byte("Success to set key"))
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}

	//send the result to client

}

// make a set request to specific server and send it a Json file
func makeSetRequest(server string, middleData []byte) *http.Response {
	req, err := http.NewRequest("PUT", server, bytes.NewBuffer(middleData))
	if err != nil {
		log.Fatalln(err)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	return resp
}

func getHandler(w http.ResponseWriter, req *http.Request) {
	//decode json file into slice
	var returnKvPairs []KvPair
	var returnNotExistedKvPairs []KvPair
	var returnJsonData JsonFile

	kvPairs := decodeJSONFromReq(req)
	multiJsons := make(map[string][]KvPair)

	for _, pair := range kvPairs {

		//use key to do the hash and determine which server it should go
		server, _ := ring.GetNode(pair.Key)

		//put all the pairs which belong to the same server into same
		//pairs slice
		_, ok := multiJsons[server]
		if ok {
			multiJsons[server] = append(multiJsons[server], KvPair{
				pair.Key,
				"",
			})
		} else {
			multiJsons[server] = make([]KvPair, 0)
			multiJsons[server] = append(multiJsons[server], KvPair{
				pair.Key,
				"",
			})
		}
	}

	//Encode every key-value data into new Json file and send
	//to each server
	for server, data := range multiJsons {

		//make a new request from proxy to server
		var middleData []byte
		middleData = encodeJSON(data)
		resp := makeGetRequest(server, middleData)
		//get the result from the server
		if resp.StatusCode == 200 {
			jsonData := decodeClientJSONFromResp(resp)
			for _, pair := range jsonData.ExistedKeys {
				returnKvPairs = append(returnKvPairs, pair)
			}
			if len(jsonData.NotExistedKeys) != 0 {
				for _, pair := range jsonData.NotExistedKeys {
					returnNotExistedKvPairs = append(returnNotExistedKvPairs, pair)
				}
			}
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}

	//send the result to client
	returnJsonData.ExistedKeys = returnKvPairs
	returnJsonData.NotExistedKeys = returnNotExistedKvPairs
	w.Write(encodeClientJSON(returnJsonData))
}

// make a get request to specific server and send it a Json file
func makeGetRequest(server string, middleData []byte) *http.Response {
	req, err := http.NewRequest("POST", server, bytes.NewBuffer(middleData))
	if err != nil {
		log.Fatalln(err)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	return resp
}

func initRing() *hashring.HashRing {
	memcacheServers := []string{
		"http://localhost:9000",
		"http://localhost:9001",
	}
	return hashring.New(memcacheServers)
}

func main() {
	fmt.Println("*** Welcome ***")
	ring = initRing()
	mux := http.NewServeMux()
	mux.HandleFunc("/set", setHandler)
	mux.HandleFunc("/get", getHandler)

	log.Println("Listening on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
