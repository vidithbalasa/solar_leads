import requests
import io
import os
from PIL import Image

def get_aerial_view(address, api_key):
    base_url = "https://maps.googleapis.com/maps/api/staticmap"
    params = {
        "center": address,
        "zoom": 20,
        "size": "600x600",
        "maptype": "satellite",
        "key": api_key
    }

    response = requests.get(base_url, params=params)

    if response.status_code == 200:
        img = Image.open(io.BytesIO(response.content))
        img.show()
    else:
        print(f"Error: {response.status_code}")

# Replace 'YOUR_API_KEY' with your actual Google Maps API key
api_key = os.environ.get("GOOGLE_MAPS_API_KEY")

# address = input("Enter the address: ")
address = "825 30th St, Oakland, CA 94608"
get_aerial_view(address, api_key)
