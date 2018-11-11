FROM ubuntu:18.04

MAINTAINER Tarasov Vladislav

ENV TZ=Europe/Moscow
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

# Обвновление списка пакетов

# Back to the root user
USER root

RUN apt-get -y update && apt-get install -y --no-install-recommends apt-utils

#
# Сборка проекта
#

# Установка golang
RUN apt install -y golang-1.10 git

# Выставляем переменную окружения для сборки проекта
ENV GOROOT /usr/lib/go-1.10
ENV GOPATH /opt/go
ENV PATH $GOROOT/bin:$GOPATH/bin:/usr/local/go/bin:$PATH

RUN go get -u github.com/golang/dep/cmd/dep

# Копируем исходный код в Docker-контейнер
WORKDIR $GOPATH/src/github.com/SinimaWath/tp-db/
ADD . $GOPATH/src/github.com/SinimaWath/tp-db/

RUN go install ./vendor/github.com/go-swagger/go-swagger/cmd/swagger

RUN make generate
RUN make install

# Объявлем порт сервера
EXPOSE 5000

#
# Установка postgresql
#
ENV PGVER 10
RUN apt-get install -y postgresql-$PGVER

# Run the rest of the commands as the ``postgres`` user created by the ``postgres-$PGVER`` package when it was ``apt-get installed``
USER postgres

# Create a PostgreSQL role named ``docker`` with ``docker`` as the password and
# then create a database `docker` owned by the ``docker`` role.
RUN /etc/init.d/postgresql start &&\
    psql --command "CREATE USER docker WITH SUPERUSER PASSWORD 'docker';" &&\
    createdb -O docker docker &&\
    psql docker -a -f assets/db/postgres/create.sql &&\
    /etc/init.d/postgresql stop

USER root
# Adjust PostgreSQL configuration so that remote connections to the
# database are possible.

RUN echo "host all all 0.0.0.0/0 md5" >> /etc/postgresql/$PGVER/main/pg_hba.conf &&\
    echo "listen_addresses='*'" >> /etc/postgresql/$PGVER/main/postgresql.conf &&\
    echo "unix_socket_directories = '/var/run/postgresql'" >> /etc/postgresql/$PGVER/main/postgresql.conf
# Expose the PostgreSQL port
EXPOSE 5432

# Add VOLUMEs to allow backup of config, logs and databases
VOLUME  ["/etc/postgresql", "/var/log/postgresql", "/var/lib/postgresql"]

#
# Запускаем PostgreSQL и сервер
#

CMD service postgresql start && forum-server --scheme=http --port=5000 --host=0.0.0.0 --database=postgres://docker:docker@localhost/docker