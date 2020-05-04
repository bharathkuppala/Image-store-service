package utility

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Shopify/sarama"
	model "github.com/image-store/webservice/models"
)

// Response ...
type Response struct {
	Ok   bool
	Data interface{}
}

// SendResponse ...
func SendResponse(w http.ResponseWriter, ok bool, data interface{}) {
	response := Response{
		Ok:   ok,
		Data: data,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Println(err)
	}
}

// ValidateConfig ...
func ValidateConfig(config *sarama.Config) error {
	err := config.Validate()
	if err != nil {
		log.Fatalln("something wrong with above config", err)
		return err
	}
	return nil
}

// ToMap ...
func ToMap(w http.ResponseWriter, r *http.Request) (map[string]interface{}, error) {
	var resp map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&resp); err != nil {
		log.Println("error in decoding request body", err)
		return nil, err
	}
	return resp, nil
}

// CheckKeyExist ...
func CheckKeyExist(m1, m2 map[string][]model.Image, k1, k2 string) (
	v1, v2 []model.Image, ok1, ok2 bool) {

	v1, ok1 = m1[k1]
	v2, ok2 = m2[k2]
	return
}
