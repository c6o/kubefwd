package proxyer
import (
    "fmt"
    "io"
    "net"
)

func Proxyer(localIp string, localPort int, remoteIp string, remotePort int) {
    listenAddress := fmt.Sprintf("%s:%d", localIp, localPort)
    remoteAddress := fmt.Sprintf("%s:%d", remoteIp, remotePort)
    ln, err := net.Listen("tcp", listenAddress)
    if err != nil {
        panic(err)
    }

    for {
        conn, err := ln.Accept()
        if err != nil {
            panic(err)
        }

        go handleRequest(conn, remoteAddress)
    }
}

func handleRequest(conn net.Conn, remoteAddress string) {
    fmt.Println("proxy to %s", remoteAddress)

    proxy, err := net.Dial("tcp", remoteAddress)
    if err != nil {
        panic(err)
    }

    go copyIO(conn, proxy)
    go copyIO(proxy, conn)
}

func copyIO(src, dest net.Conn) {
    defer dest.Close()
    io.Copy(src, dest)
}