# consterrorreturn go linter

✅ **A custom Go analyzer to detect returning constant (sentinel) errors instead of propagating the original `err` variable.**

This linter enforces better error handling practices in Go by ensuring that you don't lose the original error context.

It identifies three common anti-patterns:
1.  Returning a constant error, which discards the original error.
2.  Returning a new error created with `errors.New`, which also discards the original error.
3.  Returning an error created with `fmt.Errorf` that wraps a constant error instead of the original `err`.

---

## ✨ **Why?**

In Go, it's a common mistake to lose the original error, which makes debugging harder.

For example:
```go
// Case 1: Returning a constant error
if err != nil {
    return pkg.ErrNotFound // ❌ loses original error
}

// Case 2: Returning a new error
if err != nil {
    return errors.New("failed to marshal event") // ❌ loses original error
}

// Case 3: Misleading error chain
err := runCommand()
return fmt.Errorf("... %w", pkg.ErrNotFound) // ❌ misleading error chain and loses original error
```

✅ This linter detects such patterns and suggests returning or wrapping the **original `err` variable** instead.

It allows valid error mapping like:
```go
if errors.Is(err, sql.ErrNoRows) {
    return nil, ErrDomainNotFound // ✅ allowed
}
```

---

## **Installation**

This linter is intended to be used as a plugin for `golangci-lint`. It requires building a custom `golangci-lint` binary with the linter embedded.

### 1️⃣ **Create a custom `golangci-lint` configuration:**

Create a `.custom-gcl.yml` file in your project root with the following content:
```yaml
version: v1.60.3
plugins:
  - module: "github.com/delarean/consterrorreturn"
    import: "github.com/delarean/consterrorreturn/cmd/gclplugin"
    version: latest
```

### 2️⃣ **Build the custom `golangci-lint` binary:**

Run the following command to build the custom binary:
```bash
golangci-lint custom
```
This will create a `golangci-lint-custom` binary in your current directory.

### 3️⃣ **(Optional) Add a Makefile target:**
To simplify the process, you can add a target to your `Makefile`:
```makefile
.PHONY: build-custom-lint
build-custom-lint:
	@echo "building custom golangci-lint with consterrorreturn plugin..."
	@test -s $(GOBIN)/golangci-lint-custom || (golangci-lint custom && mv golangci-lint-custom $(GOBIN)/golangci-lint-custom)
```

---

## 📝 **Usage with golangci-lint**

1.  Enable the linter in your `.golangci.yml`:
```yaml
linters:
  enable:
    - consterrorreturn

linters-settings:
  custom:
    consterrorreturn:
      type: "module"
      description: "returning sentinel (constant) error instead of propagating original err variable"
      settings:
        include-pkgs: "your/package/prefix"
```

2.  Run the custom linter:
```bash
./golangci-lint-custom run
```

---

## 📐 **Examples of flagged code**

Here are some examples of code that will be flagged by the linter.

### Case 1: Returning a constant error
```go
func f1() error {
	err := errors.New("original error")
	if err != nil {
		return ErrNotFound // want "returning sentinel"
	}
	return nil
}
```

### Case 2: Returning a new error
```go
func f2() error {
	err := errors.New("original error")
	if err != nil {
		return errors.New("failed to marshal event") // want "returning sentinel"
	}
	return nil
}
```

### Case 3: Misleading error chain
```go
func f3() error {
	err := errors.New("original error")
	if err != nil {
		return fmt.Errorf("... %w", ErrNotFound) // want "returning sentinel"
	}
	return nil
}
```

✅ The linter is built using the `golang.org/x/tools/go/analysis` API.