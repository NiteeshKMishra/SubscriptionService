# SubscriptionService
SubscriptionService is a web app written in Golang that allows users to select different subscriptions.
This app has server side rendering using go templates. This app highlights the use of concurrency in golang using channels, waitGroups and mutexes


## Local Setup

Before running this application. Please download `Docker Desktop`, `make` and `go >= 1.20` in your system
and start docker desktop

Clone this repo locally and switch to this repo folder
```
https://github.com/NiteeshKMishra/SubscriptionService.git
```

Run below command to pull and run all the docker images required for this app
```
docker-compose up
```

Run below command to download all go packages required for running this app
```
go mod tidy
```

Create `secrets.env` at root of the directory and copy the contents of file `secrets.example.env` into `secrets.env`

Run the app with
```
make start or make restart
```

Run all tests with
```
make test
```
