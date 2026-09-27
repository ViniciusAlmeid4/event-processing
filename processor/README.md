# Processor

The processor is the responsable for understanding what the readings mean. It cheks the previous state and do a stat-skip logic if necessary, avoiding redundant data. After that it saves the data inside Mongo so the values are persisted, while also possibily posting events to the notifier if the value matches any trigger (anomaly).