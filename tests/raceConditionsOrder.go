package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

// OrderCreateRequest merepresentasikan struktur JSON payload untuk membuat order.
type OrderCreateRequest struct {
	UserID string `json:"user_id"`
	Items  []struct {
		ProductID string `json:"product_id"`
		Quantity  int    `json:"quantity"`
	} `json:"items"`
}

// OrderResponse merepresentasikan struktur respons sukses dari API Anda.
type OrderResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    struct {
		OrderID    string `json:"order_id"`
		UserID     string `json:"user_id"`
		TotalPrice float64 `json:"total_price"`
		Status     string `json:"status"`
		// ... tambahkan field lain jika ada
	} `json:"data"`
}

// ErrorResponse merepresentasikan struktur respons error dari API Anda.
type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Error   string `json:"error"` // Atau sesuai dengan nama field error Anda
}

func main() {
	// --- Konfigurasi Simulasi ---
	apiURL := "http://localhost:8080/app/api/v1/orders" // Ganti dengan URL endpoint CreateOrder Anda
	// ID Produk yang akan diuji race condition-nya (product_id_X)
	productXID := "63dbcbda-3fc4-4e41-a40c-de3131ca6a74"
	// ID Produk lain yang stoknya cukup (product_id_Y)
	productYID := "bc4c0e8a-aa24-4809-b332-76d52da338d5"

	// --- TEMPAT UNTUK JWT TOKEN ANDA ---
	// Ganti ini dengan JWT token yang valid dari sistem otentikasi Anda.
	// Anda mungkin perlu login ke aplikasi Anda untuk mendapatkan token ini.
	// Contoh: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
	const jwtToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjFkODUxODhlLWVmMTEtNDQyZS04MjYwLTJlNTg0YjIxNDExMyIsImVtYWlsIjoidXNlckBnbWFpbC5jb20iLCJyb2xlIjoidXNlciIsImlzcyI6IkRlcHVibGljIiwiZXhwIjoxNzUzMTU3MjEzfQ.jNOTL_4kiwXVWR5QcnBrDKyJqtWFwxWVcoBQwMX0BZA" // <--- GANTI INI!

	if jwtToken == "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjFkODUxODhlLWVmMTEtNDQyZS04MjYwLTJlNTg0YjIxNDExMyIsImVtYWlsIjoidXNlckBnbWFpbC5jb20iLCJyb2xlIjoidXNlciIsImlzcyI6IkRlcHVibGljIiwiZXhwIjoxNzUzMTU3MjEzfQ.jNOTL_4kiwXVWR5QcnBrDKyJqtWFwxWVcoBQwMX0BZA" || jwtToken == "" {
		fmt.Println("WARNING: Please replace 'YOUR_VALID_JWT_TOKEN_HERE' with an actual valid JWT token.")
		fmt.Println("The simulation might fail due to authentication errors.")
		// return // Anda bisa uncomment ini untuk menghentikan jika token belum diganti
	}

	// Payload JSON untuk kedua goroutine
	requestPayload := OrderCreateRequest{
		UserID: "ec4a6412-2801-4811-9a11-b3d11292708f", // Contoh User ID, bisa sama
		Items: []struct {
			ProductID string `json:"product_id"`
			Quantity  int    `json:"quantity"`
		}{
			{
				ProductID: productXID,
				Quantity:  10, // Setiap goroutine memesan 10 unit Product X
			},
			{
				ProductID: productYID,
				Quantity:  10, // Setiap goroutine memesan 10 unit Product Y
			},
		},
	}

	var wg sync.WaitGroup
	results := make(chan string, 2) // Channel untuk menampung hasil dari setiap goroutine

	fmt.Println("Starting two concurrent order creation attempts for the same products...")

	// Goroutine 1
	wg.Add(1)
	go func(id string, payload OrderCreateRequest, token string) {
		defer wg.Done()
		fmt.Printf("[%s] Attempting to create order...\n", id)
		err := sendOrderRequest(apiURL, payload, token) // Meneruskan token
		if err != nil {
			results <- fmt.Sprintf("[%s] FAILED: %v", id, err)
		} else {
			results <- fmt.Sprintf("[%s] SUCCEEDED", id)
		}
	}("Goroutine A", requestPayload, jwtToken) // Meneruskan token di sini

	// Beri sedikit jeda agar Goroutine A punya kesempatan untuk mulai mengunci
	// (meskipun dalam realitas race condition bisa terjadi sangat cepat)
	time.Sleep(50 * time.Millisecond)

	// Goroutine 2 (dengan UserID yang berbeda untuk membedakan order)
	wg.Add(1)
	requestPayload2 := requestPayload
	requestPayload2.UserID = "1d85188e-ef11-442e-8260-2e584b214113" // UserID berbeda
	go func(id string, payload OrderCreateRequest, token string) {
		defer wg.Done()
		fmt.Printf("[%s] Attempting to create order...\n", id)
		err := sendOrderRequest(apiURL, payload, token) // Meneruskan token
		if err != nil {
			results <- fmt.Sprintf("[%s] FAILED: %v", id, err)
		} else {
			results <- fmt.Sprintf("[%s] SUCCEEDED", id)
		}
	}("Goroutine B", requestPayload2, jwtToken) // Meneruskan token di sini

	wg.Wait()
	close(results)

	fmt.Println("\n--- Simulation Results ---")
	for res := range results {
		fmt.Println(res)
	}

	// --- VERIFIKASI AKHIR ---
	fmt.Println("\n--- Important: Manual Database Verification Required ---")
	fmt.Println("Please check the 'stock' of the following products in your database:")
	fmt.Printf("- Product ID: %s (Product X) -> Expected stock: (Initial Stock) - (Number of Successful Orders * 10)\n", productXID)
	fmt.Printf("- Product ID: %s (Product Y) -> Expected stock: (Initial Stock) - (Number of Successful Orders * 10)\n", productYID)
	fmt.Println("\nAlso, verify the total number of order entries for each product.")
	fmt.Println("You should see one successful order (Status 201) and one failed order (Status 400 - insufficient stock) if initial stock for Product X was, for example, 15 or 19.")
}

