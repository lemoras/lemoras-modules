package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"note"
)

func Note(w http.ResponseWriter, r *http.Request) {

	var in note.Request

	reqBody, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(reqBody, &in)

	fmt.Println("Endpoint hit: Note")

	in.Http.CustomHeader.Authorization = r.Header.Get("authorization")
	in.Http.Method = r.Method

	resp, _ := note.Invoke(in)
	w.Write([]byte(resp.Body))
}
