package middleware

import (
	"fmt"
	"net/http"

	"github.com/microcosm-cc/bluemonday"
)

func XSSMiddleware(next http.Handler) http.Handler {
	fmt.Println("INFO Initializing XSS Middleware")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("INFO XSS Middleware Run...")

		// clean url path
		sanitizedPath, err := clean(r.URL.Path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		fmt.Println("Original Path:", r.URL.Path)
		fmt.Println("Sanitized Path:", sanitizedPath)

		next.ServeHTTP(w, r)
		fmt.Println("INFO Sending responce from XSS Middleware Run")
	})
}

// Clean sanitize input data to prevent XSS attacks
func clean(data any) (any, error) {
	switch v := data.(type) {
	case map[string]any:
		for key, value := range v {
			v[key] = sanitizeValue(value)
		}
		return v, nil

	case []any:
		for key, value := range v {
			v[key] = sanitizeValue(value)
		}
		return v, nil

	case string:
		return sanitizeString(v), nil

	default:
		// error
		return nil, fmt.Errorf("unsupported type %T", data)
	}
}

func sanitizeValue(data any) any {
	switch v := data.(type) {
	case string:
		return sanitizeString(v)
	case map[string]any:
		for key, value := range v {
			v[key] = sanitizeValue(value)
		}
		return v
	case []any:
		for key, value := range v {
			v[key] = sanitizeValue(value)
		}
		return v

	default:
		return v
	}
}

func sanitizeString(data string) string {
	return bluemonday.UGCPolicy().Sanitize(data)
}
