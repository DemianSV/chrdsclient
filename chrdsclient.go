// Charybdis Monitoring System client v1.0.1

package chrdsclient

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"
)

type ConfT struct {
	ClientInSecureSkipVerify bool
	SpaceID                  string
	ModuleID                 string
	DataManagerURL           []string
}

var Conf ConfT

func statusCheck() (int, bool) {
	var index int
	status := Status()
	availability := false
	for i, item := range status {
		if item {
			index = i
			availability = true
			break
		}
	}
	return index, availability
}

// Логирование в систему CHRDS
func Log(metric string, value string) error {
	type RquestRawTextSaveT struct {
		SpaceID   string `json:"spaceid"`
		Metric    string `json:"metric"`
		Value     string `json:"value"`
		EventTime int64  `json:"eventtime"`
		Object    string `json:"object"`
	}
	type RequestRawTextSaveAT struct {
		Data []RquestRawTextSaveT `json:"data"`
	}

	var requestRawTextSave RquestRawTextSaveT
	var requestRawTextSaveA RequestRawTextSaveAT

	requestRawTextSave.SpaceID = Conf.SpaceID // Пространство UIManager
	requestRawTextSave.EventTime = MakeTimestamp()
	requestRawTextSave.Metric = metric
	requestRawTextSave.Value = value
	requestRawTextSave.Object = "uimanager"
	requestRawTextSaveA.Data = append(requestRawTextSaveA.Data, requestRawTextSave)

	index, availability := statusCheck()

	if availability {
		if len(requestRawTextSaveA.Data) > 0 {
			requestRawTextSaveJSON, err := json.Marshal(requestRawTextSaveA)
			if err != nil {
				return errors.New("JSON MARSHAL ERROR (" + err.Error() + ")")
			} else {
				req, err := http.NewRequest("POST", Conf.DataManagerURL[index]+"/api/v1/text/save", bytes.NewBuffer(requestRawTextSaveJSON))
				if err != nil {
					return errors.New("HTTP REQUEST ERROR (" + err.Error() + ")")
				}
				req.Header.Set("Content-Type", "application/json; charset=UTF-8")
				req.Header.Set("api_key", Conf.ModuleID) // Модуль UIManager

				client := &http.Client{}
				if strings.Contains(strings.ToUpper(Conf.DataManagerURL[index]), strings.ToUpper("https:")) {
					transport := &http.Transport{
						TLSClientConfig: &tls.Config{InsecureSkipVerify: Conf.ClientInSecureSkipVerify},
					}
					client = &http.Client{Transport: transport}
				}
				resp, err := client.Do(req)
				if err != nil {
					return errors.New("HTTP REQUEST ERROR (" + err.Error() + ")")
				}
				defer resp.Body.Close()

				return nil
			}
		} else {
			return errors.New("NO DATA TO SEND")
		}
	}
	return errors.New("DATAMANAGER IS NOT AVAILABLE")
}

func Status() []bool {
	var responseUp []bool

	for _, item := range Conf.DataManagerURL {
		req, err := http.NewRequest("GET", item+"/api/v1/version", nil)
		if err != nil {
			responseUp = append(responseUp, false)
		} else {
			req.Header.Set("Content-Type", "application/json; charset=UTF-8")
			req.Header.Set("api_key", Conf.ModuleID) // Модуль UIManager

			client := &http.Client{}
			if strings.Contains(strings.ToUpper(item), strings.ToUpper("https:")) {
				transport := &http.Transport{
					TLSClientConfig: &tls.Config{InsecureSkipVerify: Conf.ClientInSecureSkipVerify},
				}
				client = &http.Client{Transport: transport, Timeout: 3 * time.Second}
			}
			resp, err := client.Do(req)
			if err != nil {
				responseUp = append(responseUp, false)
			} else {
				if resp.StatusCode == 200 {
					responseUp = append(responseUp, true)
				} else {
					responseUp = append(responseUp, false)
				}
			}
		}
	}

	return responseUp
}

func MakeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func Metric(metric string, value float32) error {
	type RequestRawDataSaveT struct {
		SpaceID   string  `json:"spaceid"`
		Metric    string  `json:"metric"`
		Value     float32 `json:"value"`
		EventTime int64   `json:"eventtime"`
		Object    string  `json:"object"`
	}
	type RequestRawDataSaveAT struct {
		Data []RequestRawDataSaveT `json:"data"`
	}

	var requestRawDataSave RequestRawDataSaveT
	var requestRawDataSaveArray RequestRawDataSaveAT

	requestRawDataSave.SpaceID = Conf.SpaceID // Пространство UIManager
	requestRawDataSave.EventTime = MakeTimestamp()
	requestRawDataSave.Metric = metric
	requestRawDataSave.Value = value
	requestRawDataSave.Object = "uimanager"
	requestRawDataSaveArray.Data = append(requestRawDataSaveArray.Data, requestRawDataSave)

	index, availability := statusCheck()

	if availability {
		if len(requestRawDataSaveArray.Data) > 0 {
			requestRawDataSaveJSON, err := json.Marshal(requestRawDataSaveArray)
			if err != nil {
				return errors.New("JSON MARSHAL ERROR (" + err.Error() + ")")
			} else {
				req, err := http.NewRequest("POST", Conf.DataManagerURL[index]+"/api/v1/data/save", bytes.NewBuffer(requestRawDataSaveJSON))
				if err != nil {
					return errors.New("HTTP REQUEST ERROR (" + err.Error() + ")")
				}
				req.Header.Set("Content-Type", "application/json; charset=UTF-8")
				req.Header.Set("api_key", Conf.ModuleID) // Модуль UIManager

				client := &http.Client{}
				if strings.Contains(strings.ToUpper(Conf.DataManagerURL[index]), strings.ToUpper("https:")) {
					transport := &http.Transport{
						TLSClientConfig: &tls.Config{InsecureSkipVerify: Conf.ClientInSecureSkipVerify},
					}
					client = &http.Client{Transport: transport}
				}
				resp, err := client.Do(req)
				if err != nil {
					return errors.New("HTTP REQUEST ERROR (" + err.Error() + ")")
				}
				defer resp.Body.Close()

				return nil
			}
		} else {
			return errors.New("NO DATA TO SEND")
		}
	}
	return errors.New("DATAMANAGER IS NOT AVAILABLE")
}
