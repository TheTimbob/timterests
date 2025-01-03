package server

import (
	"encoding/json"
	"log"
	"net/http"

	"timterests/cmd/web"
	"timterests/internal/storage"

	"github.com/a-h/templ"
)

func (s *Server) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	// Register routes
	mux.Handle("/", templ.Handler(web.HomeForm()))

	mux.HandleFunc("/health", s.healthHandler)

	fileServer := http.FileServer(http.FS(web.Files))
	mux.Handle("/assets/", fileServer)
	mux.Handle("/web", templ.Handler(web.HomeForm()))
	mux.Handle("/web/home", templ.Handler(web.HomeForm()))
	mux.Handle("/web/about", templ.Handler(web.AboutForm()))
	mux.Handle("/web/letters", templ.Handler(web.Letters()))

	// Projects Routes
	mux.Handle("/web/projects", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		web.ProjectsListHandler(w, r, *s.storage)
	}))
	mux.Handle("/web/project", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		articleID := r.URL.Query().Get("id")
		web.GetProjectsHandler(w, r, *s.storage, articleID)
	}))
	mux.Handle("/projects/list", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		web.ListProjects(*s.storage)
	}))

	// Article Routes
	mux.Handle("/web/articles", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		web.ArticlesListHandler(w, r, *s.storage)
	}))
	mux.Handle("/web/article", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		articleID := r.URL.Query().Get("id")
		web.GetArticleHandler(w, r, *s.storage, articleID)
	}))
	mux.Handle("/articles/list", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		web.ListArticles(*s.storage)
	}))

	// Wrap the mux with CORS middleware
	return s.corsMiddleware(mux)
}

func (s *Server) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*") // Replace "*" with specific origins if needed
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-CSRF-Token")
		w.Header().Set("Access-Control-Allow-Credentials", "false") // Set to "true" if credentials are required

		// Handle preflight OPTIONS requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Proceed with the next handler
		next.ServeHTTP(w, r)
	})
}

func (s *Server) HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	resp := map[string]string{"message": "Hello World"}
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(jsonResp); err != nil {
		log.Printf("Failed to write response: %v", err)
	}
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	//Marshal the health check response
	resp, err := json.Marshal(storage.Health(*s.storage))
	if err != nil {
		http.Error(w, "Failed to marshal health check response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(resp); err != nil {
		log.Printf("Failed to write response: %v", err)
	}
}
