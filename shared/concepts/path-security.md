# Concept: Path Security

## ปัญหาของ Path ที่มาจากภายนอก

เวลาโปรแกรมรับ path จาก user, config, หรือ environment variable แล้วเปิดไฟล์ตาม path นั้น — มีช่องโหว่หลักสองประเภท:

1. **Symlink attack** — path ชี้ไปไฟล์ปกติ แต่จริงๆ เป็น symlink ชี้ไปที่อื่น
2. **Path traversal** — path มี `../` ที่ทำให้หลุดออกจาก directory ที่ตั้งใจ

---

## Symlink คืออะไร

**Symlink (symbolic link)** คือไฟล์พิเศษที่เป็นแค่ pointer ชี้ไปยัง path อื่น

```bash
ln -s /proc/1/environ /tmp/config.json
# /tmp/config.json ตอนนี้ชี้ไป /proc/1/environ (env vars ของ init process)
```

เวลาโปรแกรมเรียก `os.Open("/tmp/config.json")` — kernel จะ **follow** symlink และเปิด `/proc/1/environ` แทน โดยโปรแกรมไม่รู้ตัว

**Real-world incident:**
CI system อ่าน config จาก path ที่ user ส่งมา ผู้ใช้สร้าง symlink ชี้ไป `/proc/1/environ` — service โหลดเข้า memory แล้วส่งกลับใน error message ทำให้ environment variables หลุดออกไป

---

## `stat` vs `lstat` — ต่างกันยังไง

| function | ถ้า path เป็น symlink |
|---|---|
| `os.Stat(path)` | **follow** → return info ของ target |
| `os.Lstat(path)` | **ไม่ follow** → return info ของ symlink ตัวเอง |

```go
info, err := os.Lstat(path)
if err != nil {
    return nil, err
}
if info.Mode().Type() == os.ModeSymlink {
    return nil, ErrSymlink
}
```

`Lstat` ให้ตรวจสอบได้ว่า path เป็น symlink ก่อนที่จะเปิด แต่มีปัญหาคือ TOCTOU (ดูด้านล่าง)

---

## TOCTOU Race Condition

**TOCTOU = Time-of-Check to Time-of-Use**

เป็น race condition ที่เกิดระหว่างการ "ตรวจสอบ" กับ "ใช้งาน" — มีช่วงเวลาที่ state อาจเปลี่ยนได้

```
Lstat(path)           →   [attacker swaps file → symlink]   →   Open(path)
  ↑ check ตรงนี้                 ↑ race window                    ↑ use ตรงนี้
```

ถ้าใช้ `Lstat` แล้ว `Open` แยกกัน — attacker มีช่วงสั้นๆ ที่สลับไฟล์ได้

TOCTOU ไม่ใช่แค่เรื่อง symlink — เกิดกับ size check ด้วย:

```go
stat, _ := os.Stat(path)       // check: ขนาด 100 KB
// [ไฟล์โตขึ้นเป็น 10 GB ตรงนี้]
io.ReadAll(file)               // use: โหลดทั้งก้อนเข้า RAM
```

---

## `O_NOFOLLOW` — วิธีปิด Race Window

`O_NOFOLLOW` คือ flag ที่บอก kernel ว่า "ถ้า path เป็น symlink ให้ return error ทันที อย่า follow"

การตรวจสอบและเปิดไฟล์เป็น **atomic operation เดียว** ไม่มี race window

```go
import "syscall"

file, err := os.OpenFile(path, os.O_RDONLY|syscall.O_NOFOLLOW, 0)
if err != nil {
    if errors.Is(err, syscall.ELOOP) {
        return nil, fmt.Errorf("ReadConfig %s: %w", path, ErrSymlink)
    }
    return nil, fmt.Errorf("ReadConfig %s: %w", path, err)
}
defer file.Close()
```

เมื่อใช้ `O_NOFOLLOW` แล้ว `Lstat` check ก่อนหน้า **redundant** — ลบออกได้เลย

**ข้อจำกัด:** `O_NOFOLLOW` ไม่มีบน Windows — ต้องใช้ build tag:
```go
//go:build !windows
```

---

## `fstat` vs `stat` — Size Check หลัง Open

หลังจาก open file แล้ว การ check size ควรใช้ `file.Stat()` ไม่ใช่ `os.Stat(path)`:

| | ทำงานบน | TOCTOU |
|---|---|---|
| `os.Stat(path)` | path (ชื่อไฟล์) | มี — ไฟล์อาจเปลี่ยนระหว่าง Stat กับ Read |
| `file.Stat()` | fd (open file descriptor) | ไม่มี — operate บน fd ที่เปิดแล้ว |

```go
file, err := os.OpenFile(path, os.O_RDONLY|syscall.O_NOFOLLOW, 0)
// ...
info, err := file.Stat()  // ← fstat(fd) ไม่ใช่ stat(path) — ปลอดภัย
```

---

## Path Traversal

Path traversal เกิดเมื่อ user-controlled path มี `../` ที่หลุดออกจาก directory ที่ตั้งใจ:

```
base: /app/configs/
path: ../../etc/passwd
full: /app/configs/../../etc/passwd → /etc/passwd
```

**วิธีป้องกัน:**

```go
import "path/filepath"

func safePath(base, userInput string) (string, error) {
    full := filepath.Join(base, userInput)
    // filepath.Clean จัดการ .. แล้ว ตรวจว่ายังอยู่ใน base ไหม
    if !strings.HasPrefix(full, filepath.Clean(base)+string(os.PathSeparator)) {
        return "", fmt.Errorf("path traversal detected: %s", userInput)
    }
    return full, nil
}
```

---

## สรุป: เมื่อไหรควรสนใจอะไร

| สถานการณ์ | ความเสี่ยง | วิธีป้องกัน |
|---|---|---|
| path มาจาก internal config | ต่ำ | ไม่จำเป็นต้อง check |
| path มาจาก env var | กลาง | `O_NOFOLLOW` |
| path มาจาก user input | สูง | `O_NOFOLLOW` + path traversal check |
| path มาจาก untrusted source | สูงมาก | ทั้งหมด + validate allowlist |

---

## Kata ที่เกี่ยวข้อง

- `file-handling/01-read-patterns/01-whole-file-read` — symlink check + TOCTOU บน config loader
