
FROM node:17-alpine
RUN mkdir -p /home/node/app/node_modules && chown -R node:node /home/node/app
WORKDIR /home/node/app
COPY package*.json ./
USER root
RUN npm install && npm install typescript -g
RUN npm i --save-dev @types/node
COPY . .
RUN tsc index.ts
EXPOSE 8080
ENV PORT 8080
EXPOSE $PORT
CMD [ "node", "index.js" ]
