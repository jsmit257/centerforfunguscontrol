FROM nginx:bookworm AS build
COPY ./nginx/cffc.conf /cffc.conf
COPY ./www /www
ENTRYPOINT /bin/bash -c "envsubst </cffc.conf >/etc/nginx/nginx.conf && exec nginx -g 'daemon off;'"
