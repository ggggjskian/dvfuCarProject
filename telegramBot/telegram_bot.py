import asyncio
import os
import logging
from datetime import datetime, timedelta
from aiogram import Bot, Dispatcher, F, types
from aiogram.filters import CommandStart
from aiogram.fsm.context import FSMContext
from aiogram.fsm.state import State, StatesGroup
from aiogram.types import ReplyKeyboardMarkup, KeyboardButton, InlineKeyboardMarkup, InlineKeyboardButton
from aiogram.exceptions import TelegramBadRequest
import api_client

BOT_TOKEN = os.getenv("BOT_TOKEN")

logging.basicConfig(level=logging.INFO)

bot = Bot(token=BOT_TOKEN)
dp = Dispatcher()


main_kb = ReplyKeyboardMarkup(
    keyboard=[
        [KeyboardButton(text="🚗 Найти поездку"), KeyboardButton(text="➕ Создать поездку")]
    ],
    resize_keyboard=True
)


def get_trips_inline_kb():
    return InlineKeyboardMarkup(
        inline_keyboard=[
            [InlineKeyboardButton(text="🔄 Обновить список", callback_data="refresh_trips")],
            [InlineKeyboardButton(text="🔙 Назад", callback_data="close_trips")]
        ]
    )


class CreateTrip(StatesGroup):
    trip_type = State()
    point = State()
    price = State()


async def generate_trips_text():
    trips = await api_client.get_all_trips()
    

    current_time = datetime.now().strftime("%H:%M:%S")
    
    if not trips:
        return f"Пока нет доступных поездок 😔\n\n_Последнее обновление: {current_time}_"

    text = "📍 **Доступные поездки:**\n\n"
    for t in trips:
        t_type = "В кампус" if t["trip_type"] == "to_campus" else "Из кампуса"
        text += f"🚗 **{t_type}** | Точка: {t['point']}\n"
        text += f"💰 Цена: {t['price']}₽ | 💺 Мест: {t['seats_total']}\n"
        text += f"📅 Выезд: {t['departure_time'][:16].replace('T', ' ')}\n"
        text += f"Команда для брони: /book_{t['id']}\n\n"
        
    text += f"🔄 _Последнее обновление: {current_time}_"
    return text



@dp.message(CommandStart())
async def cmd_start(message: types.Message):
    await api_client.get_or_create_user(
        tg_id=message.from_user.id,
        username=message.from_user.username,
        first_name=message.from_user.first_name
    )
    await message.answer("Привет! Я бот для карпулинга. Что хочешь сделать?", reply_markup=main_kb)

@dp.message(F.text == "🚗 Найти поездку")
async def show_trips_first_time(message: types.Message):
    text = await generate_trips_text()
    await message.answer(text, parse_mode="Markdown", reply_markup=get_trips_inline_kb())

@dp.callback_query(F.data == "refresh_trips")
async def refresh_trips_handler(callback: types.CallbackQuery):
    text = await generate_trips_text()
    
    try:
        await callback.message.edit_text(
            text=text, 
            parse_mode="Markdown", 
            reply_markup=get_trips_inline_kb()
        )
    except TelegramBadRequest:

        pass 
        
    await callback.answer("Список обновлен!")

@dp.callback_query(F.data == "close_trips")
async def close_trips_handler(callback: types.CallbackQuery):
    await callback.message.delete()
    await callback.answer()


@dp.message(F.text == "➕ Создать поездку")
async def start_creating_trip(message: types.Message, state: FSMContext):
    kb = ReplyKeyboardMarkup(
        keyboard=[[KeyboardButton(text="from_campus"), KeyboardButton(text="to_campus")]],
        resize_keyboard=True,
        one_time_keyboard=True
    )
    await message.answer("Откуда едем? Выбери тип:", reply_markup=kb)
    await state.set_state(CreateTrip.trip_type)

@dp.message(CreateTrip.trip_type)
async def process_trip_type(message: types.Message, state: FSMContext):
    await state.update_data(trip_type=message.text)
    await message.answer("Напиши адрес (точку сбора/высадки):", reply_markup=types.ReplyKeyboardRemove())
    await state.set_state(CreateTrip.point)

@dp.message(CreateTrip.point)
async def process_point(message: types.Message, state: FSMContext):
    await state.update_data(point=message.text)
    await message.answer("Укажи цену (просто число, например 200):")
    await state.set_state(CreateTrip.price)

@dp.message(CreateTrip.price)
async def process_price(message: types.Message, state: FSMContext):
    data = await state.get_data()
    departure = (datetime.utcnow() + timedelta(hours=1)).isoformat() + "Z"
    
    trip_data = {
        "trip_type": data['trip_type'],
        "point": data['point'],
        "point_lat": 0.0,
        "point_lon": 0.0,
        "departure_time": departure,
        "seats_total": 4,
        "price": int(message.text),
        "comment": "Создано через бота",
        "max_deviation_km": 3,
        "time_flexibility_minutes": 30
    }

    await api_client.create_trip(message.from_user.id, trip_data)
    
    await state.clear()
    await message.answer("✅ Поездка успешно создана!", reply_markup=main_kb)


@dp.message(Command("test"))
async def cmd_test_local(message: types.Message):
    # Telegram разрешает HTTP только для 127.0.0.1 или localhost!
    local_kb = InlineKeyboardMarkup(
        inline_keyboard=[
            [InlineKeyboardButton(
                text="💻 Тест локально (Без Ngrok)", 
                web_app=WebAppInfo(url="http://127.0.0.1:5173")
            )]
        ]
    )
    await message.answer("Нажми, чтобы открыть локальный Vue:", reply_markup=local_kb)

async def main():
    await dp.start_polling(bot)

if __name__ == "__main__":
    asyncio.run(main())