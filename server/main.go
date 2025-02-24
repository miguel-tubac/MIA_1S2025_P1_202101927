package main

import (
	"bakend/src/analyzer"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type RequestData struct {
	Entrada string `json:"entrada"`
}

type ResponseData struct {
	Consola    string   `json:"consola"`
	TablaError []string `json:"tablaError"`
}

func interpretarHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var requestData RequestData
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&requestData); err != nil {
		http.Error(w, "Error al procesar la solicitud", http.StatusBadRequest)
		return
	}

	cmd, err := analyzer.Analyzer(requestData.Entrada)
	if err != nil {
		log.Println("Error al analizar la entrada:", err)
		response := ResponseData{Consola: "", TablaError: []string{err.Error()}}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	response := ResponseData{
		Consola:    fmt.Sprintf("Parsed Command: %+v", cmd),
		TablaError: []string{},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/interpretar", interpretarHandler).Methods("POST")

	// Habilitar CORS
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	originsOk := handlers.AllowedOrigins([]string{"http://localhost:3000"}) // Permite solo React
	methodsOk := handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS"})

	fmt.Println("Servidor corriendo en http://localhost:4000")
	log.Fatal(http.ListenAndServe(":4000", handlers.CORS(headersOk, originsOk, methodsOk)(router)))
}
