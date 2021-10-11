FROM node:lts as build

COPY . .

RUN npm install && \
    npm run build_worker && \
    npm run build

FROM node:alpine

ENV HOST=0.0.0.0

COPY --from=build /node_modules/ /node_modules/
COPY --from=build /public/ /public/
COPY --from=build /package.json /package.json
COPY --from=build /package-lock.json /package-lock.json
COPY --from=build /config.dev.json /config.dev.json

CMD [ "npm", "start" ]