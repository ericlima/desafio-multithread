package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type Apicep struct {
	Status   int    `json:"status"`
	Code     string `json:"code"`
	State    string `json:"state"`
	City     string `json:"city"`
	District string `json:"district"`
	Address  string `json:"address"`
}

type Viacep struct {
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

func main() {

	// le argumento passado via CLI
	if len(os.Args) < 2 {
		println("Por favor informe o CEP no formato 12345678")
		return
	}

	var cep string = os.Args[1]

	ch_apicep := make(chan Apicep)
	ch_viacep := make(chan Viacep)

	go func() {
		// consome apicep
		url := fmt.Sprintf("https://cdn.apicep.com/file/apicep/%s.json", cep)
		
		body,err := ConsomeCEP(url)
		if err != nil {
			fmt.Println("Erro na requisição HTTP:", err)
			return
		}

		var apiCEPResponse Apicep
		err = json.Unmarshal(body, &apiCEPResponse)
		if err != nil {
			fmt.Println("Erro ao fazer o parse JSON:", err)
			return
		}

		ch_apicep <-apiCEPResponse
	}()

	go func() {
		// consome viacep
		url := fmt.Sprintf("http://viacep.com.br/ws/%s/json/", cep)
		
		body,err := ConsomeCEP(url)
		if err != nil {
			fmt.Println("Erro na requisição HTTP:", err)
			return
		}

		var viaCEPResponse Viacep
		err = json.Unmarshal(body, &viaCEPResponse)
		if err != nil {
			fmt.Println("Erro ao fazer o parse JSON:", err)
			return
		}

		ch_viacep <-viaCEPResponse

	}()

	select {
	case msg := <-ch_apicep:
		fmt.Println("resultado via apicep:", msg)
	case msg := <-ch_viacep:
		fmt.Println("resultado via viacep:", msg)
	case <-time.After(time.Second):
		println("Timeout!")
	}

}

func ConsomeCEP(url string) ([]byte, error) {
	if len(url) == 0 {
		panic("url não informada")
	}
	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Erro na requisição HTTP:", err)
		return nil,err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil,errors.New(fmt.Sprintf("%s http status code %d\n",url,response.StatusCode))
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Erro no corpo da resposta:", err)
		return nil,err
	}


	return body,nil
}
