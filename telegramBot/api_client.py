import aiohttp
import os


API_URL = os.getenv("API_URL", "http://backend:8000/api")
async def get_or_create_user(tg_id: int, username: str, first_name: str):
    async with aiohttp.ClientSession() as session:
        payload = {
            "telegram_id": tg_id,
            "username": username or "",
            "first_name": first_name or ""
        }
        async with session.post(f"{API_URL}/users", json=payload) as resp:
            return await resp.json()

async def get_all_trips():
    async with aiohttp.ClientSession() as session:
        async with session.get(f"{API_URL}/trips") as resp:
            return await resp.json()

async def create_trip(driver_tg_id: int, trip_data: dict):
    async with aiohttp.ClientSession() as session:
        # Передаем driver_tg_id в параметрах, как это делает axios во Vue
        url = f"{API_URL}/trips?driver_tg_id={driver_tg_id}"
        async with session.post(url, json=trip_data) as resp:
            return await resp.json()