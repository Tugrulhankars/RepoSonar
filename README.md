<div align="center">
  <h1>📡 RepoSonar</h1>
  <p>
    <strong>High-Performance GitHub Search & Observability Engine</strong>
  </p>
</div>

---

A high-performance GitHub repository search, instant filtering, and distributed tracing (observability) engine built with **Elasticsearch**, **Go**, **OpenTelemetry**, and **React/TypeScript**.

---

## 🚀 Key Features

* **Intelligent Full-Text Search:** Weighted (`boosted`) and fault-tolerant (`fuzzy`) search across repository names, descriptions, and topics using Elasticsearch.
* **Debounced Autocomplete:** Prefix-supported search input delivering real-time suggestions without overwhelming backend resources.
* **Granular Filtering & Pagination:** Server-side pagination and real-time filtering powered by Elasticsearch Range and Term queries based on programming language (`language`) and minimum stars (`stargazers_count`).
* **End-to-End Distributed Tracing:** Full request lifecycle telemetry from the HTTP layer down to individual Elasticsearch queries exported to Jaeger via OpenTelemetry.
* **Modern User Interface:** Highly responsive and reactive client powered by React, TypeScript, Redux Toolkit, and Axios.

---

## 🛠️ Architecture & Tech Stack

| Layer | Technology | Description |
| :--- | :--- | :--- |
| **Backend** | Go (Golang) | `net/http` based RESTful API & synchronization engine |
| **Search Engine** | Elasticsearch 8+/9+ | Custom mapping and query execution via the Go Typed Client |
| **Frontend** | React, TypeScript, Redux Toolkit | Modular state management, Axios client, and custom debounce hook |
| **Observability** | OpenTelemetry, Jaeger | OTLP gRPC telemetry collection for spans and traces |

---

## ⚙️ Getting Started

### 1. Start Infrastructure Services (Docker)

Run the Jaeger and Elasticsearch containers:

```bash
# Jaeger (Distributed Tracing UI)
docker run -d --name jaeger \
  -e COLLECTOR_OTLP_ENABLED=true \
  -p 16686:16686 \
  -p 4317:4317 \
  jaegertracing/all-in-one:latest

# Elasticsearch (Default port: 9200)
docker run -d --name elasticsearch \
  -p 9200:9200 \
  -e "discovery.type=single-node" \
  -e "xpack.security.enabled=false" \
  docker.elastic.co/elasticsearch/elasticsearch:8.11.0