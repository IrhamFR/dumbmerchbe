package middleware

import (
	"context"
	dto "dumbmerch/dto/result"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

func UploadFile(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Upload file
		// FormFile returns the first file for the given key `myFile`
		// it also returns the FileHeader so we can get the Filename,
		// the Header and the size of the file
		file, _, err := r.FormFile("image")

		if err != nil && r.Method == "PATCH" {
			ctx := context.WithValue(r.Context(), "dataFile", "false")
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
			json.NewEncoder(w).Encode(response)
			return
		}
		// if err != nil {
		// 	fmt.Println(err)
		// 	json.NewEncoder(w).Encode("Error Retrieving the File")
		// 	return
		// }

		defer file.Close()

		// setup file type filtering
		buff := make([]byte, 512)
		_, err = file.Read(buff)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
			json.NewEncoder(w).Encode(response)
			return
		}

		filetype := http.DetectContentType(buff)
		if filetype != "image/jpg" {
			w.WriteHeader(http.StatusBadRequest)
			response := dto.ErrorResult{Code: http.StatusBadRequest, Message: "The provided file format is not allowed. Please upload a JPG or PNG image"}
			json.NewEncoder(w).Encode(response)
			return
		}

		_, err = file.Seek(0, io.SeekStart)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
			json.NewEncoder(w).Encode(response)
			return
		}

		// setup max-upload
		const MAX_UPLOAD_SIZE = 1 << 20 // 100 kb
		r.Body = http.MaxBytesReader(w, r.Body, MAX_UPLOAD_SIZE)
		fmt.Println(r.Body)
		r.ParseMultipartForm(MAX_UPLOAD_SIZE)
		if r.ContentLength > MAX_UPLOAD_SIZE {
			w.WriteHeader(http.StatusBadRequest)
			response := Result{Code: http.StatusBadRequest, Message: "Max size is 100 kb"}
			json.NewEncoder(w).Encode(response)
			return
		}

		// const MAX_UPLOAD_SIZE = 1 << 1 // 100 kb

		// r.Body = http.MaxBytesReader(w, r.Body, MAX_UPLOAD_SIZE)
		// fmt.Println(r.Body)

		// if err := r.ParseMultipartForm(MAX_UPLOAD_SIZE); err != nil {
		// 	http.Error(w, "Max size in 100 kb", http.StatusBadRequest)
		// 	return
		// }

		// Create a temporary file within our temp-images directory that follows
		// a particular naming pattern
		tempFile, err := ioutil.TempFile("uploads", "image-*.jpg")
		if err != nil {
			fmt.Println(err)
			fmt.Println("path upload error")
			json.NewEncoder(w).Encode(err)
			return
		}
		defer tempFile.Close()

		// read all of the contents of our uploaded file into a
		// byte array
		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			fmt.Println(err)
		}

		// write this byte array to our temporary file
		tempFile.Write(fileBytes)

		data := tempFile.Name()
		// filename := data[8:] // split uploads/

		// add filename to ctx
		ctx := context.WithValue(r.Context(), "dataFile", data)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
