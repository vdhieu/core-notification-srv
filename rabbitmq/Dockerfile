FROM rabbitmq:3.10.6-management-alpine

EXPOSE 15672

COPY defs.json /etc/definitions/defs.json
COPY rabbitmq.conf /etc/rabbitmq/conf.d/11-np-config.conf
