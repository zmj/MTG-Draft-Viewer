package main

import (
	"net/http"
	"fmt"
	"io"
	"strings"
	"os"
)

func handler(w http.ResponseWriter, r *http.Request) {
	var err error
	var event string
	defer func() {
		if err != nil {
			http.Error(w, err.Error(), 500)
			fmt.Println("error", err)
			if event != "" {
				cleanup(event)
			}
		} else {
			http.Redirect(w, r, "/draft/"+event, 302)
			fmt.Println("draft", event)
		}
	}()

	//get the draft file	
	log, _, loadLogErr := r.FormFile("log")
	if loadLogErr != nil {
		err = loadLogErr
		return
	}
	defer log.Close()

	//parse the draft
	draft, parseErr := NewDraft(log)
	if parseErr != nil {
		err = parseErr
		return
	}
	event = draft.Event

	//add the comments
	comments := strings.TrimSpace(r.FormValue("comments"))
	draft.Comments = strings.Split(comments, "\n")

	//create the directory
	mkdirErr := os.Mkdir(event, 0777)
	if mkdirErr != nil { //&& !os.IsExist(mkdirErr) {
		err = mkdirErr
		return
	}

	//write the log file
	_, seekErr := log.Seek(0, 0)
	if seekErr != nil {
		err = seekErr
		return
	}
	logWriteErr := writeToFile(event+"/"+event+".log", log)
	if logWriteErr != nil {
		err = logWriteErr
		return
	}

	//write the image
	img, imgHeader, loadImgErr := r.FormFile("image")
	if loadImgErr != nil {
		if loadImgErr != http.ErrMissingFile {
			err = loadImgErr
			return
		}
	} else {
		defer img.Close()
		draft.Image = imgHeader.Filename
		writeToFile(event+"/"+draft.Image, img)
	}

	//write the html
	f, createErr := os.Create(event+"/"+event+".html")
	if createErr != nil {
		err = createErr
		return
	}
	defer f.Close()
	makePage(draft, f)
}

func writeToFile(path string, src io.Reader) error {
	f, createErr := os.Create(path)
	if createErr != nil {
		return createErr
	}
	defer f.Close()
	_, copyErr := io.Copy(f, src)
	return copyErr
}

func cleanup(event string) {
	//delete the directory and any files in it
}

func main() {
	//start a web server that listens to posts
	http.HandleFunc("/draft/upload", handler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err)
	}
}
