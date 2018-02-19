FROM ubuntu:14.04

COPY workgen files /app/
WORKDIR "/app"
CMD ["./workgen", "-f","45User_testWorkLoad"]
