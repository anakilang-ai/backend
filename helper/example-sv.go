package main

import (
    "bufio"
    "fmt"
    "net"
    "os"
    "strings"
)

func tanganiKoneksi(koneksi net.Conn) {
    fmt.Println("Klien terhubung:", koneksi.RemoteAddr().String())
    defer koneksi.Close()
    
    for {
        pesan, _ := bufio.NewReader(koneksi).ReadString('\n')
        if strings.TrimSpace(pesan) == "KELUAR" {
            fmt.Println("Klien terputus:", koneksi.RemoteAddr().String())
            return
        }
        fmt.Print("Pesan diterima:", string(pesan))
        koneksi.Write([]byte(pesan))
    }
}

func main() {
    fmt.Println("Memulai server...")
    pendengar, err := net.Listen("tcp", "localhost:8081")
    if err != nil {
        fmt.Println("Kesalahan saat membuat server:", err)
        os.Exit(1)
    }
    defer pendengar.Close()
    
    fmt.Println("Server berjalan di localhost:8081")
    for {
        koneksi, err := pendengar.Accept()
        if err != nil {
            fmt.Println("Kesalahan saat menerima koneksi:", err)
            continue
        }
        go tanganiKoneksi(koneksi)
    }
}
