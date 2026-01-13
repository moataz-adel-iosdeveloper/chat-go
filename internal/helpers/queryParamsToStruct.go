package helpers

import (
	"encoding/json"
	"log"
	"net/http"
)

// QueryParamsToStruct parses query params (or form values) into a struct
func QueryParamsToStruct(r *http.Request, out any) error {
	// Parse multipart form (set a maxMemory for non-file parts)
	if err := r.ParseMultipartForm(10 << 20); err != nil && err != http.ErrNotMultipart {
		log.Println("Error parsing multipart form:", err)
		return err
	}

	data := make(map[string]string)

	// Read form values
	if r.MultipartForm != nil {
		for key, values := range r.MultipartForm.Value {
			if len(values) > 0 {
				log.Println("key:", key, "value:", values[0])
				data[key] = values[0]
			}
		}
	}

	// Convert map to JSON bytes
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		log.Println("Error marshaling form data:", err)
		return err
	}

	// Unmarshal JSON into the target struct
	if err := json.Unmarshal(jsonBytes, out); err != nil {
		log.Println("Error unmarshaling form data:", err)
		return err
	}

	return nil
}
