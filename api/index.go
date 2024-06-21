package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gocolly/colly"
	"github.com/gorilla/mux"
)


type Product struct {
	Title      string `json:"title"`
	Capital    string `json:"capital"`
	Population string `json:"population"`
	Area       string `json:"area"`
}


func scrapeWebsite() ([]Product, error) {
	var products []Product
	url := "https://www.scrapethissite.com/pages/simple/"
	collector := colly.NewCollector()

	collector.OnError(func(r *colly.Response, e error) {
		fmt.Println("Error occurred!:", e)
	})
	collector.OnHTML(".col-md-4.country", func(e *colly.HTMLElement) {
		product := Product{
			Title:      e.ChildText("h3"),
			Capital:    e.ChildText(".country-capital"),
			Population: e.ChildText(".country-population"),
			Area:       e.ChildText(".country-area"),
		}
		products = append(products, product)
	})

	
	err := collector.Visit(url)
	if err != nil {
		return nil, err
	}

	collector.Wait() 
	return products, nil
}


func getScrapedDataHandler(w http.ResponseWriter, r *http.Request) {
	data, err := scrapeWebsite()
	if err != nil {
		http.Error(w, "Oops, An error occured", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}


func main() {
	r := mux.NewRouter()
	
	r.HandleFunc("/api/countries", getScrapedDataHandler).Methods("GET")

	
	log.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
