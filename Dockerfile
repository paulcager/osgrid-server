FROM node:12.15.0
WORKDIR /usr/src/app
COPY package*.json ./
RUN npm --version && node --version
RUN npm install
COPY . .
EXPOSE 9090
CMD ["node", "-r", "esm", "index.js" ]

