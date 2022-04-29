package main

import (
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	p "testYoutube/api/proto"
	db "testYoutube/internal/database"
)

func main() {
	//путь до бд где хранятся уже запрашиваемые превью
	dbFile := "internal/database/previewDb"
	db, err := db.NewDB(dbFile)
	if err != nil {
		fmt.Println("ERROR: create -", err)
		return
	}
	defer db.Close()

	s := grpc.NewServer()
	srv := &GRPCServer{*db, p.UnimplementedGetterServer{}}
	p.RegisterGetterServer(s, srv)
	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	if err = s.Serve(l); err != nil {
		log.Fatal(err)
	}
	log.Println("server started")
}
