FROM alpine
ENTRYPOINT [ "/usr/bin/flagops-data-store" ]
COPY flagops-data-store /usr/bin/flagops-data-store