FROM node:18-alpine as base
WORKDIR /home/node

COPY ./package.json ./package.json

RUN yarn install --network-timeout 100000

COPY ./src ./src
COPY ./nodemon.json ./nodemon.json
COPY ./tsconfig.json ./tsconfig.json

ARG REDIS_URI_ARG

ENV NODE_ENV=production
ENV NODE_OPTIONS="--max-old-space-size=4096"
ENV PORT_ENV 80
ENV REDIS_URI_ENV $REDIS_URI_ARG

RUN yarn build

EXPOSE 80

CMD ["node", "dist/server.js"]