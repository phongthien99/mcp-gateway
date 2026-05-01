FROM mcp-gateway-hugo-site:latest AS dev

WORKDIR /src

EXPOSE 1313

CMD ["server", "--bind=0.0.0.0", "--port=1313", "--disableFastRender", "--noBuildLock", "--poll=700ms"]

FROM dev AS build

COPY . .
RUN hugo --minify

FROM nginx:alpine AS prod

COPY --from=build /src/public /usr/share/nginx/html
