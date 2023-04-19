package dao

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func GetCase(queryStr string) ([]interface{}, error) {
	// 建立 HTTP 客戶端
	httpClient := &http.Client{}

	// 設定 Elasticsearch 的主機和埠號
	esHost := "http://localhost:9200"

	// 建立搜尋請求的 JSON 樣板
	var searchBody bytes.Buffer
	searchBody.WriteString(queryStr)

	// 建立搜尋請求
	searchRequest, err := http.NewRequest("POST", esHost+"/gorm/_search", &searchBody)
	if err != nil {
		return nil, err
	}
	searchRequest.Header.Set("Content-Type", "application/json")

	// 發送搜尋請求
	searchResponse, err := httpClient.Do(searchRequest)
	if err != nil {
		return nil, err
	}
	defer searchResponse.Body.Close()

	// 解析搜尋結果
	var searchResult map[string]interface{}
	json.NewDecoder(searchResponse.Body).Decode(&searchResult)

	// 處理搜尋結果
	if hits, ok := searchResult["hits"].(map[string]interface{})["hits"].([]interface{}); ok {
		//fmt.Printf("Found a total of %d documents\n", len(hits))

		// for _, hit := range hits {
		// 	// 取得文件 ID 和內容
		// 	id := hit.(map[string]interface{})["_id"].(string)
		// 	source := hit.(map[string]interface{})["_source"].(map[string]interface{})

		// 	// 處理文件內容
		// 	fmt.Printf("Document ID: %s\n", id)
		// 	fmt.Printf("Content: %v\n", source)
		// }
		return hits, nil
	} else {
		// No documents found
		return nil, nil
	}
}
