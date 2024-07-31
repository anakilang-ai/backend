package main

import (
    "bufio"
    "fmt"package main

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
			pesan, _ := pembaca.ReadString('\n')
			fmt.Fprintf(koneksi, pesan+"\n")
			if strings.TrimSpace(pesan) == "KELUAR" {
				fmt.Println("Menutup koneksi")
				return
			}
		}
	}
	
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
            continue
        }
        go tanganiKoneksi(koneksi)
    }
}
