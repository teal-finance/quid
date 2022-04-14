// The current runtime environement types
enum EnvType {
  testRun = "testRun",
  local = "local",
  production = "production"
}

// Get the current runtime environement type
function getEnv(): string {
  switch (import.meta.env.MODE) {
    case "test":
      return EnvType.testRun;
    case "development":
      return EnvType.local;
    default:
      switch (window.location.hostname) {
        case "localhost":
          // if running localy without node
          return EnvType.local;
        default:
          return EnvType.production;
      }
  }
}

export { EnvType, getEnv };