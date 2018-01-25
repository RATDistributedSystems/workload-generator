FROM scratch

COPY workgen workload_files /app/
WORKDIR "/app"
CMD ["./workgen", "1userWorkLoad"]
