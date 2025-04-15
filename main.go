/*
 *  This file is part of BoxBox-Server (https://github.com/BrightDV/BoxBox-Server).
 *
 * BoxBox-Server is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * BoxBox-Server is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with BoxBox-Server.  If not, see <http://www.gnu.org/licenses/>.
 *
 * Copyright (c) 2022-2024, BrightDV
 */

package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/goenning/go-cache-demo/cache"
	"github.com/goenning/go-cache-demo/cache/memory"
	"github.com/gorilla/mux"
)

func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Checkout the Android app on GitHub at https://github.com/BrightDV/BoxBox")
}

var tmp, _ = base64.StdEncoding.DecodeString("eFo3QU9PRFNqaVFhZExzSVlXZWZRcnBDU1FWRGJIR0M=")
var apikey = string(tmp)
var tmp_, _ = base64.StdEncoding.DecodeString("aHR0cHM6Ly9hcGkuZm9ybXVsYTEuY29tLw==")
var f1Endpoint = string(tmp_)
var tmp__, _ = base64.StdEncoding.DecodeString("aHR0cHM6Ly9hcGkuZm9ybXVsYS1lLnB1bHNlbGl2ZS5jb20v")
var fEEndpoint = string(tmp__)
var DOMAIN string = "*"
var PORT string = "8080"

var storage cache.Storage

type Formula1 struct{}
type FormulaE struct{}

func init() {
	storage = memory.NewStorage()
}

func logger(uri string) {
	fmt.Println("["+time.Now().Format("2006-01-02 15:04:05")+"]", " GET "+uri)
}

func cached(duration string, contentType string, handler func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		content := storage.Get(r.RequestURI)
		if content != nil {
			logger(r.RequestURI)
			w.Header().Add("Content-Type", contentType)
			w.Header().Set("Access-Control-Allow-Origin", DOMAIN)
			w.Write(content)
		} else {
			c := httptest.NewRecorder()
			handler(c, r)
			for k, v := range c.Header() {
				w.Header()[k] = v
			}
			w.WriteHeader(c.Code)
			content := c.Body.Bytes()
			if d, err := time.ParseDuration(duration); err == nil {
				storage.Set(r.RequestURI, content, d)
			} else {
				fmt.Printf("Page not cached. err: %s\n", err)
			}
			w.Write(content)
		}
	})
}

func (Formula1) getArticle(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", DOMAIN)
	var article any
	articleId := mux.Vars(r)["articleId"]
	logger(r.RequestURI)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", f1Endpoint+"v1/editorial/articles/"+articleId, nil)
	req.Header.Set("User-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36")
	req.Header.Set("apikey", apikey)
	req.Header.Set("locale", "en")
	resp, _ := client.Do(req)
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &article)
	json.NewEncoder(w).Encode(article)
}

func (Formula1) getArticles(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", DOMAIN)
	var articles any
	offset := r.URL.Query().Get("offset")
	if offset == "" {
		offset = "0"
	}
	tags := r.URL.Query().Get("tags")
	articleTypes := r.URL.Query().Get("articleTypes")
	logger(r.RequestURI)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", f1Endpoint+"v1/editorial/articles?limit=16&tags="+tags+"&offset="+offset+"&articleTypes="+articleTypes, nil)
	req.Header.Set("User-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36")
	req.Header.Set("apikey", apikey)
	req.Header.Set("locale", "en")
	resp, _ := client.Do(req)
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &articles)
	json.NewEncoder(w).Encode(articles)
}

func (Formula1) getVideo(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", DOMAIN)
	var article any
	videoId := mux.Vars(r)["videoId"]
	logger(r.RequestURI)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", f1Endpoint+"v1/video-assets/videos/"+videoId, nil)
	req.Header.Set("User-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36")
	req.Header.Set("apikey", apikey)
	req.Header.Set("locale", "en")
	resp, _ := client.Do(req)
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &article)
	json.NewEncoder(w).Encode(article)
}

func (Formula1) getVideos(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", DOMAIN)
	var articles any
	limit := r.URL.Query().Get("limit")
	if limit == "" {
		limit = "16"
	}
	offset := r.URL.Query().Get("offset")
	if offset == "" {
		offset = "0"
	}
	tags := r.URL.Query().Get("tags")
	logger(r.RequestURI)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", f1Endpoint+"v1/video-assets/videos?limit="+limit+"&tag="+tags+"&offset="+offset, nil)
	req.Header.Set("User-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36")
	req.Header.Set("apikey", apikey)
	req.Header.Set("locale", "en")
	resp, _ := client.Do(req)
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &articles)
	json.NewEncoder(w).Encode(articles)
}

