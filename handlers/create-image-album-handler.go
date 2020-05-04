package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/image-store/webservice/kafka-services"
	model "github.com/image-store/webservice/models"
	utility "github.com/image-store/webservice/utilities"
)

var (
	//image *model.Image = new(model.Image)
	image model.Image
	np    kafka.NewProducer
	nc    kafka.NewConsumer
)

// ImageAlbum ...
type ImageAlbum struct {
	l *log.Logger
}

// NewImageAlbum ...
func NewImageAlbum(l *log.Logger) *ImageAlbum {
	return &ImageAlbum{l}
}

// ServeHTTP ... Acts as a main routing handler which routes request based on method
func (i *ImageAlbum) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		if r.URL.Path == "/api/v1/create-album" {
			i.createAlbum(w, r)

			return
		} else if r.URL.Path == "/api/v1/create-image" {
			i.CreateImage(w, r)
			return
		} else {
			//w.WriteHeader(http.StatusNotImplemented)
			utility.SendResponse(w, false, http.StatusNotFound)
			return
		}
	}

	if r.Method == http.MethodDelete {
		if r.URL.Path == "/api/v1/delete-album" {
			i.DeleteImageAlbum(w, r)
			return
		} else if r.URL.Path == "/api/v1/delete-image" {
			i.DeleteImage(w, r)
			return
		}
	}

	if r.Method == http.MethodGet {
		if r.URL.Path == "/api/v1/image" {
			i.GetImage(w, r)
			return
		} else if r.URL.Path == "/api/v1/images" {
			i.GetAllImages(w, r)
			return
		} else {
			//w.WriteHeader(http.StatusNotImplemented)
			utility.SendResponse(w, false, http.StatusNotFound)
			return
		}
	}

	w.WriteHeader(http.StatusNotImplemented)
}

func (i *ImageAlbum) createAlbum(w http.ResponseWriter, r *http.Request) {
	Producer, err := np.SetupProducer()
	if err != nil {
		log.Println("error in setting up producer", err)
		return
	}

	decodeRequest, err := utility.ToMap(w, r)
	if err != nil {
		return
	}
	albumName := decodeRequest["albumName"]
	// reading exsiting data from json file
	file, err := ioutil.ReadFile("image-album-db.json")
	if err != nil {
		fmt.Println("error in reading json file", err)
		utility.SendResponse(w, false, err)
		return
	}

	data := make(map[string][]model.Image, 0)
	if err = json.Unmarshal(file, &data); err != nil {
		log.Println("error with unmarshalling", err)
		utility.SendResponse(w, false, err)
		return
	}

	if _, ok := data[albumName.(string)]; ok {
		log.Println(albumName.(string) + " already exist.")
		utility.SendResponse(w, false, albumName.(string)+" already exist.")
		return
	}

	// Create new album
	jsonData := make(map[string][]model.Image, 0)
	jsonData = data
	jsonData[albumName.(string)] = []model.Image{}
	fileData, err := json.MarshalIndent(jsonData, " ", "\t")
	if err != nil {
		log.Println("marshal error", err)
		utility.SendResponse(w, false, err)
		return
	}
	ioutil.WriteFile("image-album-db.json", fileData, 0644)

	defer func() {
		np.ProducerMessage(Producer, "created new album "+albumName.(string))
	}()

	utility.SendResponse(w, true, "created new album "+albumName.(string))
}

// CreateImage ...
func (i *ImageAlbum) CreateImage(w http.ResponseWriter, r *http.Request) {
	Producer, err := np.SetupProducer()
	if err != nil {
		log.Println("error in setting up producer")
		return
	}

	var images []model.Image
	decodeRequest, err := utility.ToMap(w, r)
	if err != nil {
		utility.SendResponse(w, false, err)
		return
	}

	albumName := decodeRequest["albumName"]
	title := decodeRequest["title"]
	imageType := decodeRequest["imageType"]
	fileSize := decodeRequest["size"]
	createdOn := decodeRequest["createdOn"]

	file, _ := ioutil.ReadFile("image-album-db.json")
	data := make(map[string][]model.Image, 0)
	if err = json.Unmarshal(file, &data); err != nil {
		log.Println("error with unmarshalling", err)
		utility.SendResponse(w, false, err)
		return
	}

	if _, ok := data[title.(string)]; ok {
		log.Println(title.(string) + " already exist.")
		utility.SendResponse(w, false, title.(string)+" already exist.")
		return
	}

	defer func() {
		np.ProducerMessage(Producer, "created new image inside "+albumName.(string)+" with name "+title.(string))
	}()

	for _, v := range data[albumName.(string)] {
		image.Title = v.Title
		image.ImageType = v.ImageType
		image.Size = v.Size
		image.CreatedOn = v.CreatedOn
		images = append(images, image)
	}

	// New data from request
	image.Title = title.(string)
	image.ImageType = imageType.(string)
	image.Size = fileSize.(string)
	image.CreatedOn = createdOn.(string)
	images = append(images, image)

	if _, ok := data[albumName.(string)]; !ok {
		log.Println("key value does'nt exist!!")
		utility.SendResponse(w, false, err)
		return
	}

	data[albumName.(string)] = images

	fileData, err := json.MarshalIndent(data, " ", "\t")
	if err != nil {
		log.Println("marshal error", err)
		utility.SendResponse(w, false, err)
		return
	}
	ioutil.WriteFile("image-album-db.json", fileData, 0644)
	utility.SendResponse(w, true, "created new image inside "+albumName.(string)+" with name "+title.(string))
}

