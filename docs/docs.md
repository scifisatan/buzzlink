# Usage

```bash
buzzlink [OPTIONS] <file | directory>
```

## Options

| Flag      | Description                          |
| --------- | ------------------------------------ |
| `-h`      | Show help message                    |
| `-n NOTE` | Add a note to the upload             |
| `-p PASS` | Password-protect the file            |
| `--qr`    | Display QR code of the download link |

## Examples

### Basic File Upload

```bash
buzzlink image.jpg
```

### Upload a Directory

```bash
buzzlink folder/
```

### Upload with a note

```bash
buzzlink -n "Company invoice" invoice.pdf
```

### Encrypted upload

```bash
buzzlink -p "topsecret" design.sketch
```

## 🔐 Security

When using `-p`, your file is encrypted _locally_ before upload using AES encryption. BuzzLink does **not** store or log your password.

---

## 📁 Where are files hosted?

BuzzLink uploads files to:

```
https://www.buzzheavier.com/<FILENAME>
```

And generates a short link that directly downloads the file.

---

## 🧹 Cleanup

Temporary encrypted files are automatically deleted after upload.

---

## 🙋 FAQ

### ❓ What happens if I forget my password?

Encrypted files **cannot** be decrypted without the correct password. Store it securely.

### 📂 Can I share multiple files?

Currently, only one file at a time. Consider putting them in a directory and providing buzzlink with the directory.

### ⚙️ Can I change the host domain?

Right now it's hardcoded. Future versions may support config overrides.

---

## 📣 Contributing

Pull requests are welcome! Feel free to [open an issue](https://github.com/scifisatan/buzzlink/issues) or suggest new features.

---

## ❤️ Acknowledgments

Built with love and Go for the open-source community.
