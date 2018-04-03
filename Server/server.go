package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

type ServerIdentifier struct {
	IP   string // The IP address of the server.
	Port string // The port on which the server is listening.
}

type Server struct {
	Identifier ServerIdentifier
	Data       map[string]string
}

var server Server

func main() {
	port := os.Args[1]
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)
	fmt.Println("Listening on port ..." + port)

	server.Data = make(map[string]string)
	server.Identifier.IP = "localhost"
	server.Identifier.Port = port

	log.Fatal(http.ListenAndServe(strings.Join([]string{":", port}, ""), mux))
}

func handler(w http.ResponseWriter, r *http.Request) {
	//handle set
	if r.Method == "PUT" {
		kvPairs := decodeJSONFromReq(r)
		for _, pair := range kvPairs {
			server.Data[pair.Key] = pair.Value
			//check whether the key-value is added into Data
			_, ok := server.Data[pair.Key]
			if ok {
				w.WriteHeader(http.StatusOK)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
		}
	}

	if r.Method == "POST" {
		var jsonFile JsonFile
		kvPairs := decodeJSONFromReq(r)
		newKvPairs := make([]KvPair, 0)
		notExistedKeys := make([]KvPair, 0)

		for _, pair := range kvPairs {

			//check whether the key-value is inside Data
			value, ok := server.Data[pair.Key]
			if ok {

				//put acquired data into newKvPairs
				pair.Value = value
				newKvPairs = append(newKvPairs, pair)

				fmt.Println("Success get: ", pair)
			} else {
				notExistedKeys = append(notExistedKeys, KvPair{
					pair.Key,
					"",
				})
				fmt.Println(pair.Key, "is not existed")
			}
		}

		//encode newKvPairs into Json file and send it back to Proxy
		var middleData []byte
		jsonFile.ExistedKeys = newKvPairs
		jsonFile.NotExistedKeys = notExistedKeys
		middleData = encodeClientJSON(jsonFile)
		w.WriteHeader(http.StatusOK)
		w.Write(middleData)
	}
}
