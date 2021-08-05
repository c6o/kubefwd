package proxyer

import (
	"fmt"
	"io"
	"net"

	log "github.com/sirupsen/logrus"
)

func Proxyer(localIp string, localPort int32, remoteIp string, remotePort int32) {
	localAddress := fmt.Sprintf("%s:%d", localIp, localPort)
	remoteAddress := fmt.Sprintf("%s:%d", remoteIp, remotePort)
	listener, err := net.Listen("tcp", localAddress)
	log.Infof("Proxyer: Setting up proxy from %s to %s", localAddress, remoteAddress)
	if err != nil {
		log.Errorf("Proxyer: Proxy conflict: ", err)
		return
	}

	for {
		local, err := listener.Accept()
		if err != nil {
			log.Warnf("Proxyer: Error accepting connection: ", err)
			continue
		}

		go func() {
			remote, err := net.Dial("tcp", remoteAddress)
			if err != nil {
				log.Warnf("Proxyer: Error dialling remote address: ", err)
				local.Close()
				return
			}
			go io.Copy(remote, local)
			io.Copy(local, remote)
			remote.Close()
			local.Close()
		}()
	}
}
