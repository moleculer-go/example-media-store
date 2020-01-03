package services

import (
	"bytes"
	"encoding/base64"
	"errors"
	"image/jpeg"
	"image/png"
	"io"
	"os"

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
func saveToDisk(user, imageType, pic64, baseFolder string) (picHash, fileId string, bytesSize int, err error) {
	unbased, err := base64.StdEncoding.DecodeString(pic64)
	if err != nil {
		err = errors.New("Cannot decode base64 - source: " + err.Error())
		return picHash, fileId, bytesSize, err
	}
	picHash = hashPic(unbased)
	r := bytes.NewReader(unbased)

	path := baseFolder + user + "/"
	_ = os.MkdirAll(path, os.ModePerm)

	fileId = path + picHash + "." + imageType
	bytesSize = len(unbased)
	if imageType == "png" {
		err = savePng(r, fileId)
		if err != nil {
			return picHash, fileId, bytesSize, err
		}
	} else if imageType == "jpg" || imageType == "jpeg" {
		err = saveJpg(r, fileId)
		if err != nil {
			return picHash, fileId, bytesSize, err
		}
	} else {
		err = errors.New("Invalid imageType: " + imageType)
		return picHash, fileId, bytesSize, err
	}

	return picHash, fileId, bytesSize, nil
}

// saveToDatabase saves the image metadata and it's diskid to the database, so it can be searched and displayed in other apps.
func saveToDatabase(ctx moleculer.Context, user, fileId, picHash string, metadata map[string]interface{}) (string, error) {

	return "", nil
}

//resolvePicturesFolder return the folder where the servie will store the uploaded images
func resolvePicturesFolder(settings map[string]interface{}) string {
	pf, exists := settings["picturesFolder"]
	picturesFolder := ""
	if pfS, valid := pf.(string); exists && valid {
		picturesFolder = pfS
	} else {
		picturesFolder = "/pictures_store"
	}
	return picturesFolder
}

var settings map[string]interface{}
var Upload = moleculer.ServiceSchema{
	Name:     "upload",
	Settings: map[string]interface{}{},
	Started: func(ctx moleculer.BrokerContext, svc moleculer.ServiceSchema) {
		settings = svc.Settings
	},
	Actions: []moleculer.Action{
		{
			Name: "picture",
			Handler: func(ctx moleculer.Context, params moleculer.Payload) interface{} {
				user := params.Get("user").String()
				pic64 := params.Get("picture").String()
				metadata := params.Get("metadata").RawMap()
				imageType := metadata["imageType"].(string)

				picturesFolder := resolvePicturesFolder(settings)
				picHash, fileId, bytesSize, err := saveToDisk(user, imageType, pic64, picturesFolder)
				if err != nil {
					return err
				}

				picId, err := saveToDatabase(ctx, user, fileId, picHash, metadata)
				if err != nil {
					return err
				}

				ctx.Logger().Debug("picture uploaded succesfully! fileId: ", fileId, " bytesSize: ", bytesSize, " picId: ", picId)

				return fileId
			},
		},
	},
}
