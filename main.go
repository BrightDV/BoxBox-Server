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
 * Copyright (c) 2022-2023, BrightDV
 */

package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
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

var tmp, _ = base64.StdEncoding.DecodeString("cVBnUFBSSnlHQ0lQeEZUM2VsNE1GN3RoWEh5SkN6QVA=")
var apikey = string(tmp)
var tmp_, _ = base64.StdEncoding.DecodeString("aHR0cHM6Ly9hcGkuZm9ybXVsYTEuY29tLw==")
var endpoint = string(tmp_)
var DOMAIN string = "*"
var PORT string = "8080"

var storage cache.Storage

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
			for k, v := range c.HeaderMap {
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

func getArticle(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", DOMAIN)
	var article any
	articleId := mux.Vars(r)["articleId"]
	logger(r.RequestURI)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", endpoint+"v1/editorial/articles/"+articleId, nil)
	req.Header.Set("User-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36")
	req.Header.Set("apikey", apikey)
	req.Header.Set("locale", "en")
	resp, _ := client.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &article)
	json.NewEncoder(w).Encode(article)
}

func getArticles(w http.ResponseWriter, r *http.Request) {
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
	req, _ := http.NewRequest("GET", endpoint+"v1/editorial/articles?limit=16&tags="+tags+"&offset="+offset+"&articleTypes="+articleTypes, nil)
	req.Header.Set("User-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36")
	req.Header.Set("apikey", apikey)
	req.Header.Set("locale", "en")
	resp, _ := client.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &articles)
	json.NewEncoder(w).Encode(articles)
}

func getVideo(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", DOMAIN)
	var article any
	videoId := mux.Vars(r)["videoId"]
	logger(r.RequestURI)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", endpoint+"v1/video-assets/videos/"+videoId, nil)
	req.Header.Set("User-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36")
	req.Header.Set("apikey", apikey)
	req.Header.Set("locale", "en")
	resp, _ := client.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &article)
	json.NewEncoder(w).Encode(article)
}

func getVideos(w http.ResponseWriter, r *http.Request) {
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
	req, _ := http.NewRequest("GET", endpoint+"v1/video-assets/videos?limit="+limit+"&tag="+tags+"&offset="+offset, nil)
	req.Header.Set("User-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36")
	req.Header.Set("apikey", apikey)
	req.Header.Set("locale", "en")
	resp, _ := client.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &articles)
	json.NewEncoder(w).Encode(articles)
}

func eventTracker(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", DOMAIN)
	var event any
	logger(r.RequestURI)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", endpoint+"v1/event-tracker", nil)
	req.Header.Set("User-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36")
	req.Header.Set("apikey", apikey)
	req.Header.Set("locale", "en")
	req.Header.Set("Accept", "application/json")
	resp, _ := client.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &event)
	json.NewEncoder(w).Encode(event)
}

func eventTrackerForOneMeeting(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", DOMAIN)
	var event any
	meetingId := mux.Vars(r)["meetingId"]
	logger(r.RequestURI)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", endpoint+"v1/event-tracker/meeting/"+meetingId, nil)
	req.Header.Set("User-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36")
	req.Header.Set("apikey", apikey)
	req.Header.Set("locale", "en")
	resp, _ := client.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &event)
	json.NewEncoder(w).Encode(event)
}

func getFinishedSessions(w http.ResponseWriter, r *http.Request) {
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
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Fprint(w, string(body))
}

func getResultsForScraping(w http.ResponseWriter, r *http.Request) {
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
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Fprint(w, string(body))
}

func getCircuitDetails(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", DOMAIN)
	year := mux.Vars(r)["year"]
	circuitName := mux.Vars(r)["circuitName"]
	logger(r.RequestURI)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://www.formula1.com/en/racing/"+year+"/"+circuitName+"/Circuit.html", nil)
	req.Header.Set("User-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36")
	resp, _ := client.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Fprint(w, string(body))
}

func getDriverDetails(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", DOMAIN)
	driverId := mux.Vars(r)["driverId"]
	logger(r.RequestURI)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://www.formula1.com/en/drivers/"+driverId+".html", nil)
	req.Header.Set("User-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36")
	resp, _ := client.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Fprint(w, string(body))
}

func getTeamDetails(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", DOMAIN)
	teamId := mux.Vars(r)["teamId"]
	logger(r.RequestURI)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://www.formula1.com/en/teams/"+teamId+".html", nil)
	req.Header.Set("User-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36")
	resp, _ := client.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Fprint(w, string(body))
}

