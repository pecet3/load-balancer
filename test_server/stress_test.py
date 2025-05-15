import threading
import requests

URL = "http://localhost:8080?someQuery="""
NUM_REQUESTS = 1000

def make_request():
    try:
        resp = requests.get(URL)
        print(f"Status: {resp.status_code}")
    except Exception as e:
        print(f"Error: {e}")

threads = []

for _ in range(NUM_REQUESTS):
    t = threading.Thread(target=make_request)
    t.start()
    threads.append(t)

for t in threads:
    t.join()

print("Done.")
