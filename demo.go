package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

// 主页
func IndexView(w http.ResponseWriter, r *http.Request) {
	html := LoadHtml("index.html")
	w.Write(html)

}

// Upload
func Upload(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.Write(LoadHtml("upload.html"))
	}
	if r.Method == "POST" {
		f, h, err := r.FormFile("file")
		if err != nil {
			w.Write([]byte("文件上传失败：" + err.Error()))
			return
		}
		checkImage := h.Header.Get("Content-Type")
		if !strings.Contains(checkImage, "image") {
			w.Write([]byte("<html><a href=\"/upload\">请上传图片</a></html>"))
			return
		}
		root := "./static1/"
		// 路径，权限
		os.Mkdir(root, os.ModePerm)
		out, err := os.Create(root + h.Filename)
		if err != nil {
			io.WriteString(w, "文件创建失败："+err.Error())
			return
		}
		_, err = io.Copy(out, f)
		if err != nil {
			io.WriteString(w, "文件保存失败:"+err.Error())
			return
		}
		//io.WriteString(w, "上传成功："+root+h.Filename)
		// 重定向
		// 302:临时重定向
		http.Redirect(w, r, "/detail?name="+h.Filename, 302)
		return
	}

}

// Show
func ImageView(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() // 把 url 或 form 表单的数据解析到对应的容器 r.Form
	name := r.Form.Get("name")
	fmt.Println(name)
	root := "./static1/"
	f, err := os.Open(root + name)
	if err != nil {
		w.WriteHeader(404)
		return
	}
	defer f.Close()
	w.Header().Set("Content Type", "image")
	io.Copy(w, f)
}

func DetailView(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	name := r.Form.Get("name")
	html := LoadHtml("detail.html")
	html = bytes.Replace(html, []byte("@src"), []byte("/image?name="+name), 1)
	w.Write(html)
}
func LoadHtml(name string) []byte {
	f, err := os.Open("./views/" + name)
	if err != nil {
		return []byte("<html><head></head><body>Errors：" + err.Error() + "</body></html>")
	}

	defer f.Close()

	buf, err := ioutil.ReadAll(f)
	if err != nil {
		return []byte("<html><head></head><body>Errors：" + err.Error() + "</body></html>")
	}
	return buf
}

func ListView(w http.ResponseWriter, r *http.Request) {
	html := LoadHtml("list.html")
	names, err := ioutil.ReadDir("./static1")
	if err != nil {
		w.Write([]byte("<html><head></head><body>Errors：" + err.Error() + "</body></html>"))
	}
	tmp := ""
	for i := 0; i < len(names); i++ {
		//log.Println(names[i].Name())
		tmp += `<li><a href="/detail?name="` + names[i].Name() + `"><img src="/image?name=` + names[i].Name() + `"alt="未发现"></a> </li>`
	}
	html = bytes.Replace(html, []byte("@html"), []byte(tmp), 1)
	w.Write(html)
}

func main() {
	fmt.Println("hello")
	http.HandleFunc("/index", IndexView)
	http.HandleFunc("/upload", Upload)
	http.HandleFunc("/image", ImageView)
	http.HandleFunc("/detail", DetailView)
	http.HandleFunc("/list", ListView)
	log.Println("run...")
	http.ListenAndServe(":8080", nil)
}
