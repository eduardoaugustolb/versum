# @versum/logger

Colored console logger with sensitivity handling for Bun.

## Install

```bash
bun add @versum/logger
```

## Usage

```ts
import { logger, info, warn, error, debug, success, sensitive } from "@versum/logger";

// Convenience methods
info("User logged in");
warn("Memory usage high");
error("Failed to connect");
debug("Request details:", { url: "/api" });
success("Operation completed");

// Full logger API
logger("info", "User logged in");
logger({ level: "error", color: "red" }, "Something went wrong");

// Sensitive values - masked in production, visible in dev+debug
const password = "secret123";
info("Password:", sensitive(password));

// Multiple sensitive values
info("Credentials:", {
  user: sensitive("admin"),
  pass: sensitive("s3cr3t"),
});
```

## API

### Convenience Methods

```ts
info(...args)    // level: "info"
warn(...args)   // level: "warn"
error(...args)  // level: "error"
debug(...args) // level: "debug"
success(...args) // level: "success"
```

### `logger(config, ...args)`

**Config:**
- `level`: `"info" | "warn" | "error" | "debug" | "success" | "log"` (default: `"log"`)
- `color?`: Custom color key or ANSI code
- `icon?`: Custom icon

**Args:** Any values. Sensitive values wrapped with `sensitive()` are masked automatically.

### `sensitive<T>(value: T, shouldMask?: boolean): SensitiveWrapper`

Wraps a value to be masked in production. In development with `DEBUG=1`, values are visible.

```ts
// String
info("Token:", sensitive("abc123"));

// Object - all keys masked
info("User data:", sensitive({ name: "John", apiKey: "key" }));

// Array - each element processed
info("Secrets:", sensitive(["a", "b", "c"]));
```

### `setEnv(handlers)`

Override environment detection handlers (for testing or other runtimes).

```ts
import { setEnv } from "@versum/logger";

setEnv({
  isDev: () => process.env.NODE_ENV !== "production",
  isDebug: () => process.env.VERBOSE === "1",
});
```

### `resetEnv()`

Reset to default Bun handlers.

```ts
import { resetEnv } from "@versum/logger";
resetEnv();
```

## Environment Variables

| Variable | Description |
|----------|-------------|
| `BUN_ENV` | Set to `"production"` or `"staging"` to disable dev mode |
| `DEBUG` | Set to any value to enable debug mode (shows sensitive values) |

## Behavior

- **Dev mode** (default when `BUN_ENV` ≠ production/staging):
  - All log levels output
  - `sensitive()` values visible when `DEBUG` is set

- **Production mode** (`BUN_ENV=production`):
  - `debug` level suppressed
  - Sensitive values masked as `ab****23` or `****`

- **Invalid inputs**: Logs warning and uses fallback values

## TypeScript

Full type support out of the box. Types exported:
- `LogLevel`
- `LoggerConfig`
- `LoggerConfigObject`
- `LoggerConfigLevel`

## License

MIT
