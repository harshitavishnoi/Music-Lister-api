package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

// creating the user datastructures (User, Playlist, Song)
type User struct {
	Id         string `json:"id"`
	SecretCode string `json:"secret_code"`
	Name       string `json:"username"`
	Email      string `json:"email"`
	Playlists  []Playlist
}

type Playlist struct {
	Id    string `json:"playlistid"`
	Name  string `json:"playlistname"`
	Songs []Song
}

type Song struct {
	Id       string `json:"songid"`
	Name     string `json:"songname"`
	Composer string `json:"composer"`
	MusicURL string `json:"songurl"`
}

// creating users of user type
var (
	users          []User
	nextUserID     = 1
	nextPlaylistID = 1
	nextSongID     = 1
	mu             sync.Mutex
	getusers       = map[string]*User{}
)

// validatesecretcode function
func isValidSecretCode(inputSecretCode string) bool {
	// Iterate through the stored users and compare secret codes
	for _, user := range users {
		if inputSecretCode == user.SecretCode {
			return true // Secret code is valid
		}
	}
	return false // Secret code is not valid
}

// findUserBySecretCode retrieves a user based on their secret code
func findUserBySecretCode(inputSecretCode string) *User {
	// Iterate through the stored users
	for _, user := range users {
		if inputSecretCode == user.SecretCode {
			return &user // Return a pointer to the found user
		}
	}
	return nil // User with the given secret code not found
}

