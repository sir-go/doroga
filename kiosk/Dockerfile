#FROM node:12-alpine3.9 as builder
#WORKDIR /app
#COPY package*.json ./
#RUN yarn global add @quasar/cli
#COPY . .
#RUN yarn
#RUN quasar build

FROM nginx:stable-alpine
#COPY --from=builder /app/dist/spa /usr/share/nginx/html
COPY dist /usr/share/nginx/html
EXPOSE 8085
CMD ["nginx", "-g", "daemon off;"]
