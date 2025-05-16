export function getUploadUrl(filename: string, note?: string): string {
  /**
   * Returns the upload URL for a given filename and optional note.
   * @param filename The file name to upload.
   * @param note Optional note to append as a query parameter.
   * @returns The upload URL string.
   */

  let uploadUrl = `https://w.buzzheavier.com/${filename}`;
  if (note) uploadUrl += `?note=${encodeURIComponent(note)}`;
  return uploadUrl;
}

export function getFileUrl(id: string): string {
  const fileUrl = `https://buzzheavier.com/${id}`;
  return fileUrl;
}
