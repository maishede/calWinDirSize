package main

import (
	"fmt"
	"math"
	"os"
	"path"
	"sort"
	"strings"
)

func StaticFileSize(dir string) int64 {
	var totalSize int64 = 0
	dirList, err := os.ReadDir(dir)
	if err != nil {
		// fmt.Printf("读取文件夹: %s失败: %v\n", dir, err)
		return 0
	}
	for _, v := range dirList {
		fileName := v.Name()
		fileInfo, _ := v.Info()
		fileSize := fileInfo.Size()
		fileIsDir := v.IsDir()
		if fileIsDir {
			fileSize = StaticFileSize(path.Join(dir, fileName)) // 新路径
			totalSize += fileSize
		} else {
			totalSize += fileSize
		}
	}
	return totalSize
}

func TransByte2XB(b int64) string {
	var bFloat float64
	bFloat = float64(b)
	if bFloat < 1024 {
		return fmt.Sprintf("%.2f字节", bFloat)
	} else if bFloat >= 1024 && bFloat < math.Pow(1024, 2) {
		bFloat /= 1024
		return fmt.Sprintf("%.2fKB", bFloat)
	} else if bFloat >= math.Pow(1024, 2) && bFloat < math.Pow(1024, 3) {
		bFloat /= math.Pow(1024, 2)
		return fmt.Sprintf("%.2fMB", bFloat)
	} else if bFloat >= math.Pow(1024, 3) && bFloat < math.Pow(1024, 4) {
		bFloat /= math.Pow(1024, 3)
		return fmt.Sprintf("%.2fGB", bFloat)
	} else if bFloat >= math.Pow(1024, 4) && bFloat < math.Pow(1024, 5) {
		bFloat /= math.Pow(1024, 4)
		return fmt.Sprintf("%.2fTB", bFloat)
	} else {
		return "123"
	}
}

func main() {
	var maxNameLength int = 0
	var maxPathLength int = 0
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("获取当前路径失败:%v\n", err)
	}
	fmt.Printf("当前文件位置:%v\n", currentDir)
	var dir string
	fmt.Printf("请输入目标文件夹地址(默认为当前地址[%s]):", currentDir)
	fmt.Scanln(&dir)
	var sorted string
	fmt.Printf("是否按文件(夹)大小倒序展示(Y/N,默认N):")
	fmt.Scanln(&sorted)
	if sorted == "" {
		sorted = "N"
	}

	if dir == "" {
		dir = currentDir
	}
	fmt.Printf("查找文件地址:%s\n", dir)
	// 查看当前第一层的文件
	dirList, err := os.ReadDir(dir)
	if err != nil {
		fmt.Printf("读取当前文件夹失败:%v\n", err)
	}
	var nameSizeMap map[string]int64
	nameSizeMap = make(map[string]int64, len(dirList))
	for _, v := range dirList {
		fileName := v.Name()
		fileInfo, _ := v.Info()
		fileSize := fileInfo.Size()
		fileIsDir := v.IsDir()
		if StringLength(fileName) > maxNameLength {
			maxNameLength = StringLength(fileName)
		}
		if StringLength(path.Join(dir, fileName)) > maxPathLength {
			maxPathLength = StringLength(path.Join(dir, fileName))
		}
		if fileName == "$RECYCLE.BIN" {
			continue
		}
		if fileName == "System Volume Information" {
			continue
		}
		if fileIsDir {
			fmt.Printf("正在读取文件:%s\r", strings.Repeat(" ", maxPathLength+4))
			fmt.Printf("正在读取文件:%s\r", path.Join(dir, fileName))
			fileSize = StaticFileSize(path.Join(dir, fileName))
			nameSizeMap[fileName] = fileSize
		} else {
			nameSizeMap[fileName] = fileSize
		}
	}
	sorted = strings.ToUpper(sorted)
	fmt.Printf("正在读取文件:%s\r", strings.Repeat(" ", maxPathLength+4))
	if sorted == "Y" {
		SortPrint(nameSizeMap, maxNameLength)
	} else {
		for k, v := range nameSizeMap {
			fmt.Printf("文件名:%s%s\t文件大小:%s\n", k, strings.Repeat(" ", maxNameLength-StringLength(k)), TransByte2XB(v))
		}
	}
}

type sortStruct struct {
	fileName string
	fileSize int64
}

type bySize []sortStruct

func (b bySize) Len() int           { return len(b) }
func (b bySize) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b bySize) Less(i, j int) bool { return b[i].fileSize < b[j].fileSize }

func SortPrint(sm map[string]int64, maxNameLength int) {
	newList := []sortStruct{}
	for k, v := range sm {
		newList = append(newList, sortStruct{k, v})
	}
	// sort.Sort(bySize(newList))
	sort.Sort(sort.Reverse(bySize(newList)))

	for _, v := range newList {
		fmt.Printf("文件名:%s%s\t文件大小:%s\n", v.fileName, strings.Repeat(" ", maxNameLength-StringLength(v.fileName)), TransByte2XB(v.fileSize))
	}
}

func CalculateStringLength(s string) int64 {
	slen := len(s)
	return int64(slen)
}

func StringLength(s string) int {
	bLen := len(s)                 // 字节长度byte len
	sLen := len([]rune(s))         // 字符长度string len
	cAmount := (bLen - sLen) / 2   // 中文数量
	eAmount := (3*sLen - bLen) / 2 // 英文数量
	// fmt.Printf("字符串:%s,中文数量:%d, 英文数量:%d\n", s, cAmount, eAmount)
	return eAmount + 2*cAmount
}
