package database

import (
	"VRCHAT/src/Types"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type Avatar struct {
	ID           string
	EquippedTime *time.Time   `json:"equippedTime"`
	CachedTime   *time.Time   `json:"cachedTime"`
	Star         bool         `json:"star"`
	Info         Types.Avatar `json:"info"`
}

var Avatars *Database

func Save(avatar Avatar) error {
	return Avatars.Set(avatar.ID, avatar)
}

func GetAvatar(avatarID string) (*Avatar, error) {

	value, err := Avatars.Get(avatarID)

	if err != nil {
		return nil, err
	}

	avatarMap, ok := value.(map[string]interface{})

	if !ok {
		return nil, fmt.Errorf("failed to convert interface to map for key %s: %v", avatarID, value)
	}

	avatarJSON, err := json.Marshal(avatarMap)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal avatar map for key %s: %v", avatarID, err)
	}

	var avatar Avatar

	err = json.Unmarshal(avatarJSON, &avatar)

	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal avatar JSON for key %s: %v", avatarID, err)
	}

	return &avatar, nil
}

func GetAll() (map[string]Avatar, error) {
	avatars, err := Avatars.GetAll()

	if err != nil {
		return nil, err
	}

	result := make(map[string]Avatar)

	for key, value := range avatars {
		avatarMap, ok := value.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("failed to convert interface to map for key %s: %v", key, value)
		}

		avatarJSON, err := json.Marshal(avatarMap)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal avatar map for key %s: %v", key, err)
		}

		var avatar Avatar
		err = json.Unmarshal(avatarJSON, &avatar)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal avatar JSON for key %s: %v", key, err)
		}

		result[key] = avatar
	}

	return result, nil
}

func UpdateEquippedTime(avatarID string) error {
	avatar, err := GetAvatar(avatarID)

	if err != nil {
		return err
	}

	now := time.Now()
	avatar.EquippedTime = &now

	err = Save(*avatar)

	if err != nil {
		return err
	}

	return nil
}

func UpdateStar(avatarID string, star bool) error {
	avatar, err := GetAvatar(avatarID)

	if err != nil {
		return err
	}

	avatar.Star = star

	return Save(*avatar)
}

func Has(avatarID string) (bool, error) {
	has, err := Avatars.Has(avatarID)

	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return false, nil
		}

		return false, err
	}

	return has, nil
}

func Delete(avatarID string) error {
	return Avatars.Delete(avatarID)
}

func init() {
	var err error

	Avatars, err = New("avatars.db")

	if err != nil {
		panic(err)
	}
}
