import requests
from faker import Faker
from os import listdir
from os.path import isfile, join
import random

path = '/home/orangejohny/Downloads/images'
files = [join(path, f) for f in listdir(path) if isfile(join(path, f))]

r = requests.post('https://search-build.herokuapp.com/users/login',
                  data={'email': 'dkargashin3@gmail.com', 'password': '123456'})

tocken = r.cookies['session_id']
cookie = dict(session_id=tocken)

fake = Faker()

for i in range(50):
    title = fake.street_name()
    price = fake.building_number()
    city = fake.city()
    description = fake.text(128)
    file = random.choice(files)

    r = requests.post('https://search-build.herokuapp.com/ads/new', cookies=cookie, files={'images':open(file, 'rb')},
                      data={'title': title, 'price': price, 'city': city, 'description_ad': description})
    print('Published:', r.status_code)

requests.post('https://search-build.herokuapp.com/users/logout', cookies=cookie)