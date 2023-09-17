package crud

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
)

type Item struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

type RouteMap map[string]func(http.ResponseWriter, *http.Request)

type BasicCRUD struct {
	Routes RouteMap
	Port   string
	Items  []Item
}

func NewBasicCRUD(p ...string) BasicCRUD {
	port := "8088"
	if len(p) > 0 {
		port = p[0]
	}

	routes := RouteMap{}

	return BasicCRUD{
		Routes: routes,
		Port:   port,
	}
}

func (b BasicCRUD) Serve() {
	errLogFile, err := os.Create("BasicCRUD-error.log")
	if err != nil {
		log.Fatalln("Error creating error.log file:", err)
	}
	defer errLogFile.Close()

	errLogger := log.New(errLogFile, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

	// Set the logger for the standard log package to log to both stdout and error.log
	log.SetOutput(io.MultiWriter(os.Stdout, errLogFile))

	for route, handler := range b.Routes {
		http.HandleFunc(route, handler)
	}

	// Start the HTTP server on port 8080
	if err := http.ListenAndServe(":"+b.Port, nil); err != nil {
		errLogger.Println(err)
	}
}

func (b *BasicCRUD) GetItems(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(b.Items)
}

func (b *BasicCRUD) HandleItem(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		b.GetItem(w, r)
	} else if r.Method == http.MethodPost {
		b.CreateItem(w, r)
	} else if r.Method == http.MethodPut {
		b.UpdateItem(w, r)
	} else if r.Method == http.MethodDelete {
		b.DeleteItem(w, r)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (b *BasicCRUD) GetItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := r.URL.Path[len("/items/"):]
	for _, item := range b.Items {
		if item.ID == id {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	http.NotFound(w, r)
}

func (b *BasicCRUD) CreateItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var newItem Item
	json.NewDecoder(r.Body).Decode(&newItem)
	b.Items = append(b.Items, newItem)
	json.NewEncoder(w).Encode(newItem)
}

func (b *BasicCRUD) UpdateItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := r.URL.Path[len("/items/"):]
	for index, item := range b.Items {
		if item.ID == id {
			json.NewDecoder(r.Body).Decode(&b.Items[index])
			json.NewEncoder(w).Encode(b.Items[index])
			return
		}
	}
	http.NotFound(w, r)
}

func (b *BasicCRUD) DeleteItem(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/items/"):]
	for index, item := range b.Items {
		if item.ID == id {
			b.Items = append(b.Items[:index], b.Items[index+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	http.NotFound(w, r)
}