// sendOrderRequest mengirim permintaan POST ke API dan mengembalikan error jika ada.
// Sekarang menerima parameter jwtToken.
func sendOrderRequest(url string, payload OrderCreateRequest, jwtToken string) error {
	jsonBody, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// --- MENAMBAHKAN JWT TOKEN KE HEADER ---
	if jwtToken != "" {
		req.Header.Set("Authorization", "Bearer "+jwtToken)
	}

	client := &http.Client{Timeout: 10 * time.Second} // Tambahkan timeout
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request to API failed: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusCreated {
		var errorResp ErrorResponse
		// Coba unmarshal sebagai ErrorResponse, jika gagal, tampilkan raw body
		if err := json.Unmarshal(bodyBytes, &errorResp); err != nil {
			return fmt.Errorf("API returned status %d, but failed to parse error response. Raw body: %s", resp.StatusCode, string(bodyBytes))
		}
		// Jika unmarshal berhasil, tampilkan pesan error yang lebih spesifik
		return fmt.Errorf("API returned status %d: %s (Error: %s)", resp.StatusCode, errorResp.Message, errorResp.Error)
	}

	var successResp OrderResponse
	if err := json.Unmarshal(bodyBytes, &successResp); err != nil {
		return fmt.Errorf("API returned status %d, but failed to parse success response: %w", resp.StatusCode, err)
	}

	fmt.Printf("API Success Response (OrderID: %s, TotalPrice: %.2f, Status: %s)\n",
		successResp.Data.OrderID, successResp.Data.TotalPrice, successResp.Data.Status)
	return nil
}