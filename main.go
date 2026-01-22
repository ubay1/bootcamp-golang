package main

import (
	"encoding/json"
	"fmt"
	"kasir-api/services"
	"net/http"
)

func main() {
	services.ProdukMain()
	services.KategoriMain()

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
