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

	//********************* Esta parte es para no parar el bakend y volver a reiniciarlo***************
	// //Se reinicia la lista de los id´s montados
	// stores.ClearMountedPartitions()
	// utils.ResetMapsAndIndex()
	// //Se reinicia el login
	// comandos.SetearLogin()
	//************************** fin ******************************************************************

	//fmt.Println(requestData.Entrada)
	cmd, errs := analyzer.Analyzer(requestData.Entrada)

	//Se recorre el []error para obtener la respuesta
	var resultErrors string
	resultErrors = ""
	if len(errs) > 0 {
		//fmt.Println("Esta es la Longitud:")
		//fmt.Println(len(errs))
		for _, err := range errs {
			//fmt.Println(err)
			resultErrors += err.Error() + "\n"
		}
	}

	//Se recorre el []interface{} para obtener la respuesta
	var resultStr string
	resultStr = ""
	for _, item := range cmd {
		if item != nil && fmt.Sprintf("%+v", item) != "<nil>" {
			//fmt.Println(item)
			//resultStr += fmt.Sprintf("%+v \n", item)
		}
	}
	//Se unen los erroes al final
	resultStr += resultErrors
	response := ResponseData{
		Consola:    fmt.Sprint(resultStr),
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
