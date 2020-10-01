package intern

import (
	"bytes"
	"drops-backend/database"
	"drops-backend/models"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/Viva-con-Agua/echo-pool/resp"
	"github.com/labstack/echo"
)

func ProfileListInternal(c echo.Context) (err error) {
	body := new(models.ListRequest)
	if err = c.Bind(body); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	filter := body.Filter()
	response, api_err := database.ProfileListInternal(filter)
	if api_err.Error != nil {

		api_err.LogError(c, body)
		return c.JSON(http.StatusInternalServerError, resp.InternelServerError)
	}
	return c.JSON(http.StatusOK, response)
}

func ProfileGetRequest(u_uuid string) (p *models.Profile, err error) {
	req, err := http.NewRequest("GET", "http://"+os.Getenv("PROFILE_HOST")+"/intern/profiles/user/"+u_uuid, nil)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Print(err, " ### intern.UserRequest Step_1")
		return nil, err
	}
	defer resp.Body.Close()
	//var a []auth.User

	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	body, err := ioutil.ReadAll(resp.Body)

	err = json.Unmarshal(body, &p)
	if err != nil {
		return nil, err
	}
	return p, err
}

func ProfileCreateRequest(users *models.ProfileCreate) (p *models.Profile, err error) {
	reqBodyBytes := new(bytes.Buffer)
	json.NewEncoder(reqBodyBytes).Encode(users)
	req, err := http.NewRequest("POST", "http://"+os.Getenv("PROFILE_HOST")+"/intern/profiles", bytes.NewBuffer(reqBodyBytes.Bytes()))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Print(err, " ### intern.UserRequest Step_1")
		return nil, err
	}
	defer resp.Body.Close()
	//var a []auth.User

	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	err = json.NewDecoder(resp.Body).Decode(&p)
	if err != nil {
		return nil, err
	}
	return p, err
}
