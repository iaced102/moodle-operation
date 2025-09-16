package worker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"

	repo "moodle/internal/repository/moodle"
	"moodle/pkg/mongodbiface"
)

type MonthlyCCUReportWorker struct {
	mongoRepo *repo.MongoDB
}

// DataPoint represents a CCU sample from Prometheus
type DataPoint struct {
	Timestamp int64
	Date      string
	Month     int
	Day       int
	Hour      int
	CCU       float64
}

// CalCCUconfig minimal mapping for CalCCUconfig table
type CalCCUconfig struct {
	ID     string
	DBName string
}

func NewMonthlyCCUReportWorker(mongo mongodbiface.DB) *MonthlyCCUReportWorker {
	mongoRepo := repo.NewMongoDB(mongo)
	return &MonthlyCCUReportWorker{mongoRepo: mongoRepo}
}

// PrometheusResponse models the query_range response
type PrometheusResponse struct {
	Status string `json:"status"`
	Data   struct {
		ResultType string `json:"resultType"`
		Result     []struct {
			Metric map[string]interface{} `json:"metric"`
			Values [][]interface{}        `json:"values"`
		} `json:"result"`
	} `json:"data"`
}

func calculatePreviousMonthRange(now time.Time) (time.Time, time.Time) {
	firstOfThisMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	end := firstOfThisMonth.Add(-time.Second)
	start := time.Date(end.Year(), end.Month(), 1, 0, 0, 0, 0, end.Location())
	return start, end
}

