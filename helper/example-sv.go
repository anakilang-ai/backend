package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

// Fungsi untuk menangani koneksi dari klien
func tanganiKoneksi(koneksi net.Conn) {
	fmt.Println("Klien terhubung:", koneksi.RemoteAddr().String())
	defer koneksi.Close()

	for {
		pesan, err := bufio.NewReader(koneksi).ReadString('\n')
		if err != nil {
			fmt.Println("Kesalahan saat membaca pesan dari klien:", err)
			return
		}

		// Memeriksa apakah pesan adalah "KELUAR" untuk mengakhiri koneksi
		if strings.TrimSpace(pesan) == "KELUAR" {
			fmt.Println("Klien terputus:", koneksi.RemoteAddr().String())
			return
		}

		fmt.Print("Pesan diterima:", pesan)
		koneksi.Write([]byte(pesan))
	}
}

func main() {
	fmt.Println("Memulai server...")

	// Membuat pendengar TCP di localhost pada port 8081
	pendengar, err := net.Listen("tcp", "localhost:8081")
	if err != nil {
		fmt.Println("Kesalahan saat membuat server:", err)
		os.Exit(1)
	}
	defer pendengar.Close()

	fmt.Println("Server berjalan di localhost:8081")

	for {
		// Menerima koneksi dari klien
		koneksi, err := pendengar.Accept()
		if err != nil {
			fmt.Println("Kesalahan saat menerima koneksi:", err)
			continue
		}

		// Menangani koneksi dalam goroutine terpisah
		go tanganiKoneksi(koneksi)
	}
}
