package main

import (
	"DDSS/dataServer/internal/blobsfile"
	"DDSS/dataServer/internal/tcp"
)


func main(){

	store := blobsfile.NewStoreWorker()
	tcp.New(store).Listen()

}

