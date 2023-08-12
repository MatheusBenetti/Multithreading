package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"time"
)

type ViaCEP struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

type BrasilAPI struct {
	Cep          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
	Service      string `json:"service"`
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/cep/{cep}", handleRequest)
	http.Handle("/", r)
	log.Println("Server listening on port :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	cep := params["cep"]
	apiChan := make(chan ViaCEP)
	apiChan2 := make(chan BrasilAPI)

	go fetchViaCep(cep, apiChan)
	go fetchBrasilAPI(cep, apiChan2)

	select {
	case msg := <-apiChan:
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(msg)
		if err != nil {
			panic(err)
		}
		fmt.Println("API: Via CEP")
		fmt.Println("Resposta: ", msg)

	case msg := <-apiChan2:
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(msg)
		if err != nil {
			panic(err)
		}
		fmt.Println("API: Brasil API")
		fmt.Println("Resposta: ", msg)

	case <-time.After(time.Second * 1):
		http.Error(w, "Timeout", http.StatusRequestTimeout)
	}
}

func fetchViaCep(cep string, apiChan chan<- ViaCEP) {
	req, err := http.Get("http://viacep.com.br/ws/" + cep + "/json/")
	if err != nil {
		panic(err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(req.Body)

	res, err := io.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}

	var data ViaCEP
	err = json.Unmarshal(res, &data)
	if err != nil {
		panic(err)
	}

	apiChan <- data
}

func fetchBrasilAPI(cep string, apiChan2 chan<- BrasilAPI) {
	req, err := http.Get("https://brasilapi.com.br/api/cep/v1/" + cep)
	if err != nil {
		panic(err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(req.Body)

	res, err := io.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}

	var data BrasilAPI
	err = json.Unmarshal(res, &data)
	if err != nil {
		panic(err)
	}
	apiChan2 <- data
}
