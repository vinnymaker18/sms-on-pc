package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"os"

	"github.com/vinnymaker18/sms-on-pc/backend/common"
	"github.com/vinnymaker18/sms-on-pc/backend/storage"
)

const (
	contentTypeHeader = "Content-Type"

	jsonContentType = "application/json"

	// Accepts server control commands (shutdown, restart etc...) on this port.
	// Useful for building tools to manage servers.
	controlPort = 8001
)

func parseTextMessage(req *http.Request) (*common.SMSMessage, error) {
	form := req.Form

	// each message post request has user_id, from_addr, body fields in it.
	if _, ok := form["userid"]; !ok {
		return nil, fmt.Errorf("no user_id field found in request")
	}

	userID, err := strconv.ParseInt(form["userid"][0], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid user_id in request")
	}

	originAddress, ok := form["origin"]
	if !ok {
		return nil, fmt.Errorf("no origin address field in request")
	}

	msgBody, ok := form["msgbody"]
	if !ok {
		return nil, fmt.Errorf("no message body in request")
	}

	return &common.SMSMessage{
		UserID:        int64(userID),
		Time:          time.Now(),
		MsgBody:       msgBody[0],
		OriginAddress: originAddress[0],
	}, nil
}

func smsReadHandler(w http.ResponseWriter, req *http.Request) {
	if _, ok := req.Form["userid"]; !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userID, err := strconv.ParseInt(req.Form["userid"][0], 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// TODO - Ignore any serialization errors for now, come back to this and fix it
	serializedBytes, _ := json.Marshal(storage.FetchNewSMS(userID))
	w.Write(serializedBytes)
}

func smsWriteHandler(w http.ResponseWriter, req *http.Request) {
	newTextMsg, err := parseTextMessage(req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	storage.StoreNewSMS(newTextMsg)
}

func parseMessageIDs(req *http.Request) ([]int64, error) {
	databytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}

	var reqBody map[string][]int64
	json.NewDecoder(req.Body).Decode(&reqBody)
	msgIDParams := reqBody["msgids"]

	msgIDs := make([]int64, 0)
	for _, id := range msgIDParams {
		parsed, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("Invalid msg id parameter")
		}
		msgIDs = append(msgIDs, parsed)
	}

	fmt.Printf("%v\n", msgIDs)
	return msgIDs, nil
}

func markSmsHandler(w http.ResponseWriter, req *http.Request) {
	msgIDs, err := parseMessageIDs(req)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Invalid syntax for message Ids")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte{})
		return
	}

	err = storage.MarkAsRead(msgIDs)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	} else {
		w.Write([]byte("\"success\""))
	}
}

func main() {
	http.HandleFunc("/sms/mark", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set(contentTypeHeader, jsonContentType)
		req.ParseForm()

		if req.Method == http.MethodPost {
			fmt.Fprintln(os.Stderr, "/sms/mark post request received")
			markSmsHandler(w, req)
		} else {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte{})
		}
	})

	http.HandleFunc("/sms", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set(contentTypeHeader, jsonContentType)
		req.ParseForm()

		if req.Method == http.MethodGet {
			smsReadHandler(w, req)
		} else if req.Method == http.MethodPost {
			smsWriteHandler(w, req)
		}
	})

	go func() {
		fmt.Println("delete old sms goroutine initiated")
		for {
			storage.DeleteOldSMS()
			time.Sleep(time.Hour)
		}
	}()

	http.ListenAndServe(":8000", nil)
}
