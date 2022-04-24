FROM golang:1.18.1-bullseye

ARG USERNAME=gouser
ARG GROUPNAME=${USERNAME}
ARG UID=1000
ARG GID=1000

RUN groupadd -g ${GID} ${GROUPNAME} && \
    useradd -m -s /bin/bash -u ${UID} -g ${GID} ${USERNAME}

RUN apt-get update && apt-get install -y --no-install-recommends \
  mariadb-client \
  sudo \
  gettext \
  && apt-get -y clean \
  && rm -rf /var/lib/apt/lists/* \
  && echo "${USERNAME} ALL=(ALL) NOPASSWD:ALL" > /etc/sudoers.d/${USERNAME}

USER ${USERNAME}

RUN go install -v github.com/rubenv/sql-migrate/...@v1.1.1

WORKDIR /workspace

