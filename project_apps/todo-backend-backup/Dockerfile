FROM google/cloud-sdk:alpine

RUN apk add --no-cache postgresql-client bash

COPY backup.sh /backup.sh
RUN chmod +x /backup.sh

CMD ["/backup.sh"]