func getHallOfFame(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", DOMAIN)
	logger(r.RequestURI)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://www.formula1.com/en/drivers/hall-of-fame.html", nil)
	req.Header.Set("User-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36")
	resp, _ := client.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Fprint(w, string(body))
}

func getHallOfFameDriverDetails(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", DOMAIN)
	driver := mux.Vars(r)["driver"]
	logger(r.RequestURI)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://www.formula1.com/en/drivers/hall-of-fame/"+driver+".html", nil)
	req.Header.Set("User-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36")
	resp, _ := client.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Fprint(w, string(body))
}

func getSessionDocuments(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", DOMAIN)
	logger(r.RequestURI)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://www.fia.com/documents/championships/fia-formula-one-world-championship-14/season/season-2024-2043", nil)
	req.Header.Set("User-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36")
	resp, _ := client.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Fprint(w, string(body))
}

func getSessionDocument(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/pdf")
	w.Header().Set("Access-Control-Allow-Origin", DOMAIN)
	documentPath := mux.Vars(r)["documentPath"]
	logger(r.RequestURI)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://www.fia.com/sites/default/files/decision-document/"+documentPath, nil)
	req.Header.Set("User-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36")
	resp, _ := client.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Fprint(w, string(body))
}

func getRssFeed(w http.ResponseWriter, r *http.Request) {
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
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Fprint(w, string(body))
}

func keepAlive(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", DOMAIN)
	logger(r.RequestURI)
	fmt.Fprint(w, string("I'm still standing"))
}

func main() {
	fmt.Println("Initializing...")
	route := "/"
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homeLink)
	router.Handle(route+"v1/editorial/articles", cached("30s", "application/json", getArticles)).Methods("GET", "OPTIONS")
	router.Handle(route+"v1/editorial/articles/{articleId}", cached("5m", "application/json", getArticle)).Methods("GET", "OPTIONS")
	router.Handle(route+"v1/video-assets/videos", cached("30s", "application/json", getVideos)).Methods("GET", "OPTIONS")
	router.Handle(route+"v1/video-assets/videos/{videoId}", cached("5m", "application/json", getVideo)).Methods("GET", "OPTIONS")
	router.Handle(route+"v1/event-tracker", cached("20s", "application/json", eventTracker)).Methods("GET", "OPTIONS")
	router.Handle(route+"v1/event-tracker/meeting/{meetingId}", cached("20s", "application/json", eventTrackerForOneMeeting)).Methods("GET", "OPTIONS")
	router.Handle(route+"en/results.html/{year}/races/{fomRaceId}/{raceName}.html", cached("120s", "application/json", getFinishedSessions)).Methods("GET", "OPTIONS")
	router.Handle(route+"en/results.html/{year}/races/{fomRaceId}/{raceName}/{session}.html", cached("120s", "text/html; charset=utf-8", getResultsForScraping)).Methods("GET", "OPTIONS")
	router.Handle(route+"en/racing/{year}/{circuitName}/Circuit.html", cached("168h", "text/html; charset=utf-8", getCircuitDetails)).Methods("GET", "OPTIONS")
	router.Handle(route+"en/drivers/{driverId}.html", cached("1h", "text/html; charset=utf-8", getDriverDetails)).Methods("GET", "OPTIONS")
	router.Handle(route+"en/teams/{teamId}.html", cached("1h", "text/html; charset=utf-8", getTeamDetails)).Methods("GET", "OPTIONS")
	router.Handle(route+"en/drivers/hall-of-fame.html", cached("24h", "text/html; charset=utf-8", getHallOfFame)).Methods("GET", "OPTIONS")
	router.Handle(route+"en/drivers/hall-of-fame/{driver}.html", cached("24h", "text/html; charset=utf-8", getHallOfFameDriverDetails)).Methods("GET", "OPTIONS")
	router.Handle(route+"documents", cached("30s", "text/html; charset=utf-8", getSessionDocuments)).Methods("GET", "OPTIONS")
	router.HandleFunc(route+"documents/{documentPath}", getSessionDocument).Methods("GET", "OPTIONS")
	router.HandleFunc(route+"rss/{languageCode}", getRssFeed).Methods("GET", "OPTIONS")
	router.HandleFunc(route+"keepAlive", keepAlive).Methods("GET", "OPTIONS")
	fmt.Println("Box, Box! server running on port " + PORT)
	log.Fatal(http.ListenAndServe(":"+PORT, router))
}
