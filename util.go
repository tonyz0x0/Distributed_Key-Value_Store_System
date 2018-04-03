package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// decode Json file into pair slice
func decodeJSONFromReq(req *http.Request) []KvPair {
	kvPairs := make([]KvPair, 0)
	json.NewDecoder(req.Body).Decode(&kvPairs)
	return kvPairs
}

func decodeJSONFromResp(resp *http.Response) []KvPair {
	kvPairs := make([]KvPair, 0)
	json.NewDecoder(resp.Body).Decode(&kvPairs)
	return kvPairs
}

func decodeClientJSONFromResp(resp *http.Response) JsonFile {
	var jsonFile JsonFile
	jsonFile.ExistedKeys = make([]KvPair, 0)
	jsonFile.NotExistedKeys = make([]KvPair, 0)
	json.NewDecoder(resp.Body).Decode(&jsonFile)
	return jsonFile
}

// encode pairs into Json file
func encodeJSON(kvPairs []KvPair) []byte {
	var middleData []byte
	middleData, err := json.Marshal(kvPairs)
	if err != nil {
		log.Fatalln(err)
	}
	return middleData
}

func encodeClientJSON(jsonFile JsonFile) []byte {
	var middleData []byte
	middleData, err := json.Marshal(jsonFile)
	if err != nil {
		log.Fatalln(err)
	}
	return middleData
}
