package light

import (
	"log"
	"os/exec"
)

func (b *baize) wikipedia() {
	err := exec.Command("rundll32", "url.dll,FileProtocolHandler", b.darkBaize.Wikipedia()).Start()
	if err != nil {
		log.Println(err)
	}
}
