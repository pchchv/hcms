package hcms

import "log"

func main() {
	_, err := Load()
	if err != nil {
		log.Fatalf("config error: %v", err)
	}
}
