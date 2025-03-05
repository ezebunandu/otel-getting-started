# Instrumenting a Go WebService with Open Telemetry

This repo demos how to instrument a Go web service using the OpenTelemetry SDK for Go. It is based on the [sample dice application](https://opentelemetry.io/docs/languages/go/getting-started/) on the OpenTelemetry getting started documentation.

As an addition, I have packaged the service to run as a container, with the orchestration of the container for the service and a Prometheus instance to (hopefully) collect data from Otel (not sure how exactly this bit should fit together).
