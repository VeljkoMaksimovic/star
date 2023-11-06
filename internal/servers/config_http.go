package servers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	kuiperapi "github.com/c12s/kuiper/pkg/api"
	"github.com/c12s/star/internal/mappers/proto"
	"github.com/c12s/star/internal/services"
)

type applyConfigCommand struct {
	ID      string `json:"id"`
	Configs []struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	} `json:"configs"`
}

func StartHttpConfigServer(service *services.ConfigService) {
	http.HandleFunc("/apply-config", func(w http.ResponseWriter, r *http.Request) {
		// Only respond to POST requests
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Read the body of the request
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		// Deserialize the JSON into an ApplyConfigCommand instance
		var cmd kuiperapi.ApplyConfigCommand
		if err := json.Unmarshal(body, &cmd); err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		req, err := proto.ApplyConfigCommandToDomain(&cmd)
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Printf("Received config: %+v\n", req)

		_, err = service.Put(*req)
		if err != nil {
			log.Println(err)
		}
		// Respond to the client
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})
	})

	fmt.Println("Starting server on port 7777")
	go func() {
		if err := http.ListenAndServe(":7777", nil); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()
}
