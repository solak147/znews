package service

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"net/smtp"
	"os"
	"strings"
	"znews/app/middleware"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func Send(title string, body string, to string) error {
	envErr := godotenv.Load()
	if envErr != nil {
		panic(envErr)
	}

	domain := os.Getenv("DOMAIN")
	from := os.Getenv("GMAIL_USERNAME")
	pass := os.Getenv("GMAIL_PASSWORD")
	port := os.Getenv("GMAIL_PORT")
	server := os.Getenv("GMAIL_SERVER")

	// 讀取另一個 HTML 檔案的內容
	htmlFile := "app/service/verify.html"
	htmlContent, errFile := ioutil.ReadFile(htmlFile)
	if errFile != nil {
		middleware.Logger().WithFields(logrus.Fields{
			"title": "Read html Failed",
		}).Error(errFile.Error())
		return errors.New("發送郵件失敗" + errFile.Error())
	}

	htmlStr := strings.Replace(string(htmlContent), "xxxxxx", title, 1)

	imagePaths := []string{"body_bg.jpg", "dot_image.jpg", "header_bg.jpg", "image_1.jpg", "image_2.jpg", "image_3.jpg", "image_4.png", "image_5.png"} // 圖片路徑列表
	for index, imagePath := range imagePaths {
		htmlStr = strings.Replace(htmlStr, fmt.Sprintf("serimg%d", index), domain+"/serimg/"+imagePath, 1)
	}

	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: " + title + "\n" +
		"MIME-version: 1.0\n" +
		"Content-Type: multipart/related; boundary=boundary\n" +
		"--boundary\n" +
		"Content-Type: text/html; charset=\"UTF-8\"\n" + htmlStr

	// 添加圖片內容
	// imagePaths := []string{"body_bg.jpg", "dot_image.jpg", "header_bg.jpg", "image_1.jpg", "image_2.jpg", "image_3.jpg", "image_4.png", "image_5.png"} // 圖片路徑列表
	// for index, imagePath := range imagePaths {
	// 	imageData, err := ioutil.ReadFile("app/service/images/" + imagePath)
	// 	if err != nil {
	// 		middleware.Logger().WithFields(logrus.Fields{
	// 			"title": "無法讀取圖片:",
	// 		}).Error(err.Error())
	// 		return errors.New("無法讀取圖片:" + err.Error())
	// 	}
	// 	encodedImage := encodeImage(imageData)
	// 	imageID := fmt.Sprintf("img%d", index)

	// 	msg += "--boundary\n"
	// 	msg += "Content-Type: image/jpeg\n"
	// 	msg += "Content-Transfer-Encoding: base64\n"
	// 	msg += fmt.Sprintf("Content-ID: <%s>\n", imageID)
	// 	msg += "Content-Disposition: inline\n"
	// 	msg += encodedImage + "\n"
	// }

	err := smtp.SendMail(server+":"+port,
		smtp.PlainAuth("", from, pass, server),
		from, []string{to}, []byte(msg))
	if err != nil {
		middleware.Logger().WithFields(logrus.Fields{
			"title": "SendMail Failed",
		}).Error(err.Error())
		return errors.New("發送郵件失敗" + err.Error())
	}

	return nil
}

// 將圖片資料進行 Base64 編碼
func encodeImage(imageData []byte) string {
	encodedImage := base64.StdEncoding.EncodeToString(imageData)
	return encodedImage
}
