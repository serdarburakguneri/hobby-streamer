# Distributed Techniques

- Outbox: job requests and domain events can be persisted to Neo4j Outbox, then dispatched to Kafka.
  - Neo4j: MATCH (o:Outbox) RETURN o ORDER BY o.createdAt DESC LIMIT 50;
  - AKHQ: topics like hls.job.requested, dash.job.requested.
- Correlation/Causation: events carry correlationId; completions include jobId for tracing.
- Idempotency: UpsertVideo is the single path for create/update; safe to reprocess.
- Optimistic concurrency: version fields on aggregates; Neo4j updates compare version.
- Process manager: pipeline state (analyze, hls, dash) stored per asset/video for UI visibility.
- Retry/backoff: wrappers around ffmpeg/ffprobe and S3 uploads.

Links: [Kafka](./kafka-architecture.md), [CDN](./cdn-proposal.md)

