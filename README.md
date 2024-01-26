# Music-Lister-api

The Music Lister API is a Go-based backend service designed for managing user profiles, playlists, and song details. It provides essential functionalities for user registration, login, playlist creation, song addition, and deletion. The API leverages a secure user authentication mechanism with randomly generated secret codes.

## Features

- **User Authentication:** Users are assigned unique secret codes for secure authentication.
- **User Registration:** Register users with a unique ID, secret code, name, and email.
- **Playlist Management:** Create, view, and delete playlists associated with a user.
- **Song Operations:** Add, retrieve, and delete songs within playlists.

## API Routes

- `/login`: Authenticate users with their secret code.
- `/register`: Register new users with unique secret codes.
- `/viewProfile`: View user profiles and associated playlists.
- `/createPlaylist`: Create a new playlist for a user.
- `/getAllSongsOfPlaylist`: Retrieve all songs in a specific playlist.
- `/addSongToPlaylist`: Add a song to a playlist.
- `/deleteSongFromPlaylist`: Delete a song from a playlist.
- `/deletePlaylist`: Delete a playlist.

## Usage

1. **User Registration:**
   - **Endpoint:** `/register`
   - **Method:** `POST`
   - **Request Body Example:**
     ```json
     {
       "name": "John Doe",
       "email": "john.doe@example.com"
     }
     ```
2. **User Login:**
   - **Endpoint:** `/login`
   - **Method:** `POST`
   - **Request Body Example:**
     ```json
     {
       "secret_code": "randomly_generated_secret_code"
     }
     ```

3. **Create Playlist:**
   - **Endpoint:** `/createPlaylist?secret_code=<user_secret_code>`
   - **Method:** `POST`
   - **Request Body Example:**
     ```json
     {
       "name": "My Playlist"
     }
     ```

4. **Add Song to Playlist:**
   - **Endpoint:** `/addSongToPlaylist?secret_code=<user_secret_code>`
   - **Method:** `POST`
   - **Request Body Example:**
     ```json
     {
       "playlist_id": "playlist_id",
       "song": {
         "name": "Song Name",
         "composer": "Composer Name",
         "songurl": "https://example.com/song.mp3"
       }
     }
     ```

5. **Retrieve Songs of Playlist:**
   - **Endpoint:** `/getAllSongsOfPlaylist?secret_code=<user_secret_code>&playlist_id=<playlist_id>`
   - **Method:** `GET`

6. **Delete Song from Playlist:**
   - **Endpoint:** `/deleteSongFromPlaylist?secret_code=<user_secret_code>&playlist_id=<playlist_id>&song_id=<song_id>`
   - **Method:** `DELETE`

7. **Delete Playlist:**
   - **Endpoint:** `/deletePlaylist?secret_code=<user_secret_code>&playlist_id=<playlist_id>`
   - **Method:** `DELETE`

## Getting Started

1. Clone the repository: `git clone https://github.com/harshitavishnoi/Music-Lister-api.git`
2. Navigate to the project directory: `cd Music-Lister-api`
3. Run the server: `go run main.go`
4. The server will be accessible at `http://localhost:3000`.

Feel free to explore the API routes and integrate this backend service into your music-related applications!
