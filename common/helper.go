package common

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	//	"github.com/goph/emperror"
	"github.com/op/go-logging"
	"google.golang.org/grpc/metadata"
	"os"
)

var _logformat = logging.MustStringFormatter(
	`%{time:2006-01-02T15:04:05.000} %{module}::%{shortfunc} [%{shortfile}] > %{level:.5s} - %{message}`,
)

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}


func StringIntersect( arr1, arr2 []string ) bool {
	for _, str1 := range arr1 {
		for _, str2 := range arr2 {
			if str1 == str2 {
				return true
			}
		}
	}
	return false
}

func SingleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}

func CreateLogger(module string, logfile string, loglevel string) (log *logging.Logger, lf *os.File) {
	log = logging.MustGetLogger(module)
	var err error
	if logfile != "" {
		lf, err = os.OpenFile(logfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Errorf("Cannot open logfile %v: %v", logfile, err)
		}
		//defer lf.CloseInternal()

	} else {
		lf = os.Stderr
	}
	backend := logging.NewLogBackend(lf, "", 0)
	backendLeveled := logging.AddModuleLevel(backend)
	backendLeveled.SetLevel(logging.GetLevel(loglevel), "")

	logging.SetFormatter(_logformat)
	logging.SetBackend(backendLeveled)

	return
}

func RpcContextMetadata2(ctx context.Context) (traceId string, sourceInstance string, err error) {
	md, ok := metadata.FromIncomingContext(ctx)
	//md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		return "", "", errors.New("no metadata in context")
	}

	// check for sourceInstance Metadata
	si, exists := md["sourceinstance"]
	if !exists {
		return "", "", errors.New("no sourceinstance in context")
	}
	sourceInstance = si[0]

	// check for targetInstance Metadata
	tr, exists := md["traceid"]
	if !exists {
		return "", "", errors.New("no traceid in context")
	}
	traceId = tr[0]

	return
}
func RpcContextMetadata(ctx context.Context) (traceId string, sourceInstance string, targetInstance string, err error) {
	md, ok := metadata.FromIncomingContext(ctx)
	//md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		return "", "", "", errors.New("no metadata in context")
	}

	// check for targetInstance Metadata
	ti, exists := md["targetinstance"]
	if !exists {
		return "", "", "", errors.New("no targetinstance in context")
	}
	targetInstance = ti[0]

	// check for sourceInstance Metadata
	si, exists := md["sourceinstance"]
	if !exists {
		return "", "", "", errors.New("no sourceinstance in context")
	}
	sourceInstance = si[0]

	// check for targetInstance Metadata
	tr, exists := md["traceid"]
	if !exists {
		return "", "", "", errors.New("no traceid in context")
	}
	traceId = tr[0]

	return
}

func Openbrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
//		log.Fatal(err)
	}

}