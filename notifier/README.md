# Notifier

Notifier receives new events and makes sure they are published and/or delivered. It needs to attempt deliveries and also be sure that it retries on fails, pushing it to a dead-letter queue after some attempts