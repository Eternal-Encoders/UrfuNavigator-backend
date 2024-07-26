FROM node:18-alpine as base
WORKDIR /home/node

COPY ./package.json ./package.json

RUN yarn install

COPY ./src ./src
COPY ./nodemon.json ./nodemon.json
COPY ./tsconfig.json ./tsconfig.json

ENV NODE_ENV=production
ENV PORT_ENV 80

RUN yarn build

EXPOSE 80

CMD ["node", "dist/server.js"]