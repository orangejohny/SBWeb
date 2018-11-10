# Some examples of usage

## Curl commands for testing API

curl -v -X POST -H "Content-Type: application/x-www-form-urlencoded; charset=utf-8" --data "email=karg@gmail.com&password=123456&first_name=dmitry&last_name=kargashin" 127.0.0.1:8000/users/new

curl -v -X POST -H "Content-Type: application/x-www-form-urlencoded; charset=utf-8" --data "email=karg@gmail.com&password=123456" 127.0.0.1:8000/users/login

curl -v 127.0.0.1:8000/users/profile

curl -v --cookie "session_id=Wb8jMSneczBuOF7N54l3kAtOgNoMIkS3ms0P0hE4leU=; Expires=Mon, 05 Nov 2018 19:35:00 GMT" 127.0.0.1:8000/users/profile

curl -v -X POST --data "first_name=dmitry&last_name=kargashin&tel_number=892157655&about=heh ehehh eeh&id=5" --cookie "session_id=Wb8jMSneczBuOF7N54l3kAtOgNoMIkS3ms0P0hE4leU=; Expires=Mon, 05 Nov 2018 19:35:00 GMT" 127.0.0.1:8000/users/profile

curl -v -X POST --data "first_name=dmitry&last_name=kargashin&tel_number=892157655&about=heh ehehh eeh&id=5" --cookie "session_id=Wb8jMSneczBuOF7N54l3kAtOgNoMIkS3ms0P0hE4leU=; Expires=Mon, 05 Nov 2018 19:35:00 GMT" 127.0.0.1:8000/users/profile

curl -v -X POST --cookie "session_id=Wb8jMSneczBuOF7N54l3kAtOgNoMIkS3ms0P0hE4leU=; Expires=Mon, 05 Nov 2018 19:35:00 GMT" 127.0.0.1:8000/users/logout

## Install & Run

go install -i ./internal/service-main.go
/home/orangejohny/workspace/go/bin/service-main -cfg ./internal/config.json