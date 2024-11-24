<p align="center">
  <img width="256" src="assets/android-chrome-512x512.png">
  <h1 align="center">mouji</h1>
  <p align="center">Minimal Web Analytics</p>
</p>

<p align="center"><img src="https://github.com/sheshbabu/mouji/blob/master/docs/screenshots/home.png?raw=true" /></p>

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

This is built using Golang and uses SQLite as database. The frontend is built using vanilla HTML, CSS and JS. So far, the only dependencies it uses are [Lato](https://www.latofonts.com) (vendored), [go-sqlite](https://www.github.com/mattn/go-sqlite3) and [crypto](https://pkg.go.dev/golang.org/x/crypto).


### Installation
Build from source
```shell
$ go build
```


### Local Development
Run the application using default configuration
```shell
$ make dev
```

Run the application by customizing PORT and DATA_FOLDER
```shell
$ PORT="4000" DATA_FOLDER="./data" make dev
```

Run the application in watch mode
Install [air](https://github.com/air-verse/air)
```shell
$ make watch
```

Run the application using Docker
```shell
$ docker build -t mouji . && docker run --rm -it -v $(pwd):/data mouji
```


### Schema Migrations
* Create new migration file under `./migrations`
* Use the format `<version>_<title>.sql`
