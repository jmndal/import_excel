package api

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/jmandal/import_excel/models"
	"github.com/uadmin/uadmin"
	"github.com/xuri/excelize/v2"
)

type Response struct {
	Status      string `json:"status"`
	Data        string `json:"translated_text"`
	Description string `json:"translated_desc"`
}

const (
	URL         = "http://0.0.0.0:8833"
	contentType = "application/json"
)

var httpClient = &http.Client{
	Transport: &http.Transport{
		MaxIdleConns:    100,
		IdleConnTimeout: 90 * time.Second,
	},
}

var activeLang = "[\"en\",\"ja\",\"ko\",\"es\"]"

func UploadFile(w http.ResponseWriter, r *http.Request) {
	// Parse the form data, including the uploaded file(s)
	err := r.ParseMultipartForm(50 << 20) // 10 MB limit
	if err != nil {
		uadmin.Trail(uadmin.ERROR, "Unable to parse form. %s", err)
		uadmin.ReturnJSON(w, r, map[string]interface{}{
			"status": "error",
			"error":  err,
		})
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get the file from the request
	file, handler, err := r.FormFile("file")
	if err != nil {
		uadmin.Trail(uadmin.ERROR, "Error retrieving file from form. %s", err)
		uadmin.ReturnJSON(w, r, map[string]interface{}{
			"status": "error",
			"error":  err,
		})
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Create a new file on the server to save the uploaded file
	f, err := os.Create("media/uploads/" + handler.Filename)
	if err != nil {
		uadmin.Trail(uadmin.ERROR, "Error creating file on server. %s", err)
		uadmin.ReturnJSON(w, r, map[string]interface{}{
			"status": "error",
			"error":  err,
		})
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()

	// Copy the file content from the request to the new file
	_, err = io.Copy(f, file)
	if err != nil {
		uadmin.Trail(uadmin.ERROR, "Error copying file content. %s", err)
		uadmin.ReturnJSON(w, r, map[string]interface{}{
			"status": "error",
			"error":  err,
		})
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	now := time.Now()
	t := now.Format("2006-01-02T15:04")

	os.Mkdir("media/uploads/"+t, os.ModePerm)
	xlsx, err := unzip("media/uploads/"+handler.Filename, "media/uploads/"+t+"/")
	if err != nil {
		uadmin.Trail(uadmin.ERROR, "Unable to extract file")
		uadmin.ReturnJSON(w, r, map[string]interface{}{
			"status": "error",
			"error":  err,
		})
		w.WriteHeader(http.StatusInternalServerError)
	}

	// Respond with a success message
	uadmin.Trail(uadmin.INFO, "File uploaded successfully!")

	// Remove the zip file
	os.Remove("media/uploads/" + handler.Filename)

	// Open the excel file
	excel, err := excelize.OpenFile(xlsx)
	if err != nil {
		uadmin.Trail(uadmin.ERROR, "Unable to open excel file")
		uadmin.ReturnJSON(w, r, map[string]interface{}{
			"status": "error",
			"error":  err,
		})
		w.WriteHeader(http.StatusInternalServerError)
	}

	menuList := excel.GetSheetList()

	for _, menuSheet := range menuList {
		rows, err := excel.GetRows(menuSheet)
		if err != nil {
			uadmin.Trail(uadmin.ERROR, "Unable to get rows. %s", err)
			uadmin.ReturnJSON(w, r, map[string]interface{}{
				"status": "error",
				"error":  err,
			})
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		for i, row := range rows {

			if i != 0 {
				menu := models.Menus{}
				counter := 0
				m := ""

				for _, colCell := range row {
					colCell = strings.ReplaceAll(colCell, "\n", " ")
					if counter == 0 {
						m += colCell + ";"
						counter++
					} else {
						if counter%11 != 0 {
							m += colCell + ";"
							counter++
						}
						if counter%11 == 0 {
							m += colCell + ";"
							counter = 0

							if uadmin.Count(&menu, "name = ? AND description = ? AND category = ?", strings.Split(m, ";")[0], strings.Split(m, ";")[1], strings.Split(m, ";")[2]) == 0 {

								// Prepare the request body as a JSON payload
								payload, err := json.Marshal(map[string]interface{}{
									"activeLang":  activeLang,
									"text":        strings.Split(m, ";")[0],
									"description": strings.Split(m, ";")[1],
								})
								if err != nil {
									uadmin.Trail(uadmin.ERROR, "Error marshaling JSON:", err)
									uadmin.ReturnJSON(w, r, map[string]interface{}{
										"status": "error",
										"error":  err,
									})
									return
								}

								// Make the POST request with JSON payload
								resp, err := httpClient.Post(URL, contentType, bytes.NewBuffer(payload))
								if err != nil {
									uadmin.Trail(uadmin.ERROR, "Error making POST request:", err)
									uadmin.ReturnJSON(w, r, map[string]interface{}{
										"status": "error",
										"error":  err,
									})
									return
								}
								defer resp.Body.Close()

								// Check the response status code
								if resp.StatusCode != http.StatusOK {
									uadmin.Trail(uadmin.ERROR, "Error: Non-OK status code received: %s", resp.Status)
									uadmin.ReturnJSON(w, r, map[string]interface{}{
										"status": "error",
										"error":  err,
									})
									return
								}

								// Decode the response body into the Response struct
								var response Response
								err = json.NewDecoder(resp.Body).Decode(&response)
								if err != nil {
									uadmin.Trail(uadmin.ERROR, "Error decoding response:", err)
									uadmin.ReturnJSON(w, r, map[string]interface{}{
										"status": "error",
										"error":  err,
									})
									return
								}

								menu.Name = response.Data
								menu.Description = response.Description
								menu.Category = strings.Split(m, ";")[2]

								price, err := strconv.ParseFloat(strings.Split(m, ";")[3], 64)
								if err != nil {
									uadmin.Trail(uadmin.ERROR, "Unable to parse price")
									uadmin.ReturnJSON(w, r, map[string]interface{}{
										"status": "error",
										"error":  err,
									})
									w.WriteHeader(http.StatusInternalServerError)
									continue
								}

								menu.Price = price
								menu.Abbreviation = strings.Split(m, ";")[4]
								menu.From = strings.Split(m, ";")[5]
								menu.To = strings.Split(m, ";")[6]
								menu.Image = "/media/uploads/" + t + "/" + strings.Split(m, ";")[7]
								menu.Discount = strings.Split(m, ";")[8]
								menu.SystemCode = strings.Split(m, ";")[9]
								menu.KitchenLocation = strings.Split(m, ";")[10]

								uadmin.Save(&menu)
							}
						}
					}
				}
			}
		}
	}
}

func unzip(source, destination string) (string, error) {
	xlsxFilePath := ""
	r, err := zip.OpenReader(source)
	if err != nil {
		return "", err
	}
	defer r.Close()

	for _, file := range r.File {
		rc, err := file.Open()
		if err != nil {
			return "", err
		}
		defer rc.Close()

		// Create the destination file
		destFilepath := filepath.Join(destination, file.Name)
		if file.FileInfo().IsDir() {
			os.MkdirAll(destFilepath, os.ModePerm)
		} else {
			if err = os.MkdirAll(filepath.Dir(destFilepath), os.ModePerm); err != nil {
				return "", err
			}

			destFile, err := os.Create(destFilepath)
			if err != nil {
				return "", err
			}
			defer destFile.Close()

			_, err = io.Copy(destFile, rc)
			if err != nil {
				return "", err
			}
			if strings.HasSuffix(destFilepath, ".xlsx") {
				xlsxFilePath = destFilepath
			}
		}
	}
	return xlsxFilePath, nil
}
