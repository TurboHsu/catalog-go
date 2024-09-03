package config

type Config struct {
	Server Server
	Database Database
	Receiver Receiver
}

type Server struct {
	Listen string
}

type Database struct {
	Type string
	Path string
}

type Receiver struct {
	TelegramBot TelegramBot
}

type TelegramBot struct {
	Enable bool
	Token string
}

func DefaultConfig() Config {
	return Config{
		Server: Server{
			Listen: ":28081",
		},
		Database: Database{
			Type: "sqlite3",
			Path: "database.db",
		},
		Receiver: Receiver{
			TelegramBot: TelegramBot{
				Enable: false,
				Token: "",
			},
		},
	}
}