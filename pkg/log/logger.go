package log

import (
	"fmt"
	"github.com/chenxinya/go_common_util/pkg"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

func InitLogrus(logDir string, level logrus.Level,logFileName string) {
	logrus.SetFormatter(&logFormatter{
		TimestampFormat: pkg.TimeLayOut,
	})

	// set rotate log writer
	writer, _ := rotatelogs.New(
		filepath.Join(logDir, logFileName+".log.%Y%m%d"),
		rotatelogs.WithLinkName(logDir+logFileName+".log"),
		rotatelogs.WithRotationTime(24*time.Hour),
	)
	//设置日志输出地点
	if runtime.GOOS == "windows" {
		log.SetOutput(os.Stdout)
	}else{
		logrus.SetOutput(writer)
	}

	// set caller
	logrus.SetReportCaller(true)

	// set level
	logrus.SetLevel(level)
}

type logFormatter struct {
	TimestampFormat string
}

func (f *logFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	formattedTime := entry.Time.Format(f.TimestampFormat)
	level := strings.ToUpper(entry.Level.String())
	shortFile := TrimmedPath(entry.Caller.File)

	// 序列化map
	content := fmt.Sprintf("%s\t%s\t", formattedTime, level)

	if len(entry.Data) > 0 {
		var dataList []string
		for key, value := range entry.Data {
			dataList = append(dataList, fmt.Sprintf("%s: %+v", key, value))
		}

		content += "{" + strings.Join(dataList, ",") + "}\t"
	}

	content += fmt.Sprintf("%s:%d\t%s\n", shortFile, entry.Caller.Line, entry.Message)

	return []byte(content), nil
}

func TrimmedPath(filePath string) string {
	idx := strings.LastIndexByte(filePath, '/')
	if idx == -1 {
		return filePath
	}

	idx = strings.LastIndexByte(filePath[:idx], '/')
	if idx == -1 {
		return filePath
	}

	return filePath[idx+1:]
}

