import requests
from faker import Faker

r = requests.post('https://protected-bayou-62297.herokuapp.com/users/login',
                  data={'email': 'dkarg@gmail.com', 'password': '123456'})

print(r.status_code)
print(r.headers)

tocken = r.cookies['session_id']
cookie = dict(session_id=tocken)

fake = Faker()

for i in range(100):
    title = fake.street_name()
    price = fake.building_number()
    city = fake.city()
    description = fake.text(128)

    r = requests.post('https://protected-bayou-62297.herokuapp.com/ads/new', cookies=cookie,
                      data={'title': title, 'price': price, 'city': city, 'description_ad': description})
    print('Published:', r.status_code)
