FROM alpine

RUN apk update
RUN apk add \
    ca-certificates \
    wget
RUN update-ca-certificates
RUN apk add tzdata && \
    cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime && \
    echo "Asia/Tokyo" > /etc/timezone && \
    apk del tzdata

ARG ENTRYKIT_VERSION=0.4.0
RUN wget -O - -q https://github.com/progrium/entrykit/releases/download/v${ENTRYKIT_VERSION}/entrykit_${ENTRYKIT_VERSION}_Linux_x86_64.tgz | tar zxf - && \
    mv entrykit /bin/entrykit && \
    chmod +x /bin/entrykit && \
    entrykit --symlink

ARG SISITO_API_VERSION=0.2.6
RUN wget -O - -q https://github.com/winebarrel/sisito-api/releases/download/v${SISITO_API_VERSION}/sisito-api-v${SISITO_API_VERSION}-linux-amd64.gz \
      | gunzip > /sbin/sisito-api && \
    chmod +x /sbin/sisito-api

COPY files/ /

ENTRYPOINT [ \
  "render", \
    "/etc/sisito-api.toml", \
  "--", \
  "switch", \
    "shell=/bin/bash", \
  "--", \
  "/sbin/sisito-api", "-config", "/etc/sisito-api.toml" \
]
