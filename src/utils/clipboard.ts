import { commandExists } from "./commandExists.ts";
import { OperatingSystem } from "../types.ts";

interface ClipboardInterface {
  readText(): Promise<string>;
  writeText(data: string): Promise<void>;
}

type ClipboardMap = {
  [key in OperatingSystem]?: ClipboardImplementation;
};

const encoder = new TextEncoder();
const decoder = new TextDecoder();

class ClipboardImplementation implements ClipboardInterface {
  os: OperatingSystem;
  readCommand: string[];
  writeCommand: string[];
  postProcess?: (data: string) => string;

  constructor(
    os: OperatingSystem,
    readCommand: string[],
    writeCommand: string[],
    postProcess?: (data: string) => string
  ) {
    this.os = os;
    this.readCommand = readCommand;
    this.writeCommand = writeCommand;
    this.postProcess = postProcess;
  }

  async readText(): Promise<string> {
    const command = new Deno.Command(this.readCommand[0], {
      args: this.readCommand.slice(1),
      stdout: "piped",
      stderr: "piped",
    });

    try {
      const { stdout } = await command.output();
      let text = decoder.decode(stdout);

      if (this.postProcess) {
        text = this.postProcess(text);
      }

      return text;
    } catch (error) {
      if (error instanceof Error) {
        throw new Error(`Failed to read from clipboard: ${error.message}`);
      } else {
        throw new Error(`Failed to read from clipboard: ${String(error)}`);
      }
    }
  }

  async writeText(data: string): Promise<void> {
    const command = new Deno.Command(this.writeCommand[0], {
      args: this.writeCommand.slice(1),
      stdin: "piped",
      stderr: "piped",
    });

    try {
      const child = command.spawn();
      const writer = child.stdin.getWriter();
      await writer.write(encoder.encode(data));
      await writer.close();
      await child.status;
    } catch (error) {
      if (error instanceof Error) {
        throw new Error(`Failed to read from clipboard: ${error.message}`);
      } else {
        throw new Error(`Failed to read from clipboard: ${String(error)}`);
      }
    }
  }
}

// Define platform-specific implementations
async function getImplementations(): Promise<ClipboardMap> {
  return {
    linux: new ClipboardImplementation(
      "linux",
      (await commandExists("xsel"))
        ? ["xsel", "-b", "-o"]
        : ["xclip", "-selection", "clipboard", "-o"],
      (await commandExists("xsel"))
        ? ["xsel", "-b", "-i"]
        : ["xclip", "-selection", "clipboard", "-i"]
    ),
    darwin: new ClipboardImplementation("darwin", ["pbpaste"], ["pbcopy"]),
    windows: new ClipboardImplementation(
      "windows",
      ["powershell", "-noprofile", "-command", "Get-Clipboard"],
      ["powershell", "-noprofile", "-command", "$input|Set-Clipboard"],
      (data) => data.replace(/\r/g, "").replace(/\n$/, "")
    ),
  };
}

let clipboard: Clipboard;
let readText: Clipboard["readText"];
let writeText: Clipboard["writeText"];

(async () => {
  const implementations = await getImplementations();
  if (!implementations[Deno.build.os as OperatingSystem]) {
    console.warn(`Clipboard support not available for ${Deno.build.os}`);
  }
  clipboard = new Clipboard(Deno.build.os as OperatingSystem, implementations);
  readText = clipboard.readText.bind(clipboard);
  writeText = clipboard.writeText.bind(clipboard);
})();

class Clipboard implements ClipboardInterface {
  private impl: ClipboardImplementation;

  constructor(os: OperatingSystem, implementations: ClipboardMap) {
    const impl = implementations[os];
    if (!impl) {
      throw new Error(`Clipboard: unsupported OS: ${os}`);
    }
    this.impl = impl;
  }

  async readText(): Promise<string> {
    return await this.impl.readText();
  }

  async writeText(data: string): Promise<void> {
    return await this.impl.writeText(data);
  }
}

export { clipboard, readText, writeText };
