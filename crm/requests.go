package crm

import (
	"bytes"
	"drops-backend/models"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/Viva-con-Agua/echo-pool/api"
	"github.com/Viva-con-Agua/echo-pool/crm"
)

func IrobertCreateUser(u_crm *models.CrmUserSignUp) (err error) {
	reqBodyBytes := new(bytes.Buffer)
	json.NewEncoder(reqBodyBytes).Encode(u_crm)
	req, err := http.NewRequest("POST",
		os.Getenv("IROBERT_URL")+
			"/IB_Run4WaterImport_DE_DROPS_USER_CREATE_V01",
		bytes.NewBuffer(reqBodyBytes.Bytes()),
	)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Print(err, " ### intern.UserRequest Step_1")
		return err
	}
	defer resp.Body.Close()
	log.Print(resp)
	//var a []auth.User

	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	if resp.StatusCode == 201 {
		return nil
	}
	return api.ErrorConflict
}

func IrobertJoinEvent(u_crm *models.CrmDataBody) (err error) {
	reqBodyBytes := new(bytes.Buffer)
	json.NewEncoder(reqBodyBytes).Encode(u_crm)
	req, err := http.NewRequest("POST",
		os.Getenv("IROBERT_URL")+
			"/IB_Run4WaterImport_DE_EVENT_JOIN_V01",
		bytes.NewBuffer(reqBodyBytes.Bytes()),
	)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Print(err, " ### intern.UserRequest Step_1")
		return err
	}
	defer resp.Body.Close()
	log.Print(resp)
	//var a []auth.User

	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	if resp.StatusCode == 201 {
		return nil
	}
	return err
}

func IrobertEmail(u_crm *crm.CrmEmailBody) (err error) {
	reqBodyBytes := new(bytes.Buffer)
	json.NewEncoder(reqBodyBytes).Encode(u_crm)
	req, err := http.NewRequest("POST",
		os.Getenv("IROBERT_URL")+
			"/IB_Run4WaterImport_DE_EMAIL_SEND_V01",
		bytes.NewBuffer(reqBodyBytes.Bytes()),
	)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Print(err, " ### intern.NEWMAIL")
		return err
	}
	defer resp.Body.Close()
	log.Print(resp)
	//var a []auth.User

	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	if resp.StatusCode == 201 {
		return nil
	}
	return err
}
