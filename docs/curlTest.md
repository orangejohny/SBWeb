# Some examples of usage

## Curl commands for testing API

### User interaction

curl -v -X POST -H "Content-Type: application/x-www-form-urlencoded; charset=utf-8" --data "email=karg@gmail.com&password=123456&first_name=dmitry&last_name=kargashin" 127.0.0.1:8080/users/new

curl -v -X POST -H "Content-Type: application/x-www-form-urlencoded; charset=utf-8" --data "email=karg@gmail.com&password=123456" 127.0.0.1:8080/users/login

curl -v 127.0.0.1:8000/users/profile

curl -v --cookie "session_id=RzBgl-VlTcFOTZF5zq1JtN9tOiCnbtCHzNDGk6Yxa6I=; Expires=Mon, 05 Nov 2019 19:35:00 GMT" 127.0.0.1:8080/users/profile

curl -v -X POST --data "first_name=dmitryyyyy&last_name=kargashin&tel_number=892157655&about=heh ehehh eeh&id=21" --cookie "session_id=rGrjp9pjNic6waxXZvbTQfW-5qUMDQws7hLlyv07-c8=; Expires=Mon, 05 Nov 2018 19:35:00 GMT" 127.0.0.1:8080/users/profile

curl -v -X POST --data "first_name=dmitry&last_name=kargashin&tel_number=892157655&about=heh ehehh eeh&id=5" --cookie "session_id=Wb8jMSneczBuOF7N54l3kAtOgNoMIkS3ms0P0hE4leU=; Expires=Mon, 05 Nov 2018 19:35:00 GMT" 127.0.0.1:8000/users/profile

curl -v -X POST --cookie "session_id=Wb8jMSneczBuOF7N54l3kAtOgNoMIkS3ms0P0hE4leU=; Expires=Mon, 05 Nov 2018 19:35:00 GMT" 127.0.0.1:8000/users/logout

### Ad intercation

curl -v -X POST -H "Content-Type: application/x-www-form-urlencoded; charset=utf-8" --data 'city=London&title=House Building&description_ad=I can build any house that you want!' --cookie "session_id=a_ArmtPH3YzqBp8ZHxSab3H2ovQ2378hibe_jN1BHzY=; Expires=Sat, 10 Nov 2018 13:21:48 GMT" 127.0.0.1:8080/ads/1

curl -v -X DELETE -H "Content-Type: application/x-www-form-urlencoded; charset=utf-8" --cookie "session_id=0bfcGYzlTjgQydcSAnIWftr-bZAwa5Ki02ZDh_v3wco=; Expires=Sat, 10 Nov 2018 13:21:48 GMT" 127.0.0.1:8000/ads/delete/3

## Install & Run

go install -i ./internal/service-main.go
/home/orangejohny/workspace/go/bin/service-main -cfg ./internal/config.json

try to trigger gitlab-ci

## Create user with avatar

curl -v -F "images=@/home/orangejohny/Downloads/simpsons_PNG88.png" -F "first_name=dmitry" -F "last_name=kargashin" -F "email=dkarg@gmail.com" -F "password=123456" 127.0.0.1:8080/users/new

## Login user

curl -v -X POST -H "Content-Type: application/x-www-form-urlencoded; charset=utf-8" --data "email=dkarg@gmail.com&password=123456" 127.0.0.1:8080/users/login

## Create ad with multiple images

curl -v -F "title=Building Some Things" -F "city=New York" -F "description_ad=I can build some buildings in New York." -F "images=@/home/orangejohny/Downloads/Msd.png" -F "images=@/home/orangejohny/Downloads/Photographer_Barnstar.png" -F "images=@/home/orangejohny/Downloads/png.png" --cookie "session_id=FpXOTUwLRYFkXJbHSXayUggXTQTCOa-cHk6lsCljwp0=; Expires=Sun, 02 Dec 2018 00:24:08 GMT; HttpOnly" 127.0.0.1:8080/ads/new

## Add new image to ad

curl -v -F "title=Building Some Things" -F "city=New York" -F "description_ad=I can build some buildings in New York." -F "ad_images=/images/a545472deeee37220ae046fa2a8d4085.png" -F "ad_images=/images/c57c515b6bc91f5187b4eb77d2a4667b.png" -F "ad_images=/images/17ca761dab4009ee2f2f3c280b396703.png" -F "images=@/home/orangejohny/Downloads/images.jpeg" --cookie "session_id=FpXOTUwLRYFkXJbHSXayUggXTQTCOa-cHk6lsCljwp0=; Expires=Sun, 02 Dec 2018 00:24:08 GMT; HttpOnly" 127.0.0.1:8080/ads/edit/2

## Delete ad with images

curl -v --cookie "session_id=FpXOTUwLRYFkXJbHSXayUggXTQTCOa-cHk6lsCljwp0=; Expires=Sun, 02 Dec 2018 00:24:08 GMT; HttpOnly" 127.0.0.1:8080/ads/delete/2

## Update user with avatar

curl -v -F "images=@/home/orangejohny/Downloads/PNG_transparency_demonstration_1.png" -F "first_name=ddddmitry" -F "last_name=karg" --cookie "session_id=rGrjp9pjNic6waxXZvbTQfW-5qUMDQws7hLlyv07-c8=; Expires=Mon, 05 Nov 2018 19:35:00 GMT" 127.0.0.1:8080/users/profile