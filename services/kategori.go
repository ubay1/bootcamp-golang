package services

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type Kategori struct {
	ID          int    `json:"id"` // ini adalah backtick, digunakan untuk membuat string yang bisa diakses oleh json
	Name        string `json:"nama"`
	Description string `json:"deskripsi"`
}

var kategori = []Kategori{
	{ID: 1, Name: "Mainan Anak", Description: "sayang anak sayang anak"},
	{ID: 2, Name: "Sepatu", Description: "ada banyak jenis sepatu"},
}

func getKategoriByID(id int, w http.ResponseWriter) {
	for _, k := range kategori {
		if k.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(k)
			return
		}
	}
	w.Header().Set("Content-Type", "application/json")
	http.Error(w, `{"message": "Kategori Tidak Ditemukan"}`, http.StatusNotFound)
}

func updateKategoriById(id int, body io.Reader, w http.ResponseWriter, err error) {
	var k Kategori
	err = json.NewDecoder(body).Decode(&k)
	if err != nil {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	for i, cat := range kategori {
		if cat.ID == id {
			k.ID = id
			kategori[i] = k
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(k)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	http.Error(w, `{"message": "Kategori Tidak Ditemukan"}`, http.StatusNotFound)
}

func deleteKategoriById(id int, w http.ResponseWriter) {
	for i, cat := range kategori {
		if cat.ID == id {
			kategori = append(kategori[:i], kategori[i+1:]...)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"message": "Kategori Deleted"})
			return
		}
	}
	w.Header().Set("Content-Type", "application/json")
	http.Error(w, `{"message": "Produk Tidak Ditemukan"}`, http.StatusNotFound)
}

func KategoriMain() {
	// get, update, delete kategori by id
	http.HandleFunc("/api/categories/", func(w http.ResponseWriter, r *http.Request) {
		idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			http.Error(w, `{"message": "Invalid Kategori ID"}`, http.StatusBadRequest)
			return
		}

		switch r.Method {
		case "GET":
			getKategoriByID(id, w)
		case "PUT":
			updateKategoriById(id, r.Body, w, err)
		case "DELETE":
			deleteKategoriById(id, w)
		default:
			w.Header().Set("Content-Type", "application/json")
			http.Error(w, `{"message": "Method Not Allowed"}`, http.StatusMethodNotAllowed)
		}
	})

	// create kategori, get all kategori
	http.HandleFunc("/api/categories", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(kategori)
		default:
			// baca data dari request & masukkan data kedalam kategori
			var kategoriBaru Kategori
			err := json.NewDecoder(r.Body).Decode(&kategoriBaru)
			if err != nil {
				http.Error(w, "Invalid Request", http.StatusBadRequest)
			} else {
				kategoriBaru.ID = len(kategori) + 1
				kategori = append(kategori, kategoriBaru)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated)
				json.NewEncoder(w).Encode(kategoriBaru)
			}
		}
	})
}
