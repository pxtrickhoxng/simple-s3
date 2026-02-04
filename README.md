# Simple-S3 CLI

A minimal CLI for managing AWS S3 buckets and objects, built with Go and Cobra.  
Supports basic bucket and object operations with optional AWS profiles.

---

## Installation

### 1. Make sure Go is installed

You need [Go](https://go.dev/dl/) installed (Go 1.20+ recommended).

Check with:

```bash
go version
```

### 2. Install the CLI using go install

```bash
go install github.com/pxtrickhoxng/simple-s3@latest
```

This automatically downloads the code, fetches dependencies, compiles a binary, and installs it in `$GOPATH/bin` (or `$HOME/go/bin` by default).

### 3. Add to PATH

Make sure `$GOPATH/bin` or `$HOME/go/bin` is in your system PATH to run `simple-s3` from anywhere:

**Linux/macOS:**

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

**Windows PowerShell:**

```powershell
$env:Path += ";$env:USERPROFILE\go\bin"
```

### 4. Verify installation

```bash
simple-s3 --help
```

---

## AWS Credentials Setup

The CLI uses the AWS SDK, so credentials must be configured. Set them up using the AWS CLI:

```bash
aws configure
```

You will be prompted for:

- **AWS Access Key ID**
- **AWS Secret Access Key**
- **Default region name** (e.g., `us-east-2`)
- **Default output format** (e.g., `json`)

**Optional:** If you use multiple profiles, specify `--profile <profile-name>` with any command.

---

## Basic Usage

All commands require specifying a bucket (`--name`) and region (`--region`), except `list`, which lists all buckets.

### Bucket Commands

| Command  | Description               | Example                                                       |
| -------- | ------------------------- | ------------------------------------------------------------- |
| `create` | Create a new bucket       | `simple-s3 bucket create --name my-bucket --region us-east-2` |
| `delete` | Delete a bucket           | `simple-s3 bucket delete --name my-bucket --region us-east-2` |
| `list`   | List all buckets          | `simple-s3 bucket list --profile work`                        |
| `info`   | Show bucket info (region) | `simple-s3 bucket info --name my-bucket --region us-east-2`   |

### Object Commands

| Command         | Description              | Example                                                                                               |
| --------------- | ------------------------ | ----------------------------------------------------------------------------------------------------- |
| `upload`        | Upload a local file      | `simple-s3 bucket upload --file ./local.txt --key remote.txt --name my-bucket --region us-east-2`     |
| `download`      | Download a file          | `simple-s3 bucket download --key remote.txt --output ./local.txt --name my-bucket --region us-east-2` |
| `list-objects`  | List objects in a bucket | `simple-s3 bucket list-objects --name my-bucket --region us-east-2 --prefix logs/`                    |
| `delete-object` | Delete an object         | `simple-s3 bucket delete-object --key remote.txt --name my-bucket --region us-east-2`                 |

---

## Notes

- `--profile` is optional. If not specified, the default AWS profile is used.
- For `upload` and `download`, paths can be relative (e.g., `./file.txt`) or absolute.
- Bucket names must be globally unique across AWS.
