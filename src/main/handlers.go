package main

import (
	"encoding/json"
	"fmt"
	"github.com/DapperBlondie/user_management/src/models"
	"github.com/DapperBlondie/user_management/src/repo"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const MAX_UPLOAD_SIZE = 2048*1024
const MAX_MULTIPLE_SIZE = 5120*1024

type HandlerRepo struct {
	Config *models.AppConfig
	D	   *repo.PostgresDBRepo
}

type AppStatus struct {
	Status string
	Version string
}

type UserPayload struct {
	FirstName    string	`json:"first_name"`
	LastName     string	`json:"last_name"`
	DateOfBirth  string	`json:"date_of_birth"`
}

var handlerRepo *HandlerRepo

func NewHandlerRepo(c *models.AppConfig, d *repo.PostgresDBRepo)  {
	handlerRepo = &HandlerRepo{
		Config: c,
		D:      d,
	}
}

func (h *HandlerRepo) ImageUploadingHandler(w http.ResponseWriter, r *http.Request)  {
	if r.Method != "POST" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parsing and Setting the max request size
	r.Body = http.MaxBytesReader(w, r.Body, MAX_UPLOAD_SIZE)
	if err := r.ParseMultipartForm(MAX_UPLOAD_SIZE); err != nil {
		log.Println("Error in size of request !")
		return
	}

	// File uploading operation will be here
	file, fileHeader, err := r.FormFile("image_file")
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer func(file multipart.File) {
		err = file.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.Println(err.Error())
			return
		}
	}(file)

	dst, err := os.Create(fmt.Sprintf("./src/static/images/%s%d%s",
		r.Form.Get("first_name"),time.Now().UnixNano(), filepath.Ext(fileHeader.Filename)))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer func(dst *os.File) {
		err = dst.Close()
		if err != nil {
			log.Println(err.Error())
			return
		}
	}(dst)

	_, err = io.Copy(dst, file)
	if err != nil {
		log.Println(err.Error())
		return
	}

	log.Println("File Uploaded")
	http.Redirect(w, r, "/status", http.StatusSeeOther)
	return
}

func (h *HandlerRepo) MultipleImageUploadingHandler(w http.ResponseWriter, r *http.Request)  {
	if r.Method != "POST" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, MAX_MULTIPLE_SIZE) // Use for limiting the r.Body
	if err := r.ParseMultipartForm(MAX_MULTIPLE_SIZE); err != nil { // Parse form have multipart/form-data into it
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	files := r.MultipartForm.File["images_file"] // Use for getting the all files in the request

	for _, fileHeader := range files{
		if fileHeader.Size == MAX_UPLOAD_SIZE {
			log.Println("Uploaded file is too big")
			http.Error(w, "Uploaded file is too big", http.StatusBadRequest)
			continue
		}

		uFile, err := fileHeader.Open()
		if err != nil {
			log.Println(err.Error())
			continue
		}
		defer func(uFile multipart.File) {
			err = uFile.Close()
			if err != nil {
				log.Println(err.Error())
				return
			}
		}(uFile)

		uF, err := os.Create(fmt.Sprintf("./src/static/images/%s%d%s",
			r.Form.Get("first_name"),time.Now().UnixNano(), filepath.Ext(fileHeader.Filename)))
		defer func(uF *os.File) {
			err = uF.Close()
			if err != nil {
				log.Println(err.Error())
				return
			}
		}(uF)

		_, err = io.Copy(uF, uFile)
		if err != nil {
			log.Println(err.Error())
			return
		}
	}
	log.Println("Files Uploaded !")
	http.Redirect(w, r, "/status", http.StatusSeeOther)
	return
}

func (h *HandlerRepo) UserInfoHandler(w http.ResponseWriter, r *http.Request)  {
	if r.Method == "GET" {
		http.Redirect(w, r, "/status", http.StatusBadRequest)
		return
	}
	h.Config.SCSManager.Put(r.Context(), "url", r.UserAgent())
	var userInfo = &UserPayload{}
	err := json.NewDecoder(r.Body).Decode(&userInfo)
	if err != nil {
		log.Println("Error in decoding the user payload")
	}

	fmt.Println(userInfo)
	return
}

func (h *HandlerRepo) Status(w http.ResponseWriter, r *http.Request)  {
	currentStatus := &AppStatus{
		Status:      "Available",
		Version:     "1.0.0",
	}

	out, err := json.MarshalIndent(currentStatus, "", "\t")
	if err != nil {
		log.Fatalln("Error in marshaling currentStatus : " + err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(out)
	if err != nil {
		w.Write([]byte("Error occurred : " + err.Error() + "\n"))
		log.Println(err.Error())
		return
	}
}
