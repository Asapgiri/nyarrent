package pages

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
	"time"
)

var artifact_path string = "artifacts/"
var html_path string = "html/"
var base_template_path string = html_path + "base.html"

var file_types = map[string]string {
    "html": "text",
    "css":  "text",
}

func sizeToText(size int) string {
    const kbDiv = 1024.0

    mb := float64(size) / kbDiv / kbDiv
    gb := mb / kbDiv

    if gb >= 1.0 {
        return strconv.FormatFloat(gb, 'f', 2, 64) + " GB"
    } else {
        return strconv.FormatFloat(mb, 'f', 2, 64) + " MB"
    }
}

var funcMap = template.FuncMap {
    "inc":      func(i int) int {return i + 1},
    "dec":      func(i int) int {return i - 1},
    "size":     sizeToText,
    "timegt":   func(a time.Time, b time.Time) bool {return b.Compare(a) > 0},
    "timelt":   func(a time.Time, b time.Time) bool {return b.Compare(a) <= 0},
    "now":      time.Now,
}

func read_artifact(path string, header http.Header) (string, string) {
    var dir_path string

    ex, err := os.Executable()
    if nil != err {
        panic(err)
    }

    parts := strings.Split(path, ".")
    file_type := parts[len(parts)-1]
    if "html" == file_type {
        dir_path = filepath.Dir(ex) + "/" + html_path
    } else {
        dir_path = filepath.Dir(ex) + "/" + artifact_path
    }

    file_read, err := os.ReadFile(dir_path + path)
    if nil != err {
        not_found, _ := os.ReadFile(filepath.Dir(ex) + "/" + html_path + "not_found.html")
        return string(not_found), "text"
    }

    if nil != header {
        _, file_ok := file_types[file_type]
        if file_ok {
            header.Set("Content-Type", file_types[file_type] + "/" + file_type)
        }
    }

    return string(file_read), file_type
}

func Render(w http.ResponseWriter, temp string, dto any) {
    tmp, err := template.ParseFiles(base_template_path)
    if nil != err {
        io.WriteString(w, "Templating error!")
        return
    }

    if nil != dto {
        var tpl bytes.Buffer
        dto_tmp, err := template.New("Dto").Funcs(funcMap).Parse(temp)
        if nil != err {
            return
        }
        dto_tmp.Execute(&tpl, dto)
        temp = tpl.String()
    }

    var tpl bytes.Buffer
    tmp.Execute(&tpl, temp)
    main, err := template.New("Main").Funcs(funcMap).Parse(tpl.String())
    if nil != err {
        io.WriteString(w, "Templating error 2!" + err.Error())
        return
    }

    main.Execute(w, tpl.String())
}

// Prerender does not support session if you don't pass it...
func Pre_render(temp string, dto any) string {
    var tpl bytes.Buffer

    tmp, err := template.New("Dto").Funcs(funcMap).Parse(temp)
    if nil != err {
        return err.Error()
    }
    tmp.Execute(&tpl, dto)

    return tpl.String()
}
