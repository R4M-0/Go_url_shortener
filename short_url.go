package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

var urlMap = make(map[string]string)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func generateKey(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// Home page: form to shorten URL
func shortenHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		redirectHandler(w, r)
		return
	}

	if r.Method != http.MethodPost {
		fmt.Fprintf(w, `
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>Go URL Shortener</title>
			<style>
				body {
					font-family: Arial, sans-serif;
					background: #f0f2f5;
					display: flex;
					justify-content: center;
					align-items: center;
					height: 100vh;
					margin: 0;
				}
				.container {
					background: white;
					padding: 40px;
					border-radius: 10px;
					box-shadow: 0 4px 15px rgba(0,0,0,0.1);
					text-align: center;
					width: 350px;
				}
				input[type="text"] {
					width: 80%;
					padding: 10px;
					margin: 10px 0;
					border: 1px solid #ccc;
					border-radius: 5px;
				}
				input[type="submit"] {
					padding: 10px 20px;
					border: none;
					background: #4CAF50;
					color: white;
					border-radius: 5px;
					cursor: pointer;
					font-size: 16px;
				}
				input[type="submit"]:hover {
					background: #45a049;
				}
				h2 {
					color: #333;
				}
			</style>
		</head>
		<body>
			<div class="container">
				<h2>Go URL Shortener</h2>
				<form method="POST" action="/">
					<input type="text" name="url" placeholder="Enter your long URL" required><br>
					<input type="submit" value="Shorten URL">
				</form>
			</div>
		</body>
		</html>
		`)
		return
	}

	longURL := r.FormValue("url")
	key := generateKey(6)
	urlMap[key] = longURL
	shortURL := fmt.Sprintf("http://localhost:8000/%s", key)

	fmt.Fprintf(w, `
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>URL Shortened</title>
		<style>
			body { font-family: Arial, sans-serif; background: #f0f2f5; display:flex; justify-content:center; align-items:center; height:100vh; margin:0; }
			.container { background:white; padding:40px; border-radius:10px; box-shadow:0 4px 15px rgba(0,0,0,0.1); text-align:center; }
			a { color:#4CAF50; font-weight:bold; text-decoration:none; }
			a:hover { text-decoration:underline; }
			input[type="submit"] { margin-top:20px; padding:10px 20px; border:none; background:#4CAF50; color:white; border-radius:5px; cursor:pointer; font-size:16px; }
			input[type="submit"]:hover { background:#45a049; }
		</style>
	</head>
	<body>
		<div class="container">
			<h2>URL Shortened Successfully!</h2>
			<p>Here is your short URL:</p>
			<a href="%s">%s</a>
			<form action="/" method="get">
				<input type="submit" value="Shorten Another URL">
			</form>
		</div>
	</body>
	</html>
	`, shortURL, shortURL)
}

// Redirect handler: redirects short URL to original
func redirectHandler(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Path[1:] // remove leading "/"
	if longURL, ok := urlMap[key]; ok {
		http.Redirect(w, r, longURL, http.StatusFound)
		return
	}
	http.NotFound(w, r)
}

func main() {
	http.HandleFunc("/", shortenHandler)
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "favicon.ico")
	}) // ignore favicon

	fmt.Println("Server running at http://localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}

