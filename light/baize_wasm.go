package light

import (
	"syscall/js"
)

func (b *Baize) Wikipedia() {
	js.Global().Get("window").Call("open", b.darkBaize.Wikipedia())
}
