
package main

import (
	"fmt"
	"log"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"
"html/template"
"github.com/joho/godotenv"
"path/filepath"
"net/url"
"os"
"embed"
)
var envFile embed.FS
var staticTempContent embed.FS


func main() {

	http.HandleFunc("/", handler)
	http.HandleFunc("/update-env", updateEnvHandler)
	// Serve static files from the 'public' directory
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static", fs))
	
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, post-check=0, pre-check=0")
w.Header().Set("Pragma", "no-cache")
w.Header().Set("Expires", "0")


if r.URL.Path == "/download1" {
        tempDir := "static/temp"
        files, err := ioutil.ReadDir(tempDir)
        if err != nil {
            http.Error(w, "Error reading temp directory", http.StatusInternalServerError)
            return
        }

        // Check if there's any file in the temp folder
	var filename string
        for _, file := range files {
            if !file.IsDir() {
                filename = file.Name()
                break
            }
        }

        if filename != "" {
            // Serve the index_template page if a file is found in the temp folder
            tmpl, err := template.ParseFiles("static/index_template.html")
            if err != nil {
                http.Error(w, "Error parsing template", http.StatusInternalServerError)
                return
            }

            data := struct {
                Filename string
            }{
                Filename: filename,
            }

            tmpl.Execute(w, data)
        } else {
            // Serve the regular index.html page if no file is found in the temp folder
            http.ServeFile(w, r, "static/index.html")
        }
        return
    }

if r.URL.Path =="/download2" {
tempDir := "static/temp"
        files, err := ioutil.ReadDir(tempDir)
        if err != nil {
            http.Error(w, "Error reading temp directory", http.StatusInternalServerError)
            return
        }

var filename string
        for _, file := range files {
            if !file.IsDir() {
                filename = file.Name()
                break
            }
        }

        if filename != "" {
            // Serve the index_template page if a file is found in the temp folder
            tmpl, err := template.ParseFiles("static/index_template.html")
            if err != nil {
                http.Error(w, "Error parsing template", http.StatusInternalServerError)
                return
            }

            data := struct {
                Filename string
            }{
                Filename: filename,
            }

            tmpl.Execute(w, data)
        }
}
}


func updateEnvHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	url := r.FormValue("URL")
	if url == "" {
		http.Error(w, "URL is required", http.StatusBadRequest)
		return

	}

	// Read the .env file
	content, err := ioutil.ReadFile(".env")
	if err != nil {
		http.Error(w, "Error reading .env file", http.StatusInternalServerError)
		return
	}

	// Find and update the URL variable
	varName := "URL"
    	lines := strings.Split(string(content), "\n")
    	updated := false
    	for i, line := range lines {
        	if strings.HasPrefix(line, varName+"=") {
            	lines[i] = fmt.Sprintf("%s=\"%s\"", varName, url) // Add double quotes around the URL value
            	updated = true
            	break
        	}
    }

    // If the URL variable was not found, append it to the file
    if !updated {
        lines = append(lines, fmt.Sprintf("%s=\"%s\"", varName, url)) // Add double quotes around the URL value
    }

	// Write the updated content to the .env file
	err = ioutil.WriteFile(".env", []byte(strings.Join(lines, "\n")), 0644)
	if err != nil {
		http.Error(w, "Error updating .env file", http.StatusInternalServerError)
		return
	}
w.WriteHeader(http.StatusOK)
err = downloadMedia()

	if err != nil {

		http.Error(w, fmt.Sprintf("Error downloading media: %v", err), http.StatusInternalServerError)
		return
	}


}



func downloadMedia() error {

	_ = godotenv.Load()



	inputStr := os.Getenv("URL")

	audio := os.Getenv("AUDIO")

	audioFormat := os.Getenv("AUDIO_FORMAT")

	videoFormat := os.Getenv("VIDEO_FORMAT")



	videoCmdURL := "yt-dlp --recode-video " + videoFormat + " --output seantube_download " + inputStr

	videoCmdNoURL := "yt-dlp --recode-video " + videoFormat + " --output seantube_download ytsearch:" + `"` + inputStr + `"`

	audioCmdURL := "yt-dlp -x --audio-format " + audioFormat + " --output seantube_download " + inputStr

	audioCmdNoURL := "yt-dlp -x --audio-format " + audioFormat + " --output seantube_download ytsearch:" + `"` + inputStr + `"`

	fmt.Println(videoCmdNoURL)



	isURL := isValidURL(inputStr)



	if audio == "True" {

		if !isURL {

			fmt.Println("Mode: Audio without url")

			exec.Command("sh", "-c", audioCmdNoURL).Run()

		} else {

			fmt.Println("Mode: Audio with url")

			exec.Command("sh", "-c", audioCmdURL).Run()

		}

	} else {

		if !isURL {

			fmt.Println("Mode: Video without url")

			exec.Command("sh", "-c", videoCmdNoURL).Run()

		} else {

			fmt.Println("Mode: Video with url")

			exec.Command("sh", "-c", videoCmdURL).Run()

		}

	}



	matches, _ := filepath.Glob("seantube_*")

	for _, match := range matches {

		os.Rename(match, filepath.Join("static/temp/", filepath.Base(match)))

	}
	err := godotenv.Load()

	if err != nil {

		return fmt.Errorf("error loading .env file: %v", err)

	}
serve_files()
	return nil

}

func serve_files() {
folderToServe := "./static/temp" // Replace with the path to your folder



	// Create a file server to serve files from the specified folder

	fileServer := http.FileServer(http.Dir(folderToServe))



	// Register the file server with the desired route

	http.Handle("/dll", fileServer)



	// Set the listening address and port for the server

	address := "localhost"

	port := "8080"



	fmt.Printf("Serving files from folder: %s\n", folderToServe)

	fmt.Printf("Server listening on http://%s:%s\n", address, port)



	// Start the server

	err := http.ListenAndServe(address+":"+port, nil)

	if err != nil {

		fmt.Println("Error starting server:", err)

	}

}

func isValidURL(input string) bool {

	u, err := url.Parse(input)

	if err != nil {

		return false

	}

	return u.Scheme != "" && u.Host != ""

}

