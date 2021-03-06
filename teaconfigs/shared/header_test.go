package shared

import (
	"github.com/iwind/TeaGo/assert"
	"testing"
)

func TestHeaderConfig_Match(t *testing.T) {
	a := assert.NewAssertion(t)
	h := NewHeaderConfig()
	h.Validate()
	a.IsTrue(h.Match(200))
	a.IsFalse(h.Match(400))

	h.Status = []int{200, 201, 400}
	h.Validate()
	a.IsTrue(h.Match(400))
	a.IsFalse(h.Match(500))

	h.Always = true
	a.IsTrue(h.Match(500))
}