func (Formula1) eventTracker(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", DOMAIN)
	var event any
	logger(r.RequestURI)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", f1Endpoint+"v1/event-tracker", nil)
	req.Header.Set("User-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36")
	req.Header.Set("apikey", apikey)
	req.Header.Set("locale", "en")
	req.Header.Set("Accept", "application/json")
	resp, _ := client.Do(req)
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &event)
	json.NewEncoder(w).Encode(event)
}

func (Formula1) eventTrackerForOneMeeting(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", DOMAIN)
	var event any
	meetingId := mux.Vars(r)["meetingId"]
	logger(r.RequestURI)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", f1Endpoint+"v1/event-tracker/meeting/"+meetingId, nil)
	req.Header.Set("User-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36")
	req.Header.Set("apikey", apikey)
	req.Header.Set("locale", "en")
	resp, _ := client.Do(req)
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &event)
	json.NewEncoder(w).Encode(event)
}

func (Formula1) getRaceResults(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", DOMAIN)
	var results any
	meetingId := mux.Vars(r)["meetingId"]
	logger(r.RequestURI)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", f1Endpoint+"v1/fom-results/race?meeting="+meetingId, nil)
	req.Header.Set("User-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36")
	req.Header.Set("apikey", apikey)
	req.Header.Set("locale", "en")
	resp, _ := client.Do(req)
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &results)
	json.NewEncoder(w).Encode(results)
}

func (Formula1) getQualificationResults(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", DOMAIN)
	var results any
	meetingId := mux.Vars(r)["meetingId"]
	logger(r.RequestURI)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", f1Endpoint+"v1/fom-results/qualifying?meeting="+meetingId, nil)
	req.Header.Set("User-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36")
	req.Header.Set("apikey", apikey)
	req.Header.Set("locale", "en")
	resp, _ := client.Do(req)
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &results)
	json.NewEncoder(w).Encode(results)
}

func (Formula1) getFreePracticeResults(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", DOMAIN)
	var results any
	meetingId := mux.Vars(r)["meetingId"]
	sessionId := mux.Vars(r)["sessionId"]
	logger(r.RequestURI)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", f1Endpoint+"v1/fom-results/practice?meeting="+meetingId+"&session="+sessionId, nil)
	req.Header.Set("User-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36")
	req.Header.Set("apikey", apikey)
	req.Header.Set("locale", "en")
	resp, _ := client.Do(req)
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &results)
	json.NewEncoder(w).Encode(results)
}

func (Formula1) getSprintQualifyingResults(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", DOMAIN)
	var results any
	meetingId := mux.Vars(r)["meetingId"]
	logger(r.RequestURI)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", f1Endpoint+"v1/fom-results/sprint-shootout?meeting="+meetingId, nil)
	req.Header.Set("User-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36")
	req.Header.Set("apikey", apikey)
	req.Header.Set("locale", "en")
	resp, _ := client.Do(req)
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &results)
	json.NewEncoder(w).Encode(results)
}

func (Formula1) getSprintResults(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", DOMAIN)
	var results any
	meetingId := mux.Vars(r)["meetingId"]
	logger(r.RequestURI)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", f1Endpoint+"v1/fom-results/sprint?meeting="+meetingId, nil)
	req.Header.Set("User-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36")
	req.Header.Set("apikey", apikey)
	req.Header.Set("locale", "en")
	resp, _ := client.Do(req)
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &results)
	json.NewEncoder(w).Encode(results)
}

func (Formula1) getStartingGrid(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", DOMAIN)
	var grid any
	meetingId := mux.Vars(r)["meetingId"]
	logger(r.RequestURI)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", f1Endpoint+"v1/fom-results/starting-grid?meeting="+meetingId, nil)
	req.Header.Set("User-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36")
	req.Header.Set("apikey", apikey)
	req.Header.Set("locale", "en")
	resp, _ := client.Do(req)
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &grid)
	json.NewEncoder(w).Encode(grid)
}

func (Formula1) getRaceResultsV2(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", DOMAIN)
	var results any
	meetingId := mux.Vars(r)["meetingId"]
	logger(r.RequestURI)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", f1Endpoint+"v2/fom-results/race?meeting="+meetingId, nil)
	req.Header.Set("User-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36")
	req.Header.Set("apikey", apikey)
	req.Header.Set("locale", "en")
	resp, _ := client.Do(req)
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &results)
	json.NewEncoder(w).Encode(results)
}

