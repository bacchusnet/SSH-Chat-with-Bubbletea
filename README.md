# SSH-Chat-with-Bubbletea
## 🚀 Project Roadmap & Learning Objectives

This project serves as a sandbox for mastering Go, networking, and system security. Below are the planned features categorized by engineering domain.

---

### 🛡️ Security (Identity & Access)
* **Public Key Authentication:** Move away from passwords to purely cryptographic identity verification.
    * *Learning:* Parsing `.pub` keys and implementing `ssh.PublicKeyHandler`.
* **Rate Limiting & Brute Force Protection:** Implement a "Muzzle" system using a Token Bucket algorithm to auto-block IPs that spam or attempt rapid-fire connections.
    * *Learning:* Middleware design and `time.Ticker` logic.



---

### 🌐 Networking (Routing & State)
* **Multi-Room "MUD" Architecture:** Transition from a single global chat to a grid of "Rooms" with isolated broadcast channels.
    * *Learning:* Managing complex state across concurrent goroutines using `sync.Map`.
* **Encrypted Webhooks:** Integrate a Go HTTP server into the binary to allow external services to "Post" alerts into the chat.
    * *Learning:* Running multiple servers (SSH/HTTP) concurrently via `errgroup`.

---

### 📊 Observability & Monitoring
* **Prometheus Metrics:** Export real-time data (active connections, messages per second, memory usage) to a `/metrics` endpoint.
    * *Learning:* Time-series data visualization and the `prometheus/client_golang` library.
* **Health Checks:** Implement a "Heartbeat" mechanism to detect deadlocks or corrupted internal state.



---

### 📜 Logging & Auditing
* **Structured Logging (JSON):** Replace standard print statements with `slog` to produce machine-readable JSON logs for ingestion into ELK/Splunk.
    * *Example:* `{"level":"info", "event":"user_joined", "user":"bob", "ip":"1.2.3.4"}`
* **Session Recording:** Implement an "audit trail" that records raw terminal sessions for security playback.
    * *Learning:* I/O buffering and non-blocking filesystem writes.

---

### 🏗️ Infrastructure & Reliability
* **Graceful Shutdown:** Implement signal handling (`SIGTERM`) to notify users of maintenance and save state before the process exits.
    * *Learning:* Unix signal handling and clean resource disposal.
* **Automated Host Key Management:** Logic to check for, load, or generate persistent `id_ed25519` keys on startup to ensure a consistent host identity.