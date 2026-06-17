ARG DOCKER_PROXY_HOST_1=docker.io
FROM $DOCKER_PROXY_HOST_1/node:18 as builder

WORKDIR /app

COPY ./package*.json /app/
RUN npm install --fund false --audit false --silent
COPY . /app
RUN npm run story:build

FROM $DOCKER_PROXY_HOST_1/nginx:1.25-alpine
WORKDIR /var/www/html
COPY .nginx/nginx.conf /etc/nginx/nginx.conf
COPY --from=builder --chown=nginx /app/dist .
EXPOSE 80
