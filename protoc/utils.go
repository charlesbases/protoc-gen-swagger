package protoc

import (
	"strings"
	"time"

	"github.com/charlesbases/protoc-gen-swagger/logger"
)

// trim  prefix and suffix TODO 可优化
func trim(source string, cutsets ...string) string {
	source = strings.TrimSpace(source)

	for _, cutset := range cutsets {
		for {
			if strings.HasPrefix(source, cutset) {
				source = strings.TrimPrefix(source, cutset)
				continue
			}
			break
		}
		for {
			if strings.HasSuffix(source, cutset) {
				source = strings.TrimSuffix(source, cutset)
				continue
			}
			break
		}
	}
	return source
}

// split split by "." and return package and type name
func split(typename string) [2]string {
	list := strings.Split(typename, ".")
	if len(list) < 3 {
		logger.Fatal("split type failed. ", typename)
	}
	return [2]string{list[1], strings.Join(list[2:], "_")}
}

// nestedName message nested name
func nestedName(v ...string) string {
	return strings.Join(v, "_")
}

// methodPath .
func methodPath(v ...string) string {
	return "/" + strings.Join(v, "/")
}

// version .
func version() string {
	return time.Now().Format("20060102150405")
}

// ascending 升序
func ascending(l, r string) bool {
	var length int
	if len(l) < len(r) {
		length = len(l)
	} else {
		length = len(r)
	}

	for i := 0; i < length; i++ {
		if l[i] != r[i] {
			return l[i] < r[i]
		}
	}
	return len(l) < len(r)
}

// descending 降序
func descending(l, r string) bool {
	var length int
	if len(l) < len(r) {
		length = len(l)
	} else {
		length = len(r)
	}

	for i := 0; i < length; i++ {
		if l[i] != r[i] {
			return l[i] > r[i]
		}
	}
	return len(l) > len(r)
}
