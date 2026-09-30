# Event Contract

This document defines the JSON payloads exchanged through Kafka. Each service maintains its own Go types that follow this contract; the services do not share a code package.

All fields shown below are required. Consumers should use `schemaVersion` to identify the payload version. Timestamps use RFC 3339 in UTC. UUIDs are JSON strings. The `value` field can be a JSON string, number, or boolean; `null` is not part of the contract.

## Reading

Published by the ingestor to `sensor.readings`.

| Field | Optional | JSON type | Meaning |
| :--- | :---: | :---: | --- |
| `schemaVersion` | ( ) | integer | Contract version, currently at `1`. |
| `id` | ( ) | string (UUID) | Unique ID generated for this reading event. |
| `device_id` | ( ) | string (UUID) | Identifier of device that produced the reading. |
| `metric` | ( ) | string | Name of the measured metric. |
| `value` | ( ) | string, number, or boolean | Measured value. |
| `timestamp` | ( ) | string (RFC 3339) | When the measurement was taken, in UTC. |
| `reliability` | (x) | int (0-100) | Evaluetes how much can the system trust in the reading. |

```json
{
  "schemaVersion": 1,
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "device_id": "550e8400-e29b-41d4-a716-446655440001",
  "metric": "temperature",
  "value": 21.5,
  "timestamp": "2026-09-27T12:00:00Z",
  "reliability": 85
}
```

The HTTP `POST /readings` request does not include `id` or `schemaVersion`; the ingestor is responsible for them when it creates the Kafka event.

## NotificationRequest

Published by the processor to `notification.requested` when a reading triggers a notification.

| Field | Optional | JSON type | Meaning |
| :--- | :---: | :---: | --- |
| `schemaVersion` | ( ) | integer | Contract version; currently `1`. |
| `id` | ( ) | string (UUID) | Unique ID for this notification request. |
| `reading` | ( ) | object (`Reading`) | Reading that triggered the request. |
| `type` | ( ) | string | Notification category. |
| `severity` | ( ) | string | Severity assigned by the processor. |
| `message` | ( ) | string | Human-readable description of the notification. |
| `timestamp` | ( ) | string (RFC 3339) | Time the notification condition was detected, in UTC. |

The nested `reading` object uses the same fields and meanings as `Reading` above.

```json
{
  "schemaVersion": 1,
  "id": "550e8400-e29b-41d4-a716-446655440002",
  "reading": {
    "schemaVersion": 1,
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "device_id": "550e8400-e29b-41d4-a716-446655440001",
    "metric": "temperature",
    "value": 21.5,
    "timestamp": "2026-09-27T12:00:00Z",
    "reliability": 85
  },
  "type": "threshold_exceeded",
  "severity": "warning",
  "message": "Temperature exceeded the configured threshold.",
  "timestamp": "2026-09-27T12:00:01Z"
}
```

HTTP health responses and read/query responses are service APIs, not Kafka events, and are outside this contract.