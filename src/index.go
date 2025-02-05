package main

import (
	database "VRCHAT/src/Database"
	"VRCHAT/src/Types"
	"VRCHAT/src/VRAPI/avatars"
	"VRCHAT/src/WebPage"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	webview "github.com/webview/webview_go"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

var patternA = regexp.MustCompile("prefab-id-v1_avtr_([0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12})")
var baseDir = fmt.Sprintf("%s\\AppData\\LocalLow\\VRChat\\VRChat", os.Getenv("USERPROFILE"))

func main() {
	fmt.Print(Types.StartMessage + "\n")
	log.Printf("%s[%s Panel]%s %sStarting %s...%s\n", Types.ColorCodes["Yellow"], Types.ProjectName, Types.ColorCodes["Reset"], Types.ColorCodes["Green"], Types.ProjectName, Types.ColorCodes["Reset"])

	log.Printf("%s[%s Panel]%s %sRegistering WebPages...%s\n", Types.ColorCodes["Yellow"], Types.ProjectName, Types.ColorCodes["Reset"], Types.ColorCodes["Green"], Types.ColorCodes["Reset"])

	RegisterWebstocks()

	log.Printf("%s[%s Panel]%s %sGetting All Avatars...%s\n", Types.ColorCodes["Yellow"], Types.ProjectName, Types.ColorCodes["Reset"], Types.ColorCodes["Green"], Types.ColorCodes["Reset"])

	port, err := OpenPort()

	if err != nil {
		log.Printf("%s[%s Panel]%s %sFailed to open port: %v%s\n", Types.ColorCodes["Yellow"], Types.ProjectName, Types.ColorCodes["Reset"], Types.ColorCodes["Red"], err, Types.ColorCodes["Reset"])
		return
	}

	go func() {
		err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)

		if err != nil {
			log.Printf("%s[%s Panel]%s %sFailed to start server: %v%s\n", Types.ColorCodes["Yellow"], Types.ProjectName, Types.ColorCodes["Reset"], Types.ColorCodes["Red"], err, Types.ColorCodes["Reset"])
			return
		}
	}()

	log.Printf("%s[%s Panel]%s %sStarting Webview...%s\n", Types.ColorCodes["Yellow"], Types.ProjectName, Types.ColorCodes["Reset"], Types.ColorCodes["Green"], Types.ColorCodes["Reset"])

	w := webview.New(true)

	w.SetSize(900, 900, webview.HintMin)

	defer func() {
		w.Destroy()

		log.Printf("%s[%s Panel]%s %sWebview destroyed successfully%s\n", Types.ColorCodes["Yellow"], Types.ProjectName, Types.ColorCodes["Reset"], Types.ColorCodes["Green"], Types.ColorCodes["Reset"])
	}()

	AppHtml := WebPage.GetHomeHtml(port)

	w.SetHtml(AppHtml)

	w.SetTitle(Types.ProjectName)

	log.Printf("%s[%s Panel]%s %sWebview started successfully%s\n", Types.ColorCodes["Yellow"], Types.ProjectName, Types.ColorCodes["Reset"], Types.ColorCodes["Green"], Types.ColorCodes["Reset"])
	w.Run()
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func RegisterWebstocks() {
	http.HandleFunc("/avatars", func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")

		log.Printf("%s[%s Panel]%s %sConnection Received%s\n", Types.ColorCodes["Yellow"], Types.ProjectName, Types.ColorCodes["Reset"], Types.ColorCodes["Green"], Types.ColorCodes["Reset"])

		conn, err := upgrader.Upgrade(w, r, nil)

		if err != nil {
			log.Printf("Failed to upgrade to WebSocket: %v\n", err)
			return
		}

		defer func() {
			err := conn.Close()

			if err != nil {
				log.Printf("%s[%s Panel]%s %sFailed to close connection: %v%s\n", Types.ColorCodes["Yellow"], Types.ProjectName, Types.ColorCodes["Reset"], Types.ColorCodes["Red"], err, Types.ColorCodes["Reset"])
			}

			log.Printf("%s[%s Panel]%s %sConnection Closed%s\n", Types.ColorCodes["Yellow"], Types.ProjectName, Types.ColorCodes["Reset"], Types.ColorCodes["Green"], Types.ColorCodes["Reset"])
		}()

		cache, err := database.GetAll()

		if err != nil {
			fmt.Printf("%s[%s Panel]%s %sFailed to get all avatars: %v%s\n", Types.ColorCodes["Yellow"], Types.ProjectName, Types.ColorCodes["Reset"], Types.ColorCodes["Red"], err, Types.ColorCodes["Reset"])
			return
		}

		sortedCache := make([]database.Avatar, 0, len(cache))

		for _, avatar := range cache {
			sortedCache = append(sortedCache, avatar)
		}

		for i := 0; i < len(sortedCache); i++ {
			for j := i + 1; j < len(sortedCache); j++ {
				if (sortedCache[i].Star && !sortedCache[j].Star) ||
					(sortedCache[i].CachedTime.After(*sortedCache[j].CachedTime)) ||
					(sortedCache[i].EquippedTime != nil && sortedCache[j].EquippedTime != nil && sortedCache[i].EquippedTime.After(*sortedCache[j].EquippedTime)) ||
					(sortedCache[i].Info.Name < sortedCache[j].Info.Name) {
					sortedCache[i], sortedCache[j] = sortedCache[j], sortedCache[i]
				}
			}
		}
		for _, avatar := range sortedCache {
			if has, err := database.Has("avtr_" + avatar.ID); err != nil || has || avatar.ID == "" {
				continue
			}

			err = conn.WriteJSON(Types.WebSocketAvatar{
				Id:           avatar.ID,
				Name:         avatar.Info.Name,
				ImageUrl:     avatar.Info.ImageUrl,
				Description:  avatar.Info.Description,
				CreatedAt:    avatar.CachedTime,
				EquippedTime: avatar.EquippedTime,
			})

			if err != nil {
				log.Printf("%s[%s Panel]%s %sFailed to write to WebSocket: %v%s\n", Types.ColorCodes["Yellow"], Types.ProjectName, Types.ColorCodes["Reset"], Types.ColorCodes["Red"], err, Types.ColorCodes["Reset"])
				return
			}

		}

		log.Printf("%s[%s Panel]%s %sGetting Avatars Time Started: %s%s\n", Types.ColorCodes["Yellow"], Types.ProjectName, Types.ColorCodes["Reset"], Types.ColorCodes["Green"], time.Now().Format("15:04:05"), Types.ColorCodes["Reset"])

		avatarIds := getAvatarIds()

		log.Printf("%s[%s Panel]%s %sGetting Avatars Time Ended: %s%s\n", Types.ColorCodes["Yellow"], Types.ProjectName, Types.ColorCodes["Reset"], Types.ColorCodes["Green"], time.Now().Format("15:04:05"), Types.ColorCodes["Reset"])

		for _, avatarInfo := range avatarIds {
			if has, err := database.Has("avtr_" + avatarInfo.Id); err != nil || has || avatarInfo.Id == "" {
				continue
			}

			avatarBeforeImageLink, err := avatars.GetAvatar("avtr_" + avatarInfo.Id)

			if err != nil || avatarBeforeImageLink == nil {
				continue
			}

			avatar, err := avatars.AddImage(*avatarBeforeImageLink)

			if err != nil || avatar.Id == "" {
				continue
			}

			err = database.Save(database.Avatar{
				ID:           avatar.Id,
				EquippedTime: nil,
				CachedTime:   &avatarInfo.Time,
				Info:         avatar,
			})

			if err != nil {
				continue
			}

			err = conn.WriteJSON(Types.WebSocketAvatar{
				Id:           avatar.Id,
				Name:         avatar.Name,
				ImageUrl:     avatar.ImageUrl,
				Description:  avatar.Description,
				CreatedAt:    &avatarInfo.Time,
				EquippedTime: nil,
			})

			if err != nil {
				continue
			}
		}

		log.Printf("%s[%s Panel]%s %sAll Avatars Sent%s\n", Types.ColorCodes["Yellow"], Types.ProjectName, Types.ColorCodes["Reset"], Types.ColorCodes["Green"], Types.ColorCodes["Reset"])
	})

	http.HandleFunc("/changeAvatar", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, avatarId")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		if r.Method != http.MethodPost {
			fmt.Println("Invalid request method")
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		bodyData, err := ioutil.ReadAll(r.Body)

		if err != nil {
			fmt.Println("Failed to read body data")
			http.Error(w, "Failed to read body data", http.StatusBadRequest)
			return
		}

		var avatar struct {
			Id string `json:"avatarId"`
		}

		err = json.Unmarshal(bodyData, &avatar)

		if err != nil {
			fmt.Println("Failed to unmarshal body data")
			http.Error(w, "Failed to unmarshal body data", http.StatusBadRequest)
			return
		}

		if avatar.Id == "" {
			fmt.Println("Missing avatarId")
			http.Error(w, "Missing avatarId", http.StatusBadRequest)
			return
		}

		log.Printf("%s[%s Panel]%s %sChanging Avatar to %s%s\n", Types.ColorCodes["Yellow"], Types.ProjectName, Types.ColorCodes["Reset"], Types.ColorCodes["Green"], avatar.Id, Types.ColorCodes["Reset"])

		err = database.UpdateEquippedTime(avatar.Id)

		if err != nil {
			fmt.Println("Failed to update equipped time", err)
			http.Error(w, "Failed to update equipped time", http.StatusInternalServerError)
			return
		}

		avatars.UpdateAvatar(avatar.Id)

		_, err = w.Write([]byte("Avatar changed successfully"))

		if err != nil {
			log.Printf("%s[%s Panel]%s %sFailed to write to response: %v%s\n", Types.ColorCodes["Yellow"], Types.ProjectName, Types.ColorCodes["Reset"], Types.ColorCodes["Red"], err, Types.ColorCodes["Reset"])
			return
		}

		log.Printf("%s[%s Panel]%s %sAvatar changed successfully%s\n", Types.ColorCodes["Yellow"], Types.ProjectName, Types.ColorCodes["Reset"], Types.ColorCodes["Green"], Types.ColorCodes["Reset"])
	})

	http.HandleFunc("/starAvatar", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, avatarId")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		if r.Method != http.MethodPost {
			fmt.Println("Invalid request method")
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		bodyData, err := ioutil.ReadAll(r.Body)

		if err != nil {
			fmt.Println("Failed to read body data")
			http.Error(w, "Failed to read body data", http.StatusBadRequest)
			return
		}

		var avatar struct {
			Id   string `json:"avatarId"`
			Star bool   `json:"star"`
		}

		err = json.Unmarshal(bodyData, &avatar)

		if err != nil {
			fmt.Println("Failed to unmarshal body data")
			http.Error(w, "Failed to unmarshal body data", http.StatusBadRequest)
			return
		}

		if avatar.Id == "" {
			fmt.Println("Missing avatarId")
			http.Error(w, "Missing avatarId", http.StatusBadRequest)
			return
		}

		log.Printf("%s[%s Panel]%s %sStarring Avatar %s%s\n", Types.ColorCodes["Yellow"], Types.ProjectName, Types.ColorCodes["Reset"], Types.ColorCodes["Green"], avatar.Id, Types.ColorCodes["Reset"])

		err = database.UpdateStar(avatar.Id, avatar.Star)

		if err != nil {
			fmt.Println("Failed to update star", err)
			http.Error(w, "Failed to update star", http.StatusInternalServerError)
			return
		}

		_, err = w.Write([]byte("Avatar starred successfully"))

		if err != nil {
			log.Printf("%s[%s Panel]%s %sFailed to write to response: %v%s\n", Types.ColorCodes["Yellow"], Types.ProjectName, Types.ColorCodes["Reset"], Types.ColorCodes["Red"], err, Types.ColorCodes["Reset"])
			return
		}

		log.Printf("%s[%s Panel]%s %sAvatar starred successfully%s\n", Types.ColorCodes["Yellow"], Types.ProjectName, Types.ColorCodes["Reset"], Types.ColorCodes["Green"], Types.ColorCodes["Reset"])
	})
}

