package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

type store struct {
	mu sync.Mutex
	// maps: resource -> id -> object
	data map[string]map[string]map[string]any
	seq  map[string]int
}

func newStore() *store {
	return &store{
		data: map[string]map[string]map[string]any{},
		seq:  map[string]int{},
	}
}

func (s *store) create(resource string, obj map[string]any) map[string]any {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.data[resource]; !ok {
		s.data[resource] = map[string]map[string]any{}
	}
	s.seq[resource]++
	id := generateID(s.seq[resource])
	objCopy := map[string]any{}
	for k, v := range obj {
		objCopy[k] = v
	}
	objCopy["id"] = id
	objCopy["created_at"] = time.Now().Unix()
	s.data[resource][id] = objCopy
	return objCopy
}

// generateID creates a deterministic-ish 21-char id using the allowed charset
func generateID(seq int) string {
	charset := "useandom26T198340PX75pxJACKVERYMINDBUSHWOLFGQZbfghjklqvwyzrict-"
	l := len(charset)
	out := make([]byte, 21)
	start := seq % l
	for i := 0; i < 21; i++ {
		out[i] = charset[(start+i)%l]
	}
	return string(out)
}

func (s *store) get(resource, id string) (map[string]any, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if r, ok := s.data[resource]; ok {
		obj, found := r[id]
		return obj, found
	}
	return nil, false
}

func (s *store) update(resource, id string, obj map[string]any) (map[string]any, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if r, ok := s.data[resource]; ok {
		if _, found := r[id]; found {
			// merge
			for k, v := range obj {
				r[id][k] = v
			}
			return r[id], true
		}
	}
	return nil, false
}

func (s *store) delete(resource, id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if r, ok := s.data[resource]; ok {
		if _, found := r[id]; found {
			delete(r, id)
			return true
		}
	}
	return false
}

var s = newStore()

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", rootHandler)

	addr := ":3000"
	log.Printf("mock server listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, loggingMiddleware(mux)))
}

func loggingMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		h.ServeHTTP(w, r)
	})
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	// expected resource base paths: /api_keys, /domains, /listeners, /queues
	path := r.URL.Path
	parts := strings.Split(strings.Trim(path, "/"), "/")
	// Return 200 for root path so readiness checks succeed
	if path == "/" || path == "" {
		writeJSON(w, map[string]any{"status": "ok"})
		return
	}

	resource := parts[0]
	var id string
	if len(parts) > 1 {
		id = parts[1]
	}

	// accept any x-api-key
	_ = r.Header.Get("x-api-key")

	// support nested listener endpoints under queues: /queues/{queueId}/listeners[/{listenerId}]
	if resource == "queues" && len(parts) > 2 && parts[2] == "listeners" {
		var listenerId string
		if len(parts) > 3 {
			listenerId = parts[3]
		}
		handleListener(w, r, id, listenerId)
		return
	}

	switch resource {
	case "api_keys", "domains", "listeners", "queues":
		handleResource(w, r, resource, id)
	default:
		http.NotFound(w, r)
	}
}

func handleListener(w http.ResponseWriter, r *http.Request, queueId, id string) {
	key := "listener"
	switch r.Method {
	case http.MethodPost:
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		var obj map[string]any
		if len(body) == 0 {
			obj = map[string]any{}
		} else if err := json.Unmarshal(body, &obj); err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}
		// attach parent queue id
		obj["queue"] = queueId
		created := s.create("listeners", obj)
		resp := map[string]any{key: created}
		writeJSON(w, resp)
	case http.MethodGet:
		if id == "" {
			http.Error(w, "not implemented", http.StatusNotImplemented)
			return
		}
		if obj, ok := s.get("listeners", id); ok {
			resp := map[string]any{key: obj}
			writeJSON(w, resp)
			return
		}
		http.Error(w, "not found", http.StatusNotFound)
	case http.MethodPut:
		if id == "" {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		var obj map[string]any
		if err := json.Unmarshal(body, &obj); err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}
		if updated, ok := s.update("listeners", id, obj); ok {
			resp := map[string]any{key: updated}
			writeJSON(w, resp)
			return
		}
		http.Error(w, "not found", http.StatusNotFound)
	case http.MethodDelete:
		if id == "" {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		if s.delete("listeners", id) {
			resp := map[string]any{"_links": map[string]any{}}
			writeJSON(w, resp)
			return
		}
		http.Error(w, "not found", http.StatusNotFound)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleResource(w http.ResponseWriter, r *http.Request, resource, id string) {
	// determine singular key
	var key string
	switch resource {
	case "api_keys":
		key = "api_key"
	case "domains":
		key = "domain"
	case "listeners":
		key = "listener"
	case "queues":
		key = "queue"
	}

	switch r.Method {
	case http.MethodPost:
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		var obj map[string]any
		if len(body) == 0 {
			obj = map[string]any{}
		} else if err := json.Unmarshal(body, &obj); err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}
		created := s.create(resource, obj)
		// Ensure domains include challenge and verification objects to match client expectations
		if resource == "domains" {
			if _, ok := created["challenge"]; !ok {
				created["challenge"] = map[string]any{"https": map[string]any{"file_content": "challenge-content", "file_name": "challenge-file.txt", "context_path": "/.well-known/acme-challenge/abcd", "created_at": 0, "expires_at": 0}}
			}
			if _, ok := created["verification"]; !ok {
				created["verification"] = map[string]any{"verified": false, "verified_at": 0}
			}
		}
		resp := map[string]any{key: created}
		writeJSON(w, resp)
	case http.MethodGet:
		if id == "" {
			http.Error(w, "not implemented", http.StatusNotImplemented)
			return
		}
		if obj, ok := s.get(resource, id); ok {
			resp := map[string]any{key: obj}
			writeJSON(w, resp)
			return
		}
		http.Error(w, "not found", http.StatusNotFound)
	case http.MethodPut:
		if id == "" {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		var obj map[string]any
		if err := json.Unmarshal(body, &obj); err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}
		if updated, ok := s.update(resource, id, obj); ok {
			resp := map[string]any{key: updated}
			writeJSON(w, resp)
			return
		}
		http.Error(w, "not found", http.StatusNotFound)
	case http.MethodDelete:
		if id == "" {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		if s.delete(resource, id) {
			// return an empty _links object to match client expectations
			resp := map[string]any{"_links": map[string]any{}}
			writeJSON(w, resp)
			return
		}
		http.Error(w, "not found", http.StatusNotFound)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	_ = enc.Encode(v)
}
