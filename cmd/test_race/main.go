package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
)

const baseURL = "http://localhost:8080/api/v1"

func login(email, password string) string {
	body, _ := json.Marshal(map[string]string{
		"email":    email,
		"password": password,
	})
	resp, err := http.Post(baseURL+"/auth/login", "application/json", bytes.NewBuffer(body))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var result struct {
		Data struct {
			Token string `json:"token"`
		} `json:"data"`
	}
	_ = json.NewDecoder(resp.Body).Decode(&result)
	return result.Data.Token
}

func main() {
	// Login menggunakan 2 akun USER dari file seed
	tokenAndi := login("andi@student.univ.ac.id", "password")
	tokenSiti := login("siti@student.univ.ac.id", "password")

	payload := []byte(`{
		"room_id": 1,
		"title": "Uji Simultan Ruangan",
		"purpose": "Pengujian transaksi FOR UPDATE untuk mencegah double booking",
		"booking_date": "2026-11-25",
		"start_time": "09:00",
		"end_time": "11:00"
	}`)

	var wg sync.WaitGroup
	startSignal := make(chan struct{})

	sendBooking := func(userName, token string) {
		defer wg.Done()
		req, _ := http.NewRequest("POST", baseURL+"/bookings", bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		<-startSignal // Tunggu aba-aba agar kedua request dilepas bersamaan
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Printf("[%s] Error: %v\n", userName, err)
			return
		}
		defer resp.Body.Close()
		respBody, _ := io.ReadAll(resp.Body)
		fmt.Printf("[%s] HTTP %d -> %s\n", userName, resp.StatusCode, string(respBody))
	}

	wg.Add(2)
	go sendBooking("Andi", tokenAndi)
	go sendBooking("Siti", tokenSiti)

	close(startSignal) // Lepas kedua request secara bersamaan
	wg.Wait()
}
