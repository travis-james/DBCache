FROM redis:8

RUN mkdir -p /data

COPY data/localdev.rdb /data/dump.rdb

CMD ["redis-server", "--dir", "/data"]