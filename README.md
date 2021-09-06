# Password Manager
simple tool that uses a database file and stores your passwords in a secure encrypted format.

# Usage
## Environment Variables
- PASS_KEY: Your master encryption key
- PASS_FILE: Your database file, defaults to ~/.pass
## Commands
### Add
```bash
pass add Gmail Amirreza SecretPassword
```
### View
```bash
pass view Gmail
```
### Delete
```bash
pass delete Gmail
```

### Import
```bash
pass import [lastpass OR csv] filename
```

### Export
```bash
pass export csv filename
```