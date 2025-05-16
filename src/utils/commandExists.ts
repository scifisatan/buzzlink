export async function commandExists(command: string): Promise<boolean> {
  try {
    const checkCommand =
      Deno.build.os === "windows"
        ? new Deno.Command("where", {
            args: [command],
            stderr: "null",
            stdout: "null",
          })
        : new Deno.Command("which", {
            args: [command],
            stderr: "null",
            stdout: "null",
          });

    const { code } = await checkCommand.output();
    return code === 0;
  } catch {
    return false;
  }
}
