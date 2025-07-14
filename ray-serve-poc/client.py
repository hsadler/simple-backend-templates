import requests

text = "I love this product!"
url = "http://localhost:8000/analyze"
urlencoded_text = requests.utils.quote(text, safe='')
resp = requests.post(f"{url}/{urlencoded_text}")
print(resp.json())
