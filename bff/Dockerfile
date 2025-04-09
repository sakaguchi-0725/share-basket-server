FROM node:20.16.0-alpine AS build

WORKDIR /app

COPY package.json package-lock.json ./
RUN npm ci

COPY . .
RUN npm run build

FROM node:20.16.0-alpine AS runner

WORKDIR /app

COPY --from=build /app/node_modules ./node_modules
COPY --from=build /app/dist ./dist
COPY --from=build /app/src/**/*.graphql ./src/

EXPOSE 3000

CMD ["node", "dist/src/main"]

