package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"sync"
)

type dollars float32

func (d dollars) String() string { return fmt.Sprintf("$%.2f", d) }

type database map[string]dollars

var dbLock sync.Mutex

var listTemplate string = `
<h1>Database Items</h1>
<table>
  <tr>
    <th>Item</th>
    <th>Price</th>
  </tr>
{{range $key, $value := .}}
  <tr>
    <td>{{$key}}</td>
    <td>{{$value}}</td>
  </tr>
{{end}}
</table>
`

func main() {
	db := database{"shoes": 50, "socks": 3}
	mux := http.NewServeMux()
	mux.HandleFunc("/list", http.HandlerFunc(db.list))
	mux.HandleFunc("/price", http.HandlerFunc(db.price))
	mux.HandleFunc("/read", http.HandlerFunc(db.read))
	mux.HandleFunc("/update", http.HandlerFunc(db.update))
	mux.HandleFunc("/create", http.HandlerFunc(db.create))
	mux.HandleFunc("/delete", http.HandlerFunc(db.del))
	log.Fatal(http.ListenAndServe("localhost:8000", mux))
}

func (db database) list(w http.ResponseWriter, req *http.Request) {
	itemTable := template.Must(template.New("items").Parse(listTemplate))
	itemTable.Execute(w, db)
	// for item, price := range db {
	// 	fmt.Fprintf(w, "%s: %s\n", item, price)
	// }
}

func (db database) price(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")
	price, ok := db[item]
	if !ok {
		w.WriteHeader(http.StatusNotFound) // 404
		fmt.Fprintf(w, "no such item: %q\n", item)
		return
	}
	fmt.Fprintf(w, "%s\n", price)
}

// Exercise 7.11 Add CRUD handlers
func (db database) create(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")
	if item == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "must specify an 'item' to create")
		return
	}
	price := req.URL.Query().Get("price")
	if price == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "must specify an 'price' for the item")
		return
	}

	val, err := strconv.ParseFloat(price, 32) // ParseFloat returns a float64 regardless of specified bitsize but conversion to float32 is lossless
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Unable to parse price %v", err)
		return
	}

	dbLock.Lock()
	defer dbLock.Unlock()
	db[item] = dollars(val)
	w.WriteHeader(http.StatusOK) // not necessary, if WriteHeader is not called the first call to Write implicitly calls WriteHeader(http.StatusOk)
	fmt.Fprintf(w, "Item %q created successfully\n", item)
}
func (db database) read(w http.ResponseWriter, req *http.Request) {
	items, ok := req.URL.Query()["item"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Must specify at least one item")
		return
	}
	for _, item := range items {
		p, ok := db[item]
		if !ok {
			fmt.Fprintf(w, "%q not found in database\n", item)
		} else {
			fmt.Fprintf(w, "%s: %s\n", item, p)
		}
	}
	return
}
func (db database) update(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")
	if item == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "must specify an 'item' to update")
		return
	}
	price := req.URL.Query().Get("price")
	if price == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "must specify a 'price' to update an item")
		return
	}
	val, err := strconv.ParseFloat(price, 32)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Unable to parse price %v", err)
		return
	}

	dbLock.Lock()
	defer dbLock.Unlock()

	db[item] = dollars(val)
	fmt.Fprintf(w, "Item %q updated to price %s", item, price)
}

func (db database) del(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")
	if item == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "must specify an 'item' to update")
		return
	}

	dbLock.Lock()
	defer dbLock.Unlock()

	delete(db, item)
	fmt.Fprintf(w, "Item %q successfully deleted\n", item)
}
