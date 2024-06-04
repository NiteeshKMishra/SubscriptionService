package utils

import "net/http"

func GetFullRequestPath(r *http.Request) string {
	queryParams := r.URL.Query()
	pathWithQueryParams := r.URL.Path

	for key, values := range queryParams {
		for _, value := range values {
			pathWithQueryParams += "?" + key + "=" + value
		}
	}

	return pathWithQueryParams
}
