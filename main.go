package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"app/Godeps/_workspace/src/github.com/Azure/azure-sdk-for-go/storage"
	"app/Godeps/_workspace/src/github.com/aymerick/raymond"
	"app/Godeps/_workspace/src/github.com/gorilla/mux"
)

var (
	storageName = os.Getenv("STORAGE_NAME")
	storageKey  = os.Getenv("STORAGE_KEY")
)

func main() {
	fmt.Println("Starting Nofification Service")

	// Note; I believe we can get subrouters or setup the base api path for /template
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/template/", ListTemplates).
		Methods("GET")

	router.HandleFunc("/template/{id}", SetTemplate).
		Methods("PUT")

	router.HandleFunc("/template/{id}", GetTemplate).
		Methods("GET")

	router.HandleFunc("/template/{id}", DeleteTemplate).
		Methods("DELETE")

	router.HandleFunc("/template/{id}/view", ViewTemplate).
		Methods("POST")

	// TODO: Need to setup RabbitMQ Listener

	fmt.Println(http.ListenAndServe(":8000", router), nil)

	fmt.Println("Exited Piccolo")
}

// List templates in the repository
func ListTemplates(w http.ResponseWriter, r *http.Request) {
	// return a list of template objects from BLOB storage

}

// Generate a view of the template using the data passed
func ViewTemplate(w http.ResponseWriter, r *http.Request) {

	// Read in JSON data to bind
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if err := r.Body.Close(); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// unmarshal the json to map
	var objmap map[string]interface{}
	err = json.Unmarshal(body, &objmap)

	//msgData, err := gabs.ParseJSON([]byte(body))

	// load template from repository
	vars := mux.Vars(r)
	id := vars["id"]

	// get template from store
	tmpl, err := getTemplateById(id)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// bind the data to the template retrieved from storage
	result, err := raymond.Render(tmpl, objmap)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	fmt.Fprintln(w, result)
}

// Get a template from the repository by id
func GetTemplate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	tmpl, err := getTemplateById(id)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json;")

	fmt.Fprintf(w, tmpl)
	return
}

// Add or update a template in the repository by id
func SetTemplate(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "URL PATH - %s!", r.URL.Path[1:])
}

// Delete a template in the repository
func DeleteTemplate(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "URL PATH -  %s!", r.URL.Path[1:])
}

func ProcessMessage() {

}

// Retrieve a template from the repository
func getTemplateById(id string) (string, error) {
	cli, err := storage.NewBasicClient(storageName, storageKey)
	if err != nil {
		return "", err
	}

	blobService := cli.GetBlobService()
	blobBody, err := blobService.GetBlob("templates", id+".tmpl.html")
	if err != nil {
		return "", err
	}

	tmpl, err := ioutil.ReadAll(blobBody)
	defer blobBody.Close()

	return string(tmpl), nil
}
