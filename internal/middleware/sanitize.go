package middleware

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

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

		params := r.URL.Query()

		sanitizedQuery := make(map[string][]string)

		for key, values := range params {
			sanitizedKey, err := clean(key)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			sanitizedValues := make([]string, len(values))
			for i, value := range values {
				sanitizedValue, err := clean(value)
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
				sanitizedValues[i] = sanitizedValue.(string)
			}
			sanitizedQuery[sanitizedKey.(string)] = sanitizedValues
		}
		// Apply sanitized values back to the request
		r.URL.Path = sanitizedPath.(string)

		// Rebuild query string with sanitized parameters
		queryValues := make([]string, 0)
		for key, values := range sanitizedQuery {
			for _, value := range values {
				queryValues = append(
					queryValues,
					fmt.Sprintf("%s=%s", url.QueryEscape(key), url.QueryEscape(value)),
				)
			}
		}

		if len(queryValues) > 0 {
			r.URL.RawQuery = strings.Join(queryValues, "&")
		} else {
			r.URL.RawQuery = ""
		}
		fmt.Println("Updated URL:", r.URL.String())

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
