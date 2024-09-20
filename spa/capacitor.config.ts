import type { CapacitorConfig } from "@capacitor/cli";

const config: CapacitorConfig = {
  appId: "com.homie.app",
  appName: "spa",
  webDir: "build",
  ios: {
    contentInset: "always",
  },
};

export default config;
