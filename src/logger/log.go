package logger


import (
	"log"
	"net/http"
	"bytes"
	"encoding/json"
   "os"
)


func getApiLogURL() string{
  return os.Getenv("LOG_API_URL")
}

func PutLog(data interface {}) {

		dataJson,err:=json.Marshal(data)
		client:= http.Client{}
    logURL :=getApiLogURL()
    if logURL != "" {
      req, _ := http.NewRequest("POST", logURL, bytes.NewBuffer(dataJson))
  		req.Header.Set("Content-Type", "application/json")
  		if err != nil {
  			log.Printf("Error while preparing request for sending to logging service")
  		}
  		resp, _ := client.Do(req)
  		if err != nil {
  			log.Printf("Error while sending request to logging service")
  		}else{
  			log.Printf("Data sent successfully!")
  		}

      defer resp.Body.Close()
    }else{
      log.Printf("Valid logging URL must be provided!")
    }


}
