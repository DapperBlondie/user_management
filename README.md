First of all I want to mention that all of my applications may be are not
appropriate for this environment ((github)) because, my apps are not reusable or a pkg
that you could be able to expand it or use in other projects but I am trying and hardworking
everyday to reach that point.
This repository just a way or sign for showing my abilities in a specific language.
***


# User Management System
Someone gave me a little task to show my abilities and this repo just for that purposes.
Just a way for single file uploading, also for multiple files.
I skipped other parts because, they are sooo boring.

## Single Uploading

```
if r.Method != "POST" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parsing and Setting the max request size and then write to r.Body
	r.Body = http.MaxBytesReader(w, r.Body, MAX_UPLOAD_SIZE)
	
	// We parse our Form here, It's Required !
	if err := r.ParseMultipartForm(MAX_UPLOAD_SIZE); err != nil {
		log.Println("Error in size of request !")
		return
	}

	// We get the file by its name in form at front-end
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
```

***

## Multiple File Uploading

```
if r.Method != "POST" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	
	r.Body = http.MaxBytesReader(w, r.Body, MAX_MULTIPLE_SIZE) // Use for limiting the r.Body
	
	// Use for parsing the form we sent it as like as Single uploading file handler
	if err := r.ParseMultipartForm(MAX_MULTIPLE_SIZE); err != nil { // Parse form have multipart/form-data into it
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// At this point we get all files that associated with name 'image_file'
	files := r.MultipartForm.File["images_file"] // Use for getting the all files in the request

	// Then we iterate over them
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
```
At the end you can fork or download the code and use it :)
I provide an simple HTML file for testing it
