FROM postgres:13
COPY sql/* /docker-entrypoint-initdb.d/
COPY config/postgres.conf /etc/postgresql/postgres.conf
RUN chown postgres:postgres /var/log/postgresql/