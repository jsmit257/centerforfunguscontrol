FROM postgres:bookworm as build

COPY ./sql/create.sql /docker-entrypoint-initdb.d/01-create.sql
COPY ./sql/init.sql /docker-entrypoint-initdb.d/02-init.sql
COPY ./sql/seed.sql /docker-entrypoint-initdb.d/03-seed.sql

ENV POSTGRES_HOST_AUTH_METHOD trust

RUN grep -v 'exec "$@"' /usr/local/bin/docker-entrypoint.sh > /docker-entrypoint.sh && chmod 755 /docker-entrypoint.sh
RUN /docker-entrypoint.sh postgres

FROM postgres:bookworm
RUN mkdir /pgbackup
COPY ./bin/migration-entrypoint.sh /migration-entrypoint.sh
COPY ./bin/backup-entrypoint.sh /backup-entrypoint.sh
COPY ./bin/restore-entrypoint.sh /restore-entrypoint.sh
COPY --from=build /var/lib/postgresql/data /var/lib/postgresql/data