package main

import (
	"fmt"
	"log"
	"net/http"

	"generatePDF/handlers"

	"github.com/gorilla/mux"
)

func main() {

	// Initialize handler
	pdfHandler := handlers.NewPDFHandler()

	// Setup router
	r := mux.NewRouter()

	// Routes
	r.HandleFunc("/api/generate-pdf", pdfHandler.GeneratePDF).Methods("GET")
	r.HandleFunc("/api/generate-receipt", pdfHandler.GenerateDownloadReceiptPDF).Methods("GET")

	// Health check
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"status": "ok", "service": "pdf-generator"}`)
	}).Methods("GET")

	fmt.Println("PDF Generator API running on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
