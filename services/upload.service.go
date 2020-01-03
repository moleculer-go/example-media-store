package services

import (
	"bytes"
	"encoding/base64"
	"errors"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path/filepath"

	"github.com/moleculer-go/gateway/websocket"
	"github.com/moleculer-go/moleculer"

	"crypto/sha1"
	 
)

var websocketMixin = &websocket.WebSocketMixin{
	Mixins: []websocket.SocketMixin{
		&websocket.EventsMixin{},
	},
}

// savePng save a png image
func savePng(r io.Reader, path string) error {
	im, err := png.Decode(r)
	if err != nil {
		return errors.New("Bad png - source: " + err.Error())
	}
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		return errors.New("Cannot open file(" + path + "): " + err.Error())
	}
	png.Encode(f, im)
	return nil
}

// saveJpg save a jpg image
func saveJpg(r io.Reader, path string) error {
	im, err := jpeg.Decode(r)
	if err != nil {
		return errors.New("Bad jpeg - source: " + err.Error())
	}
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		return errors.New("Cannot open file(" + path + "): " + err.Error())
	}
	jpeg.Encode(f, im, &jpeg.Options{Quality: 100})
	return nil
}

//hashPic create a sha hash of the image bytes
//to identity the image uniquely
func hashPic(b []byte) string {
	hasher := sha1.New()
	hasher.Write(b)
	sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	return sha
}

// saveToDisk save the image to disk and return the unique file Id for the picture. So it can be retrieved from anywhere with this id.
func saveToDisk(user, imageType, pic64, baseFolder string) (fileId string, bytesSize int, err error) {
	unbased, err := base64.StdEncoding.DecodeString(pic64)
	if err != nil {
		err = errors.New("Cannot decode base64 - source: " + err.Error())
		return fileId, bytesSize, err
	}
	picHash := hashPic(unbased)
	r := bytes.NewReader(unbased)

	path := baseFolder + "/pics_store/" + user + "/"
	_ = os.MkdirAll(path, os.ModeDir)

	fileId = path + picHash + "." + imageType

	if imageType == "png" {
		err = savePng(r, fileId)
		if err != nil {
			return fileId, bytesSize, err
		}
	} else if imageType == "jpg" || imageType == "jpeg" {
		err = saveJpg(r, fileId)
		if err != nil {
			return fileId, bytesSize, err
		}
	}

	err = errors.New("Invalid imageType: " + imageType)

	return fileId, bytesSize, err
}

// saveToDatabase saves the image metadata and it's diskid to the database, so it can be searched and displayed in other apps.
func saveToDatabase(user, fileId string, metadata map[string]interface{}) string {
	return ""
}

var Upload = moleculer.ServiceSchema{
	Name: "upload",
	Actions: []moleculer.Action{
		{
			Name: "picture",
			Handler: func(ctx moleculer.Context, params moleculer.Payload) interface{} {
				user := params.Get("user").String()
				pic64 := params.Get("picture").String()
				metadata := params.Get("metadata").RawMap()
				imageType := metadata["imageType"].(string)
				baseFolder, _ := filepath.Abs("./")
				baseFolder = baseFolder + "/_test_"

				fileId, bytesSize, err := saveToDisk(user, imageType, pic64, baseFolder)
				if err != nil {
					return err
				}
				picId := saveToDatabase(user, fileId, metadata)

				ctx.Logger().Debug("picture uploaded succesfully! fileId: ", fileId, " bytesSize: ", bytesSize, " picId: ", picId)

				return fileId
			},
		},
	},
}
