# Ingestor

Ingestor service is responsible for receiving the readings, validating them and publishing them to kafka for it to be processed. The Ingestor uses a HTTP interface to accept readings, being the entry point for the pipeline.