func OpenPort() (int32, error) {
	ln, err := net.Listen("tcp", ":0")
	if err != nil {
		return 0, fmt.Errorf("failed to find an open port: %v", err)
	}
	defer ln.Close()

	port := ln.Addr().(*net.TCPAddr).Port
	return int32(port), nil
}

func getAvatarIds() []Types.FileInfo {
	var ActorIds []Types.FileInfo
	filePaths := make(chan string, 100)
	results := make(chan Types.FileInfo, 100)
	done := make(chan bool)

	worker := func() {
		for path := range filePaths {
			var fileInfo Types.FileInfo

			data, err := ioutil.ReadFile(path)
			if err != nil {
				continue
			}

			id := GetAvatarId(data)
			fileInfo.Id = id

			fileStat, err := os.Stat(path)
			if err != nil {
				continue
			}

			fileInfo.Time = fileStat.ModTime()
			results <- fileInfo
		}
		done <- true
	}

	numWorkers := 1

	for i := 0; i < numWorkers; i++ {
		go worker()
	}

	go func() {
		err := filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				log.Printf("Error accessing path %s: %v\n", path, err)
				return err
			}
			if !info.IsDir() && info.Name() == "__data" {
				filePaths <- path
			}
			return nil
		})
		if err != nil {
			log.Printf("Error walking through directory: %v\n", err)
		}
		close(filePaths)
	}()

	go func() {
		for id := range results {
			ActorIds = append(ActorIds, id)
		}
		done <- true
	}()

	for i := 0; i < numWorkers; i++ {
		<-done
	}
	close(results)
	<-done

	return ActorIds
}

func GetAvatarId(data []byte) string {
	avatarMatches := patternA.FindSubmatch(data)
	if len(avatarMatches) >= 2 {
		return string(avatarMatches[1])
	}
	return ""
}
