package worker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"moodle/internal/core/domain"
	repo "moodle/internal/repository/moodle"
	"moodle/pkg/mongodbiface"
)

type AlarmNFS struct {
	mongoRepo *repo.MongoDB
}

func NewAlarmNFSWorker(mongo mongodbiface.DB) *AlarmNFS {
	mongoRepo := repo.NewMongoDB(mongo)
	return &AlarmNFS{
		mongoRepo: mongoRepo,
	}
}

// Thêm struct và hàm lấy dữ liệu nfs_tracking
type NFSTracking struct {
	Path string `bson:"path" json:"path"`
	Size int64  `bson:"size" json:"size"`
}

type NFSTrackingDoc struct {
	Data []NFSTracking `bson:"data" json:"data"`
}

// Thêm struct cho extend_storage
type ExtendStorage struct {
	MoodleID      string `bson:"moodleId"`
	ExtendStorage string `bson:"extend_storage"`
}

func (a *AlarmNFS) Run() {
	// Lấy toàn bộ dữ liệu một lần
	moodles, err := a.mongoRepo.GetAll("")
	if err != nil {
		log.Printf("Failed to get moodles: %v", err)
		return
	}

	nfsMap, err := a.getNFSTrackingMap()
	if err != nil {
		log.Printf("Failed to get NFS tracking: %v", err)
		return
	}

	// Xử lý so sánh và gửi cảnh báo
	for _, m := range moodles {
		a.processMoodle(m, nfsMap)
	}
}

func (a *AlarmNFS) getNFSTrackingMap() (map[string]int64, error) {
	nfsTracking, err := a.mongoRepo.GetNFSTracking()
	if err != nil {
		return nil, err
	}

	result := make(map[string]int64)
	for _, d := range nfsTracking.Data {
		if len(d.Path) >= 49 {
			moodleid := d.Path[13:49]
			result[moodleid] = int64(d.Size)
		}
	}
	return result, nil
}

// Hàm lấy dung lượng mua thêm (cùng đơn vị với DocumentStorage)
func (a *AlarmNFS) getExtendStorage(moodleId string) (int64, error) {
	extra, err := a.mongoRepo.GetExtendStorage(moodleId)
	if err != nil {
		return 0, nil // lỗi coi như 0 để tránh spam lỗi
	}
	return int64(extra), nil
}

func (a *AlarmNFS) sendTelegram(m domain.Moodle, currentStorage, totalLimit int64) error {
	fmt.Println(currentStorage, totalLimit)
	apiUrl := "https://tlg.dev.bizflycloud.vn/bot248926246:AAETwv7hzpk8zv6j9aRHDITYcnIRUGydS80/sendMessage"
	text := fmt.Sprintf("[ALARM] KH: %s, Email: %s, Site: %s, Storage: %.2fGB/%.2fGB", m.Name, m.Email, m.WebSiteName, float64(currentStorage), float64(totalLimit))
	message := map[string]interface{}{
		"chat_id": -4116328390,
		"text":    text,
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

func (a *AlarmNFS) getCurrentWeekYear() (int, int) {
	now := time.Now()
	year, week := now.ISOWeek()
	return week, year
}

func (a *AlarmNFS) processMoodle(m domain.Moodle, nfsMap map[string]int64) {
	current, ok := nfsMap[m.Id]
	if !ok {
		log.Printf("skip %s: not found in nfs_tracking", m.Id)
		return
	}
	current = current / 1000000
	defaultLimit := int64(m.Packages.DocumentStorage)
	if current > defaultLimit {
		// Nếu vượt, kiểm tra extend_storage
		extBytes, _ := a.getExtendStorage(m.Id)
		totalLimit := defaultLimit + extBytes
		if current > totalLimit {
			// Kiểm tra xem đã alarm tuần này chưa
			week, year := a.getCurrentWeekYear()
			alreadySent, err := a.mongoRepo.CheckStorageAlarmSent(m.Id, week, year)
			if err != nil {
				log.Printf("check alarm sent failed for %s: %v", m.Id, err)
				return
			}

			if !alreadySent {
				// Gửi alarm và lưu vào database
				if err := a.sendTelegram(m, current, totalLimit); err != nil {
					log.Printf("send tele fail: %v", err)
					return
				}

				// Lưu alarm tracking
				alarm := domain.StorageAlarmTracking{
					MoodleId:  m.Id,
					Week:      week,
					Year:      year,
					CreatedAt: time.Now(),
				}

				if err := a.mongoRepo.CreateStorageAlarmTracking(alarm); err != nil {
					log.Printf("save alarm tracking failed for %s: %v", m.Id, err)
				} else {
					log.Printf("Alarm sent and saved for moodle %s, week %d", m.Id, week)
				}
			} else {
				log.Printf("Alarm already sent for moodle %s in week %d, skipping", m.Id, week)
			}
		}
	}
}