func (Formula1) getQualificationResultsV2(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", DOMAIN)
	var results any
	meetingId := mux.Vars(r)["meetingId"]
	logger(r.RequestURI)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", f1Endpoint+"v2/fom-results/qualifying?meeting="+meetingId, nil)
	req.Header.Set("User-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36")
	req.Header.Set("apikey", apikey)
	req.Header.Set("locale", "en")
	resp, _ := client.Do(req)
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &results)
	json.NewEncoder(w).Encode(results)
}

func (Formula1) getFreePracticeResultsV2(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", DOMAIN)
	var results any
	meetingId := mux.Vars(r)["meetingId"]
	sessionId := mux.Vars(r)["sessionId"]
	logger(r.RequestURI)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", f1Endpoint+"v2/fom-results/practice?meeting="+meetingId+"&session="+sessionId, nil)
	req.Header.Set("User-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36")
	req.Header.Set("apikey", apikey)
	req.Header.Set("locale", "en")
	resp, _ := client.Do(req)
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &results)
	json.NewEncoder(w).Encode(results)
}

func (Formula1) getSprintQualifyingResultsV2(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", DOMAIN)
	var results any
	meetingId := mux.Vars(r)["meetingId"]
	logger(r.RequestURI)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", f1Endpoint+"v2/fom-results/sprint-shootout?meeting="+meetingId, nil)
	req.Header.Set("User-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36")
	req.Header.Set("apikey", apikey)
	req.Header.Set("locale", "en")
	resp, _ := client.Do(req)
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &results)
	json.NewEncoder(w).Encode(results)
}

func (Formula1) getSprintResultsV2(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", DOMAIN)
	var results any
	meetingId := mux.Vars(r)["meetingId"]
	logger(r.RequestURI)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", f1Endpoint+"v2/fom-results/sprint?meeting="+meetingId, nil)
	req.Header.Set("User-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36")
	req.Header.Set("apikey", apikey)
	req.Header.Set("locale", "en")
	resp, _ := client.Do(req)
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &results)
	json.NewEncoder(w).Encode(results)
}

func (Formula1) getStartingGridV2(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", DOMAIN)
	var grid any
	meetingId := mux.Vars(r)["meetingId"]
	logger(r.RequestURI)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", f1Endpoint+"v2/fom-results/starting-grid?meeting="+meetingId, nil)
	req.Header.Set("User-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36")
	req.Header.Set("apikey", apikey)
	req.Header.Set("locale", "en")
	resp, _ := client.Do(req)
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &grid)
	json.NewEncoder(w).Encode(grid)
}

func (Formula1) getDriverStandings(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", DOMAIN)
	var standings any
	logger(r.RequestURI)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", f1Endpoint+"v1/editorial-driverlisting/listing", nil)
	req.Header.Set("User-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36")
	req.Header.Set("apikey", apikey)
	req.Header.Set("locale", "en")
	resp, _ := client.Do(req)
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &standings)
	json.NewEncoder(w).Encode(standings)
}

func (Formula1) getTeamStandings(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", DOMAIN)
	var standings any
	logger(r.RequestURI)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", f1Endpoint+"v1/editorial-constructorlisting/listing", nil)
	req.Header.Set("User-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36")
	req.Header.Set("apikey", apikey)
	req.Header.Set("locale", "en")
	resp, _ := client.Do(req)
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &standings)
	json.NewEncoder(w).Encode(standings)
}

func (Formula1) getSchedule(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", DOMAIN)
	var schedule any
	logger(r.RequestURI)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", f1Endpoint+"v1/editorial-eventlisting/events", nil)
	req.Header.Set("User-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36")
	req.Header.Set("apikey", apikey)
	req.Header.Set("locale", "en")
	resp, _ := client.Do(req)
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &schedule)
	json.NewEncoder(w).Encode(schedule)
}

func (Formula1) getFinishedSessions(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", DOMAIN)
	year := mux.Vars(r)["year"]
	fomRaceId := mux.Vars(r)["fomRaceId"]
	raceName := mux.Vars(r)["raceName"]
	logger(r.RequestURI)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://www.formula1.com/en/results.html/"+year+"/races/"+fomRaceId+"/"+raceName+".html", nil)
	req.Header.Set("User-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36")
	resp, _ := client.Do(req)
	body, _ := io.ReadAll(resp.Body)
	fmt.Fprint(w, string(body))
}

