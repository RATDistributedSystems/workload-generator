FROM nginx
COPY nginx.conf /etc/nginx/
COPY load-balancer.conf /etc/nginx/conf.d/
RUN mv /etc/nginx/conf.d/default.conf /etc/nginx/conf.d/default.conf.disabled
COPY workgen files create-conf.sh /app/
WORKDIR "/app"
EXPOSE 44445
ENTRYPOINT ./create-conf.sh /etc/nginx/conf.d/load-balancer.conf && tail -f /var/log/nginx/access.log