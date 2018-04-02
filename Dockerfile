FROM ubuntu:14.04
COPY workgen /app/
WORKDIR "/app"
ENTRYPOINT [ "/bin/sh" ]