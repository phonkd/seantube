package main

import (
	"fmt"
	"log"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"
)





func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/update-env", updateEnvHandler)

	// Serve static files from the 'public' directory
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static", fs))

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		http.ServeFile(w, r, "static/index.html")
		return
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
	fmt.Fprint(w, "Updated .env file successfully")
	// Run the Python script
	cmd := exec.Command("python", "main.py")
	output, err := cmd.CombinedOutput()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error running Python script: %v\nOutput: %s", err, output), http.StatusInternalServerError)
		return
	}

	// Send the script output to the client
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Updated .env file successfully\n\nPython script output:\n%s", output)
}






