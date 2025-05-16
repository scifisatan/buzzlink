import {
  red,
  green,
  yellow,
  blue,
  magenta,
  cyan,
} from "https://deno.land/std@0.224.0/fmt/colors.ts";

export const ICON = {
  success: "✅",
  error: "❌",
  warning: "⚠️",
  info: "ℹ️",
  lock: "🔒",
  link: "🔗",
  upload: "📤",
  clipboard: "📋",
  mobile: "📱",
  key: "🔑",
  package: "📦",
};

export function printStatus(
  icon: string,
  color: (msg: string) => string,
  msg: string
) {
  console.log(color(`${icon} ${msg}`));
}

export { red, green, yellow, blue, magenta, cyan };
