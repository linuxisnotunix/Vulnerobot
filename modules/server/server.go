package server

import (
	"net/http"
	"strings"

	log "github.com/Sirupsen/logrus"

	"github.com/linuxisnotunix/Vulnerobot/modules/assets"
)

//HandlePublic handlerfunc to publish assets
func HandlePublic(res http.ResponseWriter, req *http.Request) {
	log.Debug("GET " + req.URL.Path)
	path := strings.TrimPrefix(req.URL.Path, "/public")
	data, hash, contentType, err := assets.Asset("public", path)
	if err != nil {
		/*
		   log.Debug("NOT FOUND :  " + path)
		   data = []byte(err.Error())
		*/
		///*
		data, hash, contentType, err = assets.Asset("public", "/index.html") //TODO better redirect to /public/index.html
		if err != nil {
			data = []byte(err.Error())
		}
		//*/
	}
	res.Header().Set("Content-Encoding", "gzip")
	res.Header().Set("Content-Type", contentType)
	res.Header().Add("Cache-Control", "public, max-age=31536000")
	res.Header().Add("ETag", hash)
	if req.Header.Get("If-None-Match") == hash {
		res.WriteHeader(http.StatusNotModified)
	} else {
		res.WriteHeader(http.StatusOK)
		_, err := res.Write(data)
		if err != nil {
			panic(err)
		}
	}
}

//HandleAPI handlerfunc to publish assets
func HandleAPI(res http.ResponseWriter, req *http.Request) {
}