func (Formula1) getResultsForScraping(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", DOMAIN)
	year := mux.Vars(r)["year"]
	fomRaceId := mux.Vars(r)["fomRaceId"]
	raceName := mux.Vars(r)["raceName"]
	session := mux.Vars(r)["session"]
	logger(r.RequestURI)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://www.formula1.com/en/results.html/"+year+"/races/"+fomRaceId+"/"+raceName+"/"+session+".html", nil)
	req.Header.Set("User-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36")
	resp, _ := client.Do(req)
	body, _ := io.ReadAll(resp.Body)
	fmt.Fprint(w, string(body))
}

func (Formula1) getCircuitDetails(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", DOMAIN)
	year := mux.Vars(r)["year"]
	circuitName := mux.Vars(r)["circuitName"]
	logger(r.RequestURI)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://www.formula1.com/en/racing/"+year+"/"+circuitName+"/Circuit.html", nil)
	req.Header.Set("User-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36")
	resp, _ := client.Do(req)
	body, _ := io.ReadAll(resp.Body)
	fmt.Fprint(w, string(body))
}

func (Formula1) getDriverDetails(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", DOMAIN)
	driverId := mux.Vars(r)["driverId"]
	logger(r.RequestURI)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://www.formula1.com/en/drivers/"+driverId+".html", nil)
	req.Header.Set("User-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36")
	resp, _ := client.Do(req)
	body, _ := io.ReadAll(resp.Body)
	fmt.Fprint(w, string(body))
}

func (Formula1) getTeamDetails(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", DOMAIN)
	teamId := mux.Vars(r)["teamId"]
	logger(r.RequestURI)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://www.formula1.com/en/teams/"+teamId+".html", nil)
	req.Header.Set("User-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36")
	resp, _ := client.Do(req)
	body, _ := io.ReadAll(resp.Body)
	fmt.Fprint(w, string(body))
}

func (Formula1) getHallOfFame(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", DOMAIN)
	logger(r.RequestURI)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://www.formula1.com/en/drivers/hall-of-fame", nil)
	req.Header.Set("User-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36")
	resp, _ := client.Do(req)
	body, _ := io.ReadAll(resp.Body)
	fmt.Fprint(w, string(body))
}

func (Formula1) getHallOfFameDriverDetails(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", DOMAIN)
	driver := mux.Vars(r)["driver"]
	id := mux.Vars(r)["id"]
	logger(r.RequestURI)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://www.formula1.com/en/information/drivers-hall-of-fame-"+driver+"."+id, nil)
	req.Header.Set("User-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36")
	resp, _ := client.Do(req)
	body, _ := io.ReadAll(resp.Body)
	fmt.Fprint(w, string(body))
}

func (Formula1) getSessionDocuments(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", DOMAIN)
	logger(r.RequestURI)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://www.fia.com/documents/championships/fia-formula-one-world-championship-14/season/season-2025-2071", nil)
	req.Header.Set("User-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36")
	resp, _ := client.Do(req)
	body, _ := io.ReadAll(resp.Body)
	fmt.Fprint(w, string(body))
}

func (Formula1) getSessionDocument(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/pdf")
	w.Header().Set("Access-Control-Allow-Origin", DOMAIN)
	documentPath := mux.Vars(r)["documentPath"]
	logger(r.RequestURI)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://www.fia.com/system/files/decision-document/"+documentPath, nil)
	req.Header.Set("User-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36")
	resp, _ := client.Do(req)
	body, _ := io.ReadAll(resp.Body)
	fmt.Fprint(w, string(body))
}

func (Formula1) getRssFeed(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/xml; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", DOMAIN)
	languageCode := mux.Vars(r)["languageCode"]
	if languageCode == "motorsport" {
		languageCode = "br"
	}
	logger(r.RequestURI)
	client := &http.Client{}
	customFeedUrls := map[string]string{
		"fr":  "https://fr.motorsport.com",
		"es":  "https://es.motorsport.com",
		"br":  "https://motorsport.uol.com.br",
		"de":  "https://de.motorsport.com",
		"it":  "https://it.motorsport.com",
		"ru":  "https://ru.motorsport.com",
		"cn":  "https://cn.motorsport.com",
		"hu":  "https://hu.motorsport.com",
		"id":  "https://id.motorsport.com",
		"jp":  "https://jp.motorsport.com",
		"nl":  "https://nl.motorsport.com",
		"tr":  "https://tr.motorsport.com",
		"us":  "https://us.motorsport.com",
		"lat": "https://lat.motorsport.com",
		"ch":  "https://ch.motorsport.com",
		"au":  "https://au.motorsport.com",
		"pl":  "https://pl.motorsport.com",
	}
	req, _ := http.NewRequest("GET", customFeedUrls[languageCode]+"/rss/f1/news/", nil)
	req.Header.Set("User-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36")
	resp, _ := client.Do(req)
	body, _ := io.ReadAll(resp.Body)
	fmt.Fprint(w, string(body))
}

