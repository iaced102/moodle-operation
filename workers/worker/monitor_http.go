package worker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"moodle/internal/core/domain"
	repo "moodle/internal/repository/moodle"
	"moodle/pkg/mongodbiface"
)

type MonitorHTTP struct {
	mongoRepo *repo.MongoDB
}

func NewMonitorHTTPWorker(mongo mongodbiface.DB) *MonitorHTTP {
	mongoRepo := repo.NewMongoDB(mongo)
	return &MonitorHTTP{mongoRepo: mongoRepo}
}

func (w *MonitorHTTP) Run() {
	moodles, err := w.mongoRepo.GetAll("")
	if err != nil {
		log.Printf("Failed to get moodles: %v", err)
		return
	}

	for _, m := range moodles {
		w.checkAndAlertResponseTime(m)
	}
}

// Check HTTP response time for a moodle site and alert if threshold exceeded
func (w *MonitorHTTP) checkAndAlertResponseTime(m domain.Moodle) {
	if m.WebSiteName == "" {
		return
	}

	url := fmt.Sprintf("http://%s", m.WebSiteName)
	client := &http.Client{Timeout: 10 * time.Second}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("build request failed for %s: %v", url, err)
		return
	}

	start := time.Now()
	resp, err := client.Do(req)
	elapsed := time.Since(start)

	statusCode := 0
	if resp != nil {
		statusCode = resp.StatusCode
	}

	// Even on error, we still may want to alert if it took long time to fail.
	if err == nil && resp != nil && resp.Body != nil {
		_, _ = io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}

	shouldAlertSlow := elapsed.Seconds() > 1
	shouldAlertStatus := statusCode != 0 && statusCode != 200

	if shouldAlertSlow || shouldAlertStatus {
		if err := w.sendTelegramCurlAlert(m, elapsed, statusCode); err != nil {
			log.Printf("send curl-time tele fail for %s: %v", m.Id, err)
		}
	}
}

func (w *MonitorHTTP) sendTelegramCurlAlert(m domain.Moodle, elapsed time.Duration, statusCode int) error {
	apiUrl := "https://api.telegram.org/bot248926246:AAETwv7hzpk8zv6j9aRHDITYcnIRUGydS80/sendMessage"
	var text string
	if statusCode != 0 && statusCode != 200 {
		text = fmt.Sprintf("[HTTP NOK] KH: %s, Email: %s, Site: %s, Status: %d, RT: %.0fms",
			m.Name, m.Email, m.WebSiteName, statusCode, float64(elapsed.Milliseconds()))
	} else {
		text = fmt.Sprintf("[HTTP SLOW] KH: %s, Email: %s, Site: %s, RT: %.0fms (>500ms)",
			m.Name, m.Email, m.WebSiteName, float64(elapsed.Milliseconds()))
	}
	message := map[string]interface{}{
		"chat_id":           -1002960560112,
		"message_thread_id": 4,
		"text":              text,
	}
	jsonValue, _ := json.Marshal(message)
	req, err := http.NewRequest("POST", apiUrl, bytes.NewBuffer(jsonValue))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("telegram error: %s", resp.Status)
	}
	return nil
}
