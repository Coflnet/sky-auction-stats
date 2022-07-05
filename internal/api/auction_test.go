package api

import (
	"fmt"
	"github.com/Coflnet/auction-stats/internal/prometheus"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAuctionsValid(t *testing.T) {
	router := setupRouter()
	prometheus.StartPrometheus()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/new-auctions?duration=%s", url.QueryEscape("5")), nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestNewAuctionsTooSmall(t *testing.T) {
	router := setupRouter()
	prometheus.StartPrometheus()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/new-auctions?duration=%s", url.QueryEscape("-44")), nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestNewAuctionsTooBig(t *testing.T) {
	router := setupRouter()
	prometheus.StartPrometheus()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/new-auctions?duration=%s", url.QueryEscape("12345")), nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestNewAuctionsNoArgument(t *testing.T) {
	router := setupRouter()
	prometheus.StartPrometheus()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/new-auctions"), nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestNewAuctionsInvalid(t *testing.T) {
	router := setupRouter()
	prometheus.StartPrometheus()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/new-auctions?duration=%s", "hello-world"), nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestNewAuctionsOneDuration(t *testing.T) {
	router := setupRouter()
	prometheus.StartPrometheus()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/new-auctions?duration=%s", "1"), nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func FuzzNewAuctions(f *testing.F) {

	testcases := []string{"Hello, world", " ", "12345"}
	for _, tc := range testcases {
		f.Add(tc) // Use f.Add to provide a seed corpus
	}

	router := setupRouter()
	prometheus.StartPrometheus()

	f.Fuzz(func(t *testing.T, duration string) {

		fmt.Printf("fuzzing test runs with duration: %s\n", duration)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/new-auctions?duration=%s", url.QueryEscape(duration)), nil)
		router.ServeHTTP(w, req)

		if duration == "" {
			// if duration is an empty string, use default value 5
			assert.Equal(t, http.StatusOK, w.Code)
			return
		}

		v, err := strconv.Atoi(duration)

		if err != nil {
			fmt.Printf("no duration given, check if 200, actual: %v, duration: %v\n", w.Code, duration)
			assert.Equal(t, http.StatusBadRequest, w.Code)
			return
		}

		if v > 60*48 {
			fmt.Printf("v is bigger than %d: %v\n", 60*48, v)
			assert.Equal(t, http.StatusBadRequest, w.Code)
			return
		}

		if v < 1 {
			fmt.Printf("v is smaller than 1: %v\n", v)
			assert.Equal(t, http.StatusBadRequest, w.Code)
			return
		}

		fmt.Printf("check if 200, actual: %v, duration: %v\n", w.Code, duration)
		assert.Equal(t, http.StatusOK, w.Code)
	})

}
