// For more information, read [UDP Protocol Specification](https://www.bittorrent.org/beps/bep_0015.html)

package node

import (
	"crypto/sha1"
	"encoding/binary"
	"encoding/hex"
	"log"
	"net"
	"strconv"

	"github.com/Ehab-24/torrent-udp/utils"
)

type Node struct {
	ID string
	*net.UDPConn
}

// New creates a new node instance and connects it to the host over UDP
func New(host string, port uint16) (*Node, error) {

	// Dial a connection to the given host
	addr := host + ":" + strconv.Itoa(int(port))
	raddr, err := net.ResolveUDPAddr("udp4", addr)
	if err != nil {
		return nil, err
	}
	c, err := net.DialUDP("udp4", nil, raddr)
	if err != nil {
		return nil, err
	}

	// Calculate Node's ID
	buf := [2]byte{}
	binary.BigEndian.PutUint16(buf[:], port)
	sha1Sum := sha1.Sum(append([]byte(host), buf[:]...))
	id := hex.EncodeToString(sha1Sum[:])

	// Create a new Node
	n := &Node{
		ID:      id,
		UDPConn: c,
	}

	// Listen for incoming messages
	readCh := make(chan []byte)
	go n.listen(readCh)

	// process incoming messages
	go process(readCh)

	return n, nil
}

// Sends a `connect` message to the tracker
func (n *Node) SendConnect() error {
	var protocolID int64 = 0x41727101980
	var action int32 = 0
	var transactionID int32 = utils.NewTransactionID() // auto increment integer

	// Create the message as a buffer
	buf := make([]byte, 16) // `Connect` messages are of fixed length
	binary.BigEndian.PutUint64(buf[:8], uint64(protocolID))
	binary.BigEndian.PutUint32(buf[8:12], uint32(action))
	binary.BigEndian.PutUint32(buf[12:16], uint32(transactionID))

	// Send message
	_, err := n.Write(buf)

	log.Println("Sent buffer:", buf)

	return err
}

// listen listens to incoming messages and writes the buffer into `ch`
func (n *Node) listen(ch chan<- []byte) {
	buf := make([]byte, 1024)
	for {
		_, err := n.Read(buf)
		if err != nil {
			log.Printf("error while reading from %v ~ %v\n", n.RemoteAddr(), err)
		} else {
			ch <- buf
		}
	}
}

// Process the incoming messages
func process(ch <-chan []byte) {
	for {
		msg := <-ch
		log.Println(string(msg))
	}
}
