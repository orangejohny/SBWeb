# Some examples of usage

## Curl commands for testing API

### User interaction

curl -v -X POST -H "Content-Type: application/x-www-form-urlencoded; charset=utf-8" --data "email=karg@gmail.com&password=123456&first_name=dmitry&last_name=kargashin" 127.0.0.1:8080/users/new

curl -v -X POST -H "Content-Type: application/x-www-form-urlencoded; charset=utf-8" --data "email=karg@gmail.com&password=123456" 127.0.0.1:8080/users/login

curl -v 127.0.0.1:8000/users/profile

curl -v --cookie "session_id=RzBgl-VlTcFOTZF5zq1JtN9tOiCnbtCHzNDGk6Yxa6I=; Expires=Mon, 05 Nov 2019 19:35:00 GMT" 127.0.0.1:8080/users/profile

curl -v -X POST --data "first_name=dmitry&last_name=kargashin&tel_number=892157655&about=heh ehehh eeh&id=5" --cookie "session_id=; Expires=Mon, 05 Nov 2018 19:35:00 GMT" 127.0.0.1:8000/users/profile

curl -v -X POST --data "first_name=dmitry&last_name=kargashin&tel_number=892157655&about=heh ehehh eeh&id=5" --cookie "session_id=Wb8jMSneczBuOF7N54l3kAtOgNoMIkS3ms0P0hE4leU=; Expires=Mon, 05 Nov 2018 19:35:00 GMT" 127.0.0.1:8000/users/profile

curl -v -X POST --cookie "session_id=Wb8jMSneczBuOF7N54l3kAtOgNoMIkS3ms0P0hE4leU=; Expires=Mon, 05 Nov 2018 19:35:00 GMT" 127.0.0.1:8000/users/logout

### Ad intercation

curl -v -X POST -H "Content-Type: application/x-www-form-urlencoded; charset=utf-8" --data 'city=London&title=House Building&description_ad=I can build any house that you want!' --cookie "session_id=a_ArmtPH3YzqBp8ZHxSab3H2ovQ2378hibe_jN1BHzY=; Expires=Sat, 10 Nov 2018 13:21:48 GMT" 127.0.0.1:8080/ads/1

curl -v -X DELETE -H "Content-Type: application/x-www-form-urlencoded; charset=utf-8" --cookie "session_id=0bfcGYzlTjgQydcSAnIWftr-bZAwa5Ki02ZDh_v3wco=; Expires=Sat, 10 Nov 2018 13:21:48 GMT" 127.0.0.1:8000/ads/delete/3

## Install & Run

go install -i ./internal/service-main.go
/home/orangejohny/workspace/go/bin/service-main -cfg ./internal/config.json

try to trigger gitlab-ci