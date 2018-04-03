FROM ubuntu:14.04
COPY workgen files /app/
WORKDIR "/app"
ENTRYPOINT [ "tail", "-f", "/dev/null" ]