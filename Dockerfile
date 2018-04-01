FROM ubuntu:14.04
COPY workgen /app/
WORKDIR "/app"
CMD ["/bin/sh"]