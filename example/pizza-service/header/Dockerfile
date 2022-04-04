# pull official base image
FROM node:slim
COPY . ./
RUN npm install
RUN npm run build

FROM nginx:stable-alpine
COPY --from=0 ./dist /usr/share/nginx/html
COPY ./nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
