package server

import (
	"fmt"
	"github.com/smoug25/go-swagger-ui/static"
	jsoniter "github.com/json-iterator/go"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	querySwaggerURLKey  string = "url"
	querySwaggerFileKey string = "file"
	querySwaggerHost    string = "host"
)

var (
	ServerAddr  = ":8080"
	SwaggerPath = "/"
	SwaggerFile = "http://petstore.swagger.io/v2/swagger.json"
	LocalSwaggerDir = "/swagger"
	EnableTopbar = false
	IsNativeSwaggerFile   = false
	NativeSwaggerFileName = ""
)

func Serv(w http.ResponseWriter, r *http.Request) {
	source := getSource(r)
	// serve the local file
	localFile := ""
	if IsNativeSwaggerFile && source == NativeSwaggerFileName {
		localFile = SwaggerFile
	} else if strings.HasPrefix(source, "swagger/") {
		// we treat path started with swagger as a direct request of a local swagger file
		localFile = filepath.Join(LocalSwaggerDir, source[len("swagger/"):])
	}
	if len(localFile) > 0 {
		serveLocalFile(localFile, w, r)
		return
	}
	// find the in-memory static files
	staticFile, err := getFile(source)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	// set up the content type
	setContentType(w, source)
	// return back the non-index files
	if source != "index.html" {
		w.Header().Set("Content-Length", strconv.Itoa(len(staticFile)))
		w.Write(staticFile)
		return
	}
	// set up the index page
	indexHTML := prepareIndexPage(r, staticFile)
	w.Header().Set("Content-Length", strconv.Itoa(len(indexHTML)))
	fmt.Fprint(w, indexHTML)
}

func getSource(r *http.Request) string {
	source := r.URL.Path[len(SwaggerPath):]
	if len(source) == 0 {
		source = "index.html"
	}
	return source
}

func getFile(source string) ([]byte,error) {
	// find the in-memory static files
	var staticFile []byte
	file, err := static.EmbedStatic.Open(source)
	if err != nil {
		return nil, err
	}
	staticFile, err =  ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return staticFile, err
}

func serveLocalFile(localFilePath string, w http.ResponseWriter, r *http.Request) {
	newHost := r.URL.Query().Get("host")
	if len(newHost) == 0 {
		http.ServeFile(w, r, localFilePath)
		return
	}
	isJSON := false
	switch filepath.Ext(localFilePath) {
	case ".json":
		isJSON = true
	case ".yaml":
		fallthrough
	case ".yml":
		isJSON = false
	default:
		http.Error(w, "unknown swagger file: "+localFilePath, http.StatusBadRequest)
		return
	}
	// open file
	file, err := os.Open(localFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			http.Error(w, "file not exists", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()
	swg := new(map[string]interface{})
	if isJSON {
		dec := jsoniter.NewDecoder(file)
		err = dec.Decode(swg)
	} else {
		dec := yaml.NewDecoder(file)
		err = dec.Decode(swg)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	(*swg)["host"] = newHost
	var resp []byte
	if isJSON {
		resp, err = jsoniter.Marshal(swg)
	} else {
		resp, err = yaml.Marshal(swg)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Cache-Control", "no-cache, max-age=0, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Write(resp)
}

func setContentType(w http.ResponseWriter, source string) {
	switch filepath.Ext(source) {
	case ".html":
		w.Header().Set("Content-Type", "text/html")
	case ".js":
		w.Header().Set("Content-Type", "application/javascript")
	case ".css":
		w.Header().Set("Content-Type", "text/css")
	default:
		w.Header().Set("Content-Type", "application/octet-stream")
	}
}

func prepareIndexPage(r *http.Request, staticFile []byte) string {
		targetSwagger := SwaggerFile
		if f := r.URL.Query().Get(querySwaggerFileKey); len(f) > 0 {
			// requesting a local file, join it with a `swagger/` prefix
			base, err := url.Parse("swagger/")
			if err != nil {
				return ""
			}
			target, err := url.Parse(f)
			if err != nil {
				return ""
			}
			targetSwagger = base.ResolveReference(target).String()
			if h := r.URL.Query().Get(querySwaggerHost); len(h) > 0 {
				targetSwagger += "?host=" + h
			}
		} else if _url := r.URL.Query().Get(querySwaggerURLKey); len(_url) > 0 {
			// deal with the query swagger firstly
			targetSwagger = _url
		} else if IsNativeSwaggerFile {
			// for a native swagger file, use the filename directly
			targetSwagger = NativeSwaggerFileName
		}
	// replace the target swagger file in index
	indexHTML := string(staticFile)
		indexHTML = strings.Replace(indexHTML,
			"https://petstore.swagger.io/v2/swagger.json",
			targetSwagger, -1)
		if EnableTopbar {
			indexHTML = strings.Replace(indexHTML,
				"SwaggerUIBundle.plugins.DownloadUrl, HideTopbarPlugin",
				"SwaggerUIBundle.plugins.DownloadUrl", -1)
		}
	return indexHTML
}