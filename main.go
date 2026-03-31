package hcms

import (
	"log"

	"github.com/pchchv/hcms/config"
)

func main() {
	_, err := config.Load()
	if err != nil {
		log.Fatalf("config error: %v", err)
	}
}