func buildPrometheusRequest(dbname string, start, end time.Time) (*http.Request, error) {
	data := url.Values{}
	data.Set("start", strconv.FormatInt(start.Unix(), 10))
	data.Set("end", strconv.FormatInt(end.Unix(), 10))
	data.Set("step", "1200")
	data.Set("query", fmt.Sprintf(`lmsccu_online_user{dbname="%s"}`, dbname))

	u, _ := url.ParseRequestURI("https://thor-hn-metrics.bizflycloud.vn/api/datasources/proxy/11/api/v1/query_range")
	urlStr := u.String()
	req, err := http.NewRequest(http.MethodPost, urlStr, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return req, nil
}

func doHTTPRequest(req *http.Request) (string, error) {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	buf := new(strings.Builder)
	if _, err := io.Copy(buf, resp.Body); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func parsePrometheusResponse(body string) (PrometheusResponse, error) {
	var responseData PrometheusResponse
	if err := json.Unmarshal([]byte(body), &responseData); err != nil {
		return PrometheusResponse{}, err
	}
	return responseData, nil
}

func generateExcelAndSave(responseData PrometheusResponse, dbname string, start time.Time) (*excelize.File, string, string, error) {
	if len(responseData.Data.Result) == 0 || len(responseData.Data.Result[0].Values) == 0 {
		return nil, "", "", fmt.Errorf("no data to generate report")
	}

	f := excelize.NewFile()
	index, _ := f.NewSheet("RawData")

	f.SetCellValue("RawData", "A1", "Timestamp")
	f.SetCellValue("RawData", "B1", "Date")
	f.SetCellValue("RawData", "C1", "Tháng")
	f.SetCellValue("RawData", "D1", "Ngày")
	f.SetCellValue("RawData", "E1", "Giờ")
	f.SetCellValue("RawData", "F1", "CCu")

	for i, value := range responseData.Data.Result[0].Values {
		timestamp := value[0]
		dateTime := time.Unix(int64(timestamp.(float64)), 0)
		date := dateTime.Format("2/1/2006 15:04:05")
		month := int(dateTime.Month())
		day := dateTime.Day()
		hour := dateTime.Hour()

		ccuFloat, _ := strconv.ParseFloat(fmt.Sprintf("%v", value[1]), 64)

		f.SetCellValue("RawData", fmt.Sprintf("A%d", i+2), timestamp)
		f.SetCellValue("RawData", fmt.Sprintf("B%d", i+2), date)
		f.SetCellValue("RawData", fmt.Sprintf("C%d", i+2), month)
		f.SetCellValue("RawData", fmt.Sprintf("D%d", i+2), day)
		f.SetCellValue("RawData", fmt.Sprintf("E%d", i+2), hour)
		f.SetCellValue("RawData", fmt.Sprintf("F%d", i+2), ccuFloat)
	}

	pivotSheetName := "PivotTable"
	f.NewSheet(pivotSheetName)

	pivotTable := excelize.PivotTableOptions{
		DataRange:       "RawData!A1:F" + strconv.Itoa(len(responseData.Data.Result[0].Values)+1),
		PivotTableRange: pivotSheetName + "!A1:H1",
		Rows: []excelize.PivotTableField{
			{Data: "Tháng", Name: "Tháng"},
			{Data: "Ngày", Name: "Ngày"},
			{Data: "Giờ", Name: "Giờ"},
		},
		Data: []excelize.PivotTableField{
			{Data: "CCu", Name: "MAX of CCu", Subtotal: "Max"},
		},
	}

	// Additional computed columns and totals
	f.SetCellValue(pivotSheetName, "E1", "Số Block CCU")
	f.SetCellValue(pivotSheetName, "F1", "Số Block CCU tính phí")
	f.SetCellValue(pivotSheetName, "G1", "Chi phí")

	f.SetCellValue(pivotSheetName, "K1", "Tổng")
	f.SetCellValue(pivotSheetName, "L1", "Số Block CCU tính phí")
	f.SetCellValue(pivotSheetName, "M1", "Chi phí")

	for i := 0; i <= len(responseData.Data.Result[0].Values); i++ {
		f.SetCellFormula(pivotSheetName, fmt.Sprintf("E%d", i+2), fmt.Sprintf("CEILING(D%d/100, 1)", i+2))
		f.SetCellFormula(pivotSheetName, fmt.Sprintf("F%d", i+2), fmt.Sprintf("IF(E%d-1 > 0, E%d-1, 0)", i+2, i+2))
		f.SetCellFormula(pivotSheetName, fmt.Sprintf("G%d", i+2), fmt.Sprintf("F%d*3000", i+2))
	}
	f.SetCellFormula(pivotSheetName, "L2", fmt.Sprintf("SUM(F2:F%d)", len(responseData.Data.Result[0].Values)+1))
	f.SetCellFormula(pivotSheetName, "M2", fmt.Sprintf("SUM(G2:G%d)", len(responseData.Data.Result[0].Values)+1))

	if err := f.AddPivotTable(&pivotTable); err != nil {
		return nil, "", "", err
	}

	f.SetActiveSheet(index)

	filename := fmt.Sprintf("report_%s_%02d_%d.xlsx", dbname, start.Month(), start.Year())
	if err := f.SaveAs(filename); err != nil {
		return nil, "", "", err
	}

	return f, filename, pivotSheetName, nil
}

func readTotals(f *excelize.File, pivotSheetName string) (string, string) {
	totalBlockCCU, err := f.GetCellValue(pivotSheetName, "L2")
	if err != nil {
		totalBlockCCU = ""
	}
	totalCost, err := f.GetCellValue(pivotSheetName, "M2")
	if err != nil {
		totalCost = ""
	}
	return totalBlockCCU, totalCost
}

func buildTelegramCaption(customerName, siteName, totalCost, totalBlockCCU string) string {
	return fmt.Sprintf(
		"Báo cáo CCU tháng cho khách hàng: %s\nSite: %s\nTổng chi phí: %s VND\nSố block CCU tính phí: %s",
		customerName, siteName, totalCost, totalBlockCCU,
	)
}

func sendReportToTelegram(filename, caption string) error {
	telegramBotToken := "248926246:AAETwv7hzpk8zv6j9aRHDITYcnIRUGydS80"
	chatID := "-1002960560112"
	telegramURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendDocument", telegramBotToken)

	reportFile, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer reportFile.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("document", filepath.Base(filename))
	if err != nil {
		return err
	}
	if _, err = io.Copy(part, reportFile); err != nil {
		return err
	}
	_ = writer.WriteField("chat_id", chatID)
	_ = writer.WriteField("caption", caption)
	_ = writer.WriteField("parse_mode", "HTML")
	_ = writer.WriteField("message_thread_id", "2")
	writer.Close()

	req, err := http.NewRequest("POST", telegramURL, body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("telegram API error: %s", string(respBody))
	}
	return nil
}
func (w *MonthlyCCUReportWorker) Run() {
	log.Println("MonthlyCCUReportWorker: run started")
	moodles, err := w.mongoRepo.GetAll("")
	if err != nil {
		log.Printf("MonthlyCCUReportWorker: failed to get moodles: %v", err)
		return
	}
	for _, m := range moodles {
		log.Printf("MonthlyCCUReportWorker: found moodle: %+v", m)

		// For each moodle, set dbname and call to metric portal
		dbname := strings.ReplaceAll(m.Id, "-", "_")

		start, end := calculatePreviousMonthRange(time.Now())
		log.Printf("MonthlyCCUReportWorker: time range for %s -> %s - %s", dbname, start.Format(time.RFC3339), end.Format(time.RFC3339))

		req, err := buildPrometheusRequest(dbname, start, end)
		if err != nil {
			log.Printf("MonthlyCCUReportWorker: error creating request for moodle %s: %v", m.Id, err)
			continue
		}

		body, err := doHTTPRequest(req)
		if err != nil {
			log.Printf("MonthlyCCUReportWorker: error sending request for moodle %s: %v", m.Id, err)
			continue
		}
		log.Printf("MonthlyCCUReportWorker: received response for %s, size=%d bytes", dbname, len(body))

		responseData, err := parsePrometheusResponse(body)
		if err != nil {
			log.Printf("MonthlyCCUReportWorker: failed to unmarshal response for moodle %s: %v", m.Id, err)
			continue
		}

		if len(responseData.Data.Result) > 0 && len(responseData.Data.Result[0].Values) > 0 {
			log.Printf("MonthlyCCUReportWorker: %s has %d samples", dbname, len(responseData.Data.Result[0].Values))
			f, filename, pivotSheetName, err := generateExcelAndSave(responseData, dbname, start)
			if err != nil {
				log.Printf("MonthlyCCUReportWorker: failed to generate/save Excel for moodle %s: %v", m.Id, err)
				continue
			}
			log.Printf("MonthlyCCUReportWorker: report saved for moodle %s: %s", m.Id, filename)

			totalBlockCCU, totalCost := readTotals(f, pivotSheetName)
			caption := buildTelegramCaption(m.Name, m.WebSiteName, totalCost, totalBlockCCU)

			log.Printf("MonthlyCCUReportWorker: sending report to Telegram for %s", dbname)
			if err := sendReportToTelegram(filename, caption); err != nil {
				log.Printf("MonthlyCCUReportWorker: failed to send telegram request: %v", err)
				continue
			}
			log.Printf("MonthlyCCUReportWorker: report sent to Telegram for moodle %s", m.Id)
		}
	}

	// You can further process buf.String() as needed, e.g., unmarshal JSON, etc.

}
