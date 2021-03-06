package main

import (
	"fmt"
	"github.com/Zamiell/isaac-racing-server/src/log"
	"net"
	"sync"
	"time"
	"unsafe"
)

const (
	maxBufferSize = 1024
	helloSize     = 5
	port          = 9001
	clockInterval = 1
)

var (
	mux  = sync.Mutex{}
	pMap = make(PlayerMap)
)

func verifySender(mh MessageHeader, raddr net.Addr) bool {
	storedConnection := pMap.getConnection(mh)

	if storedConnection == nil {
		return false
	} else if storedConnection.(*net.UDPAddr).String() != raddr.(*net.UDPAddr).String() {
		log.Info(fmt.Sprintf("Player=%v shadow has untracked origin, recorded=%v, received=%v",
			mh.PlayerId, storedConnection, raddr))
		return false
	}
	return true
}

func handleHelloMessage(msg []byte, addr net.Addr) {
	mh := MessageHeader{}
	err := mh.Unmarshall(msg)
	if err == nil {
		pMap.update(mh, &addr)
	}
}

func handleShadowMessage(msg []byte, pc net.PacketConn, raddr net.Addr) {
	mh := MessageHeader{}
	err := mh.Unmarshall(msg)

	if err == nil {
		if !verifySender(mh, raddr) {
			return
		}

		opponent := pMap.getOpponent(mh)
		if opponent != nil {
			_, err := pc.WriteTo(msg, *opponent.ADDR)
			if err != nil {
				log.Error(fmt.Sprintf(
					"Shadow proxy failed, player=%v, msg: %v\ncause: %v", mh.PlayerId, msg, err))
			}
		}
	}
}

func shadowServer() (err error, pc net.PacketConn) {
	pc, err = net.ListenPacket("udp4", fmt.Sprintf(":%d", port))
	log.Info(fmt.Sprintf("Listening UDP connections on port %d", port))
	if err != nil {
		log.Error("Error starting UDP shadow server", err)
		return
	}

	buffer := make([]byte, maxBufferSize)

	go func() {
		for {
			n, raddr, err := pc.ReadFrom(buffer)
			if err != nil {
				log.Error(fmt.Sprintf("Error receiving UDP dgram from %v", raddr), err)
			}
			payloadSize := n - int(unsafe.Sizeof(MessageHeader{}))
			if payloadSize < helloSize {
				log.Debug("Dgram skipped, payload len=", payloadSize)
				// skip
			} else if payloadSize == helloSize {
				handleHelloMessage(buffer, raddr)
			} else {
				handleShadowMessage(buffer, pc, raddr)
			}
		}
	}()

	return
}

func sessionClock() {
	tick := time.Tick(clockInterval * time.Second)
	for {
		select {
		case <-tick:
			pMap.updateSessions()
		default:
			time.Sleep(clockInterval * time.Second)
		}
	}
}

func shadowInit() {
	go sessionClock()
	errStart, pc := shadowServer()
	if errStart != nil {
		log.Error("Exited by: ", errStart)
		if pc != nil {
			errClose := pc.Close()
			if errClose != nil {
				log.Error("Error closing connection", errClose)
			}
		}
	}
}
