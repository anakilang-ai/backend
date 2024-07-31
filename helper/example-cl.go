package main

import (
    "bufio"
    "fmt"
    "net"
    "os"
    "strings"
)

func main() {
    koneksi, err := net.Dial("tcp", "localhost:8081")
    if err != nil {
        fmt.Println("Kesalahan saat menghubungkan ke server:", err)
        os.Exit(1)
    }
    defer koneksi.Close()

    fmt.Println("Terhubung ke server")
    go func() {
        for {
            pesan, _ := bufio.NewReader(koneksi).ReadString('\n')
            fmt.Print("Pesan dari server:", string(pesan))
        }
    }()

    for {
        pembaca := bufio.NewReader(os.Stdin)
        fmt.Print("Masukkan pesan: ")
            fmt.Println("Menutup koneksi")
            return
        }
    }
}
