FROM node:15.3.0
WORKDIR /usr/src/app
COPY package*.json ./
RUN npm install
COPY . .
EXPOSE 9090
CMD ["node", "-r", "esm", "index.js" ]

