---
tier: file-handling/01-read-patterns
difficulty: 1
concepts: [whole-file-read, memory-allocation, error-handling, resource-management]
---

# Kata: Whole-File Read

## Context

ทุก service มีจุดหนึ่งที่ต้อง "อ่านไฟล์ทั้งไฟล์" — เช่น อ่าน config ตอน startup, โหลด template ก่อน render, หรืออ่าน certificate/key จาก disk
ดูเหมือนโจทย์ง่าย แต่โค้ดที่เขียนตามสัญชาตญาณมักมีรูรั่วที่ไม่โชว์ตอน dev — โชว์ตอน production เท่านั้น

## The Naive Way (และทำไมมันพัง)

**วิธีที่คนมักเขียนครั้งแรก:**
เปิดไฟล์ → อ่านทีเดียวทั้งหมด → ใช้ข้อมูล
ถ้าภาษานั้นมี one-liner (`os.ReadFile`, `ioutil.ReadFile`, `File.read()`) ก็มักใช้โดยไม่คิดต่อ

**พังตอนไหน:**
- ไฟล์หายไป / path ผิด → panic หรือ silent empty string ถ้าไม่ handle error
- file descriptor ไม่ถูกปิด (ในบางภาษา / บาง pattern) → fd leak สะสมจนถึง limit (`too many open files`)
- ไฟล์ใหญ่กว่าที่คาด (เช่น log file ถูก symlink มาแทน config) → โหลดทั้งก้อนเข้า RAM → OOM

**Root cause:**
Whole-file read โหลด content ทั้งหมดเข้า memory ในครั้งเดียว ซึ่ง OK ก็ต่อเมื่อรู้ว่าไฟล์มีขนาดจำกัดแน่นอน
ปัญหาคือคนมักไม่ validate ขนาดก่อน และไม่ handle error path ครบ

## Task

เขียนฟังก์ชัน `ReadConfig(path string) ([]byte, error)` ที่:

1. รับ path ของไฟล์ config (plaintext, ขนาดไม่เกิน 1 MB)
2. คืน content ทั้งหมดเป็น `[]byte`
3. คืน error ที่อธิบายได้ว่าเกิดอะไรขึ้น ถ้าไฟล์อ่านไม่ได้

## Requirements

- ต้องปิด file descriptor ทุกกรณี (ทั้ง success และ error path)
- ต้องปฏิเสธไฟล์ที่ใหญ่กว่า 1 MB โดย return error ที่อ่านออก (ไม่ใช่ panic)
- ห้ามใช้ `os.ReadFile` หรือ `ioutil.ReadFile` โดยตรง — ต้องเปิด/อ่าน/ปิดเอง เพื่อให้เห็น lifecycle ของ fd
- Error message ต้องบอก context ได้ เช่น "`ReadConfig: open /etc/app/config.json: no such file or directory`" ไม่ใช่แค่ `"file not found"`

## Acceptance Criteria

- [ ] อ่านไฟล์ปกติได้ถูกต้อง — content ตรงกับที่อยู่ในไฟล์ byte-for-byte
- [ ] คืน error ถ้าไฟล์ไม่มีอยู่ (`os.IsNotExist` เป็น true)
- [ ] คืน error ถ้าไฟล์ขนาดเกิน 1 MB — ไม่โหลดเข้า memory ก่อน
- [ ] คืน error ถ้าไม่มีสิทธิ์อ่าน (`permission denied`)
- [ ] ไม่มี fd leak — หลังเรียกฟังก์ชัน ไม่ว่าจะ success หรือ error fd ต้องถูกปิดแล้ว
- [ ] ฟังก์ชันทนต่อการเรียกพร้อมกัน 100 ครั้ง (concurrent-safe) — ไม่มี shared mutable state

## Concepts Involved

- `fd-lifecycle` — file descriptor คืออะไร, เปิด/ปิดอย่างไร, leak มีผลอะไร → `shared/concepts/fd-lifecycle.md`
- `error-wrapping` — การ wrap error ให้ context ไม่หาย → `shared/concepts/error-wrapping.md`
- `memory-allocation` — whole-file read กับ memory tradeoff → `shared/concepts/memory-allocation.md`
