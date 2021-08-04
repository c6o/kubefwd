package proxyer

import (
	"fmt"
	"io"
	"log"
	"net"
)

func Proxyer(localIp string, localPort int32, remoteIp string, remotePort int32) {
	listenAddress := fmt.Sprintf("%s:%d", localIp, localPort)
	remoteAddress := fmt.Sprintf("%s:%d", remoteIp, remotePort)
	listener, err := net.Listen("tcp", listenAddress)
	if err != nil {
		panic(err)
	}

	for {
		local, err := listener.Accept()
		if err != nil {
			log.Println("error accepting connection", err)
			continue
		}

		go func() {
			remote, err := net.Dial("tcp", remoteAddress)
			if err != nil {
				log.Println("error dialing remote addr", err)
				return
			}
			go io.Copy(remote, local)
			io.Copy(local, remote)
			remote.Close()
			local.Close()
		}()
	}
}
