package language

type Language int

func (l Language) Name() string {
	return names[l]
}

func (l Language) CompileCmd(srcPath, binPath string) string {
	return compileCmd[CPP](srcPath, binPath)
}

const (
	CPP Language = iota
	CPP11
	CPP14
	CPP17
	GO
	UKN
)

const Default = CPP17

var names = map[Language]string{
	CPP:   "C++",
	CPP11: "C++11",
	CPP14: "C++14",
	CPP17: "C++17",
	GO:    "go",
	UKN:   "uknown",
}

var compileCmd = map[Language]func(src, bin string) string{
	CPP:   func(src, bin string) string { return "g++ " + src + " -o " + bin },
	CPP11: func(src, bin string) string { return "g++ -std=c++11 " + src + " -o " + bin },
	CPP14: func(src, bin string) string { return "g++ -std=c++14 " + src + " -o " + bin },
	CPP17: func(src, bin string) string { return "g++ -std=c++17 " + src + " -o " + bin },
	GO:    func(src, bin string) string { return "go build -o " + bin + " " + src },
	UKN:   func(src, bin string) string { return "" },
}

func FindLangFromExt(ext string) Language {
	L := UKN
	switch ext {
	case "cpp":
		L = CPP17
	case "go":
		L = GO
	}
	return L
}
