FROM alpine
RUN apk add ca-certificates && rm -rf /var/cache/apk/*
COPY build/autocomplete-linux-amd64 /autocomplete
ENTRYPOINT [ "/autocomplete" ]
