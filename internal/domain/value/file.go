package value

var AllowedFileTypes = map[string]string{
	"jpg":                  "",
	"jpeg":                 "",
	"pdf":                  "",
	"txt":                  "",
	"text":                 "",
	"plain; charset=utf-8": "", // for unit tests
	"octet-stream":         "", // for benchmarking
}
