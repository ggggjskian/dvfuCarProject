# Используем Node.js для сборки и запуска Vue
FROM node:20-alpine

WORKDIR /app

# Копируем package.json и устанавливаем зависимости
COPY package*.json ./
RUN npm install

# Копируем весь остальной код фронтенда
COPY . .

# Открываем порт Vite
EXPOSE 5173

# Запускаем Vite в режиме разработки. 
# Флаг --host 0.0.0.0 ОБЯЗАТЕЛЕН, чтобы Docker "выпустил" сервер наружу.
CMD ["npm", "run", "dev", "--", "--host", "0.0.0.0"]