func (FormulaE) getArticle(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", DOMAIN)
	var article any
	articleId := mux.Vars(r)["articleId"]
	logger(r.RequestURI)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", fEEndpoint+"content/formula-e/text/EN/"+articleId, nil)
	req.Header.Set("User-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "application/json")
	resp, _ := client.Do(req)
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &article)
	json.NewEncoder(w).Encode(article)
}

func (FormulaE) getArticleAsHtml(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", DOMAIN)
	articleId := mux.Vars(r)["articleId"]
	logger(r.RequestURI)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://fiaformulae.com/en/news/"+articleId, nil)
	req.Header.Set("User-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36")
	resp, _ := client.Do(req)
	body, _ := io.ReadAll(resp.Body)
	fmt.Fprint(w, string(body))
}

func (FormulaE) getArticles(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", DOMAIN)
	var articles any
	page := r.URL.Query().Get("page")
	if page == "" {
		page = "0"
	}
	logger(r.RequestURI)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", fEEndpoint+"content/formula-e/text/EN/?page="+page+"&pageSize=16&tagNames=content-type%3Anews&tagExpression=&playlistTypeRestriction=&playlistId=&detail=&size=16&championshipId=&sort=", nil)
	req.Header.Set("User-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "application/json")
	resp, _ := client.Do(req)
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &articles)
	json.NewEncoder(w).Encode(articles)
}

func (FormulaE) getVideos(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", DOMAIN)
	var articles any
	limit := r.URL.Query().Get("limit")
	if limit == "" {
		limit = "16"
	}
	page := r.URL.Query().Get("page")
	if page == "" {
		page = "0"
	}
	logger(r.RequestURI)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", fEEndpoint+"content/formula-e/playlist/EN/15?page="+page+"&pageSize="+limit+"&detail=DETAILED&size="+limit, nil)
	req.Header.Set("User-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "application/json")
	resp, _ := client.Do(req)
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &articles)
	json.NewEncoder(w).Encode(articles)
}

func (FormulaE) getRaceDetails(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", DOMAIN)
	var event any
	raceId := mux.Vars(r)["raceId"]
	logger(r.RequestURI)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", fEEndpoint+"formula-e/v1/races/"+raceId, nil)
	req.Header.Set("User-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "application/json")
	resp, _ := client.Do(req)
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &event)
	json.NewEncoder(w).Encode(event)
}

func (FormulaE) getSessions(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", DOMAIN)
	var sessions any
	raceId := mux.Vars(r)["raceId"]
	logger(r.RequestURI)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", fEEndpoint+"formula-e/v1/races/"+raceId+"/sessions", nil)
	req.Header.Set("User-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "application/json")
	resp, _ := client.Do(req)
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &sessions)
	json.NewEncoder(w).Encode(sessions)
}

func (FormulaE) getRaceArticles(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", DOMAIN)
	var sessions any
	raceId := mux.Vars(r)["raceId"]
	logger(r.RequestURI)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", fEEndpoint+"content/formula-e/EN?contentTypes=video&contentTypes=news&page=0&pageSize=10&references=FORMULA_E_RACE:"+raceId+"&onlyRestrictedContent=false&detail=DETAILED", nil)
	req.Header.Set("User-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "application/json")
	resp, _ := client.Do(req)
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &sessions)
	json.NewEncoder(w).Encode(sessions)
}

func (FormulaE) getSessionResults(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", DOMAIN)
	var results any
	raceId := mux.Vars(r)["raceId"]
	sessionId := mux.Vars(r)["sessionId"]
	logger(r.RequestURI)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", fEEndpoint+"formula-e/v1/races/"+raceId+"/sessions/"+sessionId+"/results", nil)
	req.Header.Set("User-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "application/json")
	resp, _ := client.Do(req)
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &results)
	json.NewEncoder(w).Encode(results)
}

func (FormulaE) getCircuitImageDetails(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", DOMAIN)
	var results any
	raceId := mux.Vars(r)["raceId"]
	logger(r.RequestURI)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", fEEndpoint+"content/formula-e/photo/en/?references=FORMULA_E_RACE:"+raceId+"&tagNames=race:bg-image", nil)
	req.Header.Set("User-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "application/json")
	resp, _ := client.Do(req)
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &results)
	json.NewEncoder(w).Encode(results)
}

