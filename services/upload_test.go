package services

import (
	"encoding/base64"
	"io/ioutil"
	"log"
	"os"

	"github.com/moleculer-go/moleculer"
	"github.com/moleculer-go/moleculer/broker"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var logLevel = "fatal"

var _ = Describe("Upload Service", func() {

	bkr := broker.New(&moleculer.Config{
		LogLevel: logLevel,
	})
	Describe("Upload Picture", func() {

		BeforeEach(func() {
			bkr.Publish(Upload)
			bkr.Start()
		})

		AfterEach(func() {
			bkr.Stop()
		})

		It("upload.picture action should save the picture content to disk and metadata to the database", func(done Done) {

			user := "12345"
			picture := loadPic("_test_/car1.jpg")
			metadata := map[string]interface{}{
				"size":      1024,
				"imageType": "jpg",
			}

			r := <-bkr.Call("upload.picture", map[string]interface{}{
				"user":     user,
				"picture":  picture,
				"metadata": metadata,
			})
			Expect(r.Error()).Should(BeNil())

			fileId := r.String()
			Expect(fileId).ShouldNot(Equal(""))

			close(done)
		})
	})
})

//loadPic return a picture in base64 string
func loadPic(path string) string {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	imgBytes, err := ioutil.ReadAll(file)
	return base64.StdEncoding.EncodeToString(imgBytes)
}
