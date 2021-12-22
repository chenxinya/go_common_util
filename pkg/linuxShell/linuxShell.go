package pkg

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/chenxinya/go_common_util/pkg"
	log "github.com/sirupsen/logrus"
	"io"
	"os/exec"
	"strings"
	"time"
)


// RunInLinux 运行命令在Linux系统中
func RunInLinux(cmd string) (string, error) {
	startT := time.Now()		//计算当前时间
	result, err := exec.Command("/bin/sh", "-c", cmd).Output()
	if err != nil {
		return strings.TrimSpace(string(result)), err
	}
	defer func(start time.Time) {
		log.Printf("Running Linux cmd:%v, cost time : %d ms",cmd,time.Since(startT)/pkg.NanosecondToMillisecond)
	}(startT)
	return strings.TrimSpace(string(result)), err
}

// ExecCommand 运行命令在Linux系统中 返回 错误日志和正常日志
func ExecCommand(commandName string) (string,string,error) {
	var (
		outMsg []byte
		errMsg []byte
	)
	log.Infof("linux shell execCommand : %s",commandName)
	//函数返回一个*Cmd，用于使用给出的参数执行name指定的程序
	cmd := exec.Command("/bin/sh", "-c", commandName)

	//显示运行的命令
	log.Info(cmd.Args)
	//StdoutPipe方法返回一个在命令Start后与命令标准输出关联的管道。Wait方法获知命令结束后会关闭这个管道，一般不需要显式的关闭该管道。
	stdout, outErr := cmd.StdoutPipe()
	stderr, errErr := cmd.StderrPipe()

	if outErr != nil || errErr != nil{
		return "","",nil
	}

	errStart := cmd.Start()
	if errStart != nil {
		return "", "", errStart
	}

	//创建一个流来读取管道内内容，这里逻辑是通过一行一行的读取的
	reader := bufio.NewReader(stdout)
	var index int64
	//实时循环读取输出流中的一行内容
	for {
		line,_ ,err := reader.ReadLine()
		if err != nil || io.EOF == err {
			if io.EOF == err{
				log.Print("命令执行结束!")
			}else{
				log.Printf("命令out结束或者跳出循环：%v\n",err)
			}
			break
		}
		index++
		if index>5000 {
			log.Printf("输出OUT日志信息超过%d 行\n",index)
			break
		}
		outMsg=BytesCombine(outMsg,line)
	}

	var errIndex int64
	//创建一个流来读取管道内内容，这里逻辑是通过一行一行的读取的
	errReader := bufio.NewReader(stderr)
	//实时循环读取输出流中的一行内容
	for {
		errLine,_ ,err4 := errReader.ReadLine()
		//line, err2 := reader.ReadString('\n')
		if err4 != nil || io.EOF == err4 {
			fmt.Printf("命令err结束或者跳出循环：%v\n",err4)
			break
		}
		errIndex++
		if errIndex>1000 {
			log.Printf("输出ERR日志信息超过%d 行\n",errIndex)
			break
		}
		errMsg=BytesCombine(errMsg,errLine)
	}

	//阻塞直到该命令执行完成，该命令必须是被Start方法开始执行的
	err :=cmd.Wait()
	return string(outMsg),string(errMsg),err
}

//BytesCombine 多个[]byte数组合并成一个[]byte
func BytesCombine(pBytes ...[]byte) []byte {
	return bytes.Join(pBytes, []byte("\n"))
}
