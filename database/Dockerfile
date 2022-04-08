FROM postgres:13

ENV POSTGRES_PASSWORD=pass
ENV POSTGRES_USER=postgres
EXPOSE 5432

COPY ./create_database.sql /docker-entrypoint-initdb.d/


