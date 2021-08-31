package settings

type Settings struct {
	Timer     int
	Timeout   int
	File      string
	Useragent string
}

func InitSettings() Settings {
	s := Settings{
		Timer:     10,
		Timeout:   10,
		File:      "sites.txt",
		Useragent: "Mozilla/5.0 (Windows NT x.y; Win64; x64; rv:10.0) Gecko/20100101 Firefox/10.0",
	}
	return s
}
