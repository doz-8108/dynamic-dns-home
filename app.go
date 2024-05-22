package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/robfig/cron/v3"
)

var currentIp, cfZoneId, cfDnsRecordId, targetDomain, cfDnsApiKey, cfApiEmail string

type dnsUpdateReqBody struct {
	Content string   `json:"content"`
	Name    string   `json:"name"`
	Proxied bool     `json:"proxied,omitempty"`
	Type    string   `json:"type"`
	Comment string   `json:"comment,omitempty"`
	Tags    []string `json:"tags,omitempty"`
	Ttl     int32    `json:"ttl,omitempty"` // 60s ~ 86400s
}

type dnsUpdateResponseMsgItem struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type DNSRecordResponse struct {
	Errors   []dnsUpdateResponseMsgItem `json:"errors"`
	Messages []dnsUpdateResponseMsgItem `json:"messages"`
	Success  bool                       `json:"success"`
	Result   struct {
		Content   string    `json:"content"`
		Name      string    `json:"name"`
		Proxied   bool      `json:"proxied"`
		Type      string    `json:"type"`
		Comment   string    `json:"comment"`
		CreatedOn time.Time `json:"created_on"`
		ID        string    `json:"id"`
		Locked    bool      `json:"locked"`
		Meta      struct {
			AutoAdded bool   `json:"auto_added"`
			Source    string `json:"source"`
		} `json:"meta"`
		ModifiedOn time.Time `json:"modified_on"`
		Proxiable  bool      `json:"proxiable"`
		Tags       []string  `json:"tags"`
		TTL        int       `json:"ttl"`
		ZoneID     string    `json:"zone_id"`
		ZoneName   string    `json:"zone_name"`
	} `json:"result"`
}

type logWriter struct{}

func (writer logWriter) Write(bytes []byte) (int, error) {
	return fmt.Printf("[%s] %s", time.Now().UTC().Add(time.Hour*8).Format(time.UnixDate), string(bytes))
}

func main() {
	log.SetFlags(0)
	log.SetOutput((new(logWriter)))

	cfZoneId = os.Getenv("CF_ZONE_ID")
	cfDnsRecordId = os.Getenv("CF_DNS_RECORD_ID")
	cfDnsApiKey = os.Getenv("CF_DNS_API_KEY")
	cfApiEmail = os.Getenv("CF_DNS_API_EMAIL")
	targetDomain = os.Getenv("TARGET_DOMAIN")

	cronScheduler := cron.New()
	cronScheduler.AddFunc("0 * * * *", updateDnsRecord)
	log.Println("Cron job scheduled to run every hour...")
	cronScheduler.Start()

	for {
		time.Sleep(time.Hour)
	}
}

func decodeHttpRespBody[T any](resp *http.Response, output *T) {
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Failed to read response body: ", err)
		return
	}
	json.Unmarshal(respBody, output)
}

func updateDnsRecord() {
	var ipRespData struct {
		Ip string
	}
	ipResp, err := http.Get("https://api.ipify.org?format=json")
	if err != nil {
		log.Println("Failed to get public IP address: ", err)
		return
	}

	decodeHttpRespBody(ipResp, &ipRespData)
	log.Println("Public IP address: ", ipRespData.Ip)

	if currentIp == ipRespData.Ip {
		log.Println("Public IP address has not changed.")
		return
	}
	currentIp = ipRespData.Ip

	dnsUpdateReqBodyJson, err := json.Marshal(dnsUpdateReqBody{
		Content: currentIp,
		Name:    targetDomain,
		Type:    "A",
	})
	if err == nil {
		req, _ := http.NewRequest(http.MethodPut, "https://api.cloudflare.com/client/v4/zones/"+cfZoneId+"/dns_records/"+cfDnsRecordId, bytes.NewReader(dnsUpdateReqBodyJson))
		req.Header.Add("X-Auth-Key", cfDnsApiKey)
		req.Header.Add("X-Auth-Email", cfApiEmail)
		req.Header.Add("Content-Type", "application/json")

		var dnsUpdateRespData DNSRecordResponse
		dnsUpdateResp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Println("Failed to update DNS record: ", err)
			return
		}
		decodeHttpRespBody(dnsUpdateResp, &dnsUpdateRespData)
		if len(dnsUpdateRespData.Errors) != 0 {
			for _, error := range dnsUpdateRespData.Errors {
				log.Printf("(Error code: %d): %s", error.Code, error.Message)
				return
			}
		}
		log.Printf("DNS record for %s updated to %s successfully", dnsUpdateRespData.Result.Name, dnsUpdateRespData.Result.Content)
	}
}
