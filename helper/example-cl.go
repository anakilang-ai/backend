package main

import (
    "bufio"
    "fmt"
    "net"
    "os"
    "strings"
)

func main() {
    conn, err := net.Dial("tcp", "localhost:8081")
    if err != nil {
        fmt.Println("Error connecting to server:", err)
        os.Exit(1)
    }
    defer conn.Close()

    fmt.Println("Connected to server")
    go func() {
        for {
            message, _ := bufio.NewReader(conn).ReadString('\n')
            fmt.Print("Message from server:", string(message))
        }
    }()

    for {
        reader := bufio.NewReader(os.Stdin)
        fmt.Print("Enter message: ")
        message, _ := reader.ReadString('\n')
        fmt.Fprintf(conn, message+"\n")
        if strings.TrimSpace(message) == "EXIT" {
            fmt.Println("Closing connection")
            return
        }
    }
}
