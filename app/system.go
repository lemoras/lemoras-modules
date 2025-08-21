package main

import (
	"encoding/json"
	"fmt"
	"initialize"
	"io/ioutil"
	"net/http"
	"strings"
)

func Initialize(w http.ResponseWriter, r *http.Request) {

	var in initialize.Request

	reqBody, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(reqBody, &in)

	fmt.Println("Endpoint hit: Initialize")

	in.Http.CustomHeader.Authorization = r.Header.Get("authorization")
	in.Http.Path = strings.Replace(r.URL.Path, "system/init/", "", -1)
	in.Http.Method = r.Method

	resp, _ := initialize.Invoke(in)
	w.Write([]byte(resp.Body))
}
