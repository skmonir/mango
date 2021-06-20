package intel

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_whereAmI(t *testing.T) {
	Convey("Test_whereAmI", t, func() {
		Convey("valid", func() {
			p := "/codeforces/src/1545"
			c := whereAmI(p)
			So(c, ShouldNotBeNil)
			So(c.OJ, ShouldEqual, "codeforces")
			So(c.CurrentContestId, ShouldEqual, "1545")
		})
	})

}
