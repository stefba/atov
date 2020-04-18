package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"log"
	"text/template"
	"os"
	"os/exec"
	"strings"
	"time"
)

func main() {
	http.HandleFunc("/", front)
	http.HandleFunc("/tmp/", tmp)
	http.HandleFunc("/send", back)
	http.ListenAndServe(":8515", nil)
}

func front(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./main.html")
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 500)
		return
	}
	err = t.Execute(w, "")
	if err != nil {
		log.Println(err)
	}
}

func tmp(w http.ResponseWriter, r *http.Request) {
	p := "." + r.URL.Path 
	log.Println(p)
	http.ServeFile(w, r, p)
}

func back(w http.ResponseWriter, r *http.Request) {
	path, err := saveFile(r)
	if err != nil {
		log.Println(err)
		return
	}
	out, err := convertAudio(path)
	if err != nil {
		fmt.Fprint(w, out)
	}

	err = os.Remove(path)
	if err != nil {
		log.Println(err)
		return
	}

	http.Redirect(w, r, mp4Ext(path), 307)
}

func saveFile(r *http.Request) (string, error) {
	file, handler, err := r.FormFile("file")
    if err != nil {
		return "", err
    }
    defer file.Close()

	t, err := time.Parse("2006-01-02 15.04.05.mp3", handler.Filename)
	if err != nil {
		return "", err
	}

    b, err := ioutil.ReadAll(file)
    if err != nil {
		return "", err
    }

	path := "./tmp/" + t.Format("200102_150405.mp3")

	return path, ioutil.WriteFile(path, b, 0755)
}

func convertAudio(path string) (string, error) {
	cmd := exec.Command("ffmpeg", "-i", "black.mp4", "-i", path, "-shortest", mp4Ext(path))
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func mp4Ext(path string) (string) {
	return strings.Replace(path, ".mp3", ".mp4", -1)
}
