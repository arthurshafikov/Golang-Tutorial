FROM alpine:3.9

ENV GOOSE_FILE /bin/goose
ENV MIGRATIONS /migrations

ADD https://github.com/pressly/goose/releases/download/v2.6.0/goose-linux64 ${GOOSE_FILE}
RUN chmod +x ${GOOSE_FILE}

WORKDIR ${MIGRATIONS}

# Add docker-compose-wait tool -------------------
ADD https://github.com/ufoscout/docker-compose-wait/releases/download/2.7.2/wait /wait
RUN chmod +x /wait

CMD /wait && \
    ${GOOSE_FILE} postgres "${DSN}" up
