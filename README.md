# consterrorreturn go linter

✅ **A custom Go analyzer to detect returning constant (sentinel) errors or wrapping constants with `%w` instead of propagating the original `err` variable.**

This linter enforces better error handling practices in Go by ensuring:

- You return the actual `err` from function calls
- You do not mistakenly wrap sentinel errors (e.g., `pkg.ErrNotFound`) with `fmt.Errorf(... %w, ...)`
- You can intentionally map errors using `if errors.Is(...)` without false positives

---

## ✨ **Why?**

In Go, it's common to accidentally:

```go
if err != nil {
    return pkg.ErrNotFound // ❌ loses original error
}
```

or:

```go
return fmt.Errorf("... %w", pkg.ErrNotFound) // ❌ misleading error chain
```

✅ This linter detects such patterns and suggests returning or wrapping the **original `err` variable** instead.

It allows valid error mapping like:

```go
if errors.Is(err, sql.ErrNoRows) {
    return nil, ErrDomainNotFound // ✅ allowed
}
```

---

## 🚀 **Installation**

### 1️⃣ **Clone this repo:**

```bash
git clone https://github.com/delarean/consterrorreturn
cd consterrorreturn
```

### 2️⃣ **Build the plugin:**

```bash
go build -buildmode=plugin -tags=plugin -o consterrorreturn.so
```

This creates a Go plugin `consterrorreturn.so`.

✅ Requires same Go version as used by your `golangci-lint` binary.

---

## 📝 **Usage with golangci-lint**

In your project:

1. Copy `consterrorreturn.so` somewhere (e.g., `tools/linter/consterrorreturn.so`)

2. Add to `.golangci.yml`:

```yaml
linters-settings:
  custom:
    consterrorreturn:
      path: ./tools/linter/consterrorreturn.so
      description: "Checks for returning constant errors instead of err"
      original-url: "https://github.com/delarean/consterrorreturn"

linters:
  enable:
    - consterrorreturn
```

✅ Now `golangci-lint run` will invoke this linter!

---

## 📐 **Example: flagged code**

❌ **BAD:**

```go
if err != nil {
    return pkg.ErrSomething
}

return fmt.Errorf("... %w", pkg.ErrSomething)
```

✅ **GOOD:**

```go
if err != nil {
    return fmt.Errorf("something failed: %w", err)
}

if errors.Is(err, sql.ErrNoRows) {
    return nil, ErrNotFound // mapping allowed when errors.Is or error.As is used
}
```

---

## 🏗️ **Development**

If you want to modify or extend this analyzer:

1. Edit `main.go`
2. Run `go mod tidy` if dependencies change
3. Rebuild plugin:

```bash
go build -buildmode=plugin -tags=plugin -o consterrorreturn.so
```

✅ Uses `golang.org/x/tools/go/analysis` API.