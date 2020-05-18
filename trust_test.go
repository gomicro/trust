package trust

import (
	"strings"
	"testing"

	"github.com/franela/goblin"
	. "github.com/onsi/gomega"
)

func TestCACerts(t *testing.T) {
	g := goblin.Goblin(t)
	RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })

	g.Describe("Get Cert Pool", func() {
		g.It("should append added cert files into the pool", func() {
			pool := New()

			f := "./testCA.pem"
			pool.AddCAFile(f)

			p, err := pool.CACerts()
			Expect(err).To(BeNil())

			additionalFound := false
			for _, s := range p.Subjects() {
				if strings.Contains(string(s), "Gomicro") {
					additionalFound = true
					break
				}
			}
			Expect(additionalFound).To(BeTrue())
		})
	})
}
