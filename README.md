# Distributed-Key-Value-Store
This is a distributed key-value store implementing basic `get()` and `put()` functions. The whole system considers the hash consistency as load balancing. Client/server model is based on HTTP and message is packaged in JSON file. The whole system is written in Golang.

## Proposal
Golang is a relatively simple and light language. More than that, it is also a powerful language in concurrency scenarios. A distributed key-value system should handle huge amounts of requests and responses, the efficiency is totally important and fatal. Golang can guarantee that within the help of its multi threading: `'Gorutine'`. So using Golang to do this is a great choice.

## Design Strategy
The Key-value store should be able to handle data larger than any one node's memory capacity. That is, at any given time, a single node might not have all the data. The deliverables will include:
  - A server program that accepts HTTP get/post/put requests from the clients and returns a valid response. The server will communicate with it's peer processes (spread across the network) to maintain a consistent view of the key-value database. All communication between the HTTP client and this server should be in JSON format.

  - A proxy/coordinator process keeps track of available servers and data stored in those servers. A client connects to the proxy/coordinator process to learn the address of a server that it should connect for performing set/get operations. The proxy server also acts as a load-balancer and ensures a uniform workload distribution among various servers.

## Implementation
> First, the proxy/coordinate process will handle requests, contain many `get` or `set` at the same time, from client in JSON format message.

> Second, the proxy will classify all the keys and put them into different new JSON files refer to different inner servers by calculate hash.

>Third, the proxy will send these JSON files to those servers in HTTP and each server will handle different requests from proxy, whether `set` key or `get` key. When the data is acquired, each server will generate a new JSON file with the data and send back to proxy.

>Forth, proxy will cluster these different JSON files from different servers into different new JSON file by different clients and at last send one file to its belonging client.

## Example
Put `Server` folder in different directory than `proxy.go` file. 
> Starting the proxy:
```sh
$ go run proxy.go util.go
```
> Terminal will show the proxy is listening:
```sh
*** Welcome ***
Listening on port 8080...
```
> Then starting two servers with parameters `9000` and `9001` in `Server` folder:
```sh
go run server.go util.go 9000
```
```sh
go run server.go util.go 9001
```
> Two terminals will show messages meaning the servers are listening on each port now.
```sh
Listening on port ...9000
```
```sh
Listening on port ...9001
```
---
---
>Sending a `POST` request to `localhost:8080\set` to set key-values in JSON file:
```JSON
[
	{
		"key": "hi",
		"value": "hello"
	},
	{
		"key": "hi1",
		"value": "hello1"
	},
	{
		"key": "hi2",
		"value": "hello2"
	},
	{
		"key": "hi3",
		"value": "hello3"
	},
	{
		"key": "hi4",
		"value": "hello4"
	}
]
```
It will return a message:
```sh
Success to set keySuccess to set key
```
---
---
> Sending a `POST` request to `localhost:8080\get` to get key-values in JSON file:
```JSON
[
	{
		"key": "hi"
	},
	{
		"key": "hi1"
	},
	{
		"key": "hi2"
	},
	{
		"key": "hi3"
	},
	{
		"key": "hi4"
	}
]
```
>It will return a new JSON file:
```JSON
{
    "existedKey": [
        {
            "key": "hi",
            "value": "hello"
        },
        {
            "key": "hi1",
            "value": "hello1"
        },
        {
            "key": "hi2",
            "value": "hello2"
        },
        {
            "key": "hi3",
            "value": "hello3"
        },
        {
            "key": "hi4",
            "value": "hello4"
        }
    ]
}
```
---
---
> Sending a `POST` request containing those keys `hi5` and `hi6` that do not exist in key-value store to `localhost:8080\get` :
```JSON
[
	{
		"key": "hi"
	},
	{
		"key": "hi1"
	},
	{
		"key": "hi5"
	},
	{
		"key": "hi6"
	}
]
```
>It will generate a JSON file response like this:
```JSON
{
    "existedKey": [
        {
            "key": "hi",
            "value": "hello"
        },
        {
            "key": "hi1",
            "value": "hello1"
        }
    ],
    "notExistedKey": [
        {
            "key": "hi6",
            "value": ""
        },
        {
            "key": "hi5",
            "value": ""
        }
    ]
}
```
