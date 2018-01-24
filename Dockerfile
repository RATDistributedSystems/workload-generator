FROM scratch

COPY workgen workload_files /app/
WORKDIR "/app"
EXPOSE 44441
CMD ["./workgen", "1userWorkLoad"]
