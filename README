# loadws-balancer

**loadws-balancer** is a lightweight Go-based load balancer designed to distribute WebSocket and HTTP traffic intelligently across multiple backend servers. It periodically checks the health and performance status of the registered servers, ensuring optimal traffic distribution.

## 🔧 Configuration

The application requires a `config.yaml` file mounted into the container or present at the path `/app/data/config.yaml`.

### Example `config.yaml`

```yaml
port: 8080
statusInterval: 5000 # ms
servers:
  - URL: "http://localhost:8082"
    statusURL: "http://localhost:8082/api/statusz"
    isWsCandidate: true
  - URL: "http://localhost:8083"
    statusURL: "http://localhost:8083/api/statusz"
```

### Configuration Fields

| Field            | Description                                                              |
| ---------------- | ------------------------------------------------------------------------ |
| `port`           | The port the load balancer listens on.                                   |
| `statusInterval` | How often (in milliseconds) the servers are polled for CPU/memory stats. |
| `servers`        | List of backend servers with their main URL and status endpoint.         |

---

## 📁 Folder & Volume Requirements

The application expects the following structure inside the container or working directory:

```
/app
└── cfg
    └── config.yaml
```

To run in Docker with volume bindings, you can use:

```bash
docker run -p 8080:8080 \
  -v $(pwd)/data:/app/data \
  loadws-balancer
```

Make sure `data/config.yaml` exists before starting the container.

## 🚀 Features

- **Dynamic traffic distribution** using health/status checks.
- **WebSocket-ready** load balancing (if `isWsCandidate` is set).
- **Extremely lightweight** (Docker image \~25MB).

## 🧪 Sample Status Endpoint Implementations

To enable the load balancer to monitor your backend servers, each one must expose an endpoint specified in the config, that returns CPU and memory usage in JSON format like:

```json
{
  "cpu": 12.5,
  "memory": 45.3
}
```

### Don't forget to specify endpoint in `config.yaml`!

---

### 🔵 Go (Golang)

#### `main.go`

```go
package main

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/shirou/gopsutil/v3/process"
)

type Status struct {
	CPU    float64 `json:"cpu"`
	Memory float32 `json:"memory"`
}

func main() {
	http.HandleFunc("/api/statusz", func(w http.ResponseWriter, r *http.Request) {
		pid := os.Getpid()
		proc, err := process.NewProcess(int32(pid))
		if err != nil {
			http.Error(w, "failed to get process info", http.StatusInternalServerError)
			return
		}

		cpu, err := proc.CPUPercent()
		if err != nil {
			http.Error(w, "failed to get CPU usage", http.StatusInternalServerError)
			return
		}

		mem, err := proc.MemoryPercent()
		if err != nil {
			http.Error(w, "failed to get memory usage", http.StatusInternalServerError)
			return
		}

		status := Status{
			CPU:    cpu,
			Memory: mem,
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(status); err != nil {
			http.Error(w, "failed to encode JSON", http.StatusInternalServerError)
			return
		}
	})

	http.ListenAndServe(":8082", nil)
}
```

#### Install dependencies

```bash
go get github.com/shirou/gopsutil/v3/process
```

---

### 🔷 Node.js Example

#### `server.js`

```js
const express = require("express");
const pidusage = require("pidusage");
const os = require("os");

const app = express();
const port = 8082;

app.get("/api/statusz", async (req, res) => {
  try {
    const stats = await pidusage(process.pid);
    res.json({
      cpu: stats.cpu, // CPU usage %
      memory: (stats.memory / os.totalmem()) * 100, // Memory usage in %
    });
  } catch (err) {
    res.status(500).send("error");
  }
});

app.listen(port, () => {
  console.log(`Server running on port ${port}`);
});
```

#### Install dependencies

```bash
npm install express pidusage
```

---

### 🟣 C# (.NET Core)

#### `StatusController.cs`

```csharp
using Microsoft.AspNetCore.Mvc;
using System.Diagnostics;

[Route("api/statusz")]
[ApiController]
public class StatusController : ControllerBase
{
    [HttpGet]
    public IActionResult GetStatus()
    {
        var proc = Process.GetCurrentProcess();
        var cpu = 0.0; // Requires performance counter setup for accuracy
        var mem = proc.WorkingSet64 / (double)GC.GetTotalMemory(false) * 100;

        return Ok(new {
            cpu = cpu,
            memory = mem
        });
    }
}
```

#### Required in `.csproj`

```xml
<PackageReference Include="Microsoft.AspNetCore.App" />
```

_Note: Accurate CPU usage in C# requires performance counters or sampling._

---

### 🟡 Java (Spring Boot)

#### `StatusController.java`

```java
@RestController
@RequestMapping("/api")
public class StatusController {

    @GetMapping("/statusz")
    public Map<String, Object> getStatus() {
        OperatingSystemMXBean osBean = ManagementFactory.getPlatformMXBean(OperatingSystemMXBean.class);
        double cpuLoad = osBean.getSystemCpuLoad() * 100;
        double memUsage = (1 - (double) Runtime.getRuntime().freeMemory() / Runtime.getRuntime().totalMemory()) * 100;

        Map<String, Object> status = new HashMap<>();
        status.put("cpu", cpuLoad);
        status.put("memory", memUsage);
        return status;
    }
}
```

#### Dependencies (`pom.xml`)

```xml
<dependency>
    <groupId>org.springframework.boot</groupId>
    <artifactId>spring-boot-starter-web</artifactId>
</dependency>
```

---

# Examples for: node.js, C# and Java wasn't tested!
