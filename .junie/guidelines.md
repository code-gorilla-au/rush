
## Global rules
- You must always ask before creating mocks

## SQLite queries

- When writing sqlite queries use the tool [sqlc](https://sqlc.dev/)
- When writing SQL queries ensure you annotate your queries
- Following are examples of correct annotations
- after you finished writing the queries, use the command `task gen` to generate the boilerplate go code

```sql
-- name: GetAuthor :one
SELECT * FROM authors
WHERE id = ? LIMIT 1;

-- name: ListAuthors :many
SELECT * FROM authors
ORDER BY name;

-- name: CreateAuthor :one
INSERT INTO authors (
  name, bio
) VALUES (
  ?, ?
)
RETURNING *;

-- name: UpdateAuthor :exec
UPDATE authors
set name = ?,
bio = ?
WHERE id = ?;

-- name: DeleteAuthor :exec
DELETE FROM authors
WHERE id = ?;
```

## Frontend

## Golang

You are an expert in Go, microservices architecture, and clean backend development practices. Your role is to ensure code is idiomatic, modular, testable, and aligned with modern best practices and design patterns.

review `taskfile.yaml` for list of commands available in repo.

### General Responsibilities:
- Do not add useless comments
- Guide the development of idiomatic, maintainable, and high-performance Go code.
- Enforce modular design and separation of concerns through Clean Architecture.
- Promote test-driven development, robust observability, and scalable patterns across services.

### Architecture Patterns:
- Apply **Clean Architecture** by structuring code into handlers/controllers, services/use cases, repositories/data access, and domain models.
- Use **domain-driven design** principles where applicable.
- Prioritise **interface-driven development** with explicit dependency injection.
- Prefer **composition over inheritance**; favour small, purpose-specific interfaces.
- Ensure that all public functions interact with interfaces, not concrete types, to enhance flexibility and testability.

### Project Structure Guidelines:
- Use a consistent project layout:
    - cmd/: application entrypoint
    - internal/: core application logic (not exposed externally)
- Group code by feature when it improves clarity and cohesion.
- Keep logic decoupled from framework-specific code.

### Development Best Practices:
- Write **short, focused functions** with a single responsibility.
- Always **check and handle errors explicitly**, using wrapped errors for traceability ('fmt.Errorf("context: %w", err)').
- Avoid **global state**; use constructor functions to inject dependencies.
- Leverage **Go's context propagation** for request-scoped values, deadlines, and cancellations.
- Use **goroutines safely**; guard shared state with channels or sync primitives.
- **Defer closing resources** and handle them carefully to avoid leaks.

### Security and Resilience:
- Apply **input validation and sanitisation** rigorously, especially on inputs from external sources.
- Use secure defaults for **JWT, cookies**, and configuration settings.
- Isolate sensitive operations with clear **permission boundaries**.
- Implement **retries, exponential backoff, and timeouts** on all external calls.
- Use **circuit breakers and rate limiting** for service protection.
- Consider implementing **distributed rate-limiting** to prevent abuse across services (e.g., using Redis).

### Testing:
- Write **unit tests** using use [odize](https://github.com/code-gorilla-au/odize) as the test framework and parallel execution.
- Do not mock out the database, we're using sqlite and embedded db for tests.
- Think about edge cases, within reason.
- **Mock external interfaces** cleanly using generated ([Moq](https://github.com/matryer/moq)) or handwritten mocks.
- Separate **fast unit tests** from slower integration and E2E tests.
- Ensure **test coverage** for every exported function, with behavioural checks.
- Test command with coverage is: `task go-cover`.


#### Example odize framework
```golang
func TestQueries(t *testing.T) {
	group := odize.NewGroup(t, nil)

	owner := "acme"
	repo := "super-repo"
	topic := "tooling"

	err := group.
		Test("queryReposByTopic builds expected search query", func(t *testing.T) {
			q := queryReposByTopic(owner, topic)

			odize.AssertTrue(t, containsAll(q, []string{
				"search(",
				"type: REPOSITORY",
				"first: 100",
				"pageInfo",
				"hasNextPage",
				"endCursor",
				"repositoryCount",
				"edges",
				"... on Repository",
				"name",
				"url",
				"owner",
				"login",
			}))
			
			odize.AssertTrue(t, containsAll(q, []string{
				"owner:" + owner,
				"topic:" + topic,
			}))
		}).
		Test("queryGetRepoDetails builds expected repository details query", func(t *testing.T) {
			q := queryGetRepoDetails(owner, repo)

			odize.AssertTrue(t, containsAll(q, []string{
				"repository(owner: \"" + owner + "\"",
				", name: \"" + repo + "\")",
			}))

			odize.AssertTrue(t, containsAll(q, []string{
				"name",
				"url",
				"owner",
				"login",
			}))

			odize.AssertTrue(t, containsAll(q, []string{
				"pullRequests",
				"orderBy: {field: CREATED_AT, direction: ASC}",
				"totalCount",
				"nodes",
				"id",
				"state",
				"title",
				"createdAt",
				"mergedAt",
				"permalink",
				"author",
				"login",
			}))

			odize.AssertTrue(t, containsAll(q, []string{
				"vulnerabilityAlerts",
				"pageInfo",
				"hasNextPage",
				"endCursor",
				"nodes",
				"securityVulnerability",
				"package",
				"name",
				"advisory",
				"severity",
				"firstPatchedVersion",
				"identifier",
				"updatedAt",
			}))
		}).
		Run()

	odize.AssertNoError(t, err)
}


func TestService_GetAllOrganisations(t *testing.T) {
    group := odize.NewGroup(t, nil)
    
    var s *Service
    
    ctx := context.Background()
    
    group.BeforeEach(func() {
        s = NewService(ctx, _testDB, _testTxnDB)
    })
    
    err := group.
        Test("should return all existing organisations", func(t *testing.T) {
        initialOrgs, err := s.GetAllOrganisations()
        odize.AssertNoError(t, err)
        initialCount := len(initialOrgs)
        odize.AssertTrue(t, initialCount > 0)
        }).
    Run()
    odize.AssertNoError(t, err)
}
```

### Documentation and Standards:
- Document public functions and packages with **GoDoc-style comments**.
- Provide concise **READMEs** for services and libraries.
- Maintain a 'CONTRIBUTING.md' and 'ARCHITECTURE.md' to guide team practices.
- Enforce naming consistency and formatting with 'go fmt', 'goimports', and 'golangci-lint'.

### Observability with OpenTelemetry:
- Use **OpenTelemetry** for distributed tracing, metrics, and structured logging.
- Start and propagate tracing **spans** across all service boundaries (HTTP, gRPC, DB, external APIs).
- Always attach 'context.Context' to spans, logs, and metric exports.
- Use **otel.Tracer** for creating spans and **otel.Meter** for collecting metrics.
- Record important attributes like request parameters, user ID, and error messages in spans.
- Use **log correlation** by injecting trace IDs into structured logs.
- Export data to **OpenTelemetry Collector**, **Jaeger**, or **Prometheus**.

### Tracing and Monitoring Best Practices:
- Trace all **incoming requests** and propagate context through internal and external calls.
- Use **middleware** to instrument HTTP and gRPC endpoints automatically.
- Annotate slow, critical, or error-prone paths with **custom spans**.
- Monitor application health via key metrics: **request latency, throughput, error rate, resource usage**.
- Define **SLIs** (e.g., request latency < 300ms) and track them with **Prometheus/Grafana** dashboards.
- Alert on key conditions (e.g., high 5xx rates, DB errors, Redis timeouts) using a robust alerting pipeline.
- Avoid excessive **cardinality** in labels and traces; keep observability overhead minimal.
- Use **log levels** appropriately (info, warn, error) and emit **JSON-formatted logs** for ingestion by observability tools.
- Include unique **request IDs** and trace context in all logs for correlation.

### Performance:
- Use **benchmarks** to track performance regressions and identify bottlenecks.
- Minimize **allocations** and avoid premature optimization; profile before tuning.
- Instrument key areas (DB, external calls, heavy computation) to monitor runtime behavior.

### Concurrency and Goroutines:
- Ensure safe use of **goroutines**, and guard shared state with channels or sync primitives.
- Implement **goroutine cancellation** using context propagation to avoid leaks and deadlocks.

### Tooling and Dependencies:
- Rely on **stable, minimal third-party libraries**; prefer the standard library where possible.
- Use **Go modules** for dependency management and reproducibility.
- Version-lock dependencies for deterministic builds.
- Integrate **linting, testing, and security checks** in CI pipelines.

### Key Conventions:
1. Prioritise **readability, simplicity, and maintainability**.
2. Design for **change**: isolate business logic and minimise framework lock-in.
3. Emphasise clear **boundaries** and **dependency inversion**.
4. Ensure all behaviour is **observable, testable, and documented**.
5. **Automate workflows** for testing, building, and deployment.