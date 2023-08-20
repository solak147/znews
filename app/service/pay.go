package service

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/url"
	"os"
	"sort"
	"strings"
	"time"
	"znews/app/dao"
	"znews/app/middleware"
	"znews/app/model"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func CreditAll(card string) (url.Values, error){
	domain := os.Getenv("DOMAIN")

	product := model.Product{}
	if err := dao.GormSession.Select("*").Where("name=?", card).First(&product).Error; err != nil {
		return nil, err
	}

	// 獲取當前時間
	currentTime := time.Now()
	// 指定日期和時間格式
	layout := "2006/01/02 15:04:05"

	// 準備付款資料
	data := url.Values{}
	data.Set("MerchantID", "3002607")
	data.Set("MerchantTradeNo", generateOrderNumber())
	data.Set("PaymentType", "aio")
	data.Set("MerchantTradeDate", currentTime.Format(layout))
	data.Set("TotalAmount", product.Price)
	data.Set("TradeDesc", "Buddha綠界金流信用卡交易")
	data.Set("ItemName", product.ChiName)
	data.Set("ReturnURL", domain+":82/pay/result")
	data.Set("ClientBackURL", domain+":3000/deposit")
	data.Set("ChoosePayment", "Credit")
	data.Set("EncryptType", "1")

	// 計算 CheckMacValue
	checkMacValue := computeCheckMacValue(data, "pwFHCqoQZGmho4w6", "EkRm7iFT261dpevs")
	data.Set("CheckMacValue", checkMacValue)

	// // 發送 POST 請求
	// response, err := http.PostForm(apiUrl, data)
	// if err != nil {
	// 	fmt.Println("HTTP POST error:", err)
	// 	return "", err
	// }
	// defer response.Body.Close()

	// // 讀取回傳的資料
	// body, err := ioutil.ReadAll(response.Body)
	// if err != nil {
	// 	fmt.Println("Response read error:", err)
	// 	return "", err
	// }

	// 處理回傳資料
	// fmt.Println(string(body))
	return data, nil
}

// 計算 CheckMacValue
func computeCheckMacValue(data url.Values, hashKey, hashIV string) string {
    keys := []string{}
    for k := range data {
        keys = append(keys, k)
    }
    sort.Strings(keys)
    query := ""
    for _, k := range keys {
        if data.Get(k) != "" {
            query += fmt.Sprintf("%s=%s&", k, data.Get(k))
        }
    }

	query = fmt.Sprintf("HashKey=%s&", hashKey) + query
	query = query + fmt.Sprintf("HashIV=%s&", hashIV) 

    query = strings.TrimRight(query, "&")

	// 對字串進行 URL 編碼
	query = url.QueryEscape(query)
	query = strings.ToLower(query)


	// 使用 SHA-256 算法計算雜湊值
    hash := sha256.New()
    hash.Write([]byte(query))
    hashValue := hash.Sum(nil)

    // 將雜湊值轉換為十六進制表示
    hashHex := hex.EncodeToString(hashValue)
	hashHex =strings.ToUpper(hashHex)

    return hashHex
}

func generateOrderNumber() string {
	orderUUID := uuid.New()

	uuidString := orderUUID.String()

     // 去掉连字符并截取前20个字符
	 cleanedUUID := strings.ReplaceAll(uuidString, "-", "")
	 if len(cleanedUUID) > 20 {
		 cleanedUUID = cleanedUUID[:20]
	 }

    return cleanedUUID
}

func Result(c *gin.Context) error{

	// 读取整个 HTTP 请求体，这包括了表单数据
	body, err := c.GetRawData()
	if err != nil {
		return err
	}
	fmt.Println(string(body))

	middleware.Logger().WithFields(logrus.Fields{
		"title": "receive pay result error",
	}).Error(string(body))

	return nil

 
}