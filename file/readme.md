## go 文件相关操作

### 文件相关

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
      fmt.Println("文件存在")
   } else {
      f, err1 = os.Create(filename) //创建文件
      fmt.Println("文件不存在")
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
      fmt.Println("文件存在")
   } else {
      f, err1 = os.Create(filename) //创建文件
      fmt.Println("文件不存在")
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
      fmt.Println("文件存在")
   } else {
      f, err1 = os.Create(filename) //创建文件
      fmt.Println("文件不存在")
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
