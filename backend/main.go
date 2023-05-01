package main

import (
	"fmt"
	"log"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"
	"html/template"
"bufio"
"os"
)

type TemplateData struct {
    Filename string
}

func main() {
    http.HandleFunc("/", handler)
    // Add other handlers if needed

    // Serve static files from the 'static' directory
    fs := http.FileServer(http.Dir("./static"))
    http.Handle("/static/", http.StripPrefix("/static", fs))

    log.Fatal(http.ListenAndServe(":8080", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
    directoryPath := "./static/temp"
    filePrefix := "seantube_download"

    filename, err := getFileNameWithPrefix(directoryPath, filePrefix)

    if err != nil {
        // If no file is found with the given prefix, serve the regular index.html
        http.ServeFile(w, r, "./static/index.html")
    } else {
        // If a file is found, render the template with the file name
        tmpl, err := template.ParseFiles("./static/index_template.html")
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        data := TemplateData{Filename: filename}
        tmpl.Execute(w, data)
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

    err := updateEnvVariable("URL", url)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    cmd := exec.Command("python3", "your_script.py")
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    err = cmd.Run()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    fmt.Fprint(w, "Python script executed successfully")
}

func updateEnvVariable(variableName, variableValue string) error {
    inputFile, err := os.Open(".env")
    if err != nil {
        return err
    }
    defer inputFile.Close()

    outputFile, err := os.Create(".env_temp")
    if err != nil {
        return err
    }
    defer outputFile.Close()

    scanner := bufio.NewScanner(inputFile)
    updated := false

    for scanner.Scan() {
        line := scanner.Text()
        if strings.HasPrefix(line, variableName+"=") {
            _, err = outputFile.WriteString(variableName + "=" + variableValue + "\n")
            updated = true
        } else {
            _, err = outputFile.WriteString(line + "\n")
        }

        if err != nil {
            return err
        }
    }

    if !updated {
        _, err = outputFile.WriteString(variableName + "=" + variableValue + "\n")
        if err != nil {
            return err
        }
    }

    err = inputFile.Close()
    if err != nil {
        return err
    }

    err = outputFile.Close()
    if err != nil {
        return err
    }

    return os.Rename(".env_temp", ".env")
}





func getFileNameWithPrefix(directoryPath string, prefix string) (string, error) {
    files, err := ioutil.ReadDir(directoryPath)
    if err != nil {
        return "", err
    }

    for _, file := range files {
        if !file.IsDir() && strings.HasPrefix(file.Name(), prefix) {
            return file.Name(), nil
        }
    }

    return "", fmt.Errorf("no file found with prefix: %s", prefix)
}



