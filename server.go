package main

import (
	"net/http"
	"fmt"
	"io"
	"strings"
	"os"
)

const (
	draftFolder = "/home/zmj/web/draft"
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
	draft.Comments = strings.Split(r.FormValue("comments"), "\n")

	//create the directory
	mkdirErr := makeFolder(event)
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
	logWriteErr := writeToFile(event+".log", log)
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
		writeToFile(draft.Image, img)
	}
	draft.HasDeck = len(draft.Image)>0 || len(draft.Comments)>1

	//write the html
	f, createErr := os.Create(event+".html")
	if createErr != nil {
		err = createErr
		return
	}
	defer f.Close()
	makePage(draft, f)
}

func makeFolder(event string) (err error) {
	err = os.Chdir(draftFolder)
	if err != nil {
		return
	}

	err = os.Mkdir(event, 0755 | os.ModeDir)
	if err != nil {
		return
	}

	err = os.Chdir(event)
	return
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
