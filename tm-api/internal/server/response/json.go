package response

import (
	"fmt"
	"html/template"

	"github.com/gin-gonic/gin"

	"log"
	"net/http"
	"path/filepath"
	"runtime"
	"time"
)

const UserError = "USER ERROR: "

type Data struct {
	Message   string      `json:"message,omitempty"`
	Data      interface{} `json:"data,omitempty"`
	Timestamp string      `json:"timestamp,omitempty"`
	Errors    string      `json:"errors,omitempty"`
	Status    string      `json:"status,omitempty"`
}

// type Data struct {
//	Status    bool   `json:"status"`
//	Message   string `json:"message"`
//	Data      any    `json:"data"`
//	Error     string `json:"error"`
//	Timestamp string `json:"timestamp"`
//}

func HTML(c *gin.Context, path string, data interface{}) {
	// use only instance of template
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	tmpl, err := template.ParseFiles(fmt.Sprintf("%s/../../%s", basepath, path))
	if err != nil {
		log.Printf("error executing HTML template: %v", err.Error())
		http.Error(c.Writer, "internal server error", http.StatusInternalServerError)
		return
	}
	if err := tmpl.ExecuteTemplate(c.Writer, tmpl.Name(), data); err != nil {
		log.Printf("error executing HTML template: %v", err.Error())
		http.Error(c.Writer, "internal server error", http.StatusInternalServerError)
	}
}

func JSON(c *gin.Context, status int, message string, data interface{}, err error) {
	errMessage := ""
	if err != nil {
		errMessage = err.Error()
		log.Println(errMessage)
	}
	responsedata := Data{
		Message:   message,
		Data:      data,
		Errors:    errMessage,
		Status:    http.StatusText(status),
		Timestamp: time.Now().Format(time.RFC850),
	}

	c.JSON(status, responsedata)
}
