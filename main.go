package main

import (
	"net/http"
	"os"

	"github.com/jmandal/import_excel/api"
	"github.com/jmandal/import_excel/models"
	"github.com/jmandal/import_excel/views"
	"github.com/uadmin/uadmin"
)

func main() {
	uadmin.Register(
		models.Menus{},
		models.MenuItems{},
	)

	http.HandleFunc("/", uadmin.Handler(views.IndexHandler))
	http.HandleFunc("/api/upload_file", uadmin.Handler(api.UploadFile))

	// Specify the folder path
	folderPath := "media/uploads/"
	perm := os.FileMode(0777)

	err := createFolderIfNotExists(folderPath, perm)
	if err != nil {
		uadmin.Trail(uadmin.ERROR, "Error creating folder %s. %s", folderPath, err)
		return
	}

	uadmin.RootURL = "/admin/"
	uadmin.Port = 7373
	uadmin.StartServer()
}

func createFolderIfNotExists(folderPath string, perm os.FileMode) error {
	// Check if the folder exists
	_, err := os.Stat(folderPath)
	if os.IsNotExist(err) {
		// Folder does not exist, create it
		err := os.MkdirAll(folderPath, perm)
		if err != nil {
			return err
		}
		uadmin.Trail(uadmin.INFO, "Folder '%s' created successfully.", folderPath)
	} else if err != nil {
		// An error occurred while checking the folder existence
		return err
	} else {
		// Folder already exists
		uadmin.Trail(uadmin.INFO, "Folder '%s' already exists.", folderPath)
	}
	return nil
}
