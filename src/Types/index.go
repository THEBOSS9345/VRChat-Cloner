package Types

import "time"

const (
	StartMessage = " █████╗ ██╗   ██╗ █████╗ ████████╗ █████╗ ██████╗ ██╗  ██╗██╗   ██╗██████╗ \n██╔══██╗██║   ██║██╔══██╗╚══██╔══╝██╔══██╗██╔══██╗██║  ██║██║   ██║██╔══██╗\n███████║██║   ██║███████║   ██║   ███████║██████╔╝███████║██║   ██║██████╔╝\n██╔══██║╚██╗ ██╔╝██╔══██║   ██║   ██╔══██║██╔══██╗██╔══██║██║   ██║██╔══██╗\n██║  ██║ ╚████╔╝ ██║  ██║   ██║   ██║  ██║██║  ██║██║  ██║╚██████╔╝██████╔╝\n╚═╝  ╚═╝  ╚═══╝  ╚═╝  ╚═╝   ╚═╝   ╚═╝  ╚═╝╚═╝  ╚═╝╚═╝  ╚═╝ ╚═════╝ ╚═════╝ \n                                                                           \n\n"
	ProjectName  = "AvatarHub"
)

var ColorCodes = map[string]string{
	"Black":     "\033[30m",
	"Red":       "\033[31m",
	"Green":     "\033[32m",
	"Yellow":    "\033[33m",
	"Blue":      "\033[34m",
	"Magenta":   "\033[35m",
	"Cyan":      "\033[36m",
	"White":     "\033[37m",
	"Reset":     "\033[0m",
	"BlackBG":   "\033[40m",
	"RedBG":     "\033[41m",
	"GreenBG":   "\033[42m",
	"YellowBG":  "\033[43m",
	"BlueBG":    "\033[44m",
	"MagentaBG": "\033[45m",
	"CyanBG":    "\033[46m",
	"WhiteBG":   "\033[47m",
}

type FileInfo struct {
	Id   string    `json:"id"`
	Time time.Time `json:"time"`
}

type WebSocketAvatar struct {
	Id           string     `json:"id"`
	Name         string     `json:"name"`
	ImageUrl     string     `json:"imageUrl"`
	Description  string     `json:"description"`
	CreatedAt    *time.Time `json:"created_at"`
	EquippedTime *time.Time `json:"equippedTime"`
}

type Avatar struct {
	AssetUrl       string `json:"assetUrl"`
	AssetUrlObject struct {
	} `json:"assetUrlObject"`
	AuthorId          string    `json:"authorId"`
	AuthorName        string    `json:"authorName"`
	CreatedAt         time.Time `json:"created_at"`
	Description       string    `json:"description"`
	Featured          bool      `json:"featured"`
	Id                string    `json:"id"`
	ImageUrl          string    `json:"imageUrl"`
	Name              string    `json:"name"`
	ReleaseStatus     string    `json:"releaseStatus"`
	Tags              []string  `json:"tags"`
	ThumbnailImageUrl string    `json:"thumbnailImageUrl"`
	UnityPackageUrl   string    `json:"unityPackageUrl"`
	UnityPackages     []struct {
		Id             string `json:"id"`
		AssetUrl       string `json:"assetUrl"`
		AssetUrlObject struct {
		} `json:"assetUrlObject"`
		AssetVersion        int       `json:"assetVersion"`
		CreatedAt           time.Time `json:"created_at"`
		ImpostorizerVersion string    `json:"impostorizerVersion"`
		PerformanceRating   string    `json:"performanceRating"`
		Platform            string    `json:"platform"`
		PluginUrl           string    `json:"pluginUrl"`
		PluginUrlObject     struct {
		} `json:"pluginUrlObject"`
		UnitySortNumber int64  `json:"unitySortNumber"`
		UnityVersion    string `json:"unityVersion"`
		WorldSignature  string `json:"worldSignature"`
		ImpostorUrl     string `json:"impostorUrl"`
		ScanStatus      string `json:"scanStatus"`
		Variant         string `json:"variant"`
	} `json:"unityPackages"`
	UpdatedAt time.Time `json:"updated_at"`
	Version   int       `json:"version"`
	CacheId   string    `json:"cacheId,omitempty"`
}
