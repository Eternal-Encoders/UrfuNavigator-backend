FROM node:18-alpine as build-stage
WORKDIR /home/node

COPY ./package.json ./package.json
COPY ./yarn.lock ./yarn.lock

RUN yarn install

COPY ./src ./src
COPY ./nodemon.json ./nodemon.json
COPY ./tsconfig.json ./tsconfig.json

ENV NODE_ENV=production
ENV PORT_ENV 80

RUN yarn build

FROM node:18-alpine as base
WORKDIR /home/node

ENV PAYLOAD_CONFIG_PATH=dist/payload.config.js 
ENV NODE_ENV=production

COPY --from=build-stage /home/node/node_modules /home/node/node_modules
COPY --from=build-stage /home/node/dist /home/node/dist
COPY --from=build-stage /home/node/package.json /home/node/package.json
COPY --from=build-stage /home/node/build /home/node/build

EXPOSE 80

CMD ["node", "dist/server.js"]