# exporter-idrac

- Using git
  - git branch feature/readme => create new branch
  - git checkout feature/readme => switch to feature/readme branch
  - git branch => check current working branch
  - git branch delete feature/readme => delete branch feature/readme(local PC), you should `git push origin --delete feature/readme` to delete branch on github.com
  - git fetch => to synchronize your local branch with branch on github.com
  - git checkout main && git merge feature/readme => first change current branch to main and then do merge feature/readme branch to main branch
  - git push --set-upstream origin feature/readme => create feature/readme on github.com
- How to Statically compile GOLANG programs - <https://oddcode.daveamit.com/2018/08/16/statically-compile-golang-binary/>
- Follow code with ultimate service which is taught by Jacob Walker
  - <https://github.com/ardanlabs/service/tree/32c75246b11b871ca7aaf07eebb2b1ccef6ee81c>
- <https://developer.aliyun.com/article/683180>
- <https://topic.alibabacloud.com/a/write-prometheus-exporter-using-golang_8_8_10262688.html>
- Folder structure

  ```bash
  ├── cmd
  │   └── idrac
  │       └── controller.go
  ├── config
  ├── migrate
  ├── go.mod
  ├── internal
  ├── main.go
  ├── Makefile
  └── README.md
  ```

  - The **migrations** directory will contain the SQL migration files for our database.
  - The **internal** directory will contain various ancillary packages used by our API. It will contain the code for interacting with our database, doing data validation, sending emails and so on. Basically, any code which isn't application-specific and can potentially be reused will live in here. Our Go code under cmd/idrac will import the packages in the internal directory (but never the other way around).
  - The **Makefile** will contain recipes for automating common administrative tasks — like auditing our Go code, building binaries, and executing database migrations.
  - The **config** directory will contain the configuration files and setup scripts for our production server

## How to Use

- Content of Makefile

  ```bash
  run:
    go run ./cmd/idrac


  mod:
  # real tab space or get error Makefile:2: *** missing separator.  Stop.
    go mod tidy # remove unused go packages
    go mod vendor # make local copy of third party packages

  ```

- run cli:
  - make run - run main.go in idrac folder
  - make mod
    - run go mod tidy cli
    - run go mod vendor cli

### Opentelemetry

- **Specifications**: The OpenTelemetry specifications describe the requirements and expectations for all OpenTelemetry APIs, SDKs, and data protocols.
- **API**: Language-specific interfaces and implementations based on the specifications that can be used to add OpenTelemetry to an application.
- **SDK**: The concrete OpenTelemetry implementations that sit between the APIs and the Exporters, providing functionality like (for example) state tracking and batching data for transmission. An SDK also offers a number of configuration options for behaviors like request filtering and transaction sampling.
- **Exporters**: In-process SDK plug-ins that are capable of sending data to a specific destination, which may be local (such as a log file or stdout), or remote (such as Jaeger, or a commercial solution like Honeycomb or Lightstep). Exporters decouple the instrumentation from the backend, making it possible to change destinations without having to reinstrument your code.
- **Collector**: An optional, but very useful, vendor-agnostic service that can receive and process telemetry data before forwarding it to one or more destinations. It can be run either as a sidecar process alongside your application or as a standalone proxy elsewhere, providing greater flexibility for sending the application telemetry. This can be particularly useful in the kind of tightly controlled environments that are common in the enterprise.

:::image type="content" source="images/opentelemetry.png" alt-text="opentelemetry component":::

#### Tracing Concept

- **Span**: A span describes a unit of work performed by a request, such as a fork in the execution flow or hop across the network, as it propagates
through a system. Each span has an associated name, a start time, and a duration. They can be (and typically are) nested and ordered to model causal relationships.
- **Traces**: A trace represents all of the events—individually represented as spans - that make up a request as it flows through a system. A trace may be thought of as a directed acyclic graph (DAG) of spans, or more concretely as a "stack trace" in which each span represents the work done by one component.

:::image type="content" source="images/span-trace.png" alt-text="span and trace relationship":::

```bash
-  Two representations of a trace of a request as it traverses five services.
-  Resulting in five spans.
-  The full traces are visualized as a DAG (left).
-  Bar diagram (right) with a time axis illustrating start times and durations
```

#### Tracing with OpenTelemetry

- Using OpenTelemetry to instrument your code includes two phases:
  - ***Configuration***: For both tracing and metric instrumentation, the configuration phase is executed exactly once in a program, usually in the main function, and includes the following steps:
    1. The first step is to retrieve and configure the appropriate exporters for your target backends. Tracing exporters implement the
SpanExporter interface (which in OpenTelemetry v0.17.0 is located in the go.opentelemetry.io/otel/sdk/export/trace package, often aliased to export). **We'll discuss this in more detail in "Obtaining a tracer"**.
    2. Before instrumenting your code for tracing, the exporters and any other appropriate configuration options—are passed to the SDK to create the "tracer provider," which will serve as the main entry point for the OpenTelemetry tracing API for the lifetime of your program. **We'll discuss this in more detail in "Starting and ending spans"**.
    3. Once you've created your tracer provider, it's a good practice to set it as your "global" tracer provider. As we'll see in "Setting the global tracer provider", this makes it discoverable via the otel.GetTracerProvider function, which allows libraries and other dependencies that also use the OpenTelemetry API to more easily discover the SDK and emit telemetry data. **We'll discuss this in more detail in "Setting span metadata"**.
  - ***Instrumentation***:
    1. Before you can instrument an operation, you first have to obtain a Tracer, which has the central role of keeping track of trace and
span information, from the (usually global) tracer provider
    2. Once you have a handle to your Tracer you can use it to create and start the Span value that is the actual value that you'll use to instrument your code.
    3. Finally, you can also choose to add metadata to your spans, including human-readable, timestamped messages called events,
and key/value pairs called attributes.
- <https://github.com/cloud-native-go/examples/tree/main/ch11/otel-trace>
