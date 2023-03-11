package light

import (
	"log"
	"os/exec"
)

func (b *baize) wikipedia() {
	var cmd *exec.Cmd = exec.Command("xdg-open", b.darkBaize.Wikipedia())
	if cmd != nil {
		err := cmd.Start()
		if err != nil {
			log.Println(err)
		}
	}
}
