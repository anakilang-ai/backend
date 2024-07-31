package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	// Mencoba menghubungkan ke server di localhost pada port 8081
	koneksi, err := net.Dial("tcp", "localhost:8081")
	if err != nil {
		fmt.Println("Kesalahan saat menghubungkan ke server:", err)
		os.Exit(1)
	}
	defer koneksi.Close()

	fmt.Println("Terhubung ke server")

	// Goroutine untuk menerima pesan dari server
	go func() {
		for {
			pesan, err := bufio.NewReader(koneksi).ReadString('\n')
			if err != nil {
				fmt.Println("Kesalahan saat membaca pesan dari server:", err)
				return
			}
			fmt.Print("Pesan dari server:", pesan)
		}
	}()

	// Loop utama untuk mengirim pesan ke server
	for {
		pembaca := bufio.NewReader(os.Stdin)
		fmt.Print("Masukkan pesan: ")
		pesan, err := pembaca.ReadString('\n')
		if err != nil {
			fmt.Println("Kesalahan saat membaca input:", err)
			return
		}
		fmt.Fprintf(koneksi, pesan)

		// Memeriksa apakah pesan adalah "KELUAR" untuk menutup koneksi
		if strings.TrimSpace(pesan) == "KELUAR" {
			fmt.Println("Menutup koneksi")
			return
		}
	}
}
