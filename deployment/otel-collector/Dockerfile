FROM otel/opentelemetry-collector-contrib:0.123.0

# Copy the configuration file
COPY ./deployment/otel-collector/otel-collector-config.yaml /etc/otelcol-contrib/config.yaml

# Expose the port
EXPOSE 4317

ENTRYPOINT ["/otelcol-contrib"]

# Command to run the collector with the config file
CMD ["--config=/etc/otelcol-contrib/config.yaml"]
