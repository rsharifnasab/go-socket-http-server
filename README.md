## Simple Web server in Go
in order to solve Computer assignment of Computer Networks in SBU (Spring 1400) I implemented a web server in Go

more info about project could be found [here](./project-description.pdf)

this web server can only handle one specific request: request for a file in file system

for example:

```sh
curl -v 127.0.0.1/index.html
```



## how to use 

this program is written in pure go, so for compilation you need to install `go>=1.16`

after that that would be enough to run:

```sh
# compile and run 
go build 
sudo ./http-server

# just run
go run main.go
```

##### important note

don't forgot to use `sudo` because opening port 80 needs root access, of you are a normal user you can setup server on another port with this command:

```sh
go build 
./http-server -port 8080
# or
go run main.go -port 8080
```



## how to use

after running on port 80,

+ just open browser and visit `localhost`
+ add your files to `static` folder, for example add `music.mp3` to static folder and view it in address `localhost/music.mp3`
+ if you want to use curl or another friends, yo can see `requests.sh` file.



## Implementation details

you can read about details of implementation in `report.md`





