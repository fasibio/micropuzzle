# pull official base image
FROM node:10.16-slim
COPY . ./
RUN npm install
RUN npm run build


FROM nginx:stable-alpine

COPY --from=0 ./build /usr/share/nginx/html
COPY ./nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
