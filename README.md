# Goira - Go Infrared API
<p>
  <a href="https://travis-ci.com/illuminati1911/goira"><img alt="Build Status" src="https://travis-ci.com/illuminati1911/goira.svg?branch=master" target="_blank" /></a>
  <a href="https://goreportcard.com/report/github.com/illuminati1911/goira"><img alt="Go Report Card" src="https://goreportcard.com/badge/github.com/illuminati1911/goira" target="_blank" /></a>
  <img alt="Documentation" src="https://img.shields.io/badge/documentation-yes-brightgreen.svg" target="_blank" />
  <img alt="Maintenance" src="https://img.shields.io/badge/Maintained%3F-yes-green.svg" target="_blank" />
  <img alt="License: MIT" src="https://img.shields.io/badge/License-MIT-yellow.svg" target="_blank" />
</p>

Goira (Go Infrared API) is a go backend for Raspberry Pi, which exposes API routes to control GPIO connected infrared LED's. Currently Goira only supports controlling air conditioners.


## Features
- Uses [IRBlaster](https://github.com/illuminati1911/IRBlaster) linux kernel driver to manage the IR transmission
- Exposes Restful API
- Only password authentication

## Limitations
- No HTTPS yet.
- Limitations of the [IRBlaster](https://github.com/illuminati1911/IRBlaster) kernel driver
- AC Config to binary mapping has to be made specifically for each AC manufacturer. See: `goira/internal/accontrol/mappers/`
- By default only supports ChangHong AC's
 
## Tech stack
- Language: Golang
- Database: BoltDB
- UUID: Google UUID
- Architecture: [Clean Architecture by Uncle Bob](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)

## Installation
First, make sure you have installed [IRBlaster](https://github.com/illuminati1911/IRBlaster) kernel driver.

You can use the GNU make to build just the backend service or both front-end and backend.
```bash
$ make build     #build just the backend 
$ make frontend  #build just the front-end
$ make full      #build both front-end and backend
```

Resulting binary and front-end files will be placed in the `build` folder.

## Usage
After you run the binary, server will start at `localhost:8080`.

If the frontend has been installed, open `http://localhost:8080/` to access the web app.

### API
#### Login
- Type: POST
- Route: `/login`
- Parameters:
```
{
  "password": <password_here>
}
```
- Response: 200 OK for successful login
#### Get AC state
- Type: GET
- Route: `/status`
- Response 200 OK for success with body (example):
```
{
  "temp": 20,      #20 celcius
  "wind": 0,       #wind mode 0
  "mode": 0,       #ac mode 0
  "active": false  #turned off
}
```
#### Set AC state
- Type: POST
- Route: `/state`
- Parameters (example):
```
#example 1
{
  "temp": 29,
  "active": false,
  "mode": 2
}

#example 2
{
  "active": true
}
```
- Response 200 OK for success

## Contributing
- Pull requests are welcome.
- For major changes, please open an issue first to discuss what you would like to change.
- Remember to use `gofmt` to format the code before making the pull request.

## Roadmap
### Near future
- Multiple manufacturer support (Hitachi etc.)
### Some day
- Support for other kinds of IR devices.

## References
- [Clean Architecture by Uncle Bob](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
## License
[MIT](https://choosealicense.com/licenses/mit/)
