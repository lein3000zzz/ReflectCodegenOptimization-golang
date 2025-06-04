package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	_ "net/http/pprof"
	"os"
	"strings"
	"sync"
)

// вам надо написать более быструю оптимальную этой функции
// func FastSearch(out io.Writer) {
// 	SlowSearch(out)
// }

type User struct {
	Name     string   `json:"name"`
	Email    string   `json:"email"`
	Browsers []string `json:"browsers"`
}

var userPool = sync.Pool{
	New: func() interface{} {
		return &User{}
	},
}

func FastSearch(out io.Writer) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	seenBrowsers := make(map[string]struct{})
	foundUsersBuffer := bytes.NewBuffer(make([]byte, 0, 2048))

	decoder := json.NewDecoder(file)
	lineNum := 0

	for {
		user := userPool.Get().(*User)

		err := decoder.Decode(user)
		if err != nil {
			userPool.Put(user)
			if err == io.EOF {
				break
			}
			continue
		}

		isAndroid := false
		isMSIE := false

		for _, browser := range user.Browsers {
			if strings.Contains(browser, "Android") {
				isAndroid = true
				seenBrowsers[browser] = struct{}{}
			}

			if strings.Contains(browser, "MSIE") {
				isMSIE = true
				seenBrowsers[browser] = struct{}{}
			}
		}

		if isAndroid && isMSIE {
			email := strings.ReplaceAll(user.Email, "@", " [at] ")
			fmt.Fprintf(foundUsersBuffer, "[%d] %s <%s>\n", lineNum, user.Name, email)
		}

		// Возвращаем User обратно в пул
		userPool.Put(user)
		lineNum++
	}

	fmt.Fprintln(out, "found users:\n"+foundUsersBuffer.String())
	fmt.Fprintln(out, "Total unique browsers", len(seenBrowsers))
}
