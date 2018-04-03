FROM ubuntu:14.04
COPY workgen /app/
WORKDIR "/app"
ENTRYPOINT [ "tail", "-f", "/dev/null" ]