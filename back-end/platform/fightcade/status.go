package fightcade

import (
	"bytes"
	"compress/zlib"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

const statusAPI = "https://web.fightcade.com/fc2status/api/"

type statusPayload struct {
	Req        string `json:"req"`
	Token      string `json:"token"`
	UserStatus string `json:"userstatus"`
	UUID       string `json:"uuid"`
	GUID       string `json:"guid"`
	HUID       string `json:"huid"`
	Version    string `json:"version"`
	Hash       string `json:"hash"`
}

func reportStatus(token string) {
	uid := getOrCreateUID()
	guid := uid
	huid := md5hex(uid)
	status := "stcable"
	version := "6"

	hash := computeHash(token, uid, guid, version, status, huid)

	payload := statusPayload{
		Req:        "userstatus",
		Token:      token,
		UserStatus: status,
		UUID:       uid,
		GUID:       guid,
		HUID:       huid,
		Version:    version,
		Hash:       hash,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		log.Printf("[fightcade] reportStatus: marshal error: %v", err)
		return
	}

	req, err := http.NewRequest("POST", statusAPI, bytes.NewReader(body))
	if err != nil {
		log.Printf("[fightcade] reportStatus: request error: %v", err)
		return
	}
	req.Header.Set("User-Agent", "fcade")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[fightcade] reportStatus: POST error: %v", err)
		return
	}
	resp.Body.Close()
	log.Printf("[fightcade] reportStatus: done (status=%d)", resp.StatusCode)
}

func computeHash(token, uid, guid, version, status, huid string) string {
	val := "3jedoQ" + token + "qmkq0" + uid + "dsnds" + guid + "sec or" + version + "2jden3" + status + "llNjjha" + huid
	return md5hex(val)
}

func md5hex(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}

func getOrCreateUID() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return uuid.New().String()
	}
	path := filepath.Join(home, ".fcuid")

	data, err := os.ReadFile(path)
	if err == nil {
		if uid := readCompressedUID(data); uid != "" {
			return uid
		}
	}

	uid := uuid.New().String()
	if compressed, err := zlibCompress(uid); err == nil {
		_ = os.WriteFile(path, compressed, 0600)
	}
	return uid
}

func readCompressedUID(data []byte) string {
	r, err := zlib.NewReader(bytes.NewReader(data))
	if err != nil {
		return ""
	}
	defer r.Close()
	raw, err := io.ReadAll(r)
	if err != nil {
		return ""
	}
	s := string(raw)
	if len(s) >= 36 {
		s = s[:36]
	}
	if _, err := uuid.Parse(s); err != nil {
		return ""
	}
	return s
}

func zlibCompress(s string) ([]byte, error) {
	var buf bytes.Buffer
	w := zlib.NewWriter(&buf)
	_, err := w.Write([]byte(s))
	if err != nil {
		return nil, err
	}
	if err := w.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
