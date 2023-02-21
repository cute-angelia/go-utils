## go 文件相关操作

### 文件相关

[写了 30 多个 Go 常用文件操作的示例，收藏这一篇就够了](https://mp.weixin.qq.com/s/dczWeHW6JWSJMJx1nBx7rA)

- 读取-本地按行读取文件

```
file, err := os.Open("app-2019-06-01.log")
if err != nil {
    log.Fatal(err)
}
defer file.Close()
scanner := bufio.NewScanner(file)
for scanner.Scan() {
  lineText := scanner.Text()
}
```

- 读取-整个文件读取

```
b, err := ioutil.ReadFile("app-2019-06-01.log") // just pass the file name
    if err != nil {
        fmt.Print(err)
    }

str := string(b) // convert content to a 'string'
fmt.Println(str) // print the content as a 'string'
```

- 写入-文件 ioutil.WriteFile

```
// 这种方式每次都会覆盖 test.txt内容，如果test.txt文件不存在会创建。
package main

import (
   "io/ioutil"
)

func main() {

   content := []byte("测试1\n测试2\n")
   err := ioutil.WriteFile("test.txt", content, 0644)
   if err != nil {
      panic(err)
   }
}

```

- 写入-文件 os

```
// 此种方法可以在文件内容末尾添加新内容。
package main

import (
   "fmt"
   "io"
   "os"
)

func checkFileIsExist(filename string) bool {
   if _, err := os.Stat(filename); os.IsNotExist(err) {
      return false
   }
   return true
}
func main() {
   var wireteString = "测试1\n测试2\n"
   var filename = "./test.txt"
   var f *os.File
   var err1 error
   if checkFileIsExist(filename) { //如果文件存在
      f, err1 = os.OpenFile(filename, os.O_APPEND, 0666) //打开文件
      log.Println("文件存在")
   } else {
      f, err1 = os.Create(filename) //创建文件
      log.Println("文件不存在")
   }
   defer f.Close()
   n, err1 := io.WriteString(f, wireteString) //写入文件(字符串)
   if err1 != nil {
      panic(err1)
   }
   fmt.Printf("写入 %d 个字节n", n)
}
```

- 写入-文件 f.Write

```
package main

import (
   "fmt"
   "os"
)

func checkFileIsExist(filename string) bool {
   if _, err := os.Stat(filename); os.IsNotExist(err) {
      return false
   }
   return true
}
func main() {
   var str = "测试1\n测试2\n"
   var filename = "./test.txt"
   var f *os.File
   var err1 error
   if checkFileIsExist(filename) { //如果文件存在
      f, err1 = os.OpenFile(filename, os.O_APPEND, 0666) //打开文件
      log.Println("文件存在")
   } else {
      f, err1 = os.Create(filename) //创建文件
      log.Println("文件不存在")
   }
   defer f.Close()
   n, err1 := f.Write([]byte(str)) //写入文件(字节数组)

   fmt.Printf("写入 %d 个字节n", n)
   n, err1 = f.WriteString(str) //写入文件(字符串)
   if err1 != nil {
      panic(err1)
   }
   fmt.Printf("写入 %d 个字节n", n)
   f.Sync()
}
```

- 写入-文件 bufio

```
package main

import (
   "bufio"
   "fmt"
   "os"
)

func checkFileIsExist(filename string) bool {
   if _, err := os.Stat(filename); os.IsNotExist(err) {
      return false
   }
   return true
}
func main() {
   var str = "测试1\n测试2\n"
   var filename = "./test.txt"
   var f *os.File
   var err1 error
   if checkFileIsExist(filename) { //如果文件存在
      f, err1 = os.OpenFile(filename, os.O_APPEND, 0666) //打开文件
      log.Println("文件存在")
   } else {
      f, err1 = os.Create(filename) //创建文件
      log.Println("文件不存在")
   }
   defer f.Close()
   if err1 != nil {
      panic(err1)
   }
   w := bufio.NewWriter(f) //创建新的 Writer 对象
   n, _ := w.WriteString(str)
   fmt.Printf("写入 %d 个字节n", n)
   w.Flush()
}

```

### 定义的方法

- OpenLocalFile

打开已经存在的文件， 不存在会新建一个， 返回 `*os.File`

- 下载文件并保存到指定路径, more is see test.go

```
DownloadFileWithSrc("xxx.jpg", "/tmp/xxx.jpg")

```

### 测试

```

go test -v -run TestA$ file_test.go

```


### other

```
ioutil 包实现了一些 I/O 实用函数。

func NopCloser(r io.Reader) io.ReadCloser
返回一个包裹起给定 Reader r 的 ReadCloser ， 这个 ReadCloser 带有一个无参数的 Close 方法。

func ReadAll(r io.Reader) ([]byte, error)
对 r 进行读取， 直到发生错误或者遇到 EOF 为止， 然后返回被读取的数据。

func ReadDir(dirname string) ([]os.FileInfo, error)
读取 dirname 指定的目录， 并返回一个根据文件名进行排序的目录节点列表

func ReadFile(filename string) ([]byte, error)
读取名字为 filename 的文件并返回文件中的内容。

func TempDir(dir, pattern string) (name string, err error)
在目录 dir 中新创建一个带有指定前缀 prefix 的临时目录， 然后返回该目录的路径。

func TempFile(dir, pattern string) (f *os.File, err error)
在目录 dir 新创建一个名字带有指定前缀 prefix 的临时文件， 以可读写的方式打开它， 并返回一个 *os.File 指针。

func WriteFile(filename string, data []byte, perm os.FileMode) error
将给定的数据 data 写入到名字为 filename 的文件里面。
```


### http 读取文本

```

// 根据字段名获取表单文件
	formFile, _, err := r.FormFile("file")
	if err != nil {
		log.Printf("Get form file failed: %s\n", err)
		return
	}
	defer formFile.Close()

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, formFile); err != nil {
		log.Println(err)
		api.Error(w, r, "", -1)
	} else {
		scanner := bufio.NewScanner(buf)
		for scanner.Scan() {
			lineText := scanner.Text()
			if len(lineText) > 0 {
				// kindid := r.FormValue("kindid")
				// kindidii, _ := strconv.Atoi(kindid)
				// insert(lineText, int32(kindidii))
			}
		}
	}


```