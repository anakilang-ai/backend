package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	// Mencoba untuk terhubung ke server TCP di localhost pada port 8081
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
			// Membaca pesan dari server
			pesan, err := bufio.NewReader(koneksi).ReadString('\n')
			if err != nil {
				fmt.Println("Kesalahan saat membaca pesan dari server:", err)
				return
			}
			// Menampilkan pesan yang diterima dari server
			fmt.Print("Pesan dari server: ", pesan)
		}
	}()

	// Loop untuk mengirim pesan ke server
	for {
		pembaca := bufio.NewReader(os.Stdin)
		fmt.Print("Masukkan pesan: ")
		// Membaca input dari pengguna
		pesan, err := pembaca.ReadString('\n')
		if err != nil {
			fmt.Println("Kesalahan saat membaca input:", err)
			continue
		}
		// Mengirim pesan ke server
		fmt.Fprintf(koneksi, pesan+"\n")
		// Jika pesan adalah "KELUAR", menutup koneksi dan keluar
		if strings.TrimSpace(pesan) == "KELUAR" {
			fmt.Println("Menutup koneksi")
			return
		}
	}
}
