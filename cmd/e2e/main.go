package main

import (
	"context"
	"internal/app"
	"log"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("connection error: %v", err)
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {

		}
	}(conn)

	client := app.NewICacheClient(conn)
	key := "key"

	setReq := &app.SetCommand{Key: key, Value: "value"}
	setResp, err := client.Set(context.Background(), setReq)
	if err != nil {
		log.Fatalf("Error Set: %v", err)
	}
	log.Printf("SET done: %v", setResp.GetDone())

	getReq := &app.GetCommand{Key: key}
	getResp, err := client.Get(context.Background(), getReq)
	if err != nil {
		log.Fatalf("Error Get: %v", err)
	}
	log.Printf("GET %v: %v", key, getResp.GetValue())

	delReq := &app.DelCommand{Key: key}
	delResp, err := client.Del(context.Background(), delReq)
	if err != nil {
		log.Fatalf("Error Del: %v", err)
	}
	log.Printf("Del done: %v", delResp.GetDone())

	getReq = &app.GetCommand{Key: key}
	getResp, err = client.Get(context.Background(), getReq)
	if err != nil {
		log.Fatalf("Error Get: %v", err)
	}
	log.Printf("GET %v: %v", key, getResp.GetValue())
}
