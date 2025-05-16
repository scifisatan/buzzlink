import { printStatus, ICON, blue, green, red } from "./print.ts";

export async function uploadFile(uploadFile: string, uploadUrl: string) {
  printStatus(ICON.upload, blue, `Uploading ${uploadFile.split("/").pop()}...`);
  const fileData = await Deno.readFile(uploadFile);
  const resp = await fetch(uploadUrl, {
    method: "PUT",
    body: fileData,
  });
  const respText = await resp.text();
  let code = resp.status;
  let id = "";
  try {
    const json = JSON.parse(respText);
    code = json.code || code;
    if (code === 201 && json.data && json.data.id) {
      id = json.data.id;
    } else if (json.id) {
      id = json.id;
    }
  } catch (_e) {
    // Not JSON, leave id as ""
  }
  if (code !== 201 || !id) {
    printStatus(ICON.error, red, `Upload failed. Server response: ${respText}`);
    Deno.exit(1);
  }
  printStatus(ICON.success, green, `Uploaded successfully!`);
  return id;
}
