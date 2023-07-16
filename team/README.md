# Profile Interactive
go tool pprof -source_path=../../miruken/miruken http://localhost:8080/debug/pprof/profile\?seconds\=10

# Profile Web
go tool pprof -source_path=../../miruken/miruken -http=:8081 http://127.0.0.1:8080/debug/pprof/profile

# Heap Interactive
go tool pprof -source_path=../../miruken/miruken  http://localhost:8080/debug/pprof/heap\?seconds\=10

# Heap Allocations Interactive
go tool pprof -alloc_objects -source_path=../../miruken/miruken http://localhost:8080/debug/pprof/heap\?seconds\=10

