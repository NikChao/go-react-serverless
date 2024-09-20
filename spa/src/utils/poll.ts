export function poll(fn: () => Promise<void> | void, interval: number) {
  const execute = async () => {
    try {
      await fn();
    } catch (error) {
      console.error("Polling function encountered an error:", error);
    }
  };

  return setInterval(execute, interval);
}
