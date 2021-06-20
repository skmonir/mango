package intel

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_whereAmI(t *testing.T) {
	Convey("Test_whereAmI", t, func() {
		Convey("valid", func() {
			p := "/codeforces/1545"
			c := whereAmI(p)
			So(c, ShouldNotBeNil)
			So(c.Workspace, ShouldEqual, "/")
			So(c.OJ, ShouldEqual, "codeforces")
			So(c.CurrentContestId, ShouldEqual, "1545")
		})
		Convey("valid-2", func() {
			p := "/some/random/dir/cf/src/1545"
			c := whereAmI(p)
			So(c, ShouldNotBeNil)
			So(c.Workspace, ShouldEqual, "/some/random/dir/")
			So(c.OJ, ShouldEqual, "codeforces")
			So(c.CurrentContestId, ShouldEqual, "1545")
		})
		Convey("invalid", func() {
			p := "/uva/src/problem-1245"
			c := whereAmI(p)
			So(c, ShouldBeNil)
		})
	})

}
