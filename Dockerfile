FROM scratch
COPY calculator /
ENTRYPOINT ["/calculator"]