func (FormulaE) getDriverStandings(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", DOMAIN)
	var standings any
	championshipId := mux.Vars(r)["championshipId"]
	logger(r.RequestURI)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", fEEndpoint+"formula-e/v1/standings/drivers?championshipId="+championshipId, nil)
	req.Header.Set("User-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "application/json")
	resp, _ := client.Do(req)
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &standings)
	json.NewEncoder(w).Encode(standings)
}

func (FormulaE) getTeamStandings(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", DOMAIN)
	var standings any
	championshipId := mux.Vars(r)["championshipId"]
	logger(r.RequestURI)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", fEEndpoint+"formula-e/v1/standings/teams?championshipId="+championshipId, nil)
	req.Header.Set("User-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "application/json")
	resp, _ := client.Do(req)
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &standings)
	json.NewEncoder(w).Encode(standings)
}

func (FormulaE) getSchedule(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", DOMAIN)
	var schedule any
	championshipId := mux.Vars(r)["championshipId"]
	logger(r.RequestURI)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", fEEndpoint+"formula-e/v1/races?championshipId="+championshipId, nil)
	req.Header.Set("User-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "application/json")
	resp, _ := client.Do(req)
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &schedule)
	json.NewEncoder(w).Encode(schedule)
}

func (FormulaE) getLatestChampionship(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", DOMAIN)
	var schedule any
	logger(r.RequestURI)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", fEEndpoint+"formula-e/v1/championships/latest", nil)
	req.Header.Set("User-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "application/json")
	resp, _ := client.Do(req)
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &schedule)
	json.NewEncoder(w).Encode(schedule)
}

