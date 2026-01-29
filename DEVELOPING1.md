# AiredSpace 

Airedspace is a small website for reviewing movies with audio. 

This file defines the design of the project and basic guidelines for myself. 

## Developing

### Tech Stack

- Back-end [Go](https://go.dev)
- Front-end: [Vue.js](https://vuejs.org/)
  - [TailwindCSS](https://tailwindcss.com/) for some or all styling; 
  -
- Tools: 
	- Linting: [`golangci-lint`, version 2](https://golangci-lint.run/);
	- Formatting: [`gofumpt`](https://github.com/mvdan/gofumpt) and
	  [`goimports`](https://pkg.go.dev/golang.org/x/tools/cmd/goimports),
	  both already provided by `golangci-lint`;
	
### Design 

The project design is strictly made for Go Ecosystem, so modules should be made with that in mind. We are currently adopting the [Single File System](https://medium.com/@smart_byte_labs/organize-like-a-pro-a-simple-guide-to-go-project-folder-structures-e85e9c1769c2)

```
Airedspace
|- main.go        # Application entry point
|- handler.go     # HTTP handlers
|- service.go     # Business logic
|- repository.go  # Database repository
|- config.go      # Configuration settings
|- utils.go       # Utility functions
|- go.mod         # Go module file
|- go.sum         # Go module dependency file
```
