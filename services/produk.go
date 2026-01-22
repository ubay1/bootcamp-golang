package services

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
)

// data modelling like interface in typescript
type Produk struct {
	ID    int     `json:"id"` // ini adalah backtick, digunakan untuk membuat string yang bisa diakses oleh json
	Nama  string  `json:"nama"`
	Harga float64 `json:"harga"`
	Stok  int     `json:"stok"`
}

var produk = []Produk{
	{ID: 1, Nama: "Produk 1", Harga: 10000, Stok: 10},
	{ID: 2, Nama: "Produk 2", Harga: 20000, Stok: 20},
}

func getProdukByID(id int, w http.ResponseWriter) {
	for _, p := range produk {
		if p.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(p)
			return
		}
	}
	w.Header().Set("Content-Type", "application/json")
	http.Error(w, `{"message": "Produk Tidak Ditemukan"}`, http.StatusNotFound)
}

func uptadeProdukById(id int, body io.Reader, w http.ResponseWriter, err error) {
	var p Produk
	err = json.NewDecoder(body).Decode(&p)
	if err != nil {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	for i, prod := range produk {
		if prod.ID == id {
			p.ID = id
			produk[i] = p
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(p)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	http.Error(w, `{"message": "Produk Tidak Ditemukan"}`, http.StatusNotFound)
}

func deleteProdukById(id int, w http.ResponseWriter) {
	for i, prod := range produk {
		if prod.ID == id {
			produk = append(produk[:i], produk[i+1:]...)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"message": "Produk Deleted"})
			return
		}
	}
	w.Header().Set("Content-Type", "application/json")
	http.Error(w, `{"message": "Produk Tidak Ditemukan"}`, http.StatusNotFound)
}

func ProdukMain() {
	// get, update, delete produk by id
	http.HandleFunc("/api/produk/", func(w http.ResponseWriter, r *http.Request) {
		idStr := strings.TrimPrefix(r.URL.Path, "/api/produk/")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			http.Error(w, `{"message": "Invalid Produk ID"}`, http.StatusBadRequest)
			return
		}

		switch r.Method {
		case "GET":
			getProdukByID(id, w)
		case "PUT":
			uptadeProdukById(id, r.Body, w, err)
		case "DELETE":
			deleteProdukById(id, w)
		default:
			w.Header().Set("Content-Type", "application/json")
			http.Error(w, `{"message": "Method Not Allowed"}`, http.StatusMethodNotAllowed)
		}
	})

	// create produk, get all produk
	http.HandleFunc("/api/produk", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(produk)
		default:
			// baca data dari request & masukkan data kedalam produk
			var produkBaru Produk
			err := json.NewDecoder(r.Body).Decode(&produkBaru)
			if err != nil {
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
}
