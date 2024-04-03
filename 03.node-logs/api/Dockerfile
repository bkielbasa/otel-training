FROM node:16-alpine

WORKDIR /usr/src/app
COPY package*.json ./
RUN npm install
COPY . .
EXPOSE 3333

CMD ["node", "app.js"]
