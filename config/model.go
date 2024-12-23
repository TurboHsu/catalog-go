package config

type Config struct {
	Server   Server
	Database Database
	Receiver Receiver
	Store    Store
}

type Server struct {
	Listen       string
	AllowOrigins []string
}

type Database struct {
	Type             string
	Path             string
	AllowedReactions []string
}

type Store struct {
	StorePath string
	Qiniu     Qiniu
}

type Qiniu struct {
	Enable     bool
	AccessKey  string
	SecretKey  string
	Bucket     string
	UploadPath string
}

type Receiver struct {
	TelegramBot TelegramBot
}

type TelegramBot struct {
	Enable         bool
	Token          string
	PermittedUsers []int64
	ServerURL      string
}

func DefaultConfig() Config {
	return Config{
		Server: Server{
			Listen:       ":28081",
			AllowOrigins: []string{""},
		},
		Database: Database{
			Type: "sqlite3",
			Path: "database.db",
			AllowedReactions: []string{
				"😺",
				"😸",
				"😹",
				"😻",
				"😼",
				"😽",
				"🙀",
				"😿",
				"😾",
				"🐱",
			},
		},
		Receiver: Receiver{
			TelegramBot: TelegramBot{
				Enable:         false,
				Token:          "",
				PermittedUsers: []int64{},
				ServerURL:      "https://api.telegram.org",
			},
		},
		Store: Store{
			StorePath: "./static",
		},
	}
}
