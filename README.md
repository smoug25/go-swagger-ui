# Go Swagger UI

Cross-platform [Swagger UI](https://swagger.io/swagger-ui/) built with golang with only
a fat binary, which can help you serve your swagger documentation with only one command.

## Usage

### Like a module in your golang code

```shell
go install github.com/smoug25/go-swagger-ui
```
Code example

```go
package main

import "net/http"
import serv "github.com/smoug25/go-swagger-ui"

func main() {
    // Set swagger file local path or url 
	serv.SetSwaggerFile("swagger/swagger.json")
    // add swagger ui handler 
    http.HandleFunc("/", serv.Serv)
    // start you http server
    http.ListenAndServe(":8080", nil)	
}
```


### Like a server
Download the latest release for your system on [release](https://github.com/haxii/go-swagger-ui/releases) page.

Start a swagger documentation server on port `8000` for `/path/to/your/swagger.json`:

```
$ swaggerui -l "0.0.0.0:8000" -f "/path/to/your/swagger.json"
```

Everyone can read your documentation on `http://your.ip.address:8000`, you can also view other online swagger file, such as swagger's [petstore](http://petstore.swagger.io/),
by opening 

```
http://your.ip.address:8000/?url=http://petstore.swagger.io/v2/swagger.json
```

or serve a local swagger files folder defined with `-d` flag, such as `/swagger` as pre-defined, then you can view the `path/to/my-awesome-api.yaml` file in this folder by opening

```
http://your.ip.address:8000/?file=path/to/my-awesome-api.yaml
```

you can also replace the default API host endpoint with `host` query for local host swagger file

```
http://your.ip.address:8000/?file=path/to/my-awesome-api.yaml&host=www.new.host
``` 

More detailed usage:

```
Usage of swaggerui:

  -b    enable the topbar
  -d string
        swagger files vhost dir (default "/swagger")
  -f string
        swagger url or local file path (default "http://petstore.swagger.io/v2/swagger.json")
  -l string
        server's listening Address (default ":8080")
  -s string
        Send signal to a master process: install, remove, start, stop, status (default "status")
```

## Build

Source code is written in [go](https://golang.org/), [make](https://www.gnu.org/software/make/) and [xxd](https://www.systutorials.com/docs/linux/man/1-xxd/) is also needed to build the binary.

The Swagger UI's HTML pages in [dist](dist) folder is copied from the original [Swagger UI Source](https://github.com/swagger-api/swagger-ui/tree/master/dist), then converted into bytes array in file [static/static.go](static/static.go) using `xxd` command, to re-build it

```
$ make generate
```

By typing the following command, you can get the cross platform distributions of this program


```
$ make build
```