// DeleteImageAlbum ...
func (i *ImageAlbum) DeleteImageAlbum(w http.ResponseWriter, r *http.Request) {
	Producer, err := np.SetupProducer()
	if err != nil {
		log.Println("error in setting up producer")
		return
	}

	decodeRequest, err := utility.ToMap(w, r)
	if err != nil {
		utility.SendResponse(w, false, err)
		return
	}

	albumName := decodeRequest["albumName"]

	file, _ := ioutil.ReadFile("image-album-db.json")
	data := make(map[string][]model.Image, 0)
	if err = json.Unmarshal(file, &data); err != nil {
		log.Println("error with unmarshalling", err)
		utility.SendResponse(w, false, err)
		return
	}

	if _, ok := data[albumName.(string)]; !ok {
		log.Println(albumName.(string) + " does'nt exist.")
		utility.SendResponse(w, false, albumName.(string)+" does'nt exist, please provide valid album name.")
		return
	}
	delete(data, albumName.(string))

	fileData, err := json.MarshalIndent(data, " ", "\t")
	if err != nil {
		log.Println("marshal error", err)
		utility.SendResponse(w, false, err)
		return
	}
	ioutil.WriteFile("image-album-db.json", fileData, 0644)
	defer func() {
		np.ProducerMessage(Producer, albumName.(string)+" successfully deleted")
	}()

	utility.SendResponse(w, false, albumName.(string)+" successfully deleted")
}

// DeleteImage ...
func (i *ImageAlbum) DeleteImage(w http.ResponseWriter, r *http.Request) {
	Producer, err := np.SetupProducer()
	if err != nil {
		log.Println("error in setting up producer")
		return
	}

	var images []model.Image
	decodeRequest, err := utility.ToMap(w, r)
	if err != nil {
		utility.SendResponse(w, false, err)
		return
	}
	//from which album to delete
	albumName := decodeRequest["albumName"]
	//which image to delete
	imageTitle := decodeRequest["title"]

	file, _ := ioutil.ReadFile("image-album-db.json")
	data := make(map[string][]model.Image, 0)
	if err = json.Unmarshal(file, &data); err != nil {
		log.Println("error with unmarshalling", err)
		utility.SendResponse(w, false, err)
		return
	}

	if _, _, ok1, ok2 := utility.CheckKeyExist(data, data, albumName.(string), imageTitle.(string)); !ok1 && !ok2 {
		log.Println("given album name: " + albumName.(string) + " or image title: " + imageTitle.(string) + " does'nt exist or invalid")
		utility.SendResponse(w, false, "given album name: "+albumName.(string)+" or image title: "+imageTitle.(string)+" does'nt exist or invalid")
		return
	}

	for _, v := range data[albumName.(string)] {
		image.Title = v.Title
		image.ImageType = v.ImageType
		image.Size = v.Size
		image.CreatedOn = v.CreatedOn
		images = append(images, image)
	}

	for i, v := range images {
		if v.Title == imageTitle.(string) {
			images[i] = images[len(images)-1]
			images[len(images)-1] = model.Image{}
			images = images[:len(images)-1]
		}
	}
	data[albumName.(string)] = images
	fileData, err := json.MarshalIndent(data, " ", "\t")
	if err != nil {
		log.Println("marshal error", err)
		utility.SendResponse(w, false, err)
		return
	}
	ioutil.WriteFile("image-album-db.json", fileData, 0644)
	defer func() {
		np.ProducerMessage(Producer, "image "+imageTitle.(string)+" from "+albumName.(string)+" successfully deleted")
	}()
	utility.SendResponse(w, false, "image: "+imageTitle.(string)+" from "+albumName.(string)+" successfully deleted")
}

// GetAllImages ...
func (i *ImageAlbum) GetAllImages(w http.ResponseWriter, r *http.Request) {
	var images []model.Image
	decodeRequest, err := utility.ToMap(w, r)
	if err != nil {
		utility.SendResponse(w, false, err)
		return
	}

	albumName := decodeRequest["albumName"]

	file, _ := ioutil.ReadFile("image-album-db.json")
	data := make(map[string][]model.Image, 0)
	if err = json.Unmarshal(file, &data); err != nil {
		log.Println("error with unmarshalling", err)
		utility.SendResponse(w, false, err)
		return
	}

	for _, v := range data[albumName.(string)] {
		images = append(images, v)
	}

	utility.SendResponse(w, false, images)
}

// GetImage ...
func (i *ImageAlbum) GetImage(w http.ResponseWriter, r *http.Request) {
	decodeRequest, err := utility.ToMap(w, r)
	if err != nil {
		utility.SendResponse(w, false, err)
		return
	}
	//from which album to get a single image
	albumName := decodeRequest["albumName"]

	// requested image
	imageTitle := decodeRequest["title"]

	file, _ := ioutil.ReadFile("image-album-db.json")
	data := make(map[string][]model.Image, 0)
	if err = json.Unmarshal(file, &data); err != nil {
		log.Println("error with unmarshalling", err)
		utility.SendResponse(w, false, err)
		return
	}

	for _, v := range data[albumName.(string)] {
		if v.Title == imageTitle.(string) {
			utility.SendResponse(w, true, v)
			return
		}
	}
}
