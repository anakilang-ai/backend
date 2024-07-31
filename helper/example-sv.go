package main

import (
    "bufio"
    "fmt"
    "net"
    "os"
    "strings"
)

func handleConnection(conn net.Conn) {
    fmt.Println("Client connected:", conn.RemoteAddr().String())
    defer conn.Close()
    
    for {
        message, _ := bufio.NewReader(conn).ReadString('\n')
        if strings.TrimSpace(message) == "EXIT" {
            fmt.Println("Client disconnected:", conn.RemoteAddr().String())
            return
        }
        fmt.Print("Message received:", string(message))
        conn.Write([]byte(message))
    }
}

func main() {
    fmt.Println("Starting server...")
    listener, err := net.Listen("tcp", "localhost:8081")
    if err != nil {
        fmt.Println("Error creating server:", err)
        os.Exit(1)
    }
    defer listener.Close()
    
    fmt.Println("Server started on localhost:8081")
    for {
        conn, err := listener.Accept()
        if err != nil {
            fmt.Println("Error accepting connection:", err)
            continue
        }
        go handleConnection(conn)
    }
}
