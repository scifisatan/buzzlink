// buzzlink: Deno version of buzzlink.sh
// Features: file/dir upload, note, password, QR, clipboard, colored output

import { printStatus, ICON, cyan } from "./src/utils/print.ts";

import { showHelp } from "./src/flags/help.ts";
import { writeText as copyToClipboard } from "./src/utils/clipboard.ts";
import { showQRCode } from "./src/flags/qrcode.ts";
import { uploadFile as uploadFileHttp } from "./src/utils/upload.ts";
import { parseArgs } from "./src/utils/args.ts";
import { getFileUrl, getUploadUrl } from "./src/utils/getUrl.ts";

// Parse args
const { note, password, showHelpFlag, showQRFlag, fileArg } = await parseArgs();

// Step 1: Display help command
if (showHelpFlag) showHelp();

if (password) {
  console.log("We do not support password on windows");
}

const isDir = (await Deno.stat(fileArg)).isDirectory;

if (isDir) {
  console.error("Directory upload is not supported in Windows");
  Deno.exit(1);
}

// Step 2: Upload File
const uploadFile = fileArg;

const filename = uploadFile.split("/").pop()!;
const uploadUrl = getUploadUrl(filename, note);
const id = await uploadFileHttp(uploadFile, uploadUrl);
const link = getFileUrl(id);

// Step 3: Generate QR
if (showQRFlag) {
  await showQRCode(link);
}

// Step 4: Copy to clipboard
printStatus(ICON.link, cyan, link);
await copyToClipboard(link);