func main() {
	fmt.Println("Initializing...")
	route := "/"
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homeLink)
	// deprecated: not championship-specific
	router.Handle(route+"v1/editorial/articles", cached("30s", "application/json", Formula1{}.getArticles)).Methods("GET", "OPTIONS")
	router.Handle(route+"v1/editorial/articles/{articleId}", cached("5m", "application/json", Formula1{}.getArticle)).Methods("GET", "OPTIONS")
	router.Handle(route+"v1/video-assets/videos", cached("30s", "application/json", Formula1{}.getVideos)).Methods("GET", "OPTIONS")
	router.Handle(route+"v1/video-assets/videos/{videoId}", cached("5m", "application/json", Formula1{}.getVideo)).Methods("GET", "OPTIONS")
	router.Handle(route+"v1/event-tracker", cached("20s", "application/json", Formula1{}.eventTracker)).Methods("GET", "OPTIONS")
	router.Handle(route+"v1/event-tracker/meeting/{meetingId}", cached("20s", "application/json", Formula1{}.eventTrackerForOneMeeting)).Methods("GET", "OPTIONS")
	router.Handle(route+"en/results.html/{year}/races/{fomRaceId}/{raceName}.html", cached("120s", "application/json", Formula1{}.getFinishedSessions)).Methods("GET", "OPTIONS")
	router.Handle(route+"en/results.html/{year}/races/{fomRaceId}/{raceName}/{session}.html", cached("120s", "text/html; charset=utf-8", Formula1{}.getResultsForScraping)).Methods("GET", "OPTIONS")
	router.Handle(route+"en/racing/{year}/{circuitName}/Circuit.html", cached("168h", "text/html; charset=utf-8", Formula1{}.getCircuitDetails)).Methods("GET", "OPTIONS")
	router.Handle(route+"en/drivers/{driverId}.html", cached("1h", "text/html; charset=utf-8", Formula1{}.getDriverDetails)).Methods("GET", "OPTIONS")
	router.Handle(route+"en/teams/{teamId}.html", cached("1h", "text/html; charset=utf-8", Formula1{}.getTeamDetails)).Methods("GET", "OPTIONS")
	router.Handle(route+"documents", cached("30s", "text/html; charset=utf-8", Formula1{}.getSessionDocuments)).Methods("GET", "OPTIONS")
	router.HandleFunc(route+"documents/{documentPath}", Formula1{}.getSessionDocument).Methods("GET", "OPTIONS")
	router.HandleFunc(route+"rss/{languageCode}", Formula1{}.getRssFeed).Methods("GET", "OPTIONS")
	// F1 championship
	router.Handle(route+"f1/v1/editorial/articles", cached("30s", "application/json", Formula1{}.getArticles)).Methods("GET", "OPTIONS")
	router.Handle(route+"f1/v1/editorial/articles/{articleId}", cached("5m", "application/json", Formula1{}.getArticle)).Methods("GET", "OPTIONS")
	router.Handle(route+"f1/v1/video-assets/videos", cached("30s", "application/json", Formula1{}.getVideos)).Methods("GET", "OPTIONS")
	router.Handle(route+"f1/v1/video-assets/videos/{videoId}", cached("5m", "application/json", Formula1{}.getVideo)).Methods("GET", "OPTIONS")
	router.Handle(route+"f1/v1/event-tracker", cached("20s", "application/json", Formula1{}.eventTracker)).Methods("GET", "OPTIONS")
	router.Handle(route+"f1/v1/event-tracker/meeting/{meetingId}", cached("20s", "application/json", Formula1{}.eventTrackerForOneMeeting)).Methods("GET", "OPTIONS")
	router.Handle(route+"f1/v1/fom-results/race/meeting={meetingId}", cached("20s", "application/json", Formula1{}.getRaceResults)).Methods("GET", "OPTIONS")
	router.Handle(route+"f1/v1/fom-results/qualifying/meeting={meetingId}", cached("20s", "application/json", Formula1{}.getQualificationResults)).Methods("GET", "OPTIONS")
	router.Handle(route+"f1/v1/fom-results/practice/meeting={meetingId}&session={sessionId}", cached("20s", "application/json", Formula1{}.getFreePracticeResults)).Methods("GET", "OPTIONS")
	router.Handle(route+"f1/v1/fom-results/sprint-shootout/meeting={meetingId}", cached("20s", "application/json", Formula1{}.getSprintQualifyingResults)).Methods("GET", "OPTIONS")
	router.Handle(route+"f1/v1/fom-results/sprint/meeting={meetingId}", cached("20s", "application/json", Formula1{}.getSprintResults)).Methods("GET", "OPTIONS")
	router.Handle(route+"f1/v1/fom-results/starting-grid/meeting={meetingId}", cached("20s", "application/json", Formula1{}.getStartingGrid)).Methods("GET", "OPTIONS")
	router.Handle(route+"f1/v1/editorial-driverlisting/listing", cached("20s", "application/json", Formula1{}.getDriverStandings)).Methods("GET", "OPTIONS")
	router.Handle(route+"f1/v1/editorial-constructorlisting/listing", cached("20s", "application/json", Formula1{}.getTeamStandings)).Methods("GET", "OPTIONS")
	router.Handle(route+"f1/v1/editorial-eventlisting/events", cached("20s", "application/json", Formula1{}.getSchedule)).Methods("GET", "OPTIONS")
	router.Handle(route+"f1/en/results.html/{year}/races/{fomRaceId}/{raceName}.html", cached("120s", "application/json", Formula1{}.getFinishedSessions)).Methods("GET", "OPTIONS")
	router.Handle(route+"f1/en/results.html/{year}/races/{fomRaceId}/{raceName}/{session}.html", cached("120s", "text/html; charset=utf-8", Formula1{}.getResultsForScraping)).Methods("GET", "OPTIONS")
	router.Handle(route+"f1/en/racing/{year}/{circuitName}/Circuit.html", cached("168h", "text/html; charset=utf-8", Formula1{}.getCircuitDetails)).Methods("GET", "OPTIONS")
	router.Handle(route+"f1/en/drivers/{driverId}.html", cached("1h", "text/html; charset=utf-8", Formula1{}.getDriverDetails)).Methods("GET", "OPTIONS")
	router.Handle(route+"f1/en/teams/{teamId}.html", cached("1h", "text/html; charset=utf-8", Formula1{}.getTeamDetails)).Methods("GET", "OPTIONS")
	router.Handle(route+"f1/en/drivers/hall-of-fame", cached("24h", "text/html; charset=utf-8", Formula1{}.getHallOfFame)).Methods("GET", "OPTIONS")
	router.Handle(route+"f1/en/information/drivers-hall-of-fame-{driver}.{id}", cached("24h", "text/html; charset=utf-8", Formula1{}.getHallOfFameDriverDetails)).Methods("GET", "OPTIONS")
	router.Handle(route+"f1/documents", cached("30s", "text/html; charset=utf-8", Formula1{}.getSessionDocuments)).Methods("GET", "OPTIONS")
	router.HandleFunc(route+"f1/documents/{documentPath}", Formula1{}.getSessionDocument).Methods("GET", "OPTIONS")
	router.HandleFunc(route+"f1/rss/{languageCode}", Formula1{}.getRssFeed).Methods("GET", "OPTIONS")
	// v2 routes
	router.Handle(route+"f1/v2/fom-results/race/meeting={meetingId}", cached("20s", "application/json", Formula1{}.getRaceResultsV2)).Methods("GET", "OPTIONS")
	router.Handle(route+"f1/v2/fom-results/qualifying/meeting={meetingId}", cached("20s", "application/json", Formula1{}.getQualificationResultsV2)).Methods("GET", "OPTIONS")
	router.Handle(route+"f1/v2/fom-results/practice/meeting={meetingId}&session={sessionId}", cached("20s", "application/json", Formula1{}.getFreePracticeResultsV2)).Methods("GET", "OPTIONS")
	router.Handle(route+"f1/v2/fom-results/sprint-shootout/meeting={meetingId}", cached("20s", "application/json", Formula1{}.getSprintQualifyingResultsV2)).Methods("GET", "OPTIONS")
	router.Handle(route+"f1/v2/fom-results/sprint/meeting={meetingId}", cached("20s", "application/json", Formula1{}.getSprintResultsV2)).Methods("GET", "OPTIONS")
	router.Handle(route+"f1/v2/fom-results/starting-grid/meeting={meetingId}", cached("20s", "application/json", Formula1{}.getStartingGridV2)).Methods("GET", "OPTIONS")

	// FE championship
	router.Handle(route+"fe/content/formula-e/text/EN/page={page}&pageSize=16&tagNames=content-type:news&tagExpression=&playlistTypeRestriction=&playlistId=&detail=&size=16&championshipId=&sort=", cached("30s", "application/json", FormulaE{}.getArticles)).Methods("GET", "OPTIONS")
	router.Handle(route+"fe/content/formula-e/text/EN/{articleId}", cached("5m", "application/json", FormulaE{}.getArticle)).Methods("GET", "OPTIONS")
	router.Handle(route+"fe/en/news/{articleId}", cached("5m", "application/json", FormulaE{}.getArticleAsHtml)).Methods("GET", "OPTIONS")
	router.Handle(route+"fe/content/formula-e/playlist/EN/15/page={page}&pageSize={limit}&detail=DETAILED&size={limit}", cached("30s", "application/json", FormulaE{}.getVideos)).Methods("GET", "OPTIONS")
	router.Handle(route+"fe/formula-e/v1/races/championshipId={championshipId}", cached("5m", "application/json", FormulaE{}.getSchedule)).Methods("GET", "OPTIONS")
	router.Handle(route+"fe/formula-e/v1/races/{raceId}", cached("30s", "application/json", FormulaE{}.getRaceDetails)).Methods("GET", "OPTIONS")
	router.Handle(route+"fe/formula-e/v1/races/{raceId}/sessions", cached("3h", "application/json", FormulaE{}.getSessions)).Methods("GET", "OPTIONS")
	router.Handle(route+"fe/formula-e/v1/races/{raceId}/sessions/{sessionId}/results", cached("30s", "application/json", FormulaE{}.getSessionResults)).Methods("GET", "OPTIONS")
	router.Handle(route+"fe/content/formula-e/EN/contentTypes=video&contentTypes=news&page=0&pageSize=10&references=FORMULA_E_RACE:{raceId}&onlyRestrictedContent=false&detail=DETAILED", cached("30s", "application/json", FormulaE{}.getRaceArticles)).Methods("GET", "OPTIONS")
	router.Handle(route+"fe/content/formula-e/photo/en/references=FORMULA_E_RACE:{raceId}&tagNames=race:bg-image", cached("24h", "application/json", FormulaE{}.getCircuitImageDetails)).Methods("GET", "OPTIONS")
	router.Handle(route+"fe/formula-e/v1/standings/drivers/championshipId={championshipId}", cached("1m", "application/json", FormulaE{}.getDriverStandings)).Methods("GET", "OPTIONS")
	router.Handle(route+"fe/formula-e/v1/standings/teams/championshipId={championshipId}", cached("1m", "application/json", FormulaE{}.getTeamStandings)).Methods("GET", "OPTIONS")
	router.Handle(route+"fe/formula-e/v1/championships/latest", cached("168h", "application/json", FormulaE{}.getLatestChampionship)).Methods("GET", "OPTIONS")
	fmt.Println("Box, Box! server running on port " + PORT)
	log.Fatal(http.ListenAndServe(":"+PORT, router))
}
