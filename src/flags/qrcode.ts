import qrcode from "https://deno.land/x/qrcode_terminal/mod.js";
import { printStatus, ICON, magenta, yellow } from "../utils/print.ts";

export async function showQRCode(link: string) {
  printStatus(ICON.mobile, magenta, `QR Code for mobile access:`);
  try {
    qrcode.generate(link, { small: true });
  } catch (e) {
    const msg = e instanceof Error ? e.message : String(e);
    printStatus(ICON.warning, yellow, `QR code generation failed: ${msg}`);
    Deno.exit(1);
  }
}