// creating a loginuser function
func loginuser(w http.ResponseWriter, r *http.Request) {
	// parsing the json request
	var requestData struct {
		SecretCode string `json:"secret_code"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid JSON input", http.StatusBadRequest)
		return
	}

	// Validate the secret code
	if !isValidSecretCode(requestData.SecretCode) {
		http.Error(w, "Invalid secret code", http.StatusUnauthorized)
		return
	}

	// Find the user based on the secret code
	user := findUserBySecretCode(requestData.SecretCode)
	if user == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Returning the user details as JSON response
	jsonResponse, err := json.Marshal(user)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

// generateSecretCode generates a random secret code
func generateSecretCode() string {
	// Specify the desired length of the secret code
	secretCodeLength := 6 // You can adjust this as needed

	// Create a byte slice to store random bytes
	randomBytes := make([]byte, secretCodeLength)

	// Read random bytes from the crypto/rand package
	if _, err := rand.Read(randomBytes); err != nil {
		// Handle error (e.g., by using a default secret code or returning an error)
		return "Default" // Replace with your error handling logic
	}

	// Encode the random bytes as a base64 string
	secretCode := base64.RawURLEncoding.EncodeToString(randomBytes)

	// Return the generated secret code
	return secretCode
}

// writing a function to create a new user
func registernewuser(w http.ResponseWriter, r *http.Request) {
	// Parse the request JSON
	var newUser User
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		http.Error(w, "Invalid JSON input", http.StatusBadRequest)
		return
	}
	// Validate the input data
	if newUser.Name == "" || newUser.Email == "" {
		http.Error(w, "Name and email are required fields", http.StatusBadRequest)
		return
	}

	// Generate a unique ID for the user
	newUser.Id = strconv.Itoa(nextUserID)
	nextUserID++

	// Generate a unique secret code for the user
	secretCode := generateSecretCode()
	newUser.SecretCode = secretCode

	// Store the new user in your data storage
	users = append(users, newUser)

	// Returning the newly created user as JSON response
	jsonResponse, err := json.Marshal(newUser)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonResponse)
}

func viewProfile(w http.ResponseWriter, r *http.Request) {
	// Parse the request JSON
	var newUser User
	var requestData struct {
		SecretCode string `json:"secret_code"`
	}
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		http.Error(w, "Invalid JSON input", http.StatusBadRequest)
		return
	}

	// Validating the secret code and find the user
	if !isValidSecretCode(requestData.SecretCode) {
		http.Error(w, "Invalid secret code", http.StatusUnauthorized)
		return
	}

	// Find the user based on the secret code
	user := findUserBySecretCode(requestData.SecretCode)
	if user == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	jsonResponse, err := json.Marshal(user.Playlists)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func createPlaylist(w http.ResponseWriter, r *http.Request) {
	// Extracting the secret code from the query parameters
	queryValues := r.URL.Query()
	secretCode := queryValues.Get("secret_code")

	// Validate the secret code
	if secretCode == "" {
		http.Error(w, "Secret code is required", http.StatusBadRequest)
		return
	}

	mu.Lock()
	user, exists := getusers[secretCode]
	mu.Unlock()

	if !exists {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Parse the request JSON
	var newPlaylist Playlist
	err := json.NewDecoder(r.Body).Decode(&newPlaylist)
	if err != nil {
		// Handle JSON parsing error
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validating and creating the playlist
	newPlaylist.Id = strconv.Itoa(nextPlaylistID)
	nextPlaylistID++

	// Adding playlist to the user's profile
	mu.Lock()
	user.Playlists = append(user.Playlists, newPlaylist)
	mu.Unlock()

	// Return the newly created playlist as JSON
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newPlaylist)
}

func getAllSongsOfPlaylist(w http.ResponseWriter, r *http.Request) {
	// Extracting the secret code from the query parameters
	queryValues := r.URL.Query()
	secretCode := queryValues.Get("secret_code")

	// Validate the secret code
	if secretCode == "" {
		http.Error(w, "Secret code is required", http.StatusBadRequest)
		return
	}

	mu.Lock()
	user, exists := getusers[secretCode]
	mu.Unlock()

	if !exists {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Parse the request URL path to get the playlistID
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.Error(w, "Invalid URL path", http.StatusBadRequest)
		return
	}

	playlistID := parts[2]

	var playlist Playlist
	for _, pl := range user.Playlists {
		if pl.Id == playlistID {
			playlist = pl
			break
		}
	}

	// Return the playlist's songs as JSON response
	jsonResponse, err := json.Marshal(playlist.Songs)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

// adding songs to the playlist
func addSongToPlaylist(w http.ResponseWriter, r *http.Request) {
	// Parse the request JSON
	var requestData struct {
		PlaylistID string `json:"playlist_id"`
		Song       Song   `json:"song"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid JSON input", http.StatusBadRequest)
		return
	}

	// Validate the input data
	if requestData.PlaylistID == "" || requestData.Song.Name == "" || requestData.Song.Composer == "" || requestData.Song.MusicURL == "" {
		http.Error(w, "Invalid input data", http.StatusBadRequest)
		return
	}

	// Extracting the secret code from the query parameters
	queryValues := r.URL.Query()
	secretCode := queryValues.Get("secret_code")

	// Validate the secret code
	if secretCode == "" {
		http.Error(w, "Secret code is required", http.StatusBadRequest)
		return
	}

	mu.Lock()
	user, exists := getusers[secretCode]
	mu.Unlock()

	if !exists {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Find the playlist based on the playlist ID
	var targetPlaylist *Playlist
	for i := range user.Playlists {
		if user.Playlists[i].Id == requestData.PlaylistID {
			targetPlaylist = &user.Playlists[i]
			break
		}
	}

	if targetPlaylist == nil {
		http.Error(w, "Playlist not found", http.StatusNotFound)
		return
	}

	// Adding the new song to the playlist
	requestData.Song.Id = strconv.Itoa(nextSongID)
	nextSongID++

	mu.Lock()
	targetPlaylist.Songs = append(targetPlaylist.Songs, requestData.Song)
	mu.Unlock()

	// Return the updated playlist as JSON response
	jsonResponse, err := json.Marshal(targetPlaylist)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

// deleting songs from the playlist
func deleteSongFromPlaylist(w http.ResponseWriter, r *http.Request) {

	secretCode := r.URL.Query().Get("secret_code")
	mu.Lock()
	user, exists := getusers[secretCode]
	mu.Unlock()

	if !exists {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	playlistIDStr := r.URL.Query().Get("playlist_id")

	playlistID, err := strconv.Atoi(playlistIDStr)

	if err != nil {
		http.Error(w, "Invalid playlist ID", http.StatusBadRequest)
		return
	}

	songIDStr := r.URL.Query().Get("song_id")
	songID, err := strconv.Atoi(songIDStr)
	if err != nil {
		http.Error(w, "Invalid song ID", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	for i := range user.Playlists {
		if user.Playlists[i].Id == strconv.Itoa(playlistID) {
			for j := range user.Playlists[i].Songs {
				if user.Playlists[i].Songs[j].Id == strconv.Itoa(songID) {
					// Delete the song from the playlist.
					user.Playlists[i].Songs = append(user.Playlists[i].Songs[:j], user.Playlists[i].Songs[j+1:]...)

					w.WriteHeader(http.StatusNoContent)
					return
				}
			}
			http.Error(w, "Song not found in playlist", http.StatusNotFound)
			return
		}
	}

	http.Error(w, "Playlist not found", http.StatusNotFound)

}

// Delete a playlist.
func deletePlaylist(w http.ResponseWriter, r *http.Request) {
	secretCode := r.URL.Query().Get("secret_code")

	user, exists := getusers[secretCode]

	if !exists {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	playlistIDStr := r.URL.Query().Get("playlist_id")
	playlistID, err := strconv.Atoi(playlistIDStr)
	if err != nil {
		http.Error(w, "Invalid playlist ID", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	for i := range user.Playlists {
		if user.Playlists[i].Id == strconv.Itoa(playlistID) {
			// Delete the playlist.
			user.Playlists = append(user.Playlists[:i], user.Playlists[i+1:]...)

			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	http.Error(w, "Playlist not found", http.StatusNotFound)
}

// getSongDetail allows you to get all attributes of a song.
func getSongDetail(w http.ResponseWriter, r *http.Request) {
	secretCode := r.URL.Query().Get("secret_code")
	user, exists := getusers[secretCode]

	if !exists {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	playlistIDStr := r.URL.Query().Get("playlist_id")
	playlistID, err := strconv.Atoi(playlistIDStr)
	if err != nil {
		http.Error(w, "Invalid playlist ID", http.StatusBadRequest)
		return
	}

	songIDStr := r.URL.Query().Get("song_id")
	songID, err := strconv.Atoi(songIDStr)
	if err != nil {
		http.Error(w, "Invalid song ID", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	for _, playlist := range user.Playlists {
		if playlist.Id == strconv.Itoa(playlistID) {
			for _, song := range playlist.Songs {
				if song.Id == strconv.Itoa(songID) {
					w.Header().Set("Content-Type", "application/json")
					json.NewEncoder(w).Encode(song)
					return
				}
			}
			http.Error(w, "Song not found in playlist", http.StatusNotFound)
			return
		}
	}

	http.Error(w, "Playlist not found", http.StatusNotFound)
}

// defining the api routes
func main() {
	http.HandleFunc("/login", loginuser)
	http.HandleFunc("/register", registernewuser)
	http.HandleFunc("/viewProfile", viewProfile)
	http.HandleFunc("/getAllSongsOfPlaylist", getAllSongsOfPlaylist)
	http.HandleFunc("/createPlaylist", createPlaylist)
	http.HandleFunc("/addSongToPlaylist", addSongToPlaylist)
	http.HandleFunc("/deleteSongFromPlaylist", deleteSongFromPlaylist)
	http.HandleFunc("/deletePlaylist", deletePlaylist)
	// http.HandleFunc("/getSongDetail", getSongDetail)

	// Start the HTTP server
	fmt.Println("Server is listening on :3000")
	http.ListenAndServe(":3000", nil)
}
