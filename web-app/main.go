package main

import (
	"context"
	"embed"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	//"github.com/abhirockzz/amazon-bedrock-langchain-go/llm/claude"
	"github.com/build-on-aws/langchaingo-amazon-bedrock-llm/claude"
	"github.com/gorilla/mux"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/documentloaders"
	"github.com/tmc/langchaingo/schema"
)

var llm *claude.LLM
var chain chains.StuffDocuments
var docs []schema.Document

const defaultRegion = "us-east-1"

func init() {
	var err error

	region := os.Getenv("AWS_REGION")
	if region == "" {
		region = defaultRegion
	}

	llm, err = claude.New(region)
	//llm.CallbacksHandler = callbacks.LogHandler{}

	if err != nil {
		log.Fatal(err)
	}

	chain = chains.LoadStuffQA(llm)
}

func main() {

	router := mux.NewRouter()

	router.HandleFunc("/", serveStaticFiles).Methods(http.MethodGet)
	router.HandleFunc("/upload", loadData).Methods(http.MethodPost)
	router.HandleFunc("/chat", chat).Methods(http.MethodPost)

	fmt.Println("http server started...")

	log.Fatal(http.ListenAndServe(":8080", router))
}

func loadData(w http.ResponseWriter, r *http.Request) {

	b, err := io.ReadAll(r.Body)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	link := string(b)

	defer r.Body.Close()

	docs, err = getDocsFromLink(link)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

}

func chat(w http.ResponseWriter, req *http.Request) {

	body, err := io.ReadAll(req.Body)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	chatMessage := string(body)

	fmt.Println("chat message", chatMessage)

	answer, err := chains.Call(context.Background(), chain, map[string]any{
		"input_documents": docs,
		"question":        chatMessage,
	}, chains.WithMaxTokens(2048))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println("chat message response", answer["text"])

	w.Write([]byte(answer["text"].(string)))
}

func getDocsFromLink(link string) ([]schema.Document, error) {

	fmt.Println("loading data from", link)

	resp, err := http.Get(link)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	docs, err := documentloaders.NewHTML(resp.Body).Load(context.Background())

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("successfully loaded data from", link)

	return docs, nil
}

//go:embed static/*
var staticFiles embed.FS

const fileName = "index.html"

func serveStaticFiles(w http.ResponseWriter, r *http.Request) {

	path := "static" + r.URL.Path + fileName
	fmt.Println("reading from", path)

	file, err := staticFiles.ReadFile(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	contentType := http.DetectContentType(file)

	w.Header().Set("Content-Type", contentType)
	w.Write(file)
}
