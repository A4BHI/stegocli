# Stego CLI

A Go-based CLI utility for securely embedding and extracting encrypted files or raw text inside PNG images using LSB steganography.

---

## Features

- Embed files into PNG images
- Embed raw text directly from the CLI
- Extract hidden files or text from encoded images
- AES-256 GCM encryption
- Gzip compression pipeline
- PNG LSB payload embedding/extraction
- File metadata and permission preservation
- Modular Cobra-based command architecture

---

## Installation

### Clone the repository

```bash
git clone https://github.com/a4bhi/stegocli.git

cd stegocli
```

### Build the binary

```bash
go build
```

---

## Usage

### Embed a File

```bash
./stego encode \
  -i image.png \
  -f secret.txt \
  -o encoded.png \
  -p password
```

### Embed Raw Text

```bash
./stego encode \
  -i image.png \
  -t "hidden message" \
  -o encoded.png \
  -p password
```

### Decode Hidden Payload

```bash
./stego decode \
  -i encoded.png \
  -p password
```

---

## Example Output

```text
✔ Decryption completed.
✔ Decompression completed.
✔ Secret extracted successfully.

[ Secret Data ]
➜ hidden message
```

---

## How It Works

The tool follows the pipeline below:

```text
compress → encrypt → serialize metadata → embed into PNG LSBs
```

During extraction:

```text
extract bits → deserialize metadata → decrypt → decompress
```

## Technologies Used

- Go
- Cobra CLI
- AES-GCM
- Gzip
- PNG image processing

---

## Notes

- Only PNG images are currently supported
- Built for educational and research purposes
