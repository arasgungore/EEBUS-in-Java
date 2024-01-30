package ship

import (
	"crypto/sha1"
	"crypto/tls"
	"encoding/hex"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/LMF-DHBW/go_eebus/resources"

	"github.com/grandcat/zeroconf"
	"github.com/phayes/freeport"
	"golang.org/x/net/websocket"
)

type ConnectionManager func(string, *websocket.Conn)
type ConnectionManagerSpine func(*SMEInstance, string)
type CloseHandler func(*SMEInstance)
type handler func([]byte)
type dataHandler func(resources.DatagramType)

type ShipNode struct {
	serverPort            int
	hostname              string
	IsGateway             bool
	SME                   []*SMEInstance
	Requests              []*Request
	SpineConnectionNotify ConnectionManagerSpine
	SpineCloseHandler     CloseHandler
	CertName              string
	devId                 string
	brand                 string
	devType               string
}

type Request struct {
	Path string
	Id   string
	Ski  string
}

func NewShipNode(hostname string, IsGateway bool, certName string, devId string, brand string, devType string) *ShipNode {
	// Empty Ship node has empty list of clients and no server
	return &ShipNode{0, hostname, IsGateway, make([]*SMEInstance, 0), make([]*Request, 0), nil, nil, certName, devId, brand, devType}
}

func (shipNode *ShipNode) Start() {
	// ShipNode start -> assign port, create server
	port, err := freeport.GetFreePort()
	resources.CheckError(err)
	shipNode.serverPort = port
	// Start server, Register Dns and search for other DNS entries
	if !shipNode.IsGateway {
		go shipNode.StartServer()
		go shipNode.RegisterDns()

	}
	go shipNode.BrowseDns()
}

func (shipNode *ShipNode) handleFoundService(entry *zeroconf.ServiceEntry) {
	// If found service is not on same device
	if entry.Port != shipNode.serverPort {
		log.Println("Found new service", entry.HostName, entry.Port)

		skis, _ := ReadSkis()
		if resources.StringInSlice(strings.Split(entry.Text[3], "=")[1], skis) {
			// Device is trusted
			go shipNode.Connect(strings.Split(entry.Text[2], "=")[1], strings.Split(entry.Text[3], "=")[1])
		} else {
			if shipNode.IsGateway {
				requestAlreadyMade := false
				for _, req := range shipNode.Requests {
					if req.Ski == strings.Split(entry.Text[3], "=")[1] {
						requestAlreadyMade = true
					}
				}
				if !requestAlreadyMade {
					shipNode.Requests = append(shipNode.Requests, &Request{
						Path: strings.Split(entry.Text[2], "=")[1],
						Id:   strings.Split(entry.Text[6], "=")[1] + " " + strings.Split(entry.Text[5], "=")[1] + " " + strings.Split(entry.Text[1], "=")[1],
						Ski:  strings.Split(entry.Text[3], "=")[1],
					})
				}
			} else {
				ticker := time.NewTicker(3 * time.Second)
				go func() {
					for {
						select {
						case <-ticker.C:
							skis, _ := ReadSkis()
							if resources.StringInSlice(strings.Split(entry.Text[3], "=")[1], skis) {
								go shipNode.Connect(strings.Split(entry.Text[2], "=")[1], strings.Split(entry.Text[3], "=")[1])
							}
						case <-time.After(120 * time.Second):
							// Stop trying after 2 minutes
							ticker.Stop()
							return
						}
					}
				}()
			}
		}
	}
}

/* Procedure for new conncetions
1. Create SME instance and append to list from SHIP node
2. Start CME handshake
3. Start data exchange -> notify spine
(Skip Hello handshake, protocol handshake and pin exchange)
*/
func (shipNode *ShipNode) newConnection(role string, conn *websocket.Conn, ski string) {
	skiIsNew := ""
	skis, _ := ReadSkis()
	if !resources.StringInSlice(ski, skis) && shipNode.IsGateway {
		skiIsNew = ski
	}

	newSME := &SMEInstance{role, "INIT", conn, shipNode.SpineCloseHandler, ski}

	for _, e := range shipNode.SME {
		if e.Ski == ski {
			conn.Close()
			return
		}
	}

	shipNode.SME = append(shipNode.SME, newSME)

	newSME.StartCMI()
	shipNode.SpineConnectionNotify(newSME, skiIsNew)
}

func (shipNode *ShipNode) Connect(service string, ski string) {
	conf, err := websocket.NewConfig(service, "http://"+shipNode.hostname)
	resources.CheckError(err)

	var cert tls.Certificate
	cert, err = tls.LoadX509KeyPair(shipNode.CertName+".crt", shipNode.CertName+".key")
	resources.CheckError(err)

	conf.TlsConfig = &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	conn, err := websocket.DialConfig(conf)
	resources.CheckError(err)

	shipNode.newConnection("client", conn, ski)
}

func (shipNode *ShipNode) StartServer() {

	server := &http.Server{
		Addr: ":" + strconv.Itoa(shipNode.serverPort),
		TLSConfig: &tls.Config{
			ClientAuth: tls.RequestClientCert,
		},
		Handler: websocket.Handler(func(ws *websocket.Conn) {
			publickey := ws.Request().TLS.PeerCertificates[0].RawSubjectPublicKeyInfo

			hasher := sha1.New()
			hasher.Write(publickey)
			shipNode.newConnection("server", ws, hex.EncodeToString(hasher.Sum(nil)))
		}),
	}
	err := server.ListenAndServeTLS(shipNode.CertName+".crt", shipNode.CertName+".key")
	resources.CheckError(err)
}
