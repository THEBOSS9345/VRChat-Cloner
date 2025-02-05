package avatars

import (
	"VRCHAT/src/Types"
	"VRCHAT/src/VRAPI"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func GetAvatar(avatarId string) (*Types.Avatar, error) {
	url := "https://vrchat.com/api/1/avatars/" + avatarId

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil, err
	}

	VRAPI.ChangeDefaultHeader(req)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return nil, err
	}

	defer resp.Body.Close()
	var avatar Types.Avatar

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Println("Error reading body:", err)
		return nil, err
	}

	if strings.Contains(string(body), "\"message\":\"Avatar Not Found\",") {
		return nil, nil
	}

	err = json.Unmarshal(body, &avatar)

	if err != nil {
		fmt.Println("Error unmarshalling body:", err)
		return nil, err
	}

	return &avatar, nil
}

func UpdateAvatar(avatarId string) {
	url := "https://vrchat.com/api/1/avatars/" + avatarId + "/select"
	req, err := http.NewRequest("PUT", url, nil)

	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	VRAPI.ChangeDefaultHeader(req)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}

	defer resp.Body.Close()
}

func AddImage(avatar Types.Avatar) (Types.Avatar, error) {
	url := avatar.ImageUrl
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return Types.Avatar{}, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:132.0) Gecko/20100101 Firefox/132.0")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return Types.Avatar{}, err
	}
	defer resp.Body.Close()

	// Check if the response is successful
	if resp.StatusCode != http.StatusOK {
		return Types.Avatar{}, err
	}

	avatar.ImageUrl = resp.Request.URL.String()

	return avatar, nil
}
