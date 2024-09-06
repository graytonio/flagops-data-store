FROM alpine

COPY flagops-data-store /usr/bin/flagops-data-store
COPY assets/dist/ /assets/


ENTRYPOINT [ "/usr/bin/flagops-data-store" ]