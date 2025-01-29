# Используем официальный образ NGINX
FROM nginx:latest

# Копируем наш конфигурационный файл в контейнер
COPY nginx.conf /etc/nginx/nginx.conf

# Открываем порт 8080 для внешнего доступа
EXPOSE 80

# Запускаем NGINX
CMD ["nginx", "-g", "daemon off;"]