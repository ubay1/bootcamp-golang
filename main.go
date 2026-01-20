package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// data modelling like interface in typescript
type Produk struct {
	ID       int     `json:"id"` // ini adalah backtick, digunakan untuk membuat string yang bisa diakses oleh json
	Nama     string  `json:"nama"`
	Harga    float64 `json:"harga"`
	Stok 		 int     `json:"stok"`
}

var produk = []Produk{
	{ID: 1, Nama: "Produk 1", Harga: 10000, Stok: 10},
	{ID: 2, Nama: "Produk 2", Harga: 20000, Stok: 20},
}


func main() {
	// get produk by id
	http.HandleFunc("/api/produk/", func(w http.ResponseWriter, r *http.Request) {
		idStr := strings.TrimPrefix(r.URL.Path, "/api/produk/")
		id, err := strconv.Atoi(idStr)
		if (err != nil) {
		w.Header().Set("Content-Type", "application/json")
			http.Error(w, `{"message": "Invalid Produk ID"}`, http.StatusBadRequest)
			return
		}

		for _, p := range produk {
			if p.ID == id {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(p)
				return
			}
		}
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, `{"message": "Produk Tidak Ditemukan"}`, http.StatusNotFound)
	})

	// add produk
	http.HandleFunc("/api/produk", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(produk)
		case "PUT":
			idStr := strings.TrimPrefix(r.URL.Path, "/api/produk/")
			id, err := strconv.Atoi(idStr)
			if (err != nil) {
				http.Error(w, "Invalid Produk ID", http.StatusBadRequest)
				return
			}

			var p Produk
			err = json.NewDecoder(r.Body).Decode(&p)
			if (err != nil) {
				http.Error(w, "Invalid Request", http.StatusBadRequest)
				return
			}

			for i, prod := range produk {
				if prod.ID == id {
					produk[i] = p
					w.Header().Set("Content-Type", "application/json")
					json.NewEncoder(w).Encode(p)
					return
				}
			}
			
			http.Error(w, "Produk Tidak Ditemukan", http.StatusNotFound)
		default:
			// baca data dari request & masukkan data kedalam produk
			var produkBaru Produk
			err := json.NewDecoder(r.Body).Decode(&produkBaru)
			if (err != nil) {
				http.Error(w, "Invalid Request", http.StatusBadRequest)
			} else {
				produkBaru.ID = len(produk) + 1
				produk = append(produk, produkBaru)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated)
				json.NewEncoder(w).Encode(produkBaru)
			}
		}
	})

	// health check
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "OK", "message": "API Running"})
	})

	// server running
	fmt.Println("Server running on port 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("gagal running server")
	}
}
