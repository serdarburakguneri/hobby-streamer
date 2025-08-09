# Kafka

Event streaming with Apache Kafka using CloudEvents 1.0.

## Topics
`raw-video-uploaded`, `analyze.job.requested`, `hls.job.requested`, `dash.job.requested`, `analyze.job.completed`, `hls.job.completed`, `dash.job.completed`.

## Consumers
`asset-manager-group`: uploads and job completions, `transcoder-group`: analysis and transcoding jobs.

## Flows
- Upload → `raw-video-uploaded` → analyze → `analyze.job.completed`.
- HLS request → `hls.job.requested` → transcode → `hls.job.completed`.
- DASH request → `dash.job.requested` → transcode → `dash.job.completed`.

## Inspect
AKHQ: `http://localhost:8086`. Kibana: `http://localhost:5601`.

For producer/consumer code, see `backend/pkg/events/`.
