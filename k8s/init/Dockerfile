FROM alpine:latest

RUN apk --no-cache add mysql-client && apk add bash

COPY k8s/init/init.sh /scripts/

RUN chmod +x /scripts/init.sh

ENTRYPOINT ["/bin/bash", "-c", "/scripts/init.sh"]
