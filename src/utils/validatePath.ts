/**
 * Checks if the given path is a file or directory.
 * Throws an error if the path is a symlink.
 * Returns { valid: boolean}.
 */

export async function validatePath(path: string): Promise<{ valid: boolean }> {
  try {
    const stat = await Deno.stat(path);
    if (stat.isSymlink) {
      console.error("Symlinks are not allowed.");
      Deno.exit(1);
    }
    return { valid: stat.isFile || stat.isDirectory };
  } catch {
    console.error(`Error: Path '${path}' not found or not a file/directory.`);
    Deno.exit(1);
  }
}
