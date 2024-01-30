package ship

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/LMF-DHBW/go_eebus/resources"

	"github.com/grandcat/zeroconf"
)

func (shipNode *ShipNode) BrowseDns() {
	log.Println("Browsing for entries")
	// Discover all ship services on the network
	resolver, err := zeroconf.NewResolver(nil)
	resources.CheckError(err)

	entries := make(chan *zeroconf.ServiceEntry)
	go func(results <-chan *zeroconf.ServiceEntry) {
		for entry := range results {
			go shipNode.handleFoundService(entry)
		}
	}(entries)

	ctx, _ := context.WithCancel(context.Background())

	err = resolver.Browse(ctx, "_ship._tcp", "local.", entries)
	resources.CheckError(err)

	<-ctx.Done()
}

func (shipNode *ShipNode) RegisterDns() {
	// Define values for DNS entry
	port := strconv.Itoa(shipNode.serverPort)
	id := "DEVICE-EEB01-" + shipNode.devId + ".local."

	txtRecord := []string{"txtvers=1", "id=" + id, "path=wss://" + shipNode.hostname + ":" + port, "SKI=" + shipNode.getSki(), "register=true", "brand=" + shipNode.brand, "type=" + shipNode.devType}
	log.Println("Registering: ", txtRecord)
	server, err := zeroconf.Register("Device "+port, "_ship._tcp", "local.", shipNode.serverPort, txtRecord, nil)
	resources.CheckError(err)

	defer server.Shutdown()
	defer log.Println("Registering stopped")
	// Shutdown server after 2 minutes
	<-time.After(time.Second * 120)
}
