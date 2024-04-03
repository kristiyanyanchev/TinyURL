import requests
import sys
from threading import Thread


def test(url):
    for i in range(1,1000):
        x = requests.get(url)
        if x.status_code != 200:
            break
        print(i,x.status_code)

input = input("Enter the URL: ")
for i in range(1,100):
    thread = Thread(target=test(input))
    thread.start()
