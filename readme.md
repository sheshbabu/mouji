# mouji


Minimal Web Analytics
<p align="center"><img src="https://raw.githubusercontent.com/sheshbabu/mouji/master/docs/screenshots/home.png" /></p>

### Features
* Track pageviews
* Multiple projects support
* No cookies used
* No personal data collected
* Single binary deployment
* Low resource usage
* Minimal dependency footprint


### Philosophy
I built this application mainly for my own personal use. I also wanted to experiment to see if it's possible to build good software using modern technologies that consumes less resources and pull in as few dependencies as possible.

This is built using Golang and uses SQLite as database. The frontend is built using vanilla HTML and CSS. So far, the only dependencies it uses are [Lato](https://www.latofonts.com) (vendored), [go-sqlite](https://www.github.com/mattn/go-sqlite3) and [crypto](https://pkg.go.dev/golang.org/x/crypto).


### Installation
Build from source
```shell
$ go build
```


### Local Development
Run the application using default configuration
```shell
$ go run main.go
```

Run the application by customizing PORT and DATA_FOLDER
```shell
$ PORT="4000" DATA_FOLDER="./data" go run main.go
```

Run the application using Docker
```shell
$ docker build -t mouji . && docker run --rm -it -v $(pwd):/data mouji
```


### Schema Migrations
* Create new migration file under `./migrations`
* Use the format `<version>_<title>.sql`
