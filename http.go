package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"golang.org/x/talks/content/2016/applicative/google"
)

func main() {
	http.HandleFunc("/search", handleSearch)
	fmt.Println("serving on http://localhost:8080/search")
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}

func handleSearch(w http.ResponseWriter, req *http.Request) {
	log.Println("serving", req.URL)

	// Check the search query.
	query := req.FormValue("q")
	if query == "" {
		http.Error(w, `missing "q" URL parameter`, http.StatusBadRequest)
		return
	}

	// Run the Google search.
	start := time.Now()
	results, err := google.Search(query)
	elapsed := time.Since(start)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// ENDSEARCH OMIT

	// Create the structured response.
	type response struct {
		Results []google.Result
		Elapsed time.Duration
	}
	resp := response{results, elapsed}
	// ENDRESPONSE OMIT

	// Render the response.
	switch req.FormValue("output") {
	case "json":
		err = json.NewEncoder(w).Encode(resp) 
	case "prettyjson":
		var b []byte
		b, err = json.MarshalIndent(resp, "", "  ") 
		if err == nil {
			_, err = w.Write(b)
		}
	default: // HTML
		err = responseTemplate.Execute(w, resp) 
	}
	// ENDRENDER OMIT
	if err != nil {
		log.Print(err)
		return
	}
}

var responseTemplate = template.Must(template.New("results").Parse(`
<html>
<head/>
<body>
  <ol>
  {{range .Results}}
    <li>{{.Title}} - <a href="{{.URL}}">{{.URL}}</a></li>
  {{end}}
  </ol>
  <p>{{len .Results}} results in {{.Elapsed}}</p>
</body>
</html>
`))