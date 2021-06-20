package oj

type OJ string

func (o OJ) String() string {
	return string(o)
}

const (
	CF OJ = "codeforces"
)
