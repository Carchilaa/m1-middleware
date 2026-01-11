package main

import(
	"log"
	"alerter/services"
	"github.com/nats-io/nats.go"
)

func main(){
	log.Println("Demarrage du service Alerter...")

	nc, err := nats.Connect("nats://localhost:4222")
	if err != nil{
		log.Fatalf("Impossible de se connecter a NATS: %v", err)
	}

	defer nc.Close()



	services.StartAlerterConsumer(nc)
}