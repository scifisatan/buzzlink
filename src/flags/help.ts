import { cyan, yellow } from "../utils/print.ts";

export function showHelp() {
  const logo = `\n${yellow(`
  ██████╗ ██╗   ██╗███████╗███████╗██╗     ██╗███╗   ██╗██╗  ██╗
  ██╔══██╗██║   ██║╚══███╔╝╚══███╔╝██║     ██║████╗  ██║██║ ██╔╝
  ██████╔╝██║   ██║  ███╔╝   ███╔╝ ██║     ██║██╔██╗ ██║█████╔╝ 
  ██╔══██╗██║   ██║ ███╔╝   ███╔╝  ██║     ██║██║╚██╗██║██╔═██╗ 
  ██████╔╝╚██████╔╝███████╗███████╗███████╗██║██║ ╚████║██║  ██╗
  ╚═════╝  ╚═════╝ ╚══════╝╚══════╝╚══════╝╚═╝╚═╝  ╚═══╝╚═╝  ╚═╝
`)}`;

  console.log(
    `${logo}\n${cyan("Usage:")} deno task start <file> [OPTIONS]\n\n${cyan(
      "Options:"
    )}\n  -h        Show this help message\n  -n NOTE   Add a note to the upload (optional)\n  --noqr    Disable QR code display\n`
  );
  Deno.exit(0);
}
