package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

// tanganiKoneksi menangani komunikasi dengan satu klien
func tanganiKoneksi(koneksi net.Conn) {
	fmt.Println("Klien terhubung:", koneksi.RemoteAddr().String())
	defer koneksi.Close()

	for {
		// Membaca pesan dari klien
		pesan, err := bufio.NewReader(koneksi).ReadString('\n')
		if err != nil {
			fmt.Println("Kesalahan saat membaca pesan:", err)
			return
		}

		// Menghapus spasi di awal dan akhir pesan
		pesan = strings.TrimSpace(pesan)

		// Jika pesan adalah "KELUAR", tutup koneksi
		if pesan == "KELUAR" {
			fmt.Println("Klien terputus:", koneksi.RemoteAddr().String())
			return
		}

		// Menampilkan pesan yang diterima
		fmt.Println("Pesan diterima:", pesan)

		// Mengirim kembali pesan ke klien
		_, err = koneksi.Write([]byte(pesan + "\n"))
		if err != nil {
			fmt.Println("Kesalahan saat mengirim pesan:", err)
			return
		}
	}
}

func main() {
	fmt.Println("Memulai server...")

	// Membuka port untuk mendengarkan koneksi masuk
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

		// Menangani koneksi dengan goroutine
		go tanganiKoneksi(koneksi)
	}
}
