import { validatePath } from "./validatePath.ts";

export interface Args {
  note: string;
  password: string;
  showHelpFlag: boolean;
  showQRFlag: boolean;
  fileArg: string;
}

function requireArg(
  flag: string,
  next: string | undefined,
  label: string
): string {
  if (next && !next.startsWith("-")) {
    return next;
  } else {
    console.error(`Error: ${flag} requires a ${label} argument.`);
    Deno.exit(1);
  }
}

export async function parseArgs(): Promise<Args> {
  //default behavior
  const config: Args = {
    note: "",
    password: "",
    showQRFlag: true,
    showHelpFlag: false,
    fileArg: "",
  };

  const args = Deno.args;

  // Show help if requested or no arguments
  if (args.includes("-h") || args.includes("--help") || args.length === 0) {
    config.showHelpFlag = true;
    return config;
  }

  config.fileArg = args[0];

  // Validate file/dir path exists before parsing flags
  await validatePath(config.fileArg); // stops the process if not valid path/directory

  for (let i = 1; i < args.length; i++) {
    const arg = args[i];
    if (arg === "-n") {
      config.note = requireArg("-n", args[i + 1], "note");
      i++;
    } else if (arg === "-p") {
      config.password = requireArg("-p", args[i + 1], "password");
      i++;
    } else if (arg === "--noqr") {
      config.showQRFlag = false;
    }
  }

  return config;
}
