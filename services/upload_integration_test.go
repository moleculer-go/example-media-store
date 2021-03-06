package services

import (
	"encoding/base64"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/moleculer-go/moleculer"
	"github.com/moleculer-go/moleculer/broker"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var logLevel = "error"

var _ = Describe("Upload Service", func() {

	bkr := broker.New(&moleculer.Config{
		LogLevel: logLevel,
	})
	Describe("Upload Picture", func() {

		BeforeEach(func() {
			picturesFolder := os.TempDir() + "upload_test/"
			os.RemoveAll(picturesFolder)
			Upload.Settings["picturesFolder"] = picturesFolder
			bkr.Publish(Upload, MediaService)
			bkr.Start()
		})

		AfterEach(func() {
			bkr.Stop()
		})

		It("upload.picture action should save the picture content to disk and metadata to the database", func() {
			user := "12345"
			picture := loadPic("_test_/car1.jpg")
			metadata := map[string]interface{}{
				"width":     1024,
				"height":    900,
				"imageType": "jpg",
			}

			r := <-bkr.Call("upload.picture", map[string]interface{}{
				"user":     user,
				"picture":  picture,
				"metadata": metadata,
			})
			Expect(r.Error()).Should(BeNil())

			fileId := r.String()
			Expect(fileId).Should(BeARegularFile())

			//check db records
			time.Sleep(time.Second)
			um := <-bkr.Call("userMedia.find", map[string]interface{}{})
			Expect(um.Error()).Should(Succeed())
			Expect(um.Len()).Should(Equal(1))
			Expect(um.First().Get("picHash").String()).Should(Equal("YVQLpQgKGn5QOHJJLV-c39mqAhk="))
			Expect(um.First().Get("metadata").Get("imageType").String()).Should(Equal("jpg"))
			Expect(um.First().Get("metadata").Get("bytesSize").String()).Should(Equal("92285"))
			Expect(um.First().Get("metadata").Get("width").String()).Should(Equal("1024"))
			Expect(um.First().Get("metadata").Get("height").String()).Should(Equal("900"))

			am := <-bkr.Call("allMedia.find", map[string]interface{}{})
			Expect(am.Error()).Should(Succeed())
			Expect(am.Len()).Should(Equal(1))
			Expect(am.First().Get("picHash").String()).Should(Equal("YVQLpQgKGn5QOHJJLV-c39mqAhk="))
			Expect(am.First().Get("metadata").Get("imageType").String()).Should(Equal("jpg"))
			Expect(am.First().Get("metadata").Get("bytesSize").String()).Should(Equal("92285"))
			Expect(am.First().Get("metadata").Get("width").String()).Should(Equal("1024"))
			Expect(am.First().Get("metadata").Get("height").String()).Should(Equal("900"))